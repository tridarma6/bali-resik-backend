package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type RoleRepository interface {
	Create(role *models.Role) error
	FindByID(id uuid.UUID) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	FindAll() ([]models.Role, error)
	Update(role *models.Role) error
	Delete(id uuid.UUID) error
}
