package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DeploymentStatusPending   = "pending"
	DeploymentStatusRunning   = "running"
	DeploymentStatusSuccess   = "success"
	DeploymentStatusFailed    = "failed"
	DeploymentStatusCanceling = "canceling"
	DeploymentStatusCanceled  = "canceled"

	DeploymentEventLog    = "log"
	DeploymentEventStatus = "status"
	DeploymentEventResult = "result"
	DeploymentEventDone   = "done"
	DeploymentEventError  = "error"
)

var (
	ErrDeploymentNotConfigured = errors.New("ansible deployment service is not configured")
	ErrDeploymentNotFound      = errors.New("deployment not found")
	ErrDeploymentForbidden     = errors.New("deployment permission denied")
	ErrDeploymentInvalidState  = errors.New("deployment state does not allow this operation")
)

type DeploymentCreateRequest struct {
	Component      string         `json:"component"`
	BusinessLineID uint64         `json:"business_line_id"`
	Params         map[string]any `json:"params"`
}

type DeploymentEventRequest struct {
	Type    string         `json:"type"`
	Level   string         `json:"level"`
	Status  string         `json:"status"`
	Seq     uint64         `json:"seq"`
	Message string         `json:"message"`
	Payload map[string]any `json:"payload"`
}

type DeploymentFinishRequest struct {
	Status   string         `json:"status"`
	ExitCode int            `json:"exit_code"`
	Error    string         `json:"error"`
	Summary  map[string]any `json:"summary"`
}

type DeploymentEventEnvelope struct {
	Event model.DeploymentEvent
	Data  map[string]any
}

type DeploymentService struct {
	cfg    config.Config
	db     *gorm.DB
	client *http.Client

	mu          sync.Mutex
	subscribers map[string]map[chan DeploymentEventEnvelope]struct{}
}

func NewDeploymentService(cfg config.Config, db *gorm.DB) *DeploymentService {
	return &DeploymentService{
		cfg:         cfg,
		db:          db,
		client:      &http.Client{Timeout: 8 * time.Second},
		subscribers: make(map[string]map[chan DeploymentEventEnvelope]struct{}),
	}
}

func (s *DeploymentService) Create(ctx context.Context, req DeploymentCreateRequest, actorID uint64, actorName string) (model.Deployment, error) {
	component := strings.TrimSpace(strings.ToLower(req.Component))
	if !isSupportedDeploymentComponent(component) {
		return model.Deployment{}, fmt.Errorf("unsupported deployment component: %s", req.Component)
	}
	if req.BusinessLineID == 0 {
		return model.Deployment{}, errors.New("business_line_id is required")
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return model.Deployment{}, err
	}

	var businessLine model.BusinessLine
	if err := s.db.WithContext(ctx).First(&businessLine, req.BusinessLineID).Error; err != nil {
		return model.Deployment{}, err
	}

	deployment := model.Deployment{
		DeploymentID:   newDeploymentID(),
		Component:      component,
		BusinessLineID: req.BusinessLineID,
		BusinessLine:   businessLine.Name,
		Status:         DeploymentStatusPending,
		RequestPayload: string(payload),
		CreatedBy:      actorID,
		CreatedByName:  actorName,
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&deployment).Error; err != nil {
			return err
		}
		_, err := s.appendEventTx(ctx, tx, deployment.DeploymentID, DeploymentEventStatus, "info", "deployment created", map[string]any{"status": DeploymentStatusPending})
		return err
	}); err != nil {
		return model.Deployment{}, err
	}

	if err := s.startPythonDeployment(ctx, deployment, req.Params); err != nil {
		_ = s.FailStart(ctx, deployment.DeploymentID, err)
		return deployment, err
	}
	return deployment, nil
}

func (s *DeploymentService) Get(ctx context.Context, deploymentID string) (model.Deployment, error) {
	var deployment model.Deployment
	if err := s.db.WithContext(ctx).Where("deployment_id = ?", deploymentID).First(&deployment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return deployment, ErrDeploymentNotFound
		}
		return deployment, err
	}
	return deployment, nil
}

func (s *DeploymentService) Events(ctx context.Context, deploymentID string) ([]model.DeploymentEvent, error) {
	var events []model.DeploymentEvent
	err := s.db.WithContext(ctx).Where("deployment_id = ?", deploymentID).Order("seq ASC").Find(&events).Error
	return events, err
}

func (s *DeploymentService) AppendEvent(ctx context.Context, deploymentID string, req DeploymentEventRequest) (model.DeploymentEvent, error) {
	eventType := normalizeDeploymentEventType(req.Type)
	level := strings.TrimSpace(req.Level)
	if level == "" {
		level = "info"
	}
	payload := req.Payload
	if payload == nil {
		payload = map[string]any{}
	}
	if req.Status != "" {
		payload["status"] = normalizeDeploymentStatus(req.Status)
	}
	event, err := s.appendEvent(ctx, deploymentID, eventType, level, req.Message, payload)
	if err != nil {
		return event, err
	}
	if status, _ := payload["status"].(string); status != "" {
		_ = s.updateStatus(ctx, deploymentID, status, "")
	}
	return event, nil
}

