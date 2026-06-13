package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type NotificationUseCase interface {
	CreateNotification(tenantID, userID uuid.UUID, title, message, notifType string, referenceID *uuid.UUID) error
	ListNotifications(tenantID, userID uuid.UUID, query request.ListNotificationQuery) ([]response.NotificationResponse, int64, error)
	GetUnreadCount(tenantID, userID uuid.UUID) (*response.UnreadCountResponse, error)
	MarkAsRead(notificationID uuid.UUID) error
	MarkAllAsRead(tenantID, userID uuid.UUID) error
	DeleteNotification(notificationID uuid.UUID) error
}
