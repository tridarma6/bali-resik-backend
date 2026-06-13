package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/helper"
	"gorm.io/gorm"
)

type NotificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepositoryImpl {
	return &NotificationRepositoryImpl{db: db}
}

func (r *NotificationRepositoryImpl) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepositoryImpl) FindByID(id uuid.UUID) (*models.Notification, error) {
	var n models.Notification
	err := r.db.First(&n, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NotificationRepositoryImpl) FindByUser(tenantID, userID uuid.UUID, readFilter string, page, perPage int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	query := r.db.Model(&models.Notification{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	if readFilter == "unread" {
		query = query.Where("is_read = ?", false)
	} else if readFilter == "read" {
		query = query.Where("is_read = ?", true)
	}

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, meta.Total, nil
}

func (r *NotificationRepositoryImpl) CountUnread(tenantID, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("tenant_id = ? AND user_id = ? AND is_read = ?", tenantID, userID, false).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepositoryImpl) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *NotificationRepositoryImpl) MarkAllAsRead(tenantID, userID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("tenant_id = ? AND user_id = ? AND is_read = ?", tenantID, userID, false).
		Update("is_read", true).Error
}

func (r *NotificationRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Notification{}, "id = ?", id).Error
}
