package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg  config.Config
	auth *service.AuthService
}

type localLoginRequest struct {
	Username string `json:"username"`
}

func NewAuthHandler(cfg config.Config, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{cfg: cfg, auth: authService}
}

func (h *AuthHandler) Config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sso_enabled": h.cfg.SSOEnabled,
	})
}

func (h *AuthHandler) LocalLogin(c *gin.Context) {
	if h.cfg.SSOEnabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "local login is disabled when sso is enabled"})
		return
	}

	var req localLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.auth.LocalLogin(strings.TrimSpace(req.Username), c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredential):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrUserDisabled):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AuthSessionCookieName,
		Value:    result.Token,
		Path:     "/auth/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.ExpiresAt).Seconds()),
	})
	c.JSON(http.StatusOK, gin.H{
		"token":      result.Token,
		"expires_in": int(h.cfg.JWTTTL().Seconds()),
		"user": gin.H{
			"id":            result.User.ID,
			"username":      result.User.Username,
			"display_name":  result.User.DisplayName,
			"email":         result.User.Email,
			"business_line": "",
			"is_admin":      result.User.IsAdmin,
		},
	})
}
