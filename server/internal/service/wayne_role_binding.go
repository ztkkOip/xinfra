package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrWayneRoleBindingNotConfigured = errors.New("wayne native api is not configured")
	ErrWayneRoleBindingRequestFailed = errors.New("wayne native role binding request failed")
)

const wayneAdminTokenAccount = "admin"

type WayneRoleBindingRequest struct {
	Username       string   `json:"username,omitempty"`
	GroupIDs       []uint64 `json:"groupIds,omitempty"`
	OperatorUserID *uint64  `json:"operatorUserId,omitempty"`
	OperatorName   string   `json:"operatorName,omitempty"`
	Replace        *bool    `json:"replace,omitempty"`
	RequestID      string   `json:"requestId,omitempty"`
	Reason         string   `json:"reason,omitempty"`
	DryRun         bool     `json:"dryRun,omitempty"`
}

type WayneRoleBindingResult struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type WayneRoleGroup struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Type    int    `json:"type"`
}

type WayneOperatorPermissions struct {
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type WayneRoleBindingHTTPError struct {
	StatusCode int
	Body       []byte
}

func (e *WayneRoleBindingHTTPError) Error() string {
	body := strings.TrimSpace(string(e.Body))
	if body == "" {
		return fmt.Sprintf("%s: status %d", ErrWayneRoleBindingRequestFailed, e.StatusCode)
	}
	if len(body) > 512 {
		body = body[:512]
	}
	return fmt.Sprintf("%s: status %d: %s", ErrWayneRoleBindingRequestFailed, e.StatusCode, body)
}

type WayneRoleBindingService struct {
	cfg         config.Config
	db          *gorm.DB
	client      *http.Client
	now         func() time.Time
	cachedToken string
	cachedExp   time.Time
}

func NewWayneRoleBindingService(cfg config.Config, db *gorm.DB) *WayneRoleBindingService {
	return &WayneRoleBindingService{
		cfg: cfg,
		db:  db,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		now: time.Now,
	}
}

func (s *WayneRoleBindingService) BindNamespace(ctx context.Context, namespaceID uint64, username string, operatorEmail string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.bindUserRoles(ctx, "namespace", namespaceID, username, req)
}

func (s *WayneRoleBindingService) UnbindNamespace(ctx context.Context, namespaceID uint64, username string, operatorEmail string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.unbindUserRoles(ctx, "namespace", namespaceID, username, req)
}

func (s *WayneRoleBindingService) BindApp(ctx context.Context, appID uint64, username string, operatorEmail string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.bindUserRoles(ctx, "app", appID, username, req)
}

func (s *WayneRoleBindingService) UnbindApp(ctx context.Context, appID uint64, username string, operatorEmail string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.unbindUserRoles(ctx, "app", appID, username, req)
}

func (s *WayneRoleBindingService) ListNamespaces(ctx context.Context) (*WayneRoleBindingResult, error) {
	return s.callRaw(ctx, http.MethodGet, "/api/v1/namespaces", nil)
}

func (s *WayneRoleBindingService) ListGroups(ctx context.Context, groupType *int) (*WayneRoleBindingResult, error) {
	path := "/api/v1/groups"
	if groupType != nil {
		values := url.Values{}
		values.Set("type", strconv.Itoa(*groupType))
		path += "?" + values.Encode()
	}
	return s.callRaw(ctx, http.MethodGet, path, nil)
}

func (s *WayneRoleBindingService) ListNamespaceRoleGroups(ctx context.Context) ([]WayneRoleGroup, error) {
	groupType := 1
	result, err := s.ListGroups(ctx, &groupType)
	if err != nil {
		return nil, err
	}
	return parseWayneRoleGroups(result.Body)
}

func (s *WayneRoleBindingService) NamespaceVisitorGroupIDs(ctx context.Context) ([]uint64, error) {
	groups, err := s.ListNamespaceRoleGroups(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, 1)
	for _, group := range groups {
		if isWayneVisitorRoleName(group.Name) {
			ids = append(ids, group.ID)
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("wayne visitor role group not found")
	}
	return ids, nil
}

func (s *WayneRoleBindingService) GetUserRoles(ctx context.Context, username string) (*WayneRoleBindingResult, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrWayenEmailMissing
	}
	user, err := s.findUser(ctx, username)
	if err != nil {
		return nil, err
	}

	namespaces, err := s.listResourceIDs(ctx, "/api/v1/namespaces")
	if err != nil {
		return nil, err
	}
	results := make([]json.RawMessage, 0, len(namespaces))
	values := url.Values{}
	values.Set("userId", strconv.FormatUint(user.ID, 10))
	for _, namespaceID := range namespaces {
		result, err := s.callRaw(ctx, http.MethodGet, roleBindingBasePath("namespace", namespaceID)+"?"+values.Encode(), nil)
		if err != nil {
			return result, err
		}
		results = append(results, rawWayneData(result.Body))
	}
	body, err := json.Marshal(map[string]json.RawMessage{
		"namespaceRoles": mustMarshalRaw(results),
	})
	if err != nil {
		return nil, err
	}
	return &WayneRoleBindingResult{StatusCode: http.StatusOK, ContentType: "application/json", Body: body}, nil
}

func (s *WayneRoleBindingService) listResourceIDs(ctx context.Context, path string) ([]uint64, error) {
	result, err := s.callRaw(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return parseWayneResourceIDs(result.Body)
}

func (s *WayneRoleBindingService) NamespaceOperatorPermissions(ctx context.Context, namespaceID uint64, operatorEmail string) (*WayneRoleBindingResult, error) {
	return s.permissions(ctx, "namespace", namespaceID)
}

func (s *WayneRoleBindingService) NamespaceOperatorPermissionsParsed(ctx context.Context, namespaceID uint64, operatorEmail string) (*WayneOperatorPermissions, error) {
	result, err := s.NamespaceOperatorPermissions(ctx, namespaceID, operatorEmail)
	if err != nil {
		return nil, err
	}
	return parseWayneOperatorPermissions(result.Body)
}

func (s *WayneRoleBindingService) AppOperatorPermissions(ctx context.Context, appID uint64, operatorEmail string) (*WayneRoleBindingResult, error) {
	return s.permissions(ctx, "app", appID)
}

func (s *WayneRoleBindingService) bindUserRoles(ctx context.Context, scope string, resourceID uint64, username string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	if strings.TrimSpace(username) == "" {
		return nil, ErrWayenEmailMissing
	}
	user, err := s.findUser(ctx, username)
	if err != nil {
		return nil, err
	}
	bindings, err := s.listUserBindings(ctx, scope, resourceID, user.ID)
	if err != nil {
		return nil, err
	}

	body := nativeRoleBindingBody(scope, resourceID, user.ID, req.GroupIDs)
	if len(bindings) == 0 {
		return s.callRaw(ctx, http.MethodPost, roleBindingBasePath(scope, resourceID), body)
	}
	bodyWithID, err := nativeRoleBindingBodyWithID(body, bindings[0].ID)
	if err != nil {
		return nil, err
	}
	return s.callRaw(ctx, http.MethodPut, fmt.Sprintf("%s/%d", roleBindingBasePath(scope, resourceID), bindings[0].ID), bodyWithID)
}

func (s *WayneRoleBindingService) unbindUserRoles(ctx context.Context, scope string, resourceID uint64, username string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	if strings.TrimSpace(username) == "" {
		return nil, ErrWayenEmailMissing
	}
	user, err := s.findUser(ctx, username)
	if err != nil {
		return nil, err
	}
	bindings, err := s.listUserBindings(ctx, scope, resourceID, user.ID)
	if err != nil {
		return nil, err
	}
	if len(bindings) == 0 {
		return &WayneRoleBindingResult{StatusCode: http.StatusOK, ContentType: "application/json", Body: []byte(`{"data":null}`)}, nil
	}
	return s.callRaw(ctx, http.MethodDelete, fmt.Sprintf("%s/%d", roleBindingBasePath(scope, resourceID), bindings[0].ID), nil)
}

func (s *WayneRoleBindingService) permissions(ctx context.Context, scope string, resourceID uint64) (*WayneRoleBindingResult, error) {
	path := fmt.Sprintf("%s/permissions/%d", roleBindingBasePath(scope, resourceID), resourceID)
	result, err := s.callRaw(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	permissions, err := parseNativePermissions(result.Body, scope)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(map[string]WayneOperatorPermissions{"permissions": permissions})
	if err != nil {
		return nil, err
	}
	return &WayneRoleBindingResult{StatusCode: result.StatusCode, ContentType: "application/json", Body: body}, nil
}

func (s *WayneRoleBindingService) findUser(ctx context.Context, username string) (*wayneNativeUser, error) {
	values := url.Values{}
	values.Set("name", strings.TrimSpace(username))
	result, err := s.callRaw(ctx, http.MethodGet, "/api/v1/users?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	users, err := parseWayneUsers(result.Body)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if sameWayneUser(user, username) {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("wayne user %q not found", username)
}

func (s *WayneRoleBindingService) listUserBindings(ctx context.Context, scope string, resourceID uint64, userID uint64) ([]wayneRoleBinding, error) {
	values := url.Values{}
	values.Set("userId", strconv.FormatUint(userID, 10))
	result, err := s.callRaw(ctx, http.MethodGet, roleBindingBasePath(scope, resourceID)+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return parseWayneRoleBindings(result.Body)
}

func (s *WayneRoleBindingService) callRaw(ctx context.Context, method, nativePath string, body []byte) (*WayneRoleBindingResult, error) {
	token, err := s.adminToken(ctx, false)
	if err != nil {
		return nil, err
	}
	result, err := s.callRawWithToken(ctx, method, nativePath, body, token)
	if result != nil && result.StatusCode == http.StatusUnauthorized {
		token, refreshErr := s.adminToken(ctx, true)
		if refreshErr != nil {
			return result, refreshErr
		}
		return s.callRawWithToken(ctx, method, nativePath, body, token)
	}
	return result, err
}

func (s *WayneRoleBindingService) callRawWithToken(ctx context.Context, method, nativePath string, body []byte, token string) (*WayneRoleBindingResult, error) {
	if body == nil {
		body = []byte{}
	}
	target, err := s.requestURL(nativePath)
	if err != nil {
		return nil, err
	}
	log.Printf("wayne native request: method=%s target=%s body_bytes=%d", method, target, len(body))
	httpReq, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		log.Printf("wayne native request failed: method=%s target=%s error=%v", method, target, err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	result := &WayneRoleBindingResult{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        respBody,
	}
	log.Printf("wayne native response: method=%s target=%s status=%d content_type=%q body=%q", method, target, resp.StatusCode, result.ContentType, truncateForDebugLog(string(respBody), 512))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return result, &WayneRoleBindingHTTPError{StatusCode: resp.StatusCode, Body: respBody}
	}
	return result, nil
}

func (s *WayneRoleBindingService) adminToken(ctx context.Context, forceRefresh bool) (string, error) {
	if err := s.validateConfig(); err != nil {
		return "", err
	}
	now := s.now()
	if !forceRefresh && strings.TrimSpace(s.cachedToken) != "" && s.cachedExp.After(now.Add(time.Minute)) {
		return s.cachedToken, nil
	}
	if !forceRefresh && s.db != nil {
		var saved model.WayneToken
		err := s.db.WithContext(ctx).Where("account = ?", wayneAdminTokenAccount).First(&saved).Error
		if err == nil && strings.TrimSpace(saved.Token) != "" && saved.ExpiresAt.After(now.Add(time.Minute)) {
			s.cachedToken = saved.Token
			s.cachedExp = saved.ExpiresAt
			return saved.Token, nil
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}
	token, expiresAt, err := s.loginAdmin(ctx)
	if err != nil {
		return "", err
	}
	s.cachedToken = token
	s.cachedExp = expiresAt
	if s.db != nil {
		err = s.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "account"}},
			DoUpdates: clause.AssignmentColumns([]string{"token", "expires_at", "updated_at"}),
		}).Create(&model.WayneToken{
			Account:   wayneAdminTokenAccount,
			Token:     token,
			ExpiresAt: expiresAt,
		}).Error
	}
	return token, err
}

func (s *WayneRoleBindingService) loginAdmin(ctx context.Context) (string, time.Time, error) {
	body, err := json.Marshal(map[string]string{
		"username": s.cfg.WayneAdminUsername,
		"password": s.cfg.WayneAdminPassword,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	target, err := s.requestURL("/login/db")
	if err != nil {
		return "", time.Time{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", time.Time{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return "", time.Time{}, &WayneRoleBindingHTTPError{StatusCode: resp.StatusCode, Body: respBody}
	}
	token, err := parseWayneLoginToken(respBody)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := wayneTokenExpiresAt(token, s.now().Add(time.Duration(s.cfg.WayneTokenTTLMinutes)*time.Minute))
	return token, expiresAt, nil
}

func (s *WayneRoleBindingService) validateConfig() error {
	baseConfigured := strings.TrimSpace(s.cfg.WayneAPIBaseURL) != ""
	usernameConfigured := strings.TrimSpace(s.cfg.WayneAdminUsername) != ""
	passwordConfigured := strings.TrimSpace(s.cfg.WayneAdminPassword) != ""
	if !baseConfigured || !usernameConfigured || !passwordConfigured {
		log.Printf(
			"wayne native config invalid: base_url_configured=%t admin_username_configured=%t admin_password_configured=%t",
			baseConfigured,
			usernameConfigured,
			passwordConfigured,
		)
		return ErrWayneRoleBindingNotConfigured
	}
	return nil
}

func (s *WayneRoleBindingService) requestURL(nativePath string) (string, error) {
	base, err := url.Parse(strings.TrimRight(strings.TrimSpace(s.cfg.WayneAPIBaseURL), "/"))
	if err != nil {
		return "", err
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("invalid wayne api base url: %s", s.cfg.WayneAPIBaseURL)
	}
	path, rawQuery, _ := strings.Cut(nativePath, "?")
	base.Path = strings.TrimRight(base.Path, "/") + path
	base.RawQuery = rawQuery
	return base.String(), nil
}

func roleBindingBasePath(scope string, resourceID uint64) string {
	switch scope {
	case "app":
		return fmt.Sprintf("/api/v1/apps/%d/users", resourceID)
	default:
		return fmt.Sprintf("/api/v1/namespaces/%d/users", resourceID)
	}
}

func nativeRoleBindingBody(scope string, resourceID, userID uint64, groupIDs []uint64) []byte {
	groups := make([]map[string]uint64, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		groups = append(groups, map[string]uint64{"id": groupID})
	}
	body := map[string]any{
		"user":   map[string]uint64{"id": userID},
		"groups": groups,
	}
	if scope == "app" {
		body["app"] = map[string]uint64{"id": resourceID}
	} else {
		body["namespace"] = map[string]uint64{"id": resourceID}
	}
	raw, _ := json.Marshal(body)
	return raw
}

func nativeRoleBindingBodyWithID(body []byte, id uint64) ([]byte, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	payload["id"] = id
	return json.Marshal(payload)
}

func truncateForDebugLog(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit] + "...(truncated)"
}

func parseWayneLoginToken(body []byte) (string, error) {
	var wrapped struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return "", err
	}
	token := strings.TrimSpace(wrapped.Data.Token)
	if token == "" {
		token = strings.TrimSpace(wrapped.Token)
	}
	if token == "" {
		return "", fmt.Errorf("wayne login response token is empty")
	}
	return token, nil
}

func wayneTokenExpiresAt(token string, fallback time.Time) time.Time {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return fallback
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fallback
	}
	var claims struct {
		ExpiresAt int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.ExpiresAt <= 0 {
		return fallback
	}
	return time.Unix(claims.ExpiresAt, 0)
}

func parseWayneRoleGroups(body []byte) ([]WayneRoleGroup, error) {
	var wrapped struct {
		Data struct {
			List []WayneRoleGroup `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data.List, nil
}

func parseWayneOperatorPermissions(body []byte) (*WayneOperatorPermissions, error) {
	var wrapped struct {
		Permissions WayneOperatorPermissions `json:"permissions"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, err
	}
	return &wrapped.Permissions, nil
}

func parseNativePermissions(body []byte, scope string) (WayneOperatorPermissions, error) {
	var wrapped struct {
		Data map[string]map[string]bool `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return WayneOperatorPermissions{}, err
	}
	key := "namespaceUser"
	if scope == "app" {
		key = "appUser"
	}
	return WayneOperatorPermissions{
		Create: wrapped.Data[key]["create"],
		Update: wrapped.Data[key]["update"],
		Delete: wrapped.Data[key]["delete"],
	}, nil
}

func rawWayneData(body []byte) json.RawMessage {
	var wrapped struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && len(wrapped.Data) > 0 {
		return wrapped.Data
	}
	return json.RawMessage(body)
}

func mustMarshalRaw(value any) json.RawMessage {
	raw, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`null`)
	}
	return raw
}

func parseWayneResourceIDs(body []byte) ([]uint64, error) {
	var wrapped struct {
		Data struct {
			List []struct {
				ID uint64 `json:"id"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, err
	}
	return resourceIDsFromItems(wrapped.Data.List), nil
}

func resourceIDsFromItems(items []struct {
	ID uint64 `json:"id"`
}) []uint64 {
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.ID != 0 {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

type wayneNativeUser struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Display string `json:"display"`
}

func parseWayneUsers(body []byte) ([]wayneNativeUser, error) {
	var wrapped struct {
		Data struct {
			List []wayneNativeUser `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data.List, nil
}

func sameWayneUser(user wayneNativeUser, username string) bool {
	username = strings.TrimSpace(strings.ToLower(username))
	return strings.ToLower(strings.TrimSpace(user.Name)) == username ||
		strings.ToLower(strings.TrimSpace(user.Email)) == username
}

type wayneRoleBinding struct {
	ID uint64 `json:"id"`
}

func parseWayneRoleBindings(body []byte) ([]wayneRoleBinding, error) {
	var wrapped struct {
		Data struct {
			List []wayneRoleBinding `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data.List, nil
}

func isWayneVisitorRoleName(name string) bool {
	normalized := strings.ToLower(strings.TrimSpace(name))
	return normalized == "访客" || normalized == "visitor" || strings.Contains(normalized, "visitor")
}
