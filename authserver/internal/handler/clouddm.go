package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"authserver/internal/config"
	"authserver/internal/service"

	"github.com/gin-gonic/gin"
)

type CloudDMHandler struct {
	cfg   config.Config
	audit *service.AuditService
}

func NewCloudDMHandler(cfg config.Config, audit *service.AuditService) *CloudDMHandler {
	return &CloudDMHandler{cfg: cfg, audit: audit}
}

func (h *CloudDMHandler) Login(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing current user"})
		return
	}

	targetURL := strings.TrimSpace(h.cfg.CloudDMTargetURL)
	if targetURL == "" {
		h.writeAudit(c, claims.UserID, claims.Username, "deny", "clouddm target url is not configured")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "clouddm target url is not configured"})
		return
	}

	jumpURL, err := h.loginJumpURL(targetURL)
	if err != nil {
		h.writeAudit(c, claims.UserID, claims.Username, "deny", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	h.writeAudit(c, claims.UserID, claims.Username, "allow", "")
	if strings.Contains(c.GetHeader("Accept"), "application/json") {
		c.JSON(http.StatusOK, gin.H{"target_url": jumpURL})
		return
	}
	c.Redirect(http.StatusFound, jumpURL)
}

func (h *CloudDMHandler) loginJumpURL(targetURL string) (string, error) {
	requestURL := strings.TrimRight(targetURL, "/") + "/requestJumpUrl"
	body, _ := json.Marshal(gin.H{"type": "OIDC"})
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("clouddm requestJumpUrl failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var result struct {
		Success    bool   `json:"success"`
		Data       string `json:"data"`
		Msg        string `json:"msg"`
		MsgContent string `json:"msgContent"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", err
	}
	if !result.Success || strings.TrimSpace(result.Data) == "" {
		reason := strings.TrimSpace(result.MsgContent)
		if reason == "" {
			reason = strings.TrimSpace(result.Msg)
		}
		if reason == "" {
			reason = "empty clouddm jump url"
		}
		return "", fmt.Errorf("clouddm requestJumpUrl failed: %s", reason)
	}
	return result.Data, nil
}

func (h *CloudDMHandler) writeAudit(c *gin.Context, userID uint64, username, decision, reason string) {
	h.audit.Write(service.AuditEntry{
		ActorUserID:   userID,
		ActorUsername: username,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Action:        "clouddm.login",
		ResourceType:  "clouddm",
		Decision:      decision,
		Reason:        reason,
	})
}
