package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"authserver/internal/auth"
	"authserver/internal/config"
	"authserver/internal/model"
	"authserver/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const AuthSessionCookieName = "authserver_token"

type OAuthHandler struct {
	cfg   config.Config
	db    *gorm.DB
	audit *service.AuditService
}

func NewOAuthHandler(cfg config.Config, db *gorm.DB, audit *service.AuditService) *OAuthHandler {
	return &OAuthHandler{cfg: cfg, db: db, audit: audit}
}

func (h *OAuthHandler) Authorize(c *gin.Context) {
	if !h.oauthConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "oauth client is not configured"})
		return
	}

	clientID := strings.TrimSpace(c.Query("client_id"))
	redirectURI := strings.TrimSpace(c.Query("redirect_uri"))
	responseType := strings.TrimSpace(c.Query("response_type"))
	state := c.Query("state")

	if responseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported response_type"})
		return
	}
	if !h.validClientRedirect(clientID, redirectURI) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_id or redirect_uri"})
		return
	}

	user, ok := h.sessionUser(c)
	if !ok {
		loginURL := "/auth/api/v1/login/internal-sso?relay_state=" + url.QueryEscape(c.Request.URL.RequestURI())
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	code, codeID, expiresAt, err := auth.SignOAuthCode(h.cfg.JWTSecret, h.cfg.JWTIssuer, h.cfg.OAuthCodeTTL(), user.ID, clientID, redirectURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Create(&model.AccessToken{
		UserID:    user.ID,
		TokenID:   codeID,
		TokenType: "oauth_code",
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		ExpiresAt: expiresAt,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.audit.Write(service.AuditEntry{
		ActorUserID:   user.ID,
		ActorUsername: user.Username,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "oauth.authorize",
		ResourceType:  "oauth_client",
		ResourceID:    clientID,
		Decision:      "allow",
	})

	target, err := url.Parse(redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid redirect_uri"})
		return
	}
	values := target.Query()
	values.Set("code", code)
	if state != "" {
		values.Set("state", state)
	}
	target.RawQuery = values.Encode()
	c.Redirect(http.StatusFound, target.String())
}

func (h *OAuthHandler) Token(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")

	if !h.oauthConfigured() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "oauth client is not configured"})
		return
	}
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID, clientSecret := oauthClientCredentials(c)
	redirectURI := strings.TrimSpace(c.PostForm("redirect_uri"))
	code := strings.TrimSpace(c.PostForm("code"))
	if c.PostForm("grant_type") != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported grant_type"})
		return
	}
	if clientID != h.cfg.OAuthClientID || clientSecret != h.cfg.OAuthClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid client credentials"})
		return
	}

	codeClaims, err := auth.ParseOAuthCode(h.cfg.JWTSecret, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}
	if codeClaims.ClientID != clientID || codeClaims.RedirectURI != redirectURI || !h.validClientRedirect(clientID, redirectURI) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	var token string
	var expiresAt time.Time
	var user model.User
	now := time.Now()
	err = h.db.Transaction(func(tx *gorm.DB) error {
		var codeRecord model.AccessToken
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_id = ? AND token_type = ? AND revoked = ? AND expires_at > ?", codeClaims.ID, "oauth_code", false, now).
			First(&codeRecord).Error; err != nil {
			return err
		}
		if err := tx.First(&user, "id = ? AND deleted_at IS NULL", codeClaims.UserID).Error; err != nil {
			return err
		}
		if user.Status != "active" {
			return service.ErrUserDisabled
		}

		tokenID := ""
		var err error
		token, tokenID, expiresAt, err = auth.Sign(h.cfg.JWTSecret, h.cfg.JWTIssuer, h.cfg.JWTTTL(), user.ID, user.Username, user.Email, user.IsAdmin)
		if err != nil {
			return err
		}
		if err := tx.Model(&codeRecord).Updates(map[string]any{
			"revoked":    true,
			"revoked_at": &now,
		}).Error; err != nil {
			return err
		}
		return tx.Create(&model.AccessToken{
			UserID:    user.ID,
			TokenID:   tokenID,
			TokenType: "oauth_access",
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			ExpiresAt: expiresAt,
		}).Error
	})
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrUserDisabled) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": "invalid code"})
		return
	}

	h.audit.Write(service.AuditEntry{
		ActorUserID:   user.ID,
		ActorUsername: user.Username,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "oauth.token",
		ResourceType:  "oauth_client",
		ResourceID:    clientID,
		Decision:      "allow",
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *OAuthHandler) UserInfo(c *gin.Context) {
	claims, err := bearerClaims(h.cfg, c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var tokenRecord model.AccessToken
	if err := h.db.Where("token_id = ? AND revoked = ? AND expires_at > ?", claims.ID, false, time.Now()).First(&tokenRecord).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var user model.User
	if err := h.db.First(&user, "id = ? AND deleted_at IS NULL", claims.UserID).Error; err != nil || user.Status != "active" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	display := strings.TrimSpace(user.DisplayName)
	if display == "" {
		display = user.Username
	}
	c.JSON(http.StatusOK, gin.H{
		"name":    user.Username,
		"email":   user.Email,
		"display": display,
	})
}

func (h *OAuthHandler) oauthConfigured() bool {
	return h.cfg.OAuthClientID != "" && h.cfg.OAuthClientSecret != ""
}

func (h *OAuthHandler) validClientRedirect(clientID, redirectURI string) bool {
	if clientID == "" || clientID != h.cfg.OAuthClientID || redirectURI == "" {
		return false
	}
	if h.cfg.OAuthRedirectURI != "" {
		return redirectURI == h.cfg.OAuthRedirectURI
	}
	parsed, err := url.Parse(redirectURI)
	return err == nil && parsed.IsAbs() && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func (h *OAuthHandler) sessionUser(c *gin.Context) (model.User, bool) {
	tokenValue := ""
	if cookie, err := c.Cookie(AuthSessionCookieName); err == nil {
		tokenValue = cookie
	}
	if tokenValue == "" {
		tokenValue = c.Query("sso_token")
	}
	claims, err := auth.Parse(h.cfg.JWTSecret, tokenValue)
	if err != nil {
		return model.User{}, false
	}

	var user model.User
	if err := h.db.First(&user, "id = ? AND deleted_at IS NULL", claims.UserID).Error; err != nil || user.Status != "active" {
		return model.User{}, false
	}
	return user, true
}

func oauthClientCredentials(c *gin.Context) (string, string) {
	if id, secret, ok := c.Request.BasicAuth(); ok {
		return id, secret
	}
	value := c.GetHeader("Authorization")
	if strings.HasPrefix(value, "Basic ") {
		raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, "Basic "))
		if err == nil {
			id, secret, ok := strings.Cut(string(raw), ":")
			if ok {
				return id, secret
			}
		}
	}
	return strings.TrimSpace(c.PostForm("client_id")), strings.TrimSpace(c.PostForm("client_secret"))
}

func bearerClaims(cfg config.Config, value string) (*auth.Claims, error) {
	if !strings.HasPrefix(value, "Bearer ") {
		return nil, errors.New("missing bearer token")
	}
	return auth.Parse(cfg.JWTSecret, strings.TrimSpace(strings.TrimPrefix(value, "Bearer ")))
}
