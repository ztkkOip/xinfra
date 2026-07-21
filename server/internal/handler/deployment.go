package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/model"
	"github.com/1024XEngineer/xinfra/server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeploymentHandler struct {
	cfg         config.Config
	db          *gorm.DB
	deployments *service.DeploymentService
}

func NewDeploymentHandler(cfg config.Config, db *gorm.DB, deployments *service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{cfg: cfg, db: db, deployments: deployments}
}

func (h *DeploymentHandler) Create(c *gin.Context) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth claims"})
		return
	}
	var req service.DeploymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !h.ensureBusinessLineMember(c, req.BusinessLineID, claims.UserID, claims.IsAdmin) {
		return
	}
	deployment, err := h.deployments.Create(c.Request.Context(), req, claims.UserID, claims.Username)
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, service.ErrDeploymentNotConfigured) {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, gin.H{
			"error":         err.Error(),
			"deployment_id": deployment.DeploymentID,
			"status":        deployment.Status,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"deployment_id": deployment.DeploymentID, "status": deployment.Status})
}

func (h *DeploymentHandler) Get(c *gin.Context) {
	deployment, ok := h.getAuthorizedDeployment(c)
	if !ok {
		return
	}
	events, err := h.deployments.Events(c.Request.Context(), deployment.DeploymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deployment": deployment, "events": events})
}

func (h *DeploymentHandler) Events(c *gin.Context) {
	deployment, ok := h.getAuthorizedDeployment(c)
	if !ok {
		return
	}
	events, err := h.deployments.Events(c.Request.Context(), deployment.DeploymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	for _, event := range events {
		if err := writeSSE(w, event.Type, service.DeploymentEventData(event)); err != nil {
			return
		}
	}
	if isTerminalDeploymentStatus(deployment.Status) {
		_ = writeSSE(w, service.DeploymentEventDone, gin.H{"status": deployment.Status})
		return
	}

	ch, unsubscribe := h.deployments.Subscribe(deployment.DeploymentID)
	defer unsubscribe()
	flusher, _ := w.(http.Flusher)
	if flusher != nil {
		flusher.Flush()
	}
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case event := <-ch:
			if err := writeSSE(w, event.Event.Type, event.Data); err != nil {
				return
			}
			if event.Event.Type == service.DeploymentEventDone {
				return
			}
		}
	}
}

func (h *DeploymentHandler) Cancel(c *gin.Context) {
	deployment, ok := h.getAuthorizedDeployment(c)
	if !ok {
		return
	}
	if err := h.deployments.Cancel(c.Request.Context(), deployment.DeploymentID); err != nil {
		writeDeploymentError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"deployment_id": deployment.DeploymentID, "status": service.DeploymentStatusCanceling})
}

func (h *DeploymentHandler) InternalEvent(c *gin.Context) {
	if !h.authorizeInternal(c) {
		return
	}
	deploymentID := strings.TrimSpace(c.Param("id"))
	var req service.DeploymentEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := h.deployments.AppendEvent(c.Request.Context(), deploymentID, req)
	if err != nil {
		writeDeploymentError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"event_id": event.ID, "seq": event.Seq})
}

func (h *DeploymentHandler) InternalFinish(c *gin.Context) {
	if !h.authorizeInternal(c) {
		return
	}
	deploymentID := strings.TrimSpace(c.Param("id"))
	var req service.DeploymentFinishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deployment, err := h.deployments.Finish(c.Request.Context(), deploymentID, req)
	if err != nil {
		writeDeploymentError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"deployment_id": deployment.DeploymentID, "status": deployment.Status})
}

func (h *DeploymentHandler) getAuthorizedDeployment(c *gin.Context) (model.Deployment, bool) {
	claims, ok := CurrentClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth claims"})
		return model.Deployment{}, false
	}
	deploymentID := strings.TrimSpace(c.Param("id"))
	deployment, err := h.deployments.Get(c.Request.Context(), deploymentID)
	if err != nil {
		writeDeploymentError(c, err)
		return model.Deployment{}, false
	}
	if !h.ensureBusinessLineMember(c, deployment.BusinessLineID, claims.UserID, claims.IsAdmin) {
		return model.Deployment{}, false
	}
	return deployment, true
}

func (h *DeploymentHandler) ensureBusinessLineMember(c *gin.Context, businessLineID uint64, userID uint64, isAdmin bool) bool {
	if businessLineID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_line_id is required"})
		return false
	}
	if isAdmin {
		return true
	}
	var binding model.BusinessLineUser
	err := h.db.WithContext(c.Request.Context()).Where("business_line_id = ? AND user_id = ?", businessLineID, userID).First(&binding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusForbidden, gin.H{"error": "current user is not assigned to this business line"})
		return false
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (h *DeploymentHandler) authorizeInternal(c *gin.Context) bool {
	if h.cfg.AnsibleInternalToken == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ansible internal token is not configured"})
		return false
	}
	value := c.GetHeader("Authorization")
	if strings.TrimPrefix(value, "Bearer ") != h.cfg.AnsibleInternalToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid internal token"})
		return false
	}
	return true
}

func writeSSE(w gin.ResponseWriter, event string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, body); err != nil {
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

func writeDeploymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrDeploymentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrDeploymentForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrDeploymentInvalidState):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrDeploymentNotConfigured):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func isTerminalDeploymentStatus(status string) bool {
	return status == service.DeploymentStatusSuccess || status == service.DeploymentStatusFailed || status == service.DeploymentStatusCanceled
}
