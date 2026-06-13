package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type EducationRepository interface {
	Create(content *models.EducationalContent) error
	FindByID(id uuid.UUID) (*models.EducationalContent, error)
	FindByTenant(tenantID uuid.UUID, category string, publishedOnly bool) ([]models.EducationalContent, error)
	Update(content *models.EducationalContent) error
	Delete(id uuid.UUID) error
}