func (s *DeploymentService) Finish(ctx context.Context, deploymentID string, req DeploymentFinishRequest) (model.Deployment, error) {
	status := normalizeDeploymentStatus(req.Status)
	if status == "" {
		status = DeploymentStatusFailed
	}
	payload := map[string]any{
		"status":    status,
		"exit_code": req.ExitCode,
		"summary":   req.Summary,
	}
	if req.Error != "" {
		payload["error"] = req.Error
	}
	result, _ := json.Marshal(payload)
	now := time.Now()
	var deployment model.Deployment
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("deployment_id = ?", deploymentID).First(&deployment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrDeploymentNotFound
			}
			return err
		}
		updates := map[string]any{"status": status, "result_payload": string(result), "finished_at": &now}
		if err := tx.Model(&model.Deployment{}).Where("deployment_id = ?", deploymentID).Updates(updates).Error; err != nil {
			return err
		}
		eventType := DeploymentEventResult
		if status == DeploymentStatusFailed || status == DeploymentStatusCanceled {
			eventType = DeploymentEventError
		}
		if _, err := s.appendEventTx(ctx, tx, deploymentID, eventType, eventLevelForStatus(status), finishMessage(status, req.Error), payload); err != nil {
			return err
		}
		_, err := s.appendEventTx(ctx, tx, deploymentID, DeploymentEventDone, eventLevelForStatus(status), status, map[string]any{"status": status})
		return err
	})
	if err != nil {
		return deployment, err
	}
	deployment.Status = status
	deployment.ResultPayload = string(result)
	deployment.FinishedAt = &now
	return deployment, nil
}

func (s *DeploymentService) Cancel(ctx context.Context, deploymentID string) error {
	deployment, err := s.Get(ctx, deploymentID)
	if err != nil {
		return err
	}
	if !deploymentCancelable(deployment.Status) {
		return ErrDeploymentInvalidState
	}
	if err := s.updateStatus(ctx, deploymentID, DeploymentStatusCanceling, "cancel requested"); err != nil {
		return err
	}
	if s.cfg.AnsibleServiceBaseURL == "" {
		return ErrDeploymentNotConfigured
	}
	body, _ := json.Marshal(map[string]any{"deployment_id": deploymentID})
	path := s.cfg.AnsibleServiceBaseURL + "/internal/ansible/deployments/" + deploymentID + "/cancel"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if s.cfg.AnsibleInternalToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.cfg.AnsibleInternalToken)
	}
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("ansible cancel failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}

func (s *DeploymentService) Subscribe(deploymentID string) (chan DeploymentEventEnvelope, func()) {
	ch := make(chan DeploymentEventEnvelope, 64)
	s.mu.Lock()
	if s.subscribers[deploymentID] == nil {
		s.subscribers[deploymentID] = make(map[chan DeploymentEventEnvelope]struct{})
	}
	s.subscribers[deploymentID][ch] = struct{}{}
	s.mu.Unlock()
	return ch, func() {
		s.mu.Lock()
		if subscribers := s.subscribers[deploymentID]; subscribers != nil {
			delete(subscribers, ch)
			if len(subscribers) == 0 {
				delete(s.subscribers, deploymentID)
			}
		}
		s.mu.Unlock()
		close(ch)
	}
}

func (s *DeploymentService) FailStart(ctx context.Context, deploymentID string, cause error) error {
	_, err := s.Finish(ctx, deploymentID, DeploymentFinishRequest{
		Status: DeploymentStatusFailed,
		Error:  cause.Error(),
	})
	return err
}

func (s *DeploymentService) startPythonDeployment(ctx context.Context, deployment model.Deployment, params map[string]any) error {
	if s.cfg.AnsibleServiceBaseURL == "" {
		return ErrDeploymentNotConfigured
	}
	callbackBaseURL := strings.TrimRight(s.cfg.DeploymentCallbackBaseURL, "/")
	body, err := json.Marshal(map[string]any{
		"deployment_id": deployment.DeploymentID,
		"component":     deployment.Component,
		"callback_url":  callbackBaseURL + "/auth/internal/deployments/" + deployment.DeploymentID + "/events",
		"finish_url":    callbackBaseURL + "/auth/internal/deployments/" + deployment.DeploymentID + "/finish",
		"params":        params,
	})
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.AnsibleServiceBaseURL+"/internal/ansible/deploy", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if s.cfg.AnsibleInternalToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.cfg.AnsibleInternalToken)
	}
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("ansible deploy failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}

func (s *DeploymentService) appendEvent(ctx context.Context, deploymentID, eventType, level, message string, payload map[string]any) (model.DeploymentEvent, error) {
	var event model.DeploymentEvent
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		event, err = s.appendEventTx(ctx, tx, deploymentID, eventType, level, message, payload)
		return err
	})
	return event, err
}

