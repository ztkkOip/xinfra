package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
)

func TestWayneRoleBindingServiceBindNamespaceUsesNativeAPI(t *testing.T) {
	var got []string
	var updateBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Method+" "+r.URL.RequestURI())
		switch r.URL.RequestURI() {
		case "/login/db":
			assertLoginBody(t, r)
			_, _ = w.Write([]byte(`{"data":{"token":"` + testJWT(time.Now().Add(time.Hour)) + `"}}`))
		case "/api/v1/users?name=target%40example.com":
			assertBearer(t, r)
			_, _ = w.Write([]byte(`{"data":{"list":[{"id":7,"name":"target@example.com","email":"target@example.com"}]}}`))
		case "/api/v1/namespaces/1/users?userId=7":
			assertBearer(t, r)
			_, _ = w.Write([]byte(`{"data":{"list":[{"id":99}]}}`))
		case "/api/v1/namespaces/1/users/99":
			assertBearer(t, r)
			body := readTestBody(t, r)
			if err := json.Unmarshal(body, &updateBody); err != nil {
				t.Fatalf("invalid update body: %v", err)
			}
			_, _ = w.Write([]byte(`{"data":{"changed":true}}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(testWayneConfig(server.URL), nil)
	result, err := svc.BindNamespace(context.Background(), 1, "target@example.com", "operator@example.com", WayneRoleBindingRequest{GroupIDs: []uint64{10, 11}})
	if err != nil {
		t.Fatalf("BindNamespace() error = %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want 200", result.StatusCode)
	}
	want := []string{
		"POST /login/db",
		"GET /api/v1/users?name=target%40example.com",
		"GET /api/v1/namespaces/1/users?userId=7",
		"PUT /api/v1/namespaces/1/users/99",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if updateBody["id"].(float64) != 99 {
		t.Fatalf("update id = %v, want 99", updateBody["id"])
	}
	if updateBody["namespace"].(map[string]any)["id"].(float64) != 1 {
		t.Fatalf("namespace body = %#v", updateBody["namespace"])
	}
	if updateBody["user"].(map[string]any)["id"].(float64) != 7 {
		t.Fatalf("user body = %#v", updateBody["user"])
	}
}

func TestParseWayneUsersSupportsNativePageList(t *testing.T) {
	users, err := parseWayneUsers([]byte(`{"data":{"pageNo":1,"pageSize":10,"totalCount":1,"list":[{"id":7,"name":"target@example.com","email":"target@example.com"}]}}`))
	if err != nil {
		t.Fatalf("parseWayneUsers() error = %v", err)
	}
	if len(users) != 1 || users[0].ID != 7 || users[0].Name != "target@example.com" {
		t.Fatalf("users = %#v", users)
	}
}

func TestWayneRoleBindingServiceCreateAndDeleteNativeBindings(t *testing.T) {
	tests := []struct {
		name string
		call func(*WayneRoleBindingService) (*WayneRoleBindingResult, error)
		want string
	}{
		{
			name: "create app binding",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.BindApp(context.Background(), 3, "target@example.com", "operator@example.com", WayneRoleBindingRequest{GroupIDs: []uint64{20}})
			},
			want: "POST /api/v1/apps/3/users",
		},
		{
			name: "delete namespace binding",
			call: func(s *WayneRoleBindingService) (*WayneRoleBindingResult, error) {
				return s.UnbindNamespace(context.Background(), 1, "target@example.com", "operator@example.com", WayneRoleBindingRequest{})
			},
			want: "DELETE /api/v1/namespaces/1/users/44",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var finalRequest string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case r.URL.RequestURI() == "/login/db":
					_, _ = w.Write([]byte(`{"data":{"token":"` + testJWT(time.Now().Add(time.Hour)) + `"}}`))
				case strings.HasPrefix(r.URL.RequestURI(), "/api/v1/users?"):
					_, _ = w.Write([]byte(`{"data":{"list":[{"id":7,"name":"target@example.com"}]}}`))
				case strings.Contains(r.URL.RequestURI(), "userId=7"):
					if strings.Contains(tt.want, "POST ") {
						_, _ = w.Write([]byte(`{"data":{"list":[]}}`))
					} else {
						_, _ = w.Write([]byte(`{"data":{"list":[{"id":44}]}}`))
					}
				default:
					finalRequest = r.Method + " " + r.URL.RequestURI()
					_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
				}
			}))
			defer server.Close()

			svc := NewWayneRoleBindingService(testWayneConfig(server.URL), nil)
			if _, err := tt.call(svc); err != nil {
				t.Fatalf("call error = %v", err)
			}
			if finalRequest != tt.want {
				t.Fatalf("final request = %q, want %q", finalRequest, tt.want)
			}
		})
	}
}

func TestWayneRoleBindingServiceRefreshesTokenAfterUnauthorized(t *testing.T) {
	loginCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/login/db":
			loginCount++
			_, _ = w.Write([]byte(`{"data":{"token":"` + testJWT(time.Now().Add(time.Hour)) + `"}}`))
		case "/api/v1/groups?type=1":
			if loginCount == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"msg":"expired"}`))
				return
			}
			_, _ = w.Write([]byte(`{"data":{"list":[]}}`))
		default:
			t.Fatalf("unexpected request %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(testWayneConfig(server.URL), nil)
	groupType := 1
	if _, err := svc.ListGroups(context.Background(), &groupType); err != nil {
		t.Fatalf("ListGroups() error = %v", err)
	}
	if loginCount != 2 {
		t.Fatalf("loginCount = %d, want 2", loginCount)
	}
}

