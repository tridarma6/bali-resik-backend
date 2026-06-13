package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reward struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	PointsCost  int            `gorm:"not null" json:"points_cost"`
	Stock       *int           `gorm:"default:null" json:"stock,omitempty"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	IsActive    bool           `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Reward) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (Reward) TableName() string {
	return "rewards"
}

type RewardTransaction struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Points      int            `gorm:"not null" json:"points"`
	Type        string         `gorm:"type:varchar(10);not null" json:"type"`
	ReferenceID *uuid.UUID     `gorm:"type:uuid" json:"reference_id,omitempty"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (rt *RewardTransaction) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

func (RewardTransaction) TableName() string {
	return "reward_transactions"
}
