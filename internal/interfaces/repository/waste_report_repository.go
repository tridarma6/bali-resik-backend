package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type WasteReportRepository interface {
	Create(report *models.WasteReport) error
	FindByID(id uuid.UUID) (*models.WasteReport, error)
	FindByTenant(tenantID uuid.UUID, status string, page, perPage int) ([]models.WasteReport, int64, error)
	FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.WasteReport, int64, error)
	FindNearby(tenantID uuid.UUID, lat, lng float64, radiusKm float64) ([]models.WasteReport, error)
	Update(report *models.WasteReport) error
	UpdateStatus(id uuid.UUID, status string) error
	AddImage(image *models.ReportImage) error
	GetImages(reportID uuid.UUID) ([]models.ReportImage, error)
}
