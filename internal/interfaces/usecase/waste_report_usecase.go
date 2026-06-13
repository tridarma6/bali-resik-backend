package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type WasteReportUseCase interface {
	CreateReport(tenantID, userID uuid.UUID, req *request.CreateReportRequest) (*response.WasteReportResponse, error)
	GetReport(tenantID, reportID uuid.UUID) (*response.WasteReportResponse, error)
	ListReports(tenantID uuid.UUID, query request.ListReportQuery) ([]response.WasteReportResponse, int64, error)
	ListMyReports(tenantID, userID uuid.UUID, query request.ListReportQuery) ([]response.WasteReportResponse, int64, error)
	FindNearby(tenantID uuid.UUID, query request.NearbyQuery) ([]response.WasteReportResponse, error)
	UpdateReportStatus(tenantID, reportID uuid.UUID, status string) (*response.WasteReportResponse, error)
	AddReportImage(tenantID, reportID uuid.UUID, imageURL string) error
}
