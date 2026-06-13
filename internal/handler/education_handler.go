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

type EducationHandler struct {
	eduUseCase ucase.EducationUseCase
	log        *logrus.Logger
}

func NewEducationHandler(eduUseCase ucase.EducationUseCase, log *logrus.Logger) *EducationHandler {
	return &EducationHandler{
		eduUseCase: eduUseCase,
		log:        log,
	}
}

func (h *EducationHandler) Create(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var req request.CreateContentRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.eduUseCase.CreateContent(tenantID, userID, &req)
	if err != nil {
		return h.handleEduError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Content created")
}

func (h *EducationHandler) GetByID(c echo.Context) error {
	contentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid content ID")
	}

	resp, err := h.eduUseCase.GetContent(contentID)
	if err != nil {
		return h.handleEduError(c, err)
	}

	return helper.SuccessOK(c, resp, "Content retrieved")
}

func (h *EducationHandler) List(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	var query request.ListContentQuery
	c.Bind(&query)

	resp, err := h.eduUseCase.ListContents(tenantID, query)
	if err != nil {
		return h.handleEduError(c, err)
	}

	return helper.SuccessOK(c, resp, "Content list retrieved")
}

func (h *EducationHandler) Update(c echo.Context) error {
	contentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid content ID")
	}

	var req request.UpdateContentRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.eduUseCase.UpdateContent(contentID, &req)
	if err != nil {
		return h.handleEduError(c, err)
	}

	return helper.SuccessOK(c, resp, "Content updated")
}

func (h *EducationHandler) Delete(c echo.Context) error {
	contentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid content ID")
	}

	if err := h.eduUseCase.DeleteContent(contentID); err != nil {
		return h.handleEduError(c, err)
	}

	return helper.SuccessOK(c, nil, "Content deleted")
}

func (h *EducationHandler) handleEduError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled education error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
