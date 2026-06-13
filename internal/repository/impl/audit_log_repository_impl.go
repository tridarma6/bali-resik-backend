package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/helper"
	"gorm.io/gorm"
)

type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepositoryImpl {
	return &AuditLogRepositoryImpl{db: db}
}

func (r *AuditLogRepositoryImpl) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepositoryImpl) FindByTenant(tenantID uuid.UUID, page, perPage int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	query := r.db.Where("tenant_id = ?", tenantID)

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, meta.Total, nil
}
