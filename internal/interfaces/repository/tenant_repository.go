package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type TenantRepository interface {
	Create(tenant *models.Tenant) error
	FindByID(id uuid.UUID) (*models.Tenant, error)
	FindBySlug(slug string) (*models.Tenant, error)
	FindAll() ([]models.Tenant, error)
	Update(tenant *models.Tenant) error
	Delete(id uuid.UUID) error
}
