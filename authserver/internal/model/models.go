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
