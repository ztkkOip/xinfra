package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/1024XEngineer/xinfra/server/internal/service"

	"github.com/gin-gonic/gin"
)

type WayneRoleBindingHandler struct {
	wayne *service.WayneRoleBindingService
	audit *service.AuditService
}

func NewWayneRoleBindingHandler(wayne *service.WayneRoleBindingService, audit *service.AuditService) *WayneRoleBindingHandler {
	return &WayneRoleBindingHandler{wayne: wayne, audit: audit}
}

func (h *WayneRoleBindingHandler) BindNamespace(c *gin.Context) {
	h.handle(c, "namespace", http.MethodPut)
}

func (h *WayneRoleBindingHandler) UnbindNamespace(c *gin.Context) {
	h.handle(c, "namespace", http.MethodDelete)
}

func (h *WayneRoleBindingHandler) BindApp(c *gin.Context) {
	h.handle(c, "app", http.MethodPut)
}

func (h *WayneRoleBindingHandler) UnbindApp(c *gin.Context) {
	h.handle(c, "app", http.MethodDelete)
}

func (h *WayneRoleBindingHandler) ListNamespaces(c *gin.Context) {
	h.handleQuery(c, "namespaces", 0, "")
}

func (h *WayneRoleBindingHandler) ListGroups(c *gin.Context) {
	var groupType *int
	if raw := strings.TrimSpace(c.Query("type")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || (value != 0 && value != 1) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
			return
		}
		groupType = &value
	}
	result, err := h.wayne.ListGroups(c.Request.Context(), groupType)
	if err != nil {
		h.writeQueryAudit(c, "groups", 0, "deny", err.Error())
		writeWayneRoleBindingError(c, result, err)
		return
	}
	h.writeQueryAudit(c, "groups", 0, "allow", "")
	writeWayneRoleBindingResult(c, result)
}

func (h *WayneRoleBindingHandler) GetCurrentUserRoles(c *gin.Context) {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	h.handleQuery(c, "user_roles", 0, username)
}

func (h *WayneRoleBindingHandler) NamespaceOperatorPermissions(c *gin.Context) {
	namespaceID, ok := parseUintPathParam(c, "namespaceid")
	if !ok {
		return
	}
	h.handleOperatorPermissions(c, "namespace", namespaceID)
}

func (h *WayneRoleBindingHandler) AppOperatorPermissions(c *gin.Context) {
	appID, ok := parseUintPathParam(c, "appid")
	if !ok {
		return
	}
	h.handleOperatorPermissions(c, "app", appID)
}

func (h *WayneRoleBindingHandler) handle(c *gin.Context, scope string, method string) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}
	operatorEmail := strings.TrimSpace(claims.Email)
	if operatorEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is missing in token"})
		return
	}

	resourceParam := "namespaceid"
	if scope == "app" {
		resourceParam = "appid"
	}
	resourceID, ok := parseUintPathParam(c, resourceParam)
	if !ok {
		return
	}
	req, ok := parseRoleBindingRequest(c)
	if !ok {
		return
	}
	targetUsername, ok := roleBindingTargetUsername(c, req)
	if !ok {
		return
	}
	if method == http.MethodPut && len(req.GroupIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupIds is required"})
		return
	}

	result, err := h.call(c, scope, method, resourceID, targetUsername, operatorEmail, req)
	if err != nil {
		h.writeAudit(c, claims.UserID, operatorEmail, scope, resourceID, targetUsername, "deny", req.RequestID, err.Error())
		writeWayneRoleBindingError(c, result, err)
		return
	}

	h.writeAudit(c, claims.UserID, operatorEmail, scope, resourceID, targetUsername, "allow", req.RequestID, "")
	writeWayneRoleBindingResult(c, result)
}

func (h *WayneRoleBindingHandler) handleQuery(c *gin.Context, resourceType string, resourceID uint64, username string) {
	var result *service.WayneRoleBindingResult
	var err error
	switch resourceType {
	case "namespaces":
		result, err = h.wayne.ListNamespaces(c.Request.Context())
	case "user_roles":
		result, err = h.wayne.GetUserRoles(c.Request.Context(), username)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported query resource"})
		return
	}
	if err != nil {
		h.writeQueryAudit(c, resourceType, resourceID, "deny", err.Error())
		writeWayneRoleBindingError(c, result, err)
		return
	}
	h.writeQueryAudit(c, resourceType, resourceID, "allow", "")
	writeWayneRoleBindingResult(c, result)
}

