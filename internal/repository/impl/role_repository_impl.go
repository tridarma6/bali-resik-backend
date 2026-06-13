package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepositoryImpl {
	return &RoleRepositoryImpl{db: db}
}

func (r *RoleRepositoryImpl) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepositoryImpl) FindByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepositoryImpl) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepositoryImpl) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Order("name ASC").Find(&roles).Error
	return roles, err
}

func (r *RoleRepositoryImpl) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Role{}, "id = ?", id).Error
}
