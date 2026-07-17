package handler

import (
	"net/http"

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
	var users []model.User
	if err := h.db.Where("deleted_at IS NULL").Order("id ASC").Find(&users).Error; err != nil {
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