func (h *WayneRoleBindingHandler) handleOperatorPermissions(c *gin.Context, scope string, resourceID uint64) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}
	operatorEmail := strings.TrimSpace(claims.Email)
	if operatorEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is missing in token"})
		return
	}

	var result *service.WayneRoleBindingResult
	var err error
	if scope == "namespace" {
		result, err = h.wayne.NamespaceOperatorPermissions(c.Request.Context(), resourceID, operatorEmail)
	} else {
		result, err = h.wayne.AppOperatorPermissions(c.Request.Context(), resourceID, operatorEmail)
	}
	resourceType := scope + "_operator_permissions"
	if err != nil {
		h.writeQueryAudit(c, resourceType, resourceID, "deny", err.Error())
		writeWayneRoleBindingError(c, result, err)
		return
	}
	h.writeQueryAudit(c, resourceType, resourceID, "allow", "")
	writeWayneRoleBindingResult(c, result)
}

func (h *WayneRoleBindingHandler) call(c *gin.Context, scope, method string, resourceID uint64, targetUsername string, operatorEmail string, req service.WayneRoleBindingRequest) (*service.WayneRoleBindingResult, error) {
	if scope == "namespace" {
		if method == http.MethodPut {
			return h.wayne.BindNamespace(c.Request.Context(), resourceID, targetUsername, operatorEmail, req)
		}
		return h.wayne.UnbindNamespace(c.Request.Context(), resourceID, targetUsername, operatorEmail, req)
	}
	if method == http.MethodPut {
		return h.wayne.BindApp(c.Request.Context(), resourceID, targetUsername, operatorEmail, req)
	}
	return h.wayne.UnbindApp(c.Request.Context(), resourceID, targetUsername, operatorEmail, req)
}

func parseRoleBindingRequest(c *gin.Context) (service.WayneRoleBindingRequest, bool) {
	var req service.WayneRoleBindingRequest
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return req, false
	}
	if strings.TrimSpace(string(body)) == "" {
		return req, true
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return req, false
	}
	return req, true
}

func roleBindingTargetUsername(c *gin.Context, req service.WayneRoleBindingRequest) (string, bool) {
	if username := strings.TrimSpace(req.Username); username != "" {
		return username, true
	}
	for _, key := range []string{"username", "userName", "user_name"} {
		raw := strings.TrimSpace(c.Query(key))
		if raw == "" {
			continue
		}
		return raw, true
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
	return "", false
}

func parseUintPathParam(c *gin.Context, name string) (uint64, bool) {
	raw := strings.TrimSpace(c.Param(name))
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid %s", name)})
		return 0, false
	}
	return value, true
}

func writeWayneRoleBindingResult(c *gin.Context, result *service.WayneRoleBindingResult) {
	if result == nil {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
		return
	}
	contentType := result.ContentType
	if contentType == "" {
		contentType = "application/json"
	}
	if len(result.Body) == 0 {
		c.Status(result.StatusCode)
		return
	}
	c.Data(result.StatusCode, contentType, result.Body)
}

func writeWayneRoleBindingError(c *gin.Context, result *service.WayneRoleBindingResult, err error) {
	if result != nil && len(result.Body) > 0 {
		contentType := result.ContentType
		if contentType == "" {
			contentType = "application/json"
		}
		c.Data(result.StatusCode, contentType, result.Body)
		return
	}
	var httpErr *service.WayneRoleBindingHTTPError
	switch {
	case errors.Is(err, service.ErrWayneRoleBindingNotConfigured):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrWayenEmailMissing):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.As(err, &httpErr):
		c.JSON(httpErr.StatusCode, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
}

func (h *WayneRoleBindingHandler) writeAudit(c *gin.Context, userID uint64, operatorEmail, scope string, resourceID uint64, targetUsername string, decision, requestID, reason string) {
	h.audit.Write(service.AuditEntry{
		RequestID:     requestID,
		ActorUserID:   userID,
		ActorUsername: operatorEmail,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "wayne.role_binding." + decision,
		ResourceType:  "wayne_" + scope,
		ResourceID:    strconv.FormatUint(resourceID, 10),
		ScopeType:     scope,
		ScopeID:       resourceID,
		Decision:      decision,
		Reason:        truncateAuditReason(reason),
		Metadata: map[string]any{
			"targetUsername": targetUsername,
		},
	})
}

func (h *WayneRoleBindingHandler) writeQueryAudit(c *gin.Context, resourceType string, resourceID uint64, decision, reason string) {
	claims, _ := CurrentClaims(c)
	var actorUserID uint64
	var actorUsername string
	if claims != nil {
		actorUserID = claims.UserID
		actorUsername = strings.TrimSpace(claims.Email)
		if actorUsername == "" {
			actorUsername = claims.Username
		}
	}
	h.audit.Write(service.AuditEntry{
		ActorUserID:   actorUserID,
		ActorUsername: actorUsername,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "wayne.role_query." + decision,
		ResourceType:  "wayne_" + resourceType,
		ResourceID:    strconv.FormatUint(resourceID, 10),
		Decision:      decision,
		Reason:        truncateAuditReason(reason),
	})
}

func truncateAuditReason(reason string) string {
	const maxReasonBytes = 512
	if len(reason) <= maxReasonBytes {
		return reason
	}
	return reason[:maxReasonBytes]
}
