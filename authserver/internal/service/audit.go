package service

import (
	"encoding/json"

	"authserver/internal/model"

	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

type AuditEntry struct {
	RequestID      string
	ActorUserID    uint64
	ActorUsername  string
	ClientIP       string
	UserAgent      string
	Action         string
	ResourceType   string
	ResourceID     string
	ScopeType      string
	ScopeID        uint64
	BusinessLineID uint64
	NamespaceID    uint64
	ClusterID      uint64
	EnvironmentID  uint64
	Decision       string
	Reason         string
	Metadata       map[string]any
}

func (s *AuditService) Write(entry AuditEntry) {
	raw, _ := json.Marshal(entry.Metadata)
	_ = s.db.Create(&model.AuditLog{
		RequestID:      entry.RequestID,
		ActorUserID:    entry.ActorUserID,
		ActorUsername:  entry.ActorUsername,
		ClientIP:       entry.ClientIP,
		UserAgent:      entry.UserAgent,
		Action:         entry.Action,
		ResourceType:   entry.ResourceType,
		ResourceID:     entry.ResourceID,
		ScopeType:      entry.ScopeType,
		ScopeID:        entry.ScopeID,
		BusinessLineID: entry.BusinessLineID,
		NamespaceID:    entry.NamespaceID,
		ClusterID:      entry.ClusterID,
		EnvironmentID:  entry.EnvironmentID,
		Decision:       entry.Decision,
		Reason:         entry.Reason,
		Metadata:       string(raw),
	}).Error
}
