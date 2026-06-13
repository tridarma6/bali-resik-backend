package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WasteReport struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Latitude    float64        `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude   float64        `gorm:"type:decimal(10,7);not null" json:"longitude"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Status      string         `gorm:"type:varchar(20);index;not null;default:reported" json:"status"`
	Severity    string         `gorm:"type:varchar(20);not null;default:medium" json:"severity"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Tenant *Tenant      `gorm:"foreignKey:TenantID" json:"-"`
	User   *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Images []ReportImage `gorm:"foreignKey:ReportID" json:"images,omitempty"`
}

func (w *WasteReport) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (WasteReport) TableName() string {
	return "waste_reports"
}

type ReportImage struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ReportID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"report_id"`
	ImageURL  string         `gorm:"type:varchar(500);not null" json:"image_url"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *ReportImage) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (ReportImage) TableName() string {
	return "report_images"
}
