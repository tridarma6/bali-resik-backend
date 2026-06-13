package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/helper"
	"github.com/indim/bali-resik-backend/internal/middleware"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type NotificationHandler struct {
	notifUseCase ucase.NotificationUseCase
	log          *logrus.Logger
}

func NewNotificationHandler(notifUseCase ucase.NotificationUseCase, log *logrus.Logger) *NotificationHandler {
	return &NotificationHandler{
		notifUseCase: notifUseCase,
		log:          log,
	}
}

func (h *NotificationHandler) List(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var query request.ListNotificationQuery
	c.Bind(&query)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	notifications, total, err := h.notifUseCase.ListNotifications(tenantID, userID, query)
	if err != nil {
		return h.handleNotifError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:    query.Page,
		PerPage: query.PerPage,
		Total:   total,
	}

	return helper.SuccessPaginated(c, notifications, meta, "Notifications retrieved")
}

func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	resp, err := h.notifUseCase.GetUnreadCount(tenantID, userID)
	if err != nil {
		return h.handleNotifError(c, err)
	}

	return helper.SuccessOK(c, resp, "Unread count retrieved")
}

func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid notification ID")
	}

	if err := h.notifUseCase.MarkAsRead(notificationID); err != nil {
		return h.handleNotifError(c, err)
	}

	return helper.SuccessOK(c, nil, "Notification marked as read")
}

func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	if err := h.notifUseCase.MarkAllAsRead(tenantID, userID); err != nil {
		return h.handleNotifError(c, err)
	}

	return helper.SuccessOK(c, nil, "All notifications marked as read")
}

func (h *NotificationHandler) Delete(c echo.Context) error {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid notification ID")
	}

	if err := h.notifUseCase.DeleteNotification(notificationID); err != nil {
		return h.handleNotifError(c, err)
	}

	return helper.SuccessOK(c, nil, "Notification deleted")
}

func (h *NotificationHandler) handleNotifError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled notification error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