func (s *DeploymentService) appendEventTx(ctx context.Context, tx *gorm.DB, deploymentID, eventType, level, message string, payload map[string]any) (model.DeploymentEvent, error) {
	var deployment model.Deployment
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("deployment_id = ?", deploymentID).First(&deployment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.DeploymentEvent{}, ErrDeploymentNotFound
		}
		return model.DeploymentEvent{}, err
	}
	var maxSeq uint64
	if err := tx.Model(&model.DeploymentEvent{}).Where("deployment_id = ?", deploymentID).Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq).Error; err != nil {
		return model.DeploymentEvent{}, err
	}
	payloadBody, _ := json.Marshal(payload)
	event := model.DeploymentEvent{
		DeploymentID: deploymentID,
		Seq:          maxSeq + 1,
		Type:         eventType,
		Level:        level,
		Message:      message,
		Payload:      string(payloadBody),
	}
	if err := tx.Create(&event).Error; err != nil {
		return model.DeploymentEvent{}, err
	}
	go s.broadcast(event, payload)
	return event, nil
}

func (s *DeploymentService) updateStatus(ctx context.Context, deploymentID string, status string, message string) error {
	status = normalizeDeploymentStatus(status)
	if status == "" {
		return nil
	}
	now := time.Now()
	updates := map[string]any{"status": status}
	if status == DeploymentStatusRunning {
		updates["started_at"] = &now
	}
	if isDeploymentTerminal(status) {
		updates["finished_at"] = &now
	}
	if err := s.db.WithContext(ctx).Model(&model.Deployment{}).Where("deployment_id = ?", deploymentID).Updates(updates).Error; err != nil {
		return err
	}
	if message != "" {
		_, err := s.appendEvent(ctx, deploymentID, DeploymentEventStatus, eventLevelForStatus(status), message, map[string]any{"status": status})
		return err
	}
	return nil
}

func (s *DeploymentService) broadcast(event model.DeploymentEvent, data map[string]any) {
	envelope := DeploymentEventEnvelope{Event: event, Data: eventData(event, data)}
	s.mu.Lock()
	subscribers := make([]chan DeploymentEventEnvelope, 0, len(s.subscribers[event.DeploymentID]))
	for ch := range s.subscribers[event.DeploymentID] {
		subscribers = append(subscribers, ch)
	}
	s.mu.Unlock()
	for _, ch := range subscribers {
		select {
		case ch <- envelope:
		default:
		}
	}
}

func eventData(event model.DeploymentEvent, data map[string]any) map[string]any {
	out := map[string]any{
		"deployment_id": event.DeploymentID,
		"seq":           event.Seq,
		"type":          event.Type,
		"level":         event.Level,
		"message":       event.Message,
		"created_at":    event.CreatedAt,
	}
	for key, value := range data {
		out[key] = value
	}
	return out
}

func DeploymentEventData(event model.DeploymentEvent) map[string]any {
	payload := map[string]any{}
	if strings.TrimSpace(event.Payload) != "" {
		_ = json.Unmarshal([]byte(event.Payload), &payload)
	}
	return eventData(event, payload)
}

func isSupportedDeploymentComponent(component string) bool {
	switch component {
	case "mysql", "openresty":
		return true
	default:
		return false
	}
}

func normalizeDeploymentEventType(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case DeploymentEventStatus:
		return DeploymentEventStatus
	case DeploymentEventResult:
		return DeploymentEventResult
	case DeploymentEventDone:
		return DeploymentEventDone
	case DeploymentEventError:
		return DeploymentEventError
	default:
		return DeploymentEventLog
	}
}

func normalizeDeploymentStatus(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case DeploymentStatusPending, DeploymentStatusRunning, DeploymentStatusSuccess, DeploymentStatusFailed, DeploymentStatusCanceling, DeploymentStatusCanceled:
		return strings.TrimSpace(strings.ToLower(value))
	default:
		return ""
	}
}

func isDeploymentTerminal(status string) bool {
	return status == DeploymentStatusSuccess || status == DeploymentStatusFailed || status == DeploymentStatusCanceled
}

func deploymentCancelable(status string) bool {
	return status == DeploymentStatusPending || status == DeploymentStatusRunning || status == DeploymentStatusCanceling
}

func eventLevelForStatus(status string) string {
	if status == DeploymentStatusFailed || status == DeploymentStatusCanceled {
		return "error"
	}
	return "info"
}

func finishMessage(status, fallback string) string {
	if fallback != "" {
		return fallback
	}
	switch status {
	case DeploymentStatusSuccess:
		return "deployment completed"
	case DeploymentStatusCanceled:
		return "deployment canceled"
	default:
		return "deployment failed"
	}
}

func newDeploymentID() string {
	now := time.Now()
	buf := make([]byte, 3)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("CMP-%s-%d", now.Format("20060102"), now.UnixNano()%1000000)
	}
	return fmt.Sprintf("CMP-%s-%s", now.Format("20060102"), strings.ToUpper(hex.EncodeToString(buf)))
}

func SortedDeploymentEvents(events []model.DeploymentEvent) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Seq < events[j].Seq
	})
}
