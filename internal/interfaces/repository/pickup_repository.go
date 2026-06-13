package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type PickupRepository interface {
	Create(pickup *models.PickupRequest) error
	FindByID(id uuid.UUID) (*models.PickupRequest, error)
	FindByTenant(tenantID uuid.UUID, status string, page, perPage int) ([]models.PickupRequest, int64, error)
	FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.PickupRequest, int64, error)
	FindByCollector(tenantID, collectorID uuid.UUID, page, perPage int) ([]models.PickupRequest, int64, error)
	Update(pickup *models.PickupRequest) error
	AssignCollector(id, collectorID uuid.UUID) error
	UpdateStatus(id uuid.UUID, status string) error
	CancelPickup(id uuid.UUID) error
}
