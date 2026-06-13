package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Action     string         `gorm:"type:varchar(50);index;not null" json:"action"`
	EntityType string         `gorm:"type:varchar(50);not null" json:"entity_type"`
	EntityID   uuid.UUID      `gorm:"type:uuid" json:"entity_id"`
	OldValues  string         `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues  string         `gorm:"type:jsonb" json:"new_values,omitempty"`
	IPAddress  string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string         `gorm:"type:varchar(500)" json:"user_agent"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
