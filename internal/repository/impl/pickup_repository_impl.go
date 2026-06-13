package impl

import (
	"time"

	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/constants"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/helper"
	"gorm.io/gorm"
)

type PickupRepositoryImpl struct {
	db *gorm.DB
}

func NewPickupRepository(db *gorm.DB) *PickupRepositoryImpl {
	return &PickupRepositoryImpl{db: db}
}

func (r *PickupRepositoryImpl) Create(pickup *models.PickupRequest) error {
	return r.db.Create(pickup).Error
}

func (r *PickupRepositoryImpl) FindByID(id uuid.UUID) (*models.PickupRequest, error) {
	var pickup models.PickupRequest
	err := r.db.
		Preload("User").
		Preload("Collector").
		First(&pickup, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pickup, nil
}

func (r *PickupRepositoryImpl) FindByTenant(tenantID uuid.UUID, status string, page, perPage int) ([]models.PickupRequest, int64, error) {
	var pickups []models.PickupRequest
	query := r.db.Where("tenant_id = ?", tenantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	param := helper.PaginationParam{
		Page:    page,
		PerPage: perPage,
	}

	paginatedQuery, meta := helper.Paginate(query, param)
	err := paginatedQuery.
		Preload("User").
		Preload("Collector").
		Order("created_at DESC").
		Find(&pickups).Error

	if err != nil {
		return nil, 0, err
	}

	return pickups, meta.Total, nil
}

func (r *PickupRepositoryImpl) FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.PickupRequest, int64, error) {
	var pickups []models.PickupRequest
	query := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	param := helper.PaginationParam{
		Page:    page,
		PerPage: perPage,
	}

	paginatedQuery, meta := helper.Paginate(query, param)
	err := paginatedQuery.
		Preload("Collector").
		Order("created_at DESC").
		Find(&pickups).Error

	if err != nil {
		return nil, 0, err
	}

	return pickups, meta.Total, nil
}

func (r *PickupRepositoryImpl) FindByCollector(tenantID, collectorID uuid.UUID, page, perPage int) ([]models.PickupRequest, int64, error) {
	var pickups []models.PickupRequest
	query := r.db.Where("tenant_id = ? AND collector_id = ?", tenantID, collectorID)

	param := helper.PaginationParam{
		Page:    page,
		PerPage: perPage,
	}

	paginatedQuery, meta := helper.Paginate(query, param)
	err := paginatedQuery.
		Preload("User").
		Order("created_at DESC").
		Find(&pickups).Error

	if err != nil {
		return nil, 0, err
	}

	return pickups, meta.Total, nil
}

func (r *PickupRepositoryImpl) Update(pickup *models.PickupRequest) error {
	return r.db.Save(pickup).Error
}

func (r *PickupRepositoryImpl) AssignCollector(id, collectorID uuid.UUID) error {
	return r.db.Model(&models.PickupRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"collector_id": collectorID,
			"status":       constants.PickupStatusAssigned,
			"updated_at":   time.Now(),
		}).Error
}

func (r *PickupRepositoryImpl) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.PickupRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *PickupRepositoryImpl) CancelPickup(id uuid.UUID) error {
	return r.db.Model(&models.PickupRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     constants.PickupStatusCancelled,
			"updated_at": time.Now(),
		}).Error
}
