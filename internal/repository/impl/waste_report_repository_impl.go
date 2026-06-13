package impl

import (
	"time"

	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/helper"
	"gorm.io/gorm"
)

type WasteReportRepositoryImpl struct {
	db *gorm.DB
}

func NewWasteReportRepository(db *gorm.DB) *WasteReportRepositoryImpl {
	return &WasteReportRepositoryImpl{db: db}
}

func (r *WasteReportRepositoryImpl) Create(report *models.WasteReport) error {
	return r.db.Create(report).Error
}

func (r *WasteReportRepositoryImpl) FindByID(id uuid.UUID) (*models.WasteReport, error) {
	var report models.WasteReport
	err := r.db.
		Preload("User").
		Preload("Images").
		First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *WasteReportRepositoryImpl) FindByTenant(tenantID uuid.UUID, status string, page, perPage int) ([]models.WasteReport, int64, error) {
	var reports []models.WasteReport
	query := r.db.Where("tenant_id = ?", tenantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.
		Preload("User").
		Preload("Images").
		Order("created_at DESC").
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, meta.Total, nil
}

func (r *WasteReportRepositoryImpl) FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.WasteReport, int64, error) {
	var reports []models.WasteReport
	query := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.
		Preload("Images").
		Order("created_at DESC").
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, meta.Total, nil
}

func (r *WasteReportRepositoryImpl) FindNearby(tenantID uuid.UUID, lat, lng float64, radiusKm float64) ([]models.WasteReport, error) {
	var reports []models.WasteReport

	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / (111.0 * 0.766)

	err := r.db.
		Where("tenant_id = ? AND latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
			tenantID, lat-latDelta, lat+latDelta, lng-lngDelta, lng+lngDelta).
		Preload("Images").
		Order("created_at DESC").
		Limit(50).
		Find(&reports).Error

	return reports, err
}

func (r *WasteReportRepositoryImpl) Update(report *models.WasteReport) error {
	return r.db.Save(report).Error
}

func (r *WasteReportRepositoryImpl) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.WasteReport{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *WasteReportRepositoryImpl) AddImage(image *models.ReportImage) error {
	return r.db.Create(image).Error
}

func (r *WasteReportRepositoryImpl) GetImages(reportID uuid.UUID) ([]models.ReportImage, error) {
	var images []models.ReportImage
	err := r.db.Where("report_id = ?", reportID).Order("created_at ASC").Find(&images).Error
	return images, err
}
