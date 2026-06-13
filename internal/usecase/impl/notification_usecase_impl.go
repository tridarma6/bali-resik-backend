package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
)

type NotificationUseCaseImpl struct {
	notifRepo repo.NotificationRepository
	log       *logrus.Logger
}

func NewNotificationUseCase(
	notifRepo repo.NotificationRepository,
	log *logrus.Logger,
) *NotificationUseCaseImpl {
	return &NotificationUseCaseImpl{
		notifRepo: notifRepo,
		log:       log,
	}
}

func (uc *NotificationUseCaseImpl) CreateNotification(tenantID, userID uuid.UUID, title, message, notifType string, referenceID *uuid.UUID) error {
	notification := &models.Notification{
		TenantID:    tenantID,
		UserID:      userID,
		Title:       title,
		Message:     message,
		Type:        notifType,
		ReferenceID: referenceID,
	}

	if err := uc.notifRepo.Create(notification); err != nil {
		uc.log.WithError(err).Error("failed to create notification")
		return errors.NewDomainError("INTERNAL_ERROR", "Failed to create notification", err)
	}

	return nil
}

func (uc *NotificationUseCaseImpl) ListNotifications(tenantID, userID uuid.UUID, query request.ListNotificationQuery) ([]response.NotificationResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	notifications, total, err := uc.notifRepo.FindByUser(tenantID, userID, query.Read, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list notifications")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list notifications", err)
	}

	result := make([]response.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, response.NotificationResponse{
			ID:          n.ID,
			Title:       n.Title,
			Message:     n.Message,
			Type:        n.Type,
			ReferenceID: n.ReferenceID,
			IsRead:      n.IsRead,
			CreatedAt:   n.CreatedAt,
		})
	}

	return result, total, nil
}

func (uc *NotificationUseCaseImpl) GetUnreadCount(tenantID, userID uuid.UUID) (*response.UnreadCountResponse, error) {
	count, err := uc.notifRepo.CountUnread(tenantID, userID)
	if err != nil {
		uc.log.WithError(err).Error("failed to count unread notifications")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to count unread notifications", err)
	}

	return &response.UnreadCountResponse{Count: count}, nil
}

func (uc *NotificationUseCaseImpl) MarkAsRead(notificationID uuid.UUID) error {
	notif, err := uc.notifRepo.FindByID(notificationID)
	if err != nil {
		return errors.NewDomainError("NOT_FOUND", "Notification not found", errors.ErrNotFound)
	}

	if notif.IsRead {
		return nil
	}

	return uc.notifRepo.MarkAsRead(notificationID)
}

func (uc *NotificationUseCaseImpl) MarkAllAsRead(tenantID, userID uuid.UUID) error {
	return uc.notifRepo.MarkAllAsRead(tenantID, userID)
}

func (uc *NotificationUseCaseImpl) DeleteNotification(notificationID uuid.UUID) error {
	if err := uc.notifRepo.Delete(notificationID); err != nil {
		return errors.NewDomainError("NOT_FOUND", "Notification not found", errors.ErrNotFound)
	}
	return nil
}
