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
	scope := strings.TrimSpace(c.Query("scope"))
	nonce := strings.TrimSpace(c.Query("nonce"))

	if responseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported response_type"})
		return
	}
	client, ok := h.client(clientID)
	if !ok || !validClientRedirect(client, redirectURI) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_id or redirect_uri"})
		return
	}

	user, ok := h.sessionUser(c)
	if !ok {
		loginURL := "/auth/api/v1/login/internal-sso?relay_state=" + url.QueryEscape(c.Request.URL.RequestURI())
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	code, codeID, expiresAt, err := auth.SignOAuthCode(h.cfg.JWTSecret, h.cfg.JWTIssuer, h.cfg.OAuthCodeTTL(), user.ID, clientID, redirectURI, scope, nonce)
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
	client, ok := h.client(clientID)
	if !ok || clientSecret != client.Secret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid client credentials"})
		return
	}

	codeClaims, err := auth.ParseOAuthCode(h.cfg.JWTSecret, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}
	if codeClaims.ClientID != clientID || codeClaims.RedirectURI != redirectURI || !validClientRedirect(client, redirectURI) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	var token string
	var idToken string
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

		display := strings.TrimSpace(user.DisplayName)
		if display == "" {
			display = user.Username
		}
		tokenID := ""
		var err error
		token, tokenID, expiresAt, err = auth.Sign(h.cfg.JWTSecret, h.cfg.JWTIssuer, h.cfg.JWTTTL(), user.ID, user.Username, user.Email, user.IsAdmin)
		if err != nil {
			return err
		}
		privateKey, err := auth.LoadRSAPrivateKey(h.cfg.SAMLSPKey)
		if err != nil {
			return err
		}
		keyID := auth.RSAKeyID(&privateKey.PublicKey)
		idToken, _, err = auth.SignIDToken(privateKey, keyID, h.cfg.OIDCIssuer, clientID, h.cfg.JWTTTL(), user.ID, user.Username, user.Email, display, codeClaims.Nonce)
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
		"id_token":     idToken,
		"token_type":   "Bearer",
		"expires_in":   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *OAuthHandler) Discovery(c *gin.Context) {
	issuer := strings.TrimRight(h.cfg.OIDCIssuer, "/")
	c.JSON(http.StatusOK, gin.H{
		"issuer":                                issuer,
		"authorization_endpoint":                h.cfg.OIDCAuthorizeURL,
		"token_endpoint":                        h.cfg.OIDCTokenURL,
		"userinfo_endpoint":                     h.cfg.OIDCUserInfoURL,
		"jwks_uri":                              h.cfg.OIDCJWKSURL,
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"claims_supported":                      []string{"sub", "name", "preferred_username", "email", "email_verified"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
	})
}

func (h *OAuthHandler) JWKS(c *gin.Context) {
	privateKey, err := auth.LoadRSAPrivateKey(h.cfg.SAMLSPKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"keys": []auth.JWK{auth.PublicJWKFromKey(privateKey)},
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
		"sub":                user.Email,
		"name":               user.Username,
		"preferred_username": user.Username,
		"email":              user.Email,
		"email_verified":     user.Email != "",
		"display":            display,
	})
}

func (h *OAuthHandler) oauthConfigured() bool {
	return len(h.cfg.OAuthClients()) > 0
}

func (h *OAuthHandler) client(clientID string) (config.OAuthClient, bool) {
	client, ok := h.cfg.OAuthClients()[strings.TrimSpace(clientID)]
	return client, ok
}

func validClientRedirect(client config.OAuthClient, redirectURI string) bool {
	if redirectURI == "" {
		return false
	}
	if len(client.RedirectURIs) > 0 {
		for _, allowed := range client.RedirectURIs {
			if redirectURI == allowed {
				return true
			}
		}
		return false
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
