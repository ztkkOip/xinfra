package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type OAuthClient struct {
	ID           string
	Secret       string
	RedirectURIs []string
}

type Config struct {
	AppEnv                   string
	HTTPAddr                 string
	PublicBaseURL            string
	MySQLDSN                 string
	AutoMigrate              bool
	SSOEnabled               bool
	JWTSecret                string
	JWTIssuer                string
	JWTTTLMinutes            int
	SAMLEntityID             string
	SAMLACSURL               string
	SAMLSPCert               string
	SAMLSPKey                string
	SAMLIDPMetaURL           string
	WayenLoginURL            string
	WayenTargetURL           string
	WayenUsernameKey         string
	WayenPasswordKey         string
	WayenLoginFormat         string
	WayenLoginValue          string
	WayenOAuthRef            string
	WayenOAuthLoginURL       string
	WayneInternalAPIBaseURL  string
	WayneServiceName         string
	WayneServiceAPISecretKey string
	OAuthClientID            string
	OAuthClientSecret        string
	OAuthRedirectURI         string
	OAuthCodeTTLSeconds      int
	OIDCIssuer               string
	OIDCAuthorizeURL         string
	OIDCTokenURL             string
	OIDCUserInfoURL          string
	OIDCJWKSURL              string
	CloudDMClientID          string
	CloudDMClientSecret      string
	CloudDMRedirectURI       string
	CloudDMTargetURL         string
}

func Load() Config {
	loadDotEnv(".env")
	httpAddr := env("HTTP_ADDR", ":8080")
	publicBaseURL := env("PUBLIC_BASE_URL", defaultPublicBaseURL(httpAddr))
	samlEntityID := env("SAML_ENTITY_ID", strings.TrimRight(publicBaseURL, "/")+"/auth/api/v1/saml/metadata")
	samlACSURL := env("SAML_ACS_URL", strings.TrimRight(publicBaseURL, "/")+"/auth/api/v1/saml/acs")
	oidcIssuer := env("OIDC_ISSUER", strings.TrimRight(publicBaseURL, "/")+"/auth")
	oidcIssuer = strings.TrimRight(oidcIssuer, "/")

	return Config{
		AppEnv:                   env("APP_ENV", "dev"),
		HTTPAddr:                 httpAddr,
		PublicBaseURL:            publicBaseURL,
		MySQLDSN:                 env("MYSQL_DSN", "auth:auth@tcp(127.0.0.1:3306)/authserver?charset=utf8mb4&parseTime=True&loc=Local"),
		AutoMigrate:              envBool("AUTO_MIGRATE", true),
		SSOEnabled:               envBool("SSO_ENABLED", true),
		JWTSecret:                env("JWT_SECRET", "change-this-secret"),
		JWTIssuer:                env("JWT_ISSUER", "authserver"),
		JWTTTLMinutes:            envInt("JWT_TTL_MINUTES", 120),
		SAMLEntityID:             samlEntityID,
		SAMLACSURL:               samlACSURL,
		SAMLSPCert:               env("SAML_SP_CERT_FILE", "certs/sp.crt"),
		SAMLSPKey:                env("SAML_SP_KEY_FILE", "certs/sp.key"),
		SAMLIDPMetaURL:           env("SAML_IDP_METADATA_URL", "http://sso-internal.dev.qiniu.io/saml2/meta"),
		WayenLoginURL:            env("WAYEN_LOGIN_URL", ""),
		WayenTargetURL:           env("WAYEN_TARGET_URL", ""),
		WayenUsernameKey:         env("WAYEN_USERNAME_KEY", "email"),
		WayenPasswordKey:         env("WAYEN_PASSWORD_KEY", "password"),
		WayenLoginFormat:         env("WAYEN_LOGIN_FORMAT", "form"),
		WayenLoginValue:          env("WAYEN_LOGIN_VALUE", "email"),
		WayenOAuthRef:            env("WAYEN_OAUTH_REF", "/portal/namespace/1/app"),
		WayenOAuthLoginURL:       trimURL(env("WAYEN_OAUTH_LOGIN_URL", "")),
		WayneInternalAPIBaseURL:  trimURL(env("WAYNE_INTERNAL_API_BASE_URL", "")),
		WayneServiceName:         env("WAYNE_SERVICE_NAME", "xinfra"),
		WayneServiceAPISecretKey: env("WAYNE_SERVICE_API_SECRET_KEY", ""),
		OAuthClientID:            env("OAUTH_WAYNE_CLIENT_ID", "wayne"),
		OAuthClientSecret:        env("OAUTH_WAYNE_CLIENT_SECRET", "wayne-secret"),
		OAuthRedirectURI:         env("OAUTH_WAYNE_REDIRECT_URI", ""),
		OAuthCodeTTLSeconds:      envInt("OAUTH_CODE_TTL_SECONDS", 120),
		OIDCIssuer:               oidcIssuer,
		OIDCAuthorizeURL:         trimURL(env("OIDC_AUTHORIZATION_ENDPOINT", oidcIssuer+"/oauth/authorize")),
		OIDCTokenURL:             trimURL(env("OIDC_TOKEN_ENDPOINT", oidcIssuer+"/oauth/token")),
		OIDCUserInfoURL:          trimURL(env("OIDC_USERINFO_ENDPOINT", oidcIssuer+"/oauth/userinfo")),
		OIDCJWKSURL:              trimURL(env("OIDC_JWKS_URI", oidcIssuer+"/oauth/jwks")),
		CloudDMClientID:          env("OIDC_CLOUDDM_CLIENT_ID", "clouddm"),
		CloudDMClientSecret:      env("OIDC_CLOUDDM_CLIENT_SECRET", ""),
		CloudDMRedirectURI:       env("OIDC_CLOUDDM_REDIRECT_URI", ""),
		CloudDMTargetURL:         env("CLOUDDM_TARGET_URL", ""),
	}
}

func loadDotEnv(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func (c Config) JWTTTL() time.Duration {
	return time.Duration(c.JWTTTLMinutes) * time.Minute
}

func (c Config) OAuthCodeTTL() time.Duration {
	return time.Duration(c.OAuthCodeTTLSeconds) * time.Second
}

func (c Config) OAuthClients() map[string]OAuthClient {
	clients := make(map[string]OAuthClient)
	addOAuthClient(clients, c.OAuthClientID, c.OAuthClientSecret, c.OAuthRedirectURI)
	addOAuthClient(clients, c.CloudDMClientID, c.CloudDMClientSecret, c.CloudDMRedirectURI)
	return clients
}

func addOAuthClient(clients map[string]OAuthClient, id, secret, redirectURIs string) {
	id = strings.TrimSpace(id)
	secret = strings.TrimSpace(secret)
	if id == "" || secret == "" {
		return
	}
	clients[id] = OAuthClient{
		ID:           id,
		Secret:       secret,
		RedirectURIs: splitCSV(redirectURIs),
	}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}

func trimURL(value string) string {
	return strings.TrimRight(strings.TrimSpace(value), "/")
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func defaultPublicBaseURL(httpAddr string) string {
	addr := strings.TrimSpace(httpAddr)
	if addr == "" {
		return "http://localhost:8080"
	}
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}
	return "http://" + addr
}
