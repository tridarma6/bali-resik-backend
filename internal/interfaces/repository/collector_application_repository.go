package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type CollectorApplicationRepository interface {
	Create(app *models.CollectorApplication) error
	FindByID(id uuid.UUID) (*models.CollectorApplication, error)
	FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.CollectorApplication, int64, error)
	FindByTenant(tenantID uuid.UUID, status string, page, perPage int) ([]models.CollectorApplication, int64, error)
	FindPendingByUser(tenantID, userID uuid.UUID) (*models.CollectorApplication, error)
	Update(app *models.CollectorApplication) error
}
