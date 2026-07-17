package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/wayne"
)

func TestWayneRoleBindingServiceBindNamespaceSignsAndOverridesOperator(t *testing.T) {
	var requestPath string
	var payload WayneRoleBindingRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath = r.URL.RequestURI()
		body := readTestBody(t, r)
		if !wayne.Verify("service-secret", r.Header.Get(wayne.HeaderSignature), r.Method, r.URL.RequestURI(), r.Header.Get(wayne.HeaderTimestamp), r.Header.Get(wayne.HeaderNonce), body) {
			t.Fatalf("invalid signature headers: %#v body=%s", r.Header, string(body))
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("invalid request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"changed":true}}`))
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(config.Config{
		WayneInternalAPIBaseURL:  server.URL,
		WayneServiceName:         "xinfra",
		WayneServiceAPISecretKey: "service-secret",
	})
	svc.now = func() time.Time { return time.Unix(1721000000, 0) }

	operatorUserID := uint64(123)
	replace := false
	result, err := svc.BindNamespace(context.Background(), 1, "target@example.com", "eastsales@qiniu.com", WayneRoleBindingRequest{
		Username:       "target@example.com",
		GroupIDs:       []uint64{10, 11},
		OperatorUserID: &operatorUserID,
		OperatorName:   "attacker@example.com",
		Replace:        &replace,
		RequestID:      "req-001",
		Reason:         "grant",
	})
	if err != nil {
		t.Fatalf("BindNamespace() error = %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want 200", result.StatusCode)
	}
	if requestPath != "/api/v1/internal/namespaces/1/users/target@example.com/roles" {
		t.Fatalf("requestPath = %q", requestPath)
	}
	if payload.OperatorName != "eastsales@qiniu.com" {
		t.Fatalf("OperatorName = %q, want token email", payload.OperatorName)
	}
	if payload.OperatorUserID != nil {
		t.Fatalf("OperatorUserID should be omitted, got %v", *payload.OperatorUserID)
	}
	if payload.Username != "" {
		t.Fatalf("Username should be omitted from Wayne body, got %q", payload.Username)
	}
	if payload.Replace == nil || *payload.Replace {
		t.Fatalf("Replace = %v, want false", payload.Replace)
	}
}

func TestWayneRoleBindingServiceCallsAllDocumentedEndpoints(t *testing.T) {
	tests := []struct {
		name string
		call func(*WayneRoleBindingService) (*WayneRoleBindingResult, error)
		want string
	}{
		{
			name: "unbind namespace",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.UnbindNamespace(context.Background(), 1, "target@example.com", "eastsales@qiniu.com", WayneRoleBindingRequest{GroupIDs: []uint64{10}})
			},
			want: "DELETE /api/v1/internal/namespaces/1/users/target@example.com/roles",
		},
		{
			name: "bind app",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.BindApp(context.Background(), 3, "target@example.com", "eastsales@qiniu.com", WayneRoleBindingRequest{GroupIDs: []uint64{20}})
			},
			want: "PUT /api/v1/internal/apps/3/users/target@example.com/roles",
		},
		{
			name: "unbind app",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.UnbindApp(context.Background(), 3, "target@example.com", "eastsales@qiniu.com", WayneRoleBindingRequest{GroupIDs: []uint64{20}})
			},
			want: "DELETE /api/v1/internal/apps/3/users/target@example.com/roles",
		},
		{
			name: "list namespaces",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.ListNamespaces(context.Background())
			},
			want: "GET /api/v1/internal/namespaces",
		},
		{
			name: "list namespace groups",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				groupType := 1
				return s.ListGroups(context.Background(), &groupType)
			},
			want: "GET /api/v1/internal/groups?type=1",
		},
		{
			name: "list all groups",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.ListGroups(context.Background(), nil)
			},
			want: "GET /api/v1/internal/groups",
		},
		{
			name: "get user roles",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.GetUserRoles(context.Background(), "target@example.com")
			},
			want: "GET /api/v1/internal/users/target@example.com/roles",
		},
		{
			name: "namespace operator permissions",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.NamespaceOperatorPermissions(context.Background(), 1, "eastsales@qiniu.com")
			},
			want: "GET /api/v1/internal/namespaces/1/operator-permissions?operatorName=eastsales%40qiniu.com",
		},
		{
			name: "app operator permissions",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.AppOperatorPermissions(context.Background(), 3, "eastsales@qiniu.com")
			},
			want: "GET /api/v1/internal/apps/3/operator-permissions?operatorName=eastsales%40qiniu.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.Method + " " + r.URL.RequestURI()
				_, _ = w.Write([]byte(`{"data":{"changed":true}}`))
			}))
			defer server.Close()

			svc := NewWayneRoleBindingService(config.Config{
				WayneInternalAPIBaseURL:  server.URL,
				WayneServiceName:         "xinfra",
				WayneServiceAPISecretKey: "service-secret",
			})
			if _, err := tt.call(svc); err != nil {
				t.Fatalf("call error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("got endpoint %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWayneRoleBindingServiceQuerySignsEmptyBodyAndQueryURI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := readTestBody(t, r)
		if len(body) != 0 {
			t.Fatalf("GET body length = %d, want 0", len(body))
		}
		if r.URL.RequestURI() != "/api/v1/internal/groups?type=1" {
			t.Fatalf("RequestURI = %q", r.URL.RequestURI())
		}
		if !wayne.Verify("service-secret", r.Header.Get(wayne.HeaderSignature), r.Method, r.URL.RequestURI(), r.Header.Get(wayne.HeaderTimestamp), r.Header.Get(wayne.HeaderNonce), body) {
			t.Fatalf("invalid GET signature headers: %#v", r.Header)
		}
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(config.Config{
		WayneInternalAPIBaseURL:  server.URL,
		WayneServiceName:         "xinfra",
		WayneServiceAPISecretKey: "service-secret",
	})
	groupType := 1
	if _, err := svc.ListGroups(context.Background(), &groupType); err != nil {
		t.Fatalf("ListGroups() error = %v", err)
	}
}

func TestWayneRoleBindingServiceOperatorPermissionsSignsQueryURI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := readTestBody(t, r)
		if r.URL.RequestURI() != "/api/v1/internal/namespaces/1/operator-permissions?operatorName=eastsales%40qiniu.com" {
			t.Fatalf("RequestURI = %q", r.URL.RequestURI())
		}
		if !wayne.Verify("service-secret", r.Header.Get(wayne.HeaderSignature), r.Method, r.URL.RequestURI(), r.Header.Get(wayne.HeaderTimestamp), r.Header.Get(wayne.HeaderNonce), body) {
			t.Fatalf("invalid operator permissions signature headers: %#v", r.Header)
		}
		_, _ = w.Write([]byte(`{"data":{"permissions":{"create":true,"update":true,"delete":false}}}`))
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(config.Config{
		WayneInternalAPIBaseURL:  server.URL,
		WayneServiceName:         "xinfra",
		WayneServiceAPISecretKey: "service-secret",
	})
	if _, err := svc.NamespaceOperatorPermissions(context.Background(), 1, "eastsales@qiniu.com"); err != nil {
		t.Fatalf("NamespaceOperatorPermissions() error = %v", err)
	}
}

func TestWayneRoleBindingServiceHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"code":403,"msg":"denied"}`))
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(config.Config{
		WayneInternalAPIBaseURL:  server.URL,
		WayneServiceName:         "xinfra",
		WayneServiceAPISecretKey: "service-secret",
	})
	result, err := svc.BindApp(context.Background(), 3, "target@example.com", "eastsales@qiniu.com", WayneRoleBindingRequest{GroupIDs: []uint64{20}})
	if err == nil {
		t.Fatal("expected error")
	}
	if result == nil || result.StatusCode != http.StatusForbidden {
		t.Fatalf("result = %#v, want 403", result)
	}
	if !strings.Contains(err.Error(), "denied") {
		t.Fatalf("error = %q, want denied body", err.Error())
	}
}

func readTestBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return body
}
