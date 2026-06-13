package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type AnalyticsUseCase interface {
	GetDashboard(tenantID uuid.UUID) (*response.AnalyticsResponse, error)
	GetRegionalStats() ([]response.RegionalStatResponse, error)
}
