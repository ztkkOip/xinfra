package handler

import "testing"

func TestSSORedirectURLAllowsFrontendRoutes(t *testing.T) {
	tests := []struct {
		name  string
		relay string
		want  string
	}{
		{
			name:  "root",
			relay: "/",
			want:  "/?sso_token=token-001",
		},
		{
			name:  "frontend route",
			relay: "/subsystem?open_app=wayne#section",
			want:  "/subsystem?open_app=wayne&sso_token=token-001#section",
		},
		{
			name:  "oauth authorize route",
			relay: "/auth/oauth/authorize?client_id=wayne",
			want:  "/auth/oauth/authorize?client_id=wayne&sso_token=token-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ssoRedirectURL(tt.relay, "token-001"); got != tt.want {
				t.Fatalf("ssoRedirectURL(%q) = %q, want %q", tt.relay, got, tt.want)
			}
		})
	}
}

func TestSSORedirectURLRejectsUnsafeRelayState(t *testing.T) {
	tests := []string{
		"",
		"https://evil.example.com/",
		"//evil.example.com/",
		"dashboard",
		"/auth/api/v1/users/me",
	}

	for _, relay := range tests {
		t.Run(relay, func(t *testing.T) {
			if got := ssoRedirectURL(relay, "token-001"); got != "/?sso_token=token-001" {
				t.Fatalf("ssoRedirectURL(%q) = %q, want root fallback", relay, got)
			}
		})
	}
}
