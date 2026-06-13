package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectorApplication struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"user_id"`
	Status     string     `gorm:"type:varchar(20);default:pending;not null;index" json:"status"`
	AdminNotes string     `gorm:"type:text" json:"admin_notes,omitempty"`
	ReviewedBy *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Admin *User  `gorm:"foreignKey:ReviewedBy" json:"admin,omitempty"`
}

func (c *CollectorApplication) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (CollectorApplication) TableName() string {
	return "collector_applications"
}
