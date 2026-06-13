package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex:idx_users_email_tenant;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string         `gorm:"type:varchar(150);not null" json:"full_name"`
	Phone        string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	AvatarURL    string         `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	IsActive     bool           `gorm:"default:true;not null" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Roles  []Role  `gorm:"-" json:"roles,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
