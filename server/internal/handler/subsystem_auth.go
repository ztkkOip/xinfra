package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/auth"
	"github.com/1024XEngineer/xinfra/server/internal/model"
	"github.com/1024XEngineer/xinfra/server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubsystemAuthHandler struct {
	db    *gorm.DB
	wayne *service.WayneRoleBindingService
	audit *service.AuditService
}

func NewSubsystemAuthHandler(db *gorm.DB, wayne *service.WayneRoleBindingService, audit *service.AuditService) *SubsystemAuthHandler {
	return &SubsystemAuthHandler{db: db, wayne: wayne, audit: audit}
}

func (h *SubsystemAuthHandler) ListSystems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"items": []gin.H{
			{"key": "wayne", "name": "Wayne", "enabled": true},
			{"key": "clouddm", "name": "CloudDM", "enabled": false},
		},
	})
}

func (h *SubsystemAuthHandler) ListWayneNamespaceRoles(c *gin.Context) {
	groups, err := h.wayne.ListNamespaceRoleGroups(c.Request.Context())
	if err != nil {
		writeWayneRoleBindingError(c, nil, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": groups})
}

func (h *SubsystemAuthHandler) ListWayneBusinessLineNamespaces(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}
	businessLineID, ok := parseBusinessLineID(c)
	if !ok {
		return
	}
	if !h.canManageBusinessLine(c, claims, businessLineID) {
		return
	}

	rows, ok := h.listBusinessLineWayneNamespaces(c, businessLineID)
	if !ok {
		return
	}
	operatorEmail, ok := subsystemOperatorEmail(c, claims)
	if !ok {
		return
	}

	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		item := gin.H{
			"id":            row.WayneNamespaceID,
			"name":          row.WayneNamespaceName,
			"kubeNamespace": row.KubeNamespace,
		}
		if claims.IsAdmin {
			item["permissions"] = &service.WayneOperatorPermissions{Create: true, Update: true, Delete: true}
			item["can_bind"] = true
			item["can_unbind"] = true
		} else {
			permissions, err := h.wayne.NamespaceOperatorPermissionsParsed(c.Request.Context(), row.WayneNamespaceID, operatorEmail)
			if err != nil {
				item["permission_error"] = err.Error()
			} else {
				item["permissions"] = permissions
				item["can_bind"] = permissions.Create || permissions.Update
				item["can_unbind"] = permissions.Delete
			}
		}
		items = append(items, item)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SubsystemAuthHandler) GetWayneUserRoles(c *gin.Context) {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	result, err := h.wayne.GetUserRoles(c.Request.Context(), username)
	if err != nil {
		writeWayneRoleBindingError(c, result, err)
		return
	}
	writeWayneRoleBindingResult(c, result)
}

func (h *SubsystemAuthHandler) BindWayneNamespaceRoles(c *gin.Context) {
	h.handleWayneNamespaceRoles(c, http.MethodPut)
}

func (h *SubsystemAuthHandler) UnbindWayneNamespaceRoles(c *gin.Context) {
	h.handleWayneNamespaceRoles(c, http.MethodDelete)
}

func (h *SubsystemAuthHandler) InitWayneBusinessLineUser(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}
	businessLineID, ok := parseBusinessLineID(c)
	if !ok {
		return
	}
	targetUserID, ok := parseUintPathParam(c, "userid")
	if !ok {
		return
	}
	if !h.canManageBusinessLine(c, claims, businessLineID) {
		return
	}

	var target model.User
	if err := h.db.Where("id = ? AND deleted_at IS NULL", targetUserID).First(&target).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "target user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	operatorEmail, ok := subsystemOperatorEmail(c, claims)
	if !ok {
		return
	}
	targetUsername := wayneUsernameForUser(target)
	result, ok := h.initializeWayneVisitor(c, businessLineID, targetUsername, operatorEmail, claims.IsAdmin, service.WayneRoleBindingRequest{
		RequestID: "business-line-user-init-" + strconv.FormatUint(targetUserID, 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10),
		Reason:    "初始化业务线 Wayne 访客角色",
	})
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "items": result})
}

func (h *SubsystemAuthHandler) handleWayneNamespaceRoles(c *gin.Context, method string) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}
	businessLineID, ok := parseBusinessLineID(c)
	if !ok {
		return
	}
	namespaceID, ok := parseUintPathParam(c, "namespaceid")
	if !ok {
		return
	}
	targetUsername := strings.TrimSpace(c.Param("username"))
	if targetUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	if !h.canManageBusinessLine(c, claims, businessLineID) {
		return
	}
	if !h.ensureNamespaceBelongsToBusinessLine(c, businessLineID, namespaceID) {
		return
	}
	if !h.ensureLocalUserExists(c, targetUsername) {
		return
	}
	operatorEmail, ok := subsystemOperatorEmail(c, claims)
	if !ok {
		return
	}
	if !claims.IsAdmin && !h.ensureWayneOperatorPermission(c, namespaceID, operatorEmail, method) {
		return
	}

	req, ok := parseRoleBindingRequest(c)
	if !ok {
		return
	}
	req.Username = ""
	if method == http.MethodPut && len(req.GroupIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupIds is required"})
		return
	}

	var result *service.WayneRoleBindingResult
	var err error
	if method == http.MethodPut {
		result, err = h.wayne.BindNamespace(c.Request.Context(), namespaceID, targetUsername, operatorEmail, req)
	} else {
		result, err = h.wayne.UnbindNamespace(c.Request.Context(), namespaceID, targetUsername, operatorEmail, req)
	}
	if err != nil {
		h.writeAudit(c, claims, businessLineID, namespaceID, targetUsername, method, "deny", req.RequestID, err.Error())
		writeWayneRoleBindingError(c, result, err)
		return
	}

	h.writeAudit(c, claims, businessLineID, namespaceID, targetUsername, method, "allow", req.RequestID, "")
	writeWayneRoleBindingResult(c, result)
}

