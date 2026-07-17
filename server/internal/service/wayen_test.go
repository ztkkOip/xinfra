package service

import (
	"net/url"
	"testing"

	"github.com/1024XEngineer/xinfra/server/internal/config"
)

func TestWayenLoginUsesDedicatedOAuthLoginURL(t *testing.T) {
	service := NewWayenService(config.Config{
		WayenOAuthLoginURL: "http://218.11.5.223:32000/login/oauth2/oauth2",
		OAuthRedirectURI:   "http://218.11.5.223:30008/login/oauth2/oauth2",
		WayenTargetURL:     "http://218.11.5.223:32000/",
		WayenOAuthRef:      "/portal/namespace/1/app",
	}, nil)

	result, err := service.Login("eastsales@qiniu.com", "eastsales@qiniu.com", "")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	parsed, err := url.Parse(result.TargetURL)
	if err != nil {
		t.Fatalf("invalid target url: %v", err)
	}
	if parsed.Host != "218.11.5.223:32000" {
		t.Fatalf("target host = %q, want Wayne frontend 32000", parsed.Host)
	}
	if parsed.Path != "/login/oauth2/oauth2" {
		t.Fatalf("target path = %q", parsed.Path)
	}

	next := parsed.Query().Get("next")
	if next == "" {
		t.Fatal("next is empty")
	}
	parsedNext, err := url.Parse(next)
	if err != nil {
		t.Fatalf("invalid next url: %v", err)
	}
	if parsedNext.Host != "218.11.5.223:32000" {
		t.Fatalf("next host = %q, want Wayne frontend 32000", parsedNext.Host)
	}
	if parsedNext.Path != "/sign-in" {
		t.Fatalf("next path = %q, want /sign-in", parsedNext.Path)
	}
	if parsedNext.Query().Get("ref") != "/portal/namespace/1/app" {
		t.Fatalf("next ref = %q", parsedNext.Query().Get("ref"))
	}
}

func TestWayenLoginFallsBackToOAuthRedirectURI(t *testing.T) {
	service := NewWayenService(config.Config{
		OAuthRedirectURI: "http://218.11.5.223:32000/login/oauth2/oauth2",
		WayenTargetURL:   "http://218.11.5.223:32000/",
	}, nil)

	result, err := service.Login("eastsales@qiniu.com", "eastsales@qiniu.com", "")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	parsed, err := url.Parse(result.TargetURL)
	if err != nil {
		t.Fatalf("invalid target url: %v", err)
	}
	if parsed.Host != "218.11.5.223:32000" {
		t.Fatalf("target host = %q, want fallback redirect host", parsed.Host)
	}
}

func TestWayenLoginUsesRefOverride(t *testing.T) {
	service := NewWayenService(config.Config{
		WayenOAuthLoginURL: "http://218.11.5.223:32000/login/oauth2/oauth2",
		WayenTargetURL:     "http://218.11.5.223:32000/",
		WayenOAuthRef:      "/portal/namespace/1/app",
	}, nil)

	result, err := service.Login("eastsales@qiniu.com", "eastsales@qiniu.com", "/portal/namespace/3/app")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	parsed, err := url.Parse(result.TargetURL)
	if err != nil {
		t.Fatalf("invalid target url: %v", err)
	}
	parsedNext, err := url.Parse(parsed.Query().Get("next"))
	if err != nil {
		t.Fatalf("invalid next url: %v", err)
	}
	if parsedNext.Query().Get("ref") != "/portal/namespace/3/app" {
		t.Fatalf("next ref = %q", parsedNext.Query().Get("ref"))
	}
}
