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

type CollectorApplicationHandler struct {
	appUseCase ucase.CollectorApplicationUseCase
	log        *logrus.Logger
}

func NewCollectorApplicationHandler(appUseCase ucase.CollectorApplicationUseCase, log *logrus.Logger) *CollectorApplicationHandler {
	return &CollectorApplicationHandler{
		appUseCase: appUseCase,
		log:        log,
	}
}

func (h *CollectorApplicationHandler) Submit(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	resp, err := h.appUseCase.Submit(tenantID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Collector application submitted")
}

func (h *CollectorApplicationHandler) ListMine(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var query request.ListCollectorAppQuery
	c.Bind(&query)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	apps, total, err := h.appUseCase.ListMyApplications(tenantID, userID, query)
	if err != nil {
		return h.handleError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:    query.Page,
		PerPage: query.PerPage,
		Total:   total,
	}

	return helper.SuccessPaginated(c, apps, meta, "Applications retrieved")
}

func (h *CollectorApplicationHandler) ListAll(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	var query request.ListCollectorAppQuery
	c.Bind(&query)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	apps, total, err := h.appUseCase.ListAll(tenantID, query)
	if err != nil {
		return h.handleError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:    query.Page,
		PerPage: query.PerPage,
		Total:   total,
	}

	return helper.SuccessPaginated(c, apps, meta, "Applications retrieved")
}

func (h *CollectorApplicationHandler) Approve(c echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid application ID")
	}

	resp, err := h.appUseCase.Approve(adminID, appID)
	if err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessOK(c, resp, "Application approved")
}

func (h *CollectorApplicationHandler) Reject(c echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid application ID")
	}

	var req request.ReviewCollectorAppRequest
	if err := c.Bind(&req); err != nil {
		req = request.ReviewCollectorAppRequest{}
	}

	resp, err := h.appUseCase.Reject(adminID, appID, req.Notes)
	if err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessOK(c, resp, "Application rejected")
}

func (h *CollectorApplicationHandler) handleError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled collector application error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "DUPLICATE_APPLICATION", "ALREADY_COLLECTOR":
		return helper.Conflict(c, domainErr.Message)
	case "INVALID_STATUS":
		return helper.BadRequest(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