func (h *SubsystemAuthHandler) canManageBusinessLine(c *gin.Context, claims *auth.Claims, businessLineID uint64) bool {
	var businessLine model.BusinessLine
	if err := h.db.First(&businessLine, "id = ?", businessLineID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business line not found"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if claims.IsAdmin {
		return true
	}

	var binding model.BusinessLineUser
	if err := h.db.Where("business_line_id = ? AND user_id = ? AND permission = ?", businessLineID, claims.UserID, 0).
		First(&binding).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "current user is not platform admin or business line admin"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (h *SubsystemAuthHandler) ensureNamespaceBelongsToBusinessLine(c *gin.Context, businessLineID, namespaceID uint64) bool {
	var row model.BusinessLineWayneNamespace
	if err := h.db.Where("business_line_id = ? AND wayne_namespace_id = ?", businessLineID, namespaceID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "wayne namespace is not bound to current business line"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (h *SubsystemAuthHandler) ensureLocalUserExists(c *gin.Context, username string) bool {
	var user model.User
	if err := h.db.Where("(username = ? OR email = ?) AND deleted_at IS NULL", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "target user not found"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (h *SubsystemAuthHandler) ensureWayneOperatorPermission(c *gin.Context, namespaceID uint64, operatorEmail string, method string) bool {
	permissions, err := h.wayne.NamespaceOperatorPermissionsParsed(c.Request.Context(), namespaceID, operatorEmail)
	if err != nil {
		writeWayneRoleBindingError(c, nil, err)
		return false
	}
	if method == http.MethodDelete {
		if permissions.Delete {
			return true
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "current user does not have Wayne namespace role delete permission"})
		return false
	}
	if permissions.Create || permissions.Update {
		return true
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "current user does not have Wayne namespace role create or update permission"})
	return false
}

func (h *SubsystemAuthHandler) listBusinessLineWayneNamespaces(c *gin.Context, businessLineID uint64) ([]model.BusinessLineWayneNamespace, bool) {
	var rows []model.BusinessLineWayneNamespace
	if err := h.db.Where("business_line_id = ?", businessLineID).Order("wayne_namespace_id ASC").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}
	return rows, true
}

func (h *SubsystemAuthHandler) initializeWayneVisitor(c *gin.Context, businessLineID uint64, targetUsername string, operatorEmail string, skipWaynePermissionCheck bool, req service.WayneRoleBindingRequest) ([]gin.H, bool) {
	namespaces, ok := h.listBusinessLineWayneNamespaces(c, businessLineID)
	if !ok {
		return nil, false
	}
	if len(namespaces) == 0 {
		return []gin.H{}, true
	}
	groupIDs, err := h.wayne.NamespaceVisitorGroupIDs(c.Request.Context())
	if err != nil {
		writeWayneRoleBindingError(c, nil, err)
		return nil, false
	}
	req.GroupIDs = groupIDs
	replace := true
	req.Replace = &replace

	items := make([]gin.H, 0, len(namespaces))
	for _, namespace := range namespaces {
		if !skipWaynePermissionCheck && !h.ensureWayneOperatorPermission(c, namespace.WayneNamespaceID, operatorEmail, http.MethodPut) {
			return nil, false
		}
		result, err := h.wayne.BindNamespace(c.Request.Context(), namespace.WayneNamespaceID, targetUsername, operatorEmail, req)
		if err != nil {
			writeWayneRoleBindingError(c, result, err)
			return nil, false
		}
		items = append(items, gin.H{
			"namespace_id":   namespace.WayneNamespaceID,
			"namespace_name": namespace.WayneNamespaceName,
			"group_ids":      groupIDs,
		})
	}
	return items, true
}

func (h *SubsystemAuthHandler) writeAudit(c *gin.Context, claims *auth.Claims, businessLineID, namespaceID uint64, targetUsername, method, decision, requestID, reason string) {
	h.audit.Write(service.AuditEntry{
		RequestID:      requestID,
		ActorUserID:    claims.UserID,
		ActorUsername:  actorNameFromClaims(claims),
		ClientIP:       c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
		Action:         "subsystem_auth.wayne." + strings.ToLower(method) + "." + decision,
		ResourceType:   "wayne_namespace",
		ResourceID:     strconv.FormatUint(namespaceID, 10),
		ScopeType:      "business_line",
		ScopeID:        businessLineID,
		BusinessLineID: businessLineID,
		NamespaceID:    namespaceID,
		Decision:       decision,
		Reason:         truncateAuditReason(reason),
		Metadata: map[string]any{
			"targetUsername": targetUsername,
		},
	})
}

func subsystemOperatorEmail(c *gin.Context, claims *auth.Claims) (string, bool) {
	operatorEmail := strings.TrimSpace(claims.Email)
	if operatorEmail == "" {
		operatorEmail = strings.TrimSpace(claims.Username)
	}
	if operatorEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is missing in token"})
		return "", false
	}
	return operatorEmail, true
}

func actorNameFromClaims(claims *auth.Claims) string {
	if claims == nil {
		return ""
	}
	if value := strings.TrimSpace(claims.Email); value != "" {
		return value
	}
	return claims.Username
}

func wayneUsernameForUser(user model.User) string {
	if value := strings.TrimSpace(user.Email); value != "" {
		return value
	}
	return strings.TrimSpace(user.Username)
}
