package handler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/service"
	"github.com/1024XEngineer/xinfra/server/internal/sso"

	"github.com/gin-gonic/gin"
)

type SAMLHandler struct {
	cfg  config.Config
	auth *service.AuthService
}

func NewSAMLHandler(cfg config.Config, authService *service.AuthService) *SAMLHandler {
	return &SAMLHandler{cfg: cfg, auth: authService}
}

func (h *SAMLHandler) Metadata(c *gin.Context) {
	data, err := sso.BuildSPMetadata(sso.MetadataConfig{
		EntityID: h.cfg.SAMLEntityID,
		ACSURL:   h.cfg.SAMLACSURL,
		CertFile: h.cfg.SAMLSPCert,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/samlmetadata+xml; charset=utf-8", data)
}

func (h *SAMLHandler) MetadataConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"entity_id":          h.cfg.SAMLEntityID,
		"acs_url":            h.cfg.SAMLACSURL,
		"sp_cert_file":       h.cfg.SAMLSPCert,
		"sp_key_file":        h.cfg.SAMLSPKey,
		"idp_metadata_url":   h.cfg.SAMLIDPMetaURL,
		"sp_metadata_url":    h.cfg.SAMLEntityID,
		"metadata_generated": true,
	})
}

func (h *SAMLHandler) Login(c *gin.Context) {
	redirectURL, err := sso.BuildLoginRedirect(sso.LoginConfig{
		EntityID:       h.cfg.SAMLEntityID,
		ACSURL:         h.cfg.SAMLACSURL,
		IDPMetadataURL: h.cfg.SAMLIDPMetaURL,
		RelayState:     c.Query("relay_state"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *SAMLHandler) Logout(c *gin.Context) {
	expired := time.Unix(0, 0)
	for _, path := range []string{"/auth/", "/"} {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     AuthSessionCookieName,
			Value:    "",
			Path:     path,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Expires:  expired,
		})
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *SAMLHandler) ACS(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	info, err := sso.DecodeSAMLResponse(c.PostForm("SAMLResponse"), h.cfg.SAMLSPKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("SAML ACS debug summary:\n%s", info.JSON())
	if info.DecryptedAssertionXML != "" {
		log.Printf("SAML ACS decrypted assertion:\n%s", info.DecryptedAssertionXML)
	}

	result, err := h.auth.SAMLLogin(info, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AuthSessionCookieName,
		Value:    result.Token,
		Path:     "/auth/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.ExpiresAt).Seconds()),
	})
	c.Redirect(http.StatusFound, ssoRedirectURL(c.PostForm("RelayState"), result.Token))
}

func ssoRedirectURL(relayState, token string) string {
	target := strings.TrimSpace(relayState)
	if target == "" {
		target = "/"
	}
	parsed, err := url.Parse(target)
	if err != nil || parsed.IsAbs() || !strings.HasPrefix(target, "/") || strings.HasPrefix(target, "//") || strings.HasPrefix(parsed.Path, "/auth/api/") {
		parsed = &url.URL{Path: "/"}
	}
	values := parsed.Query()
	values.Set("sso_token", token)
	parsed.RawQuery = values.Encode()
	return parsed.String()
}
