package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByID(id uuid.UUID) (*models.Notification, error)
	FindByUser(tenantID, userID uuid.UUID, readFilter string, page, perPage int) ([]models.Notification, int64, error)
	CountUnread(tenantID, userID uuid.UUID) (int64, error)
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(tenantID, userID uuid.UUID) error
	Delete(id uuid.UUID) error
}
