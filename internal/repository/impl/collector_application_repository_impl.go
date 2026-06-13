package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/helper"
	"gorm.io/gorm"
)

type CollectorApplicationRepositoryImpl struct {
	db *gorm.DB
}

func NewCollectorApplicationRepository(db *gorm.DB) *CollectorApplicationRepositoryImpl {
	return &CollectorApplicationRepositoryImpl{db: db}
}

func (r *CollectorApplicationRepositoryImpl) Create(app *models.CollectorApplication) error {
	return r.db.Create(app).Error
}

func (r *CollectorApplicationRepositoryImpl) FindByID(id uuid.UUID) (*models.CollectorApplication, error) {
	var app models.CollectorApplication
	err := r.db.
		Preload("User").
		First(&app, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *CollectorApplicationRepositoryImpl) FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.CollectorApplication, int64, error) {
	var apps []models.CollectorApplication
	query := r.db.Model(&models.CollectorApplication{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.Order("created_at DESC").Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, meta.Total, nil
}

func (r *CollectorApplicationRepositoryImpl) FindByTenant(tenantID uuid.UUID, status string, page, perPage int) ([]models.CollectorApplication, int64, error) {
	var apps []models.CollectorApplication
	query := r.db.Model(&models.CollectorApplication{}).Where("tenant_id = ?", tenantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.
		Preload("User").
		Order("created_at DESC").
		Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	return apps, meta.Total, nil
}

func (r *CollectorApplicationRepositoryImpl) FindPendingByUser(tenantID, userID uuid.UUID) (*models.CollectorApplication, error) {
	var app models.CollectorApplication
	err := r.db.Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, "pending").
		First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *CollectorApplicationRepositoryImpl) Update(app *models.CollectorApplication) error {
	return r.db.Save(app).Error
}
