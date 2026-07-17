package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/model"

	"gorm.io/gorm"
)

var (
	ErrWayenNotConfigured      = errors.New("wayen login is not configured")
	ErrWayenEmailMissing       = errors.New("email is missing in token")
	ErrWayenCredentialNotFound = errors.New("wayen credential not found")
	ErrWayenLoginFailed        = errors.New("wayen login failed")
)

type WayenService struct {
	cfg    config.Config
	db     *gorm.DB
	client *http.Client
}

type WayenLoginResult struct {
	TargetURL string
	Cookies   []*http.Cookie
}

func NewWayenService(cfg config.Config, db *gorm.DB) *WayenService {
	return &WayenService{
		cfg: cfg,
		db:  db,
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (s *WayenService) Login(email, username, refOverride string) (*WayenLoginResult, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, ErrWayenEmailMissing
	}
	oauthLoginURL := strings.TrimSpace(s.cfg.WayenOAuthLoginURL)
	if oauthLoginURL == "" {
		oauthLoginURL = strings.TrimSpace(s.cfg.OAuthRedirectURI)
	}
	if oauthLoginURL != "" && strings.TrimSpace(s.cfg.WayenTargetURL) != "" {
		target, err := s.oauthLoginURL(oauthLoginURL, s.cfg.WayenTargetURL, refOverride)
		if err != nil {
			return nil, err
		}
		return &WayenLoginResult{TargetURL: target}, nil
	}
	if strings.TrimSpace(s.cfg.WayenLoginURL) == "" || strings.TrimSpace(s.cfg.WayenTargetURL) == "" {
		return nil, ErrWayenNotConfigured
	}

	var credential model.WayenCredential
	if err := s.db.Where("email = ? AND enabled = ?", email, true).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWayenCredentialNotFound
		}
		return nil, err
	}

	loginName := email
	if strings.EqualFold(s.cfg.WayenLoginValue, "username") && strings.TrimSpace(username) != "" {
		loginName = strings.TrimSpace(username)
	}

	loginURL, body, contentType, err := s.loginRequest(s.cfg.WayenLoginURL, loginName, credential.Password)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, loginURL, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json, text/html;q=0.9, */*;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("%w: status %d", ErrWayenLoginFailed, resp.StatusCode)
	}

	return &WayenLoginResult{
		TargetURL: s.cfg.WayenTargetURL,
		Cookies:   resp.Cookies(),
	}, nil
}

func (s *WayenService) oauthLoginURL(redirectURI, targetURL, refOverride string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(redirectURI))
	if err != nil {
		return "", err
	}
	next, err := url.Parse(strings.TrimSpace(targetURL))
	if err != nil {
		return "", err
	}
	if next.Path == "" || next.Path == "/" {
		next.Path = "/sign-in"
	}
	values := next.Query()
	values.Set("ref", defaultConfigValue(refOverride, defaultConfigValue(s.cfg.WayenOAuthRef, "/portal/namespace/1/app")))
	next.RawQuery = values.Encode()

	query := parsed.Query()
	query.Set("next", next.String())
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (s *WayenService) loginRequest(loginURL, email, password string) (string, io.Reader, string, error) {
	usernameKey := defaultConfigValue(s.cfg.WayenUsernameKey, "email")
	passwordKey := defaultConfigValue(s.cfg.WayenPasswordKey, "password")

	if strings.EqualFold(s.cfg.WayenLoginFormat, "query") {
		parsed, err := url.Parse(loginURL)
		if err != nil {
			return "", nil, "", err
		}
		values := parsed.Query()
		values.Set(usernameKey, email)
		values.Set(passwordKey, password)
		parsed.RawQuery = values.Encode()
		return parsed.String(), nil, "", nil
	}

	if strings.EqualFold(s.cfg.WayenLoginFormat, "json") {
		payload := map[string]string{
			usernameKey: email,
			passwordKey: password,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return "", nil, "", err
		}
		return loginURL, bytes.NewReader(data), "application/json", nil
	}

	values := url.Values{}
	values.Set(usernameKey, email)
	values.Set(passwordKey, password)
	return loginURL, strings.NewReader(values.Encode()), "application/x-www-form-urlencoded", nil
}

func defaultConfigValue(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
