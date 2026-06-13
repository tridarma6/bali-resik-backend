package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index:idx_notif_user_read;not null" json:"user_id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Type        string         `gorm:"type:varchar(30);index;not null" json:"type"`
	ReferenceID *uuid.UUID     `gorm:"type:uuid" json:"reference_id,omitempty"`
	IsRead      bool           `gorm:"index:idx_notif_user_read;default:false;not null" json:"is_read"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (Notification) TableName() string {
	return "notifications"
}
