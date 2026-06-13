package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/helper"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type AdminHandler struct {
	adminUseCase ucase.AdminUseCase
	log          *logrus.Logger
}

func NewAdminHandler(adminUseCase ucase.AdminUseCase, log *logrus.Logger) *AdminHandler {
	return &AdminHandler{
		adminUseCase: adminUseCase,
		log:          log,
	}
}

func (h *AdminHandler) CreateTenant(c echo.Context) error {
	var req request.CreateTenantRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.adminUseCase.CreateTenant(&req)
	if err != nil {
		return h.handleAdminError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Tenant created successfully")
}

func (h *AdminHandler) ListTenants(c echo.Context) error {
	resp, err := h.adminUseCase.ListTenants()
	if err != nil {
		return h.handleAdminError(c, err)
	}

	return helper.SuccessOK(c, resp, "Tenants retrieved")
}

func (h *AdminHandler) CreateAdmin(c echo.Context) error {
	var req request.CreateAdminRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.adminUseCase.CreateAdmin(&req)
	if err != nil {
		return h.handleAdminError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Admin user created successfully")
}

func (h *AdminHandler) handleAdminError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled admin error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "SLUG_EXISTS":
		return helper.Conflict(c, domainErr.Message)
	case "TENANT_NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "EMAIL_EXISTS":
		return helper.Conflict(c, domainErr.Message)
	case "INVALID_TENANT_ID":
		return helper.BadRequest(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