func TestWayneRoleBindingServiceOperatorPermissionsFromNativePermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/login/db":
			_, _ = w.Write([]byte(`{"data":{"token":"` + testJWT(time.Now().Add(time.Hour)) + `"}}`))
		case "/api/v1/namespaces/1/users/permissions/1":
			_, _ = w.Write([]byte(`{"data":{"namespaceUser":{"create":true,"update":true,"delete":false}}}`))
		default:
			t.Fatalf("unexpected request %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(testWayneConfig(server.URL), nil)
	permissions, err := svc.NamespaceOperatorPermissionsParsed(context.Background(), 1, "operator@example.com")
	if err != nil {
		t.Fatalf("NamespaceOperatorPermissionsParsed() error = %v", err)
	}
	if !permissions.Create || !permissions.Update || permissions.Delete {
		t.Fatalf("permissions = %#v", permissions)
	}
}

func TestWayneRoleBindingServiceHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/login/db":
			_, _ = w.Write([]byte(`{"data":{"token":"` + testJWT(time.Now().Add(time.Hour)) + `"}}`))
		case "/api/v1/groups":
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"code":403,"msg":"denied"}`))
		default:
			t.Fatalf("unexpected request %s", r.URL.RequestURI())
		}
	}))
	defer server.Close()

	svc := NewWayneRoleBindingService(testWayneConfig(server.URL), nil)
	result, err := svc.ListGroups(context.Background(), nil)
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

func testWayneConfig(baseURL string) config.Config {
	return config.Config{
		WayneAPIBaseURL:      baseURL,
		WayneAdminUsername:   "wayne-admin",
		WayneAdminPassword:   "wayne-password",
		WayneTokenTTLMinutes: 60,
	}
}

func assertLoginBody(t *testing.T, r *http.Request) {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(readTestBody(t, r), &body); err != nil {
		t.Fatalf("invalid login body: %v", err)
	}
	if body["username"] != "wayne-admin" || body["password"] != "wayne-password" {
		t.Fatalf("login body = %#v", body)
	}
}

func assertBearer(t *testing.T, r *http.Request) {
	t.Helper()
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		t.Fatalf("missing bearer header: %#v", r.Header)
	}
}

func testJWT(exp time.Time) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"exp":` + strconv.FormatInt(exp.Unix(), 10) + `}`))
	return header + "." + payload + ".sig"
}

func readTestBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return body
}
