package handler

import (
	"net/http"
	"strconv"

	"github.com/1024XEngineer/xinfra/server/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Me(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
		"is_admin": claims.IsAdmin,
	})
}

func (h *UserHandler) List(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	query := h.db.Model(&model.User{}).Where("users.deleted_at IS NULL")
	if businessLineID := c.Query("business_line_id"); businessLineID != "" {
		if !claims.IsAdmin {
			id, err := strconv.ParseUint(businessLineID, 10, 64)
			if err != nil || id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business line id"})
				return
			}
			var binding model.BusinessLineUser
			if err := h.db.Where("business_line_id = ? AND user_id = ?", id, claims.UserID).First(&binding).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "current user is not in business line"})
				return
			}
		}
		query = query.
			Joins("JOIN business_line_users ON business_line_users.user_id = users.id").
			Where("business_line_users.business_line_id = ?", businessLineID)
	}

	var users []model.User
	if err := query.Order("users.id ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]gin.H, 0, len(users))
	for _, user := range users {
		items = append(items, gin.H{
			"uid":      user.ID,
			"username": user.Username,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
