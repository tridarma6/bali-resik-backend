package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/helper"
	"github.com/indim/bali-resik-backend/internal/middleware"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type WasteReportHandler struct {
	reportUseCase ucase.WasteReportUseCase
	log           *logrus.Logger
}

func NewWasteReportHandler(reportUseCase ucase.WasteReportUseCase, log *logrus.Logger) *WasteReportHandler {
	return &WasteReportHandler{
		reportUseCase: reportUseCase,
		log:           log,
	}
}

func (h *WasteReportHandler) Create(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var req request.CreateReportRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.reportUseCase.CreateReport(tenantID, userID, &req)
	if err != nil {
		return h.handleReportError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Waste report created")
}

func (h *WasteReportHandler) GetByID(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid report ID")
	}

	resp, err := h.reportUseCase.GetReport(tenantID, reportID)
	if err != nil {
		return h.handleReportError(c, err)
	}

	return helper.SuccessOK(c, resp, "Waste report retrieved")
}

func (h *WasteReportHandler) List(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	var query request.ListReportQuery
	c.Bind(&query)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	reports, total, err := h.reportUseCase.ListReports(tenantID, query)
	if err != nil {
		return h.handleReportError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:    query.Page,
		PerPage: query.PerPage,
		Total:   total,
	}

	return helper.SuccessPaginated(c, reports, meta, "Waste reports retrieved")
}

func (h *WasteReportHandler) ListMine(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var query request.ListReportQuery
	c.Bind(&query)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	reports, total, err := h.reportUseCase.ListMyReports(tenantID, userID, query)
	if err != nil {
		return h.handleReportError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:    query.Page,
		PerPage: query.PerPage,
		Total:   total,
	}

	return helper.SuccessPaginated(c, reports, meta, "My waste reports retrieved")
}

func (h *WasteReportHandler) FindNearby(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	var query request.NearbyQuery
	if err := c.Bind(&query); err != nil {
		return helper.BadRequest(c, "Invalid query parameters")
	}

	reports, err := h.reportUseCase.FindNearby(tenantID, query)
	if err != nil {
		return h.handleReportError(c, err)
	}

	return helper.SuccessOK(c, reports, "Nearby waste reports retrieved")
}

func (h *WasteReportHandler) UpdateStatus(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid report ID")
	}

	var req request.UpdateReportStatusRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.reportUseCase.UpdateReportStatus(tenantID, reportID, req.Status)
	if err != nil {
		return h.handleReportError(c, err)
	}

	return helper.SuccessOK(c, resp, "Report status updated")
}

func (h *WasteReportHandler) UploadImage(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid report ID")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return helper.BadRequest(c, "Image file is required")
	}

	if file.Size > 10*1024*1024 {
		return helper.BadRequest(c, "Image file must be less than 10MB")
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return helper.BadRequest(c, "Only JPG, PNG, and WebP images are allowed")
	}

	filename := fmt.Sprintf("%s_%d%s", reportID.String(), time.Now().UnixNano(), ext)
	uploadDir := "uploads/reports"
	os.MkdirAll(uploadDir, 0755)

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		h.log.WithError(err).Error("failed to create upload file")
		return helper.InternalError(c, "Failed to upload image")
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return helper.InternalError(c, "Failed to read uploaded file")
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		h.log.WithError(err).Error("failed to save uploaded file")
		return helper.InternalError(c, "Failed to save image")
	}

	imageURL := fmt.Sprintf("/uploads/reports/%s", filename)

	if err := h.reportUseCase.AddReportImage(tenantID, reportID, imageURL); err != nil {
		return h.handleReportError(c, err)
	}

	return helper.SuccessCreated(c, map[string]string{
		"image_url": imageURL,
	}, "Image uploaded successfully")
}

func (h *WasteReportHandler) handleReportError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled report error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "FORBIDDEN":
		return helper.Forbidden(c, domainErr.Message)
	case "INVALID_STATUS", "INVALID_TRANSITION":
		return helper.BadRequest(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
