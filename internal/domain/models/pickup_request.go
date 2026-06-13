package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickupRequest struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;index:idx_pickup_tenant_status;not null" json:"tenant_id"`
	UserID        uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	CollectorID   *uuid.UUID     `gorm:"type:uuid;index" json:"collector_id,omitempty"`
	WasteType     string         `gorm:"type:varchar(30);not null" json:"waste_type"`
	Latitude      float64        `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude     float64        `gorm:"type:decimal(10,7);not null" json:"longitude"`
	Address       string         `gorm:"type:text;not null" json:"address"`
	Status        string         `gorm:"type:varchar(20);index:idx_pickup_tenant_status;not null;default:pending" json:"status"`
	ScheduledDate *time.Time     `json:"scheduled_date,omitempty"`
	Notes         string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Tenant    *Tenant `gorm:"foreignKey:TenantID" json:"-"`
	User      *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Collector *User   `gorm:"foreignKey:CollectorID" json:"collector,omitempty"`
}

func (p *PickupRequest) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (PickupRequest) TableName() string {
	return "pickup_requests"
}
