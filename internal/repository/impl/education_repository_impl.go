package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"gorm.io/gorm"
)

type EducationRepositoryImpl struct {
	db *gorm.DB
}

func NewEducationRepository(db *gorm.DB) *EducationRepositoryImpl {
	return &EducationRepositoryImpl{db: db}
}

func (r *EducationRepositoryImpl) Create(content *models.EducationalContent) error {
	return r.db.Create(content).Error
}

func (r *EducationRepositoryImpl) FindByID(id uuid.UUID) (*models.EducationalContent, error) {
	var content models.EducationalContent
	err := r.db.Preload("Author").First(&content, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *EducationRepositoryImpl) FindByTenant(tenantID uuid.UUID, category string, publishedOnly bool) ([]models.EducationalContent, error) {
	query := r.db.Where("tenant_id = ?", tenantID)

	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var contents []models.EducationalContent
	err := query.Preload("Author").
		Order("created_at DESC").
		Find(&contents).Error

	return contents, err
}

func (r *EducationRepositoryImpl) Update(content *models.EducationalContent) error {
	return r.db.Save(content).Error
}

func (r *EducationRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.EducationalContent{}, "id = ?", id).Error
}
