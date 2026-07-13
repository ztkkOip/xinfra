package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv              string
	HTTPAddr            string
	PublicBaseURL       string
	MySQLDSN            string
	AutoMigrate         bool
	JWTSecret           string
	JWTIssuer           string
	JWTTTLMinutes       int
	SAMLEntityID        string
	SAMLACSURL          string
	SAMLSPCert          string
	SAMLSPKey           string
	SAMLIDPMetaURL      string
	WayenLoginURL       string
	WayenTargetURL      string
	WayenUsernameKey    string
	WayenPasswordKey    string
	WayenLoginFormat    string
	WayenLoginValue     string
	OAuthClientID       string
	OAuthClientSecret   string
	OAuthRedirectURI    string
	OAuthCodeTTLSeconds int
}

func Load() Config {
	loadDotEnv(".env")
	httpAddr := env("HTTP_ADDR", ":8080")
	publicBaseURL := env("PUBLIC_BASE_URL", defaultPublicBaseURL(httpAddr))
	samlEntityID := env("SAML_ENTITY_ID", strings.TrimRight(publicBaseURL, "/")+"/auth/api/v1/saml/metadata")
	samlACSURL := env("SAML_ACS_URL", strings.TrimRight(publicBaseURL, "/")+"/auth/api/v1/saml/acs")

	return Config{
		AppEnv:              env("APP_ENV", "dev"),
		HTTPAddr:            httpAddr,
		PublicBaseURL:       publicBaseURL,
		MySQLDSN:            env("MYSQL_DSN", "auth:auth@tcp(127.0.0.1:3306)/authserver?charset=utf8mb4&parseTime=True&loc=Local"),
		AutoMigrate:         envBool("AUTO_MIGRATE", true),
		JWTSecret:           env("JWT_SECRET", "change-this-secret"),
		JWTIssuer:           env("JWT_ISSUER", "authserver"),
		JWTTTLMinutes:       envInt("JWT_TTL_MINUTES", 120),
		SAMLEntityID:        samlEntityID,
		SAMLACSURL:          samlACSURL,
		SAMLSPCert:          env("SAML_SP_CERT_FILE", "certs/sp.crt"),
		SAMLSPKey:           env("SAML_SP_KEY_FILE", "certs/sp.key"),
		SAMLIDPMetaURL:      env("SAML_IDP_METADATA_URL", "http://sso-internal.dev.qiniu.io/saml2/meta"),
		WayenLoginURL:       env("WAYEN_LOGIN_URL", ""),
		WayenTargetURL:      env("WAYEN_TARGET_URL", ""),
		WayenUsernameKey:    env("WAYEN_USERNAME_KEY", "email"),
		WayenPasswordKey:    env("WAYEN_PASSWORD_KEY", "password"),
		WayenLoginFormat:    env("WAYEN_LOGIN_FORMAT", "form"),
		WayenLoginValue:     env("WAYEN_LOGIN_VALUE", "email"),
		OAuthClientID:       env("OAUTH_WAYNE_CLIENT_ID", "wayne"),
		OAuthClientSecret:   env("OAUTH_WAYNE_CLIENT_SECRET", "wayne-secret"),
		OAuthRedirectURI:    env("OAUTH_WAYNE_REDIRECT_URI", ""),
		OAuthCodeTTLSeconds: envInt("OAUTH_CODE_TTL_SECONDS", 120),
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
