package service

import (
	"bytes"
	"context"
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
	"github.com/1024XEngineer/xinfra/server/internal/wayne"
)

var (
	ErrWayneRoleBindingNotConfigured = errors.New("wayne internal role binding api is not configured")
	ErrWayneRoleBindingRequestFailed = errors.New("wayne internal role binding request failed")
)

type WayneRoleBindingRequest struct {
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
	cfg    config.Config
	client *http.Client
	now    func() time.Time
}

func NewWayneRoleBindingService(cfg config.Config) *WayneRoleBindingService {
	return &WayneRoleBindingService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		now: time.Now,
	}
}

func (s *WayneRoleBindingService) BindNamespace(ctx context.Context, namespaceID uint64, username string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.call(ctx, http.MethodPut, fmt.Sprintf("/api/v1/internal/namespaces/%d/users/%s/roles", namespaceID, url.PathEscape(username)), username, req)
}

func (s *WayneRoleBindingService) UnbindNamespace(ctx context.Context, namespaceID uint64, username string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.call(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/internal/namespaces/%d/users/%s/roles", namespaceID, url.PathEscape(username)), username, req)
}

func (s *WayneRoleBindingService) BindApp(ctx context.Context, appID uint64, username string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.call(ctx, http.MethodPut, fmt.Sprintf("/api/v1/internal/apps/%d/users/%s/roles", appID, url.PathEscape(username)), username, req)
}

func (s *WayneRoleBindingService) UnbindApp(ctx context.Context, appID uint64, username string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	return s.call(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/internal/apps/%d/users/%s/roles", appID, url.PathEscape(username)), username, req)
}

func (s *WayneRoleBindingService) ListNamespaces(ctx context.Context) (*WayneRoleBindingResult, error) {
	return s.callRaw(ctx, http.MethodGet, "/api/v1/internal/namespaces", nil)
}

func (s *WayneRoleBindingService) ListGroups(ctx context.Context, groupType *int) (*WayneRoleBindingResult, error) {
	internalPath := "/api/v1/internal/groups"
	if groupType != nil {
		values := url.Values{}
		values.Set("type", strconv.Itoa(*groupType))
		internalPath += "?" + values.Encode()
	}
	return s.callRaw(ctx, http.MethodGet, internalPath, nil)
}

func (s *WayneRoleBindingService) GetUserRoles(ctx context.Context, username string) (*WayneRoleBindingResult, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrWayenEmailMissing
	}
	return s.callRaw(ctx, http.MethodGet, fmt.Sprintf("/api/v1/internal/users/%s/roles", url.PathEscape(username)), nil)
}

func (s *WayneRoleBindingService) NamespaceOperatorPermissions(ctx context.Context, namespaceID uint64, operatorEmail string) (*WayneRoleBindingResult, error) {
	return s.operatorPermissions(ctx, fmt.Sprintf("/api/v1/internal/namespaces/%d/operator-permissions", namespaceID), operatorEmail)
}

func (s *WayneRoleBindingService) AppOperatorPermissions(ctx context.Context, appID uint64, operatorEmail string) (*WayneRoleBindingResult, error) {
	return s.operatorPermissions(ctx, fmt.Sprintf("/api/v1/internal/apps/%d/operator-permissions", appID), operatorEmail)
}

func (s *WayneRoleBindingService) call(ctx context.Context, method, internalPath, operatorEmail string, req WayneRoleBindingRequest) (*WayneRoleBindingResult, error) {
	if err := s.validateConfig(); err != nil {
		return nil, err
	}
	operatorEmail = strings.TrimSpace(operatorEmail)
	if operatorEmail == "" {
		return nil, ErrWayenEmailMissing
	}

	req.OperatorUserID = nil
	req.OperatorName = operatorEmail

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return s.callRaw(ctx, method, internalPath, body)
}

func (s *WayneRoleBindingService) operatorPermissions(ctx context.Context, internalPath, operatorEmail string) (*WayneRoleBindingResult, error) {
	operatorEmail = strings.TrimSpace(operatorEmail)
	if operatorEmail == "" {
		return nil, ErrWayenEmailMissing
	}
	values := url.Values{}
	values.Set("operatorName", operatorEmail)
	return s.callRaw(ctx, http.MethodGet, internalPath+"?"+values.Encode(), nil)
}

func (s *WayneRoleBindingService) callRaw(ctx context.Context, method, internalPath string, body []byte) (*WayneRoleBindingResult, error) {
	if err := s.validateConfig(); err != nil {
		return nil, err
	}
	if body == nil {
		body = []byte{}
	}
	target, signingURI, err := s.requestURL(internalPath)
	if err != nil {
		return nil, err
	}
	log.Printf(
		"wayne role binding request: method=%s target=%s signing_uri=%s body_bytes=%d service_name=%s secret_configured=%t",
		method,
		target,
		signingURI,
		len(body),
		s.cfg.WayneServiceName,
		strings.TrimSpace(s.cfg.WayneServiceAPISecretKey) != "",
	)
	httpReq, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")

	headers, err := wayne.BuildSignedHeaders(s.cfg.WayneServiceName, s.cfg.WayneServiceAPISecretKey, method, signingURI, body, s.now())
	if err != nil {
		return nil, err
	}
	headers.Apply(httpReq)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		log.Printf("wayne role binding request failed: method=%s target=%s error=%v", method, target, err)
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
	log.Printf(
		"wayne role binding response: method=%s target=%s status=%d content_type=%q body=%q",
		method,
		target,
		resp.StatusCode,
		result.ContentType,
		truncateForDebugLog(string(respBody), 512),
	)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return result, &WayneRoleBindingHTTPError{StatusCode: resp.StatusCode, Body: respBody}
	}
	return result, nil
}

func (s *WayneRoleBindingService) validateConfig() error {
	baseConfigured := strings.TrimSpace(s.cfg.WayneInternalAPIBaseURL) != ""
	serviceConfigured := strings.TrimSpace(s.cfg.WayneServiceName) != ""
	secretConfigured := strings.TrimSpace(s.cfg.WayneServiceAPISecretKey) != ""
	if !baseConfigured || !serviceConfigured || !secretConfigured {
		log.Printf(
			"wayne role binding config invalid: base_url_configured=%t service_name_configured=%t secret_configured=%t",
			baseConfigured,
			serviceConfigured,
			secretConfigured,
		)
		return ErrWayneRoleBindingNotConfigured
	}
	return nil
}

func (s *WayneRoleBindingService) requestURL(internalPath string) (string, string, error) {
	base, err := url.Parse(strings.TrimRight(strings.TrimSpace(s.cfg.WayneInternalAPIBaseURL), "/"))
	if err != nil {
		return "", "", err
	}
	if base.Scheme == "" || base.Host == "" {
		return "", "", fmt.Errorf("invalid wayne internal api base url: %s", s.cfg.WayneInternalAPIBaseURL)
	}
	path, rawQuery, _ := strings.Cut(internalPath, "?")
	base.Path = strings.TrimRight(base.Path, "/") + path
	base.RawQuery = rawQuery
	return base.String(), base.RequestURI(), nil
}

func truncateForDebugLog(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit] + "...(truncated)"
}
