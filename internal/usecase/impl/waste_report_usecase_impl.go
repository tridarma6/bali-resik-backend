package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/constants"
	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
)

type WasteReportUseCaseImpl struct {
	reportRepo repo.WasteReportRepository
	userRepo   repo.UserRepository
	log        *logrus.Logger
}

func NewWasteReportUseCase(
	reportRepo repo.WasteReportRepository,
	userRepo repo.UserRepository,
	log *logrus.Logger,
) *WasteReportUseCaseImpl {
	return &WasteReportUseCaseImpl{
		reportRepo: reportRepo,
		userRepo:   userRepo,
		log:        log,
	}
}

func (uc *WasteReportUseCaseImpl) CreateReport(tenantID, userID uuid.UUID, req *request.CreateReportRequest) (*response.WasteReportResponse, error) {
	report := &models.WasteReport{
		TenantID:    tenantID,
		UserID:      userID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		Status:      constants.ReportStatusReported,
		Severity:    req.Severity,
	}

	if err := uc.reportRepo.Create(report); err != nil {
		uc.log.WithError(err).Error("failed to create waste report")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create report", err)
	}

	return uc.toReportResponse(report), nil
}

func (uc *WasteReportUseCaseImpl) GetReport(tenantID, reportID uuid.UUID) (*response.WasteReportResponse, error) {
	report, err := uc.reportRepo.FindByID(reportID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Waste report not found", errors.ErrNotFound)
	}

	if report.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied to this report", errors.ErrForbidden)
	}

	return uc.toReportResponse(report), nil
}

func (uc *WasteReportUseCaseImpl) ListReports(tenantID uuid.UUID, query request.ListReportQuery) ([]response.WasteReportResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	reports, total, err := uc.reportRepo.FindByTenant(tenantID, query.Status, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list reports")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list reports", err)
	}

	result := make([]response.WasteReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, *uc.toReportResponse(&r))
	}

	return result, total, nil
}

func (uc *WasteReportUseCaseImpl) ListMyReports(tenantID, userID uuid.UUID, query request.ListReportQuery) ([]response.WasteReportResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	reports, total, err := uc.reportRepo.FindByUser(tenantID, userID, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list user reports")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list reports", err)
	}

	result := make([]response.WasteReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, *uc.toReportResponse(&r))
	}

	return result, total, nil
}

func (uc *WasteReportUseCaseImpl) FindNearby(tenantID uuid.UUID, query request.NearbyQuery) ([]response.WasteReportResponse, error) {
	reports, err := uc.reportRepo.FindNearby(tenantID, query.Latitude, query.Longitude, query.RadiusKm)
	if err != nil {
		uc.log.WithError(err).Error("failed to find nearby reports")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to find nearby reports", err)
	}

	result := make([]response.WasteReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, *uc.toReportResponse(&r))
	}

	return result, nil
}

func (uc *WasteReportUseCaseImpl) UpdateReportStatus(tenantID, reportID uuid.UUID, status string) (*response.WasteReportResponse, error) {
	report, err := uc.reportRepo.FindByID(reportID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Waste report not found", errors.ErrNotFound)
	}

	if report.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied", errors.ErrForbidden)
	}

	if err := uc.validateReportTransition(report.Status, status); err != nil {
		return nil, err
	}

	if err := uc.reportRepo.UpdateStatus(reportID, status); err != nil {
		uc.log.WithError(err).Error("failed to update report status")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to update status", err)
	}

	report, err = uc.reportRepo.FindByID(reportID)
	if err != nil {
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve updated report", err)
	}

	return uc.toReportResponse(report), nil
}

func (uc *WasteReportUseCaseImpl) AddReportImage(tenantID, reportID uuid.UUID, imageURL string) error {
	report, err := uc.reportRepo.FindByID(reportID)
	if err != nil {
		return errors.NewDomainError("NOT_FOUND", "Waste report not found", errors.ErrNotFound)
	}

	if report.TenantID != tenantID {
		return errors.NewDomainError("FORBIDDEN", "Access denied", errors.ErrForbidden)
	}

	image := &models.ReportImage{
		TenantID: tenantID,
		ReportID: reportID,
		ImageURL: imageURL,
	}

	return uc.reportRepo.AddImage(image)
}

func (uc *WasteReportUseCaseImpl) validateReportTransition(current, next string) error {
	transitions := map[string][]string{
		constants.ReportStatusReported: {constants.ReportStatusVerified, constants.ReportStatusRejected},
		constants.ReportStatusVerified: {constants.ReportStatusCleaning, constants.ReportStatusRejected},
		constants.ReportStatusCleaning: {constants.ReportStatusResolved},
	}

	allowed, ok := transitions[current]
	if !ok {
		return errors.NewDomainError("INVALID_STATUS", "Cannot change status from "+current, nil)
	}

	for _, s := range allowed {
		if s == next {
			return nil
		}
	}

	return errors.NewDomainError("INVALID_TRANSITION",
		"Cannot transition from "+current+" to "+next, nil)
}

func (uc *WasteReportUseCaseImpl) toReportResponse(r *models.WasteReport) *response.WasteReportResponse {
	resp := &response.WasteReportResponse{
		ID:          r.ID,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Description: r.Description,
		Status:      r.Status,
		Severity:    r.Severity,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	if r.User != nil {
		resp.User = &response.UserRef{
			ID:       r.User.ID,
			FullName: r.User.FullName,
			Email:    r.User.Email,
			Phone:    r.User.Phone,
		}
	}

	if len(r.Images) > 0 {
		images := make([]response.ReportImageResp, 0, len(r.Images))
		for _, img := range r.Images {
			images = append(images, response.ReportImageResp{
				ID:       img.ID,
				ImageURL: img.ImageURL,
			})
		}
		resp.Images = images
	}

	return resp
}
