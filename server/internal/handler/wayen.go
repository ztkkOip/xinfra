package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/1024XEngineer/xinfra/server/internal/auth"
	"github.com/1024XEngineer/xinfra/server/internal/model"
	"github.com/1024XEngineer/xinfra/server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type selfWayenCredentialPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type WayenHandler struct {
	db    *gorm.DB
	wayen *service.WayenService
	audit *service.AuditService
}

func NewWayenHandler(db *gorm.DB, wayen *service.WayenService, audit *service.AuditService) *WayenHandler {
	return &WayenHandler{db: db, wayen: wayen, audit: audit}
}

func (h *WayenHandler) Login(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	email, _, ok := h.currentUser(c, claims)
	if !ok {
		return
	}

	refOverride, ok := h.resolveWayneRef(c, claims)
	if !ok {
		return
	}

	result, err := h.wayen.Login(email, claims.Username, refOverride)
	if err != nil {
		h.writeAudit(c, claims.UserID, claims.Username, "deny", err.Error())
		switch {
		case errors.Is(err, service.ErrWayenEmailMissing):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrWayenNotConfigured):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrWayenCredentialNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrWayenLoginFailed):
			status := http.StatusBadGateway
			if strings.Contains(err.Error(), "status 400") {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	for _, cookie := range result.Cookies {
		http.SetCookie(c.Writer, cookie)
	}
	h.writeAudit(c, claims.UserID, claims.Username, "allow", "")
	if strings.Contains(c.GetHeader("Accept"), "application/json") {
		c.JSON(http.StatusOK, gin.H{"target_url": result.TargetURL})
		return
	}
	c.Redirect(http.StatusFound, result.TargetURL)
}

func (h *WayenHandler) resolveWayneRef(c *gin.Context, claims *auth.Claims) (string, bool) {
	rawBusinessLineID := strings.TrimSpace(c.Query("business_line_id"))
	if rawBusinessLineID == "" {
		return "", true
	}

	businessLineID, err := strconv.ParseUint(rawBusinessLineID, 10, 64)
	if err != nil || businessLineID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business_line_id"})
		return "", false
	}

	var businessLine model.BusinessLine
	if err := h.db.First(&businessLine, "id = ?", businessLineID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business line not found"})
			return "", false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", false
	}

	if !claims.IsAdmin {
		var binding model.BusinessLineUser
		if err := h.db.Where("business_line_id = ? AND user_id = ?", businessLineID, claims.UserID).First(&binding).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusForbidden, gin.H{"error": "current user does not have business line permission"})
				return "", false
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return "", false
		}
	}

	var mapping model.BusinessLineWayneNamespace
	if err := h.db.Where("business_line_id = ?", businessLineID).Order("wayne_namespace_id ASC").First(&mapping).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", true
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", false
	}
	if mapping.WayneNamespaceID == 0 {
		return "", true
	}
	return fmt.Sprintf("/portal/namespace/%d/app", mapping.WayneNamespaceID), true
}

func (h *WayenHandler) GetCredential(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	email, _, ok := h.currentUser(c, claims)
	if !ok {
		return
	}

	response := gin.H{
		"email":      email,
		"configured": false,
		"enabled":    false,
	}
	if email == "" {
		c.JSON(http.StatusOK, response)
		return
	}

	var item model.WayenCredential
	if err := h.db.Where("email = ?", email).First(&item).Error; err == nil {
		response["configured"] = true
		response["enabled"] = item.Enabled
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *WayenHandler) SaveCredential(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	var req selfWayenCredentialPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, user, ok := h.currentUser(c, claims)
	if !ok {
		return
	}
	if trimmed := strings.TrimSpace(req.Email); trimmed != "" {
		email = trimmed
	}
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	password := strings.TrimSpace(req.Password)
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if user.ID != 0 && strings.TrimSpace(user.Email) != email {
			if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Update("email", email).Error; err != nil {
				return err
			}
		}

		var item model.WayenCredential
		err := tx.Where("email = ?", email).First(&item).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(&model.WayenCredential{
				Email:    email,
				Password: password,
				Enabled:  true,
			}).Error
		}
		if err != nil {
			return err
		}
		return tx.Model(&item).Updates(map[string]any{
			"password": password,
			"enabled":  true,
		}).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.audit.Write(service.AuditEntry{
		ActorUserID:   claims.UserID,
		ActorUsername: claims.Username,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "wayen.credential.save",
		ResourceType:  "wayen_credential",
		ResourceID:    email,
		Decision:      "allow",
	})
	c.JSON(http.StatusOK, gin.H{"email": email, "configured": true, "enabled": true})
}

func (h *WayenHandler) currentUser(c *gin.Context, claims *auth.Claims) (string, model.User, bool) {
	var user model.User
	email := strings.TrimSpace(claims.Email)
	if claims.UserID != 0 {
		if err := h.db.Select("id", "email").Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "current user not found"})
			return "", user, false
		}
		if email == "" {
			email = strings.TrimSpace(user.Email)
		}
	}
	return email, user, true
}

func (h *WayenHandler) writeAudit(c *gin.Context, userID uint64, username, decision, reason string) {
	h.audit.Write(service.AuditEntry{
		ActorUserID:   userID,
		ActorUsername: username,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "wayen.login",
		ResourceType:  "wayen",
		Decision:      decision,
		Reason:        reason,
	})
}
