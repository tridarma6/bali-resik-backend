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

type PickupHandler struct {
	pickupUseCase ucase.PickupUseCase
	log           *logrus.Logger
}

func NewPickupHandler(pickupUseCase ucase.PickupUseCase, log *logrus.Logger) *PickupHandler {
	return &PickupHandler{
		pickupUseCase: pickupUseCase,
		log:           log,
	}
}

func (h *PickupHandler) Create(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var req request.CreatePickupRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.pickupUseCase.CreatePickup(tenantID, userID, &req)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Pickup request created")
}

func (h *PickupHandler) GetByID(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	pickupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid pickup ID")
	}

	resp, err := h.pickupUseCase.GetPickup(tenantID, pickupID)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	return helper.SuccessOK(c, resp, "Pickup request retrieved")
}

func (h *PickupHandler) List(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	var query request.ListPickupQuery
	if err := c.Bind(&query); err != nil {
		query = request.ListPickupQuery{Page: 1, PerPage: 20}
	}

	pickups, total, err := h.pickupUseCase.ListPickups(tenantID, query)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:      query.Page,
		PerPage:   query.PerPage,
		Total:     total,
	}

	return helper.SuccessPaginated(c, pickups, meta, "Pickup requests retrieved")
}

func (h *PickupHandler) ListMine(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var query request.ListPickupQuery
	if err := c.Bind(&query); err != nil {
		query = request.ListPickupQuery{Page: 1, PerPage: 20}
	}

	pickups, total, err := h.pickupUseCase.ListMyPickups(tenantID, userID, query)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:      query.Page,
		PerPage:   query.PerPage,
		Total:     total,
	}

	return helper.SuccessPaginated(c, pickups, meta, "My pickup requests retrieved")
}

func (h *PickupHandler) ListAssigned(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var query request.ListPickupQuery
	if err := c.Bind(&query); err != nil {
		query = request.ListPickupQuery{Page: 1, PerPage: 20}
	}

	pickups, total, err := h.pickupUseCase.ListCollectorPickups(tenantID, userID, query)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:      query.Page,
		PerPage:   query.PerPage,
		Total:     total,
	}

	return helper.SuccessPaginated(c, pickups, meta, "Assigned pickups retrieved")
}

func (h *PickupHandler) AssignCollector(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	pickupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid pickup ID")
	}

	var req request.AssignCollectorRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	collectorID, err := uuid.Parse(req.CollectorID)
	if err != nil {
		return helper.BadRequest(c, "Invalid collector ID")
	}

	resp, err := h.pickupUseCase.AssignCollector(tenantID, pickupID, collectorID)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	return helper.SuccessOK(c, resp, "Collector assigned")
}

func (h *PickupHandler) UpdateStatus(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	pickupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid pickup ID")
	}

	var req request.UpdatePickupStatusRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.pickupUseCase.UpdatePickupStatus(tenantID, pickupID, req.Status)
	if err != nil {
		return h.handlePickupError(c, err)
	}

	return helper.SuccessOK(c, resp, "Pickup status updated")
}

func (h *PickupHandler) Cancel(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	pickupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid pickup ID")
	}

	if err := h.pickupUseCase.CancelPickup(tenantID, pickupID, userID); err != nil {
		return h.handlePickupError(c, err)
	}

	return helper.SuccessOK(c, nil, "Pickup request cancelled")
}

func (h *PickupHandler) handlePickupError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled pickup error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "FORBIDDEN":
		return helper.Forbidden(c, domainErr.Message)
	case "INVALID_STATUS", "INVALID_TRANSITION":
		return helper.BadRequest(c, domainErr.Message)
	case "COLLECTOR_NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
