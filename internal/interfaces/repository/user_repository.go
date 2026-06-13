package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByEmailAndTenant(email string, tenantID uuid.UUID) (*models.User, error)
	FindByTenant(tenantID uuid.UUID) ([]models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	AssignRole(userID, roleID, tenantID uuid.UUID) error
	GetRoles(userID uuid.UUID) ([]models.Role, error)
}
