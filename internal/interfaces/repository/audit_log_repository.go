package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	FindByTenant(tenantID uuid.UUID, page, perPage int) ([]models.AuditLog, int64, error)
}
