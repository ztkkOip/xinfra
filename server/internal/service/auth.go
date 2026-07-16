package service

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/auth"
	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/model"
	"github.com/1024XEngineer/xinfra/server/internal/sso"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredential  = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user disabled")
	ErrSAMLSubjectMissing = errors.New("saml subject is missing")
)

type AuthService struct {
	cfg   config.Config
	db    *gorm.DB
	audit *AuditService
}

type LoginResult struct {
	Token     string     `json:"token"`
	TokenID   string     `json:"token_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	User      model.User `json:"user"`
}

func NewAuthService(cfg config.Config, db *gorm.DB, audit *AuditService) *AuthService {
	return &AuthService{cfg: cfg, db: db, audit: audit}
}

func (s *AuthService) LocalLogin(username, clientIP, userAgent string) (*LoginResult, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrInvalidCredential
	}

	var result *LoginResult
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user, err := findOrCreateLocalUser(tx, username)
		if err != nil {
			return err
		}
		if user.Status != "active" {
			s.audit.Write(AuditEntry{ActorUserID: user.ID, ActorUsername: user.Username, ClientIP: clientIP, UserAgent: userAgent, Action: "local.login.failed", Decision: "deny", Reason: "user_disabled"})
			return ErrUserDisabled
		}

		token, tokenID, expiresAt, err := auth.Sign(s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTTTL(), user.ID, user.Username, user.Email, user.IsAdmin)
		if err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Update("last_login_at", now).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.AccessToken{
			UserID:    user.ID,
			TokenID:   tokenID,
			TokenType: "access",
			ClientIP:  clientIP,
			UserAgent: userAgent,
			ExpiresAt: expiresAt,
		}).Error; err != nil {
			return err
		}

		user.LastLoginAt = &now
		result = &LoginResult{Token: token, TokenID: tokenID, ExpiresAt: expiresAt, User: user}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.audit.Write(AuditEntry{ActorUserID: result.User.ID, ActorUsername: result.User.Username, ClientIP: clientIP, UserAgent: userAgent, Action: "local.login.success", Decision: "allow"})
	return result, nil
}

func (s *AuthService) SAMLLogin(info *sso.SAMLDebugInfo, clientIP, userAgent string) (*LoginResult, error) {
	subject := strings.TrimSpace(info.NameID)
	email := firstSAMLAttribute(info.Attributes,
		"email",
		"mail",
		"Email",
		"EmailAddress",
		"emailAddress",
		"urn:oid:0.9.2342.19200300.100.1.3",
		"urn:oid:1.3.6.1.4.1.5923.1.1.1.6",
	)
	if email == "" && looksLikeEmail(subject) {
		email = subject
	}
	if subject == "" {
		subject = email
	}
	if subject == "" {
		reason := fmt.Sprintf("subject_missing name_id=%q attribute_keys=%v", info.NameID, samlAttributeKeys(info.Attributes))
		s.audit.Write(AuditEntry{ClientIP: clientIP, UserAgent: userAgent, Action: "saml.login.failed", Decision: "deny", Reason: reason})
		return nil, ErrSAMLSubjectMissing
	}

	displayName := firstSAMLAttribute(info.Attributes,
		"displayName",
		"name",
		"cn",
		"username",
		"uid",
		"urn:oid:2.5.4.3",
		"urn:oid:0.9.2342.19200300.100.1.1",
	)

	var result *LoginResult
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user, err := findOrCreateSAMLUser(tx, subject, email, displayName)
		if err != nil {
			return err
		}
		if user.Status != "active" {
			s.audit.Write(AuditEntry{ActorUserID: user.ID, ActorUsername: user.Username, ClientIP: clientIP, UserAgent: userAgent, Action: "saml.login.failed", Decision: "deny", Reason: "user_disabled"})
			return ErrUserDisabled
		}

		token, tokenID, expiresAt, err := auth.Sign(s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTTTL(), user.ID, user.Username, user.Email, user.IsAdmin)
		if err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]any{
			"last_login_at": now,
			"email":         user.Email,
			"display_name":  user.DisplayName,
		}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.AccessToken{
			UserID:    user.ID,
			TokenID:   tokenID,
			TokenType: "access",
			ClientIP:  clientIP,
			UserAgent: userAgent,
			ExpiresAt: expiresAt,
		}).Error; err != nil {
			return err
		}

		user.LastLoginAt = &now
		result = &LoginResult{Token: token, TokenID: tokenID, ExpiresAt: expiresAt, User: user}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.audit.Write(AuditEntry{ActorUserID: result.User.ID, ActorUsername: result.User.Username, ClientIP: clientIP, UserAgent: userAgent, Action: "saml.login.success", Decision: "allow"})
	return result, nil
}

func findOrCreateLocalUser(tx *gorm.DB, username string) (model.User, error) {
	email := username

	var user model.User
	query := tx.Where("username = ? AND deleted_at IS NULL", username)
	query = query.Or("email = ? AND deleted_at IS NULL", email)
	if err := query.First(&user).Error; err == nil {
		return user, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}

	var count int64
	if err := tx.Model(&model.User{}).Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return model.User{}, err
	}
	user = model.User{
		Username:    username,
		DisplayName: username,
		Email:       email,
		Source:      "local",
		Status:      "active",
		IsAdmin:     count == 0,
	}
	if err := tx.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func findOrCreateSAMLUser(tx *gorm.DB, subject, email, displayName string) (model.User, error) {
	username := samlUsername(email, subject)
	if displayName == "" {
		displayName = username
	}

	var user model.User
	if err := findExistingSAMLUser(tx, &user, username, email, subject); err == nil {
		return updateSAMLUser(tx, user, username, email, displayName, subject)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}

	user = model.User{
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Source:      "saml",
		ExternalID:  subject,
		Status:      "active",
		IsAdmin:     false,
	}
	if err := tx.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func findExistingSAMLUser(tx *gorm.DB, user *model.User, username, email, subject string) error {
	query := tx.Where("username = ? AND deleted_at IS NULL", username)
	if email != "" {
		query = query.Or("email = ? AND deleted_at IS NULL", email)
	}
	if subject != "" {
		query = query.Or("source = ? AND external_id = ? AND deleted_at IS NULL", "saml", subject)
	}
	return query.First(user).Error
}

func updateSAMLUser(tx *gorm.DB, user model.User, username, email, displayName, subject string) (model.User, error) {
	updates := map[string]any{
		"source":      "saml",
		"external_id": subject,
	}
	user.Source = "saml"
	user.ExternalID = subject

	if email != "" && user.Email != email {
		updates["email"] = email
		user.Email = email
	}
	if displayName != "" && user.DisplayName != displayName {
		updates["display_name"] = displayName
		user.DisplayName = displayName
	}
	if username != "" && user.Username != username {
		var count int64
		if err := tx.Model(&model.User{}).Where("username = ? AND id <> ? AND deleted_at IS NULL", username, user.ID).Count(&count).Error; err != nil {
			return user, err
		}
		if count == 0 {
			updates["username"] = username
			user.Username = username
		}
	}
	if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		return user, err
	}
	return user, nil
}

func firstSAMLAttribute(attrs map[string][]string, names ...string) string {
	for _, name := range names {
		for _, value := range attrs[name] {
			value = strings.TrimSpace(value)
			if value != "" {
				return value
			}
		}
	}
	for key, values := range attrs {
		lower := strings.ToLower(key)
		if strings.Contains(lower, "email") || strings.Contains(lower, "mail") {
			for _, value := range values {
				value = strings.TrimSpace(value)
				if value != "" {
					return value
				}
			}
		}
	}
	return ""
}

func looksLikeEmail(value string) bool {
	if value == "" {
		return false
	}
	_, err := mail.ParseAddress(value)
	return err == nil && strings.Contains(value, "@")
}

var usernameCleaner = regexp.MustCompile(`[^a-zA-Z0-9@._-]+`)

func samlUsername(email, subject string) string {
	base := strings.TrimSpace(email)
	if base == "" {
		base = strings.TrimSpace(subject)
	}
	base = usernameCleaner.ReplaceAllString(base, "-")
	base = strings.Trim(base, ".-_")
	if base == "" {
		return "saml-user"
	}
	return base
}

func samlAttributeKeys(attrs map[string][]string) []string {
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}
	return keys
}
