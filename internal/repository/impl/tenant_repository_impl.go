package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"gorm.io/gorm"
)

type TenantRepositoryImpl struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepositoryImpl {
	return &TenantRepositoryImpl{db: db}
}

func (r *TenantRepositoryImpl) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *TenantRepositoryImpl) FindByID(id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.First(&tenant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepositoryImpl) FindBySlug(slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.First(&tenant, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepositoryImpl) FindAll() ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.db.Order("name ASC").Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepositoryImpl) Update(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *TenantRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Tenant{}, "id = ?", id).Error
}
