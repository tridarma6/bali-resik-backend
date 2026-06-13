package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/helper"
	"github.com/indim/bali-resik-backend/internal/middleware"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type AnalyticsHandler struct {
	analyticsUseCase ucase.AnalyticsUseCase
	log              *logrus.Logger
}

func NewAnalyticsHandler(analyticsUseCase ucase.AnalyticsUseCase, log *logrus.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsUseCase: analyticsUseCase,
		log:              log,
	}
}

func (h *AnalyticsHandler) GetDashboard(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	resp, err := h.analyticsUseCase.GetDashboard(tenantID)
	if err != nil {
		return h.handleAnalyticsError(c, err)
	}

	return helper.SuccessOK(c, resp, "Dashboard data retrieved")
}

func (h *AnalyticsHandler) GetRegionalStats(c echo.Context) error {
	resp, err := h.analyticsUseCase.GetRegionalStats()
	if err != nil {
		return h.handleAnalyticsError(c, err)
	}

	return helper.SuccessOK(c, resp, "Regional statistics retrieved")
}

func (h *AnalyticsHandler) handleAnalyticsError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled analytics error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	return helper.InternalError(c, domainErr.Message)
}
