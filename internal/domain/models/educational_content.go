package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EducationalContent struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	ContentType string         `gorm:"type:varchar(20);not null;default:article" json:"content_type"`
	Category    string         `gorm:"type:varchar(50);index;not null" json:"category"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	IsPublished bool           `gorm:"default:false;not null" json:"is_published"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Author *User `gorm:"foreignKey:CreatedBy" json:"author,omitempty"`
}

func (e *EducationalContent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (EducationalContent) TableName() string {
	return "educational_contents"
}
