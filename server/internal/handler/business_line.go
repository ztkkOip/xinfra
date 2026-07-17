package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessLineHandler struct {
	db *gorm.DB
}

type BusinessLineWithPermission struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Permission int    `json:"permission"`
}

type GrantBusinessLinePermissionRequest struct {
	BusinessLineID       uint64 `json:"business_line_id" binding:"required"`
	TargetUserID         uint64 `json:"target_user_id" binding:"required"`
	TargetBusinessLineID uint64 `json:"target_business_line_id" binding:"required"`
	Permission           int    `json:"permission"`
}

type BusinessLinePayload struct {
	Name string `json:"name" binding:"required"`
}

type WayneNamespaceBindingPayload struct {
	Namespaces []WayneNamespaceBindingItem `json:"namespaces"`
}

type WayneNamespaceBindingItem struct {
	ID            uint64 `json:"id" binding:"required"`
	Name          string `json:"name"`
	KubeNamespace string `json:"kubeNamespace"`
}

func NewBusinessLineHandler(db *gorm.DB) *BusinessLineHandler {
	return &BusinessLineHandler{db: db}
}

func (h *BusinessLineHandler) ListCurrentUserBusinessLines(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	if claims.IsAdmin {
		var rows []model.BusinessLine
		if err := h.db.Order("id ASC").Find(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items := make([]BusinessLineWithPermission, 0, len(rows))
		for _, row := range rows {
			items = append(items, BusinessLineWithPermission{
				ID:         row.ID,
				Name:       row.Name,
				CreatedAt:  row.CreatedAt.Format(time.RFC3339),
				UpdatedAt:  row.UpdatedAt.Format(time.RFC3339),
				Permission: 0,
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
		return
	}

	var rows []struct {
		ID         uint64
		Name       string
		CreatedAt  time.Time
		UpdatedAt  time.Time
		Permission int
	}
	if err := h.db.
		Table("business_line_users").
		Select("business_lines.id, business_lines.name, business_lines.created_at, business_lines.updated_at, business_line_users.permission").
		Joins("JOIN business_lines ON business_lines.id = business_line_users.business_line_id").
		Where("business_line_users.user_id = ?", claims.UserID).
		Order("business_lines.id ASC").
		Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]BusinessLineWithPermission, 0, len(rows))
	for _, row := range rows {
		items = append(items, BusinessLineWithPermission{
			ID:         row.ID,
			Name:       row.Name,
			CreatedAt:  row.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  row.UpdatedAt.Format(time.RFC3339),
			Permission: row.Permission,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BusinessLineHandler) ListAll(c *gin.Context) {
	if !requirePlatformAdmin(c) {
		return
	}

	var rows []model.BusinessLine
	if err := h.db.Order("id ASC").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"id":         row.ID,
			"name":       row.Name,
			"created_at": row.CreatedAt.Format(time.RFC3339),
			"updated_at": row.UpdatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BusinessLineHandler) Create(c *gin.Context) {
	if !requirePlatformAdmin(c) {
		return
	}

	var req BusinessLinePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := model.BusinessLine{Name: req.Name}
	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	writeBusinessLine(c, item)
}

func (h *BusinessLineHandler) Update(c *gin.Context) {
	if !requirePlatformAdmin(c) {
		return
	}

	var req BusinessLinePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item model.BusinessLine
	if err := h.db.First(&item, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business line not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Model(&item).Update("name", req.Name).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	item.Name = req.Name
	writeBusinessLine(c, item)
}

func (h *BusinessLineHandler) Delete(c *gin.Context) {
	if !requirePlatformAdmin(c) {
		return
	}

	var item model.BusinessLine
	if err := h.db.First(&item, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business line not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("business_line_id = ?", item.ID).Delete(&model.BusinessLineUser{}).Error; err != nil {
			return err
		}
		if err := tx.Where("business_line_id = ?", item.ID).Delete(&model.BusinessLineWayneNamespace{}).Error; err != nil {
			return err
		}
		return tx.Delete(&item).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *BusinessLineHandler) GrantPermission(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	var req GrantBusinessLinePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Permission != 0 && req.Permission != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "permission must be 0 or 1"})
		return
	}

	if !claims.IsAdmin {
		var currentBinding model.BusinessLineUser
		if err := h.db.Where("business_line_id = ? AND user_id = ? AND permission = ?", req.BusinessLineID, claims.UserID, 0).
			First(&currentBinding).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusForbidden, gin.H{"error": "current user is not platform admin or business line admin"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var businessLine model.BusinessLine
	if err := h.db.First(&businessLine, "id = ?", req.BusinessLineID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business line not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var targetUser model.User
	if err := h.db.Where("id = ? AND deleted_at IS NULL", req.TargetUserID).First(&targetUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "target user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var targetBusinessLine model.BusinessLine
	if err := h.db.First(&targetBusinessLine, "id = ?", req.TargetBusinessLineID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "target business line not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var binding model.BusinessLineUser
	err := h.db.Where("business_line_id = ? AND user_id = ?", req.TargetBusinessLineID, req.TargetUserID).First(&binding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		binding = model.BusinessLineUser{
			BusinessLineID: req.TargetBusinessLineID,
			UserID:         req.TargetUserID,
			Permission:     req.Permission,
		}
		if err := h.db.Create(&binding).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if binding.Permission != req.Permission {
		if err := h.db.Model(&binding).Update("permission", req.Permission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		binding.Permission = req.Permission
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               binding.ID,
		"business_line_id": binding.BusinessLineID,
		"user_id":          binding.UserID,
		"permission":       binding.Permission,
		"created_at":       binding.CreatedAt.Format(time.RFC3339),
		"updated_at":       binding.UpdatedAt.Format(time.RFC3339),
	})
}

func (h *BusinessLineHandler) ListWayneNamespaces(c *gin.Context) {
	businessLineID, ok := parseBusinessLineID(c)
	if !ok {
		return
	}
	if !h.canManageBusinessLine(c, businessLineID) {
		return
	}

	var rows []model.BusinessLineWayneNamespace
	if err := h.db.Where("business_line_id = ?", businessLineID).Order("wayne_namespace_id ASC").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"id":            row.WayneNamespaceID,
			"name":          row.WayneNamespaceName,
			"kubeNamespace": row.KubeNamespace,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BusinessLineHandler) ReplaceWayneNamespaces(c *gin.Context) {
	businessLineID, ok := parseBusinessLineID(c)
	if !ok {
		return
	}
	if !h.canManageBusinessLine(c, businessLineID) {
		return
	}

	var req WayneNamespaceBindingPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("business_line_id = ?", businessLineID).Delete(&model.BusinessLineWayneNamespace{}).Error; err != nil {
			return err
		}
		for _, item := range req.Namespaces {
			row := model.BusinessLineWayneNamespace{
				BusinessLineID:     businessLineID,
				WayneNamespaceID:   item.ID,
				WayneNamespaceName: item.Name,
				KubeNamespace:      item.KubeNamespace,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *BusinessLineHandler) canManageBusinessLine(c *gin.Context, businessLineID uint64) bool {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return false
	}

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

func parseBusinessLineID(c *gin.Context) (uint64, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business line id"})
		return 0, false
	}
	return value, true
}

func requirePlatformAdmin(c *gin.Context) bool {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return false
	}
	if !claims.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "platform admin required"})
		return false
	}
	return true
}

func writeBusinessLine(c *gin.Context, item model.BusinessLine) {
	c.JSON(http.StatusOK, gin.H{
		"id":         item.ID,
		"name":       item.Name,
		"created_at": item.CreatedAt.Format(time.RFC3339),
		"updated_at": item.UpdatedAt.Format(time.RFC3339),
	})
}
