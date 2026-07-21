package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"size:128;not null;uniqueIndex" json:"username"`
	DisplayName string         `gorm:"size:128;not null;default:''" json:"display_name"`
	Email       string         `gorm:"size:255;not null;default:'';index" json:"email"`
	Phone       string         `gorm:"size:32;not null;default:''" json:"phone"`
	Source      string         `gorm:"size:32;not null;default:'local';index:idx_users_external,priority:1" json:"source"`
	ExternalID  string         `gorm:"size:128;not null;default:'';index:idx_users_external,priority:2" json:"external_id"`
	Status      string         `gorm:"size:32;not null;default:'active';index" json:"status"`
	IsAdmin     bool           `gorm:"not null;default:false" json:"is_admin"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WayenCredential struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:512;not null" json:"-"`
	Enabled   bool           `gorm:"not null;default:true;index" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type BusinessLine struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BusinessLineUser struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	BusinessLineID uint64    `gorm:"not null;uniqueIndex:idx_business_line_users_unique,priority:1;index" json:"business_line_id"`
	UserID         uint64    `gorm:"not null;uniqueIndex:idx_business_line_users_unique,priority:2;index" json:"user_id"`
	Permission     int       `gorm:"not null;default:1" json:"permission"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type BusinessLineWayneNamespace struct {
	ID                 uint64    `gorm:"primaryKey" json:"id"`
	BusinessLineID     uint64    `gorm:"not null;uniqueIndex:idx_business_line_wayne_namespaces_unique,priority:1;index" json:"business_line_id"`
	WayneNamespaceID   uint64    `gorm:"not null;uniqueIndex:idx_business_line_wayne_namespaces_unique,priority:2" json:"wayne_namespace_id"`
	WayneNamespaceName string    `gorm:"size:128;not null;default:''" json:"wayne_namespace_name"`
	KubeNamespace      string    `gorm:"size:128;not null;default:''" json:"kube_namespace"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type AccessToken struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	UserID    uint64     `gorm:"not null;index" json:"user_id"`
	TokenID   string     `gorm:"size:128;not null;uniqueIndex" json:"token_id"`
	TokenType string     `gorm:"size:32;not null;default:'access'" json:"token_type"`
	ClientIP  string     `gorm:"size:64;not null;default:''" json:"client_ip"`
	UserAgent string     `gorm:"size:512;not null;default:''" json:"user_agent"`
	Revoked   bool       `gorm:"not null;default:false;index" json:"revoked"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type WayneToken struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Account   string    `gorm:"size:128;not null;uniqueIndex" json:"account"`
	Token     string    `gorm:"type:text;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Deployment struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	DeploymentID   string     `gorm:"size:64;not null;uniqueIndex" json:"deployment_id"`
	Component      string     `gorm:"size:64;not null;index" json:"component"`
	BusinessLineID uint64     `gorm:"not null;index" json:"business_line_id"`
	BusinessLine   string     `gorm:"size:128;not null;default:''" json:"business_line"`
	Status         string     `gorm:"size:32;not null;index" json:"status"`
	RequestPayload string     `gorm:"type:longtext" json:"request_payload"`
	ResultPayload  string     `gorm:"type:longtext" json:"result_payload"`
	CreatedBy      uint64     `gorm:"not null;index" json:"created_by"`
	CreatedByName  string     `gorm:"size:128;not null;default:''" json:"created_by_name"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	CreatedAt      time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type DeploymentEvent struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	DeploymentID string    `gorm:"size:64;not null;index:idx_deployment_events_deployment_seq,priority:1" json:"deployment_id"`
	Seq          uint64    `gorm:"not null;index:idx_deployment_events_deployment_seq,priority:2" json:"seq"`
	Type         string    `gorm:"size:32;not null;index" json:"type"`
	Level        string    `gorm:"size:32;not null;default:''" json:"level"`
	Message      string    `gorm:"type:longtext" json:"message"`
	Payload      string    `gorm:"type:longtext" json:"payload"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

type AuditLog struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	RequestID      string    `gorm:"size:128;not null;default:''" json:"request_id"`
	ActorUserID    uint64    `gorm:"not null;default:0;index:idx_actor_time,priority:1" json:"actor_user_id"`
	ActorUsername  string    `gorm:"size:128;not null;default:''" json:"actor_username"`
	ClientIP       string    `gorm:"size:64;not null;default:''" json:"client_ip"`
	UserAgent      string    `gorm:"size:512;not null;default:''" json:"user_agent"`
	Action         string    `gorm:"size:128;not null;index:idx_action_time,priority:1" json:"action"`
	ResourceType   string    `gorm:"size:64;not null;default:'';index:idx_resource,priority:1" json:"resource_type"`
	ResourceID     string    `gorm:"size:128;not null;default:'';index:idx_resource,priority:2" json:"resource_id"`
	ScopeType      string    `gorm:"size:32;not null;default:'';index:idx_scope,priority:1" json:"scope_type"`
	ScopeID        uint64    `gorm:"not null;default:0;index:idx_scope,priority:2" json:"scope_id"`
	BusinessLineID uint64    `gorm:"not null;default:0" json:"business_line_id"`
	NamespaceID    uint64    `gorm:"not null;default:0;index:idx_ns_cluster_env_time,priority:1" json:"namespace_id"`
	ClusterID      uint64    `gorm:"not null;default:0;index:idx_ns_cluster_env_time,priority:2" json:"cluster_id"`
	EnvironmentID  uint64    `gorm:"not null;default:0;index:idx_ns_cluster_env_time,priority:3" json:"environment_id"`
	Decision       string    `gorm:"size:32;not null;default:'';index:idx_decision_time,priority:1" json:"decision"`
	Reason         string    `gorm:"size:512;not null;default:''" json:"reason"`
	Metadata       string    `gorm:"type:longtext" json:"metadata"`
	CreatedAt      time.Time `gorm:"index:idx_actor_time,priority:2;index:idx_action_time,priority:2;index:idx_ns_cluster_env_time,priority:4;index:idx_decision_time,priority:2" json:"created_at"`
}
