package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type CollectorApplicationUseCase interface {
	Submit(tenantID, userID uuid.UUID) (*response.CollectorApplicationResponse, error)
	ListMyApplications(tenantID, userID uuid.UUID, query request.ListCollectorAppQuery) ([]response.CollectorApplicationResponse, int64, error)
	ListAll(tenantID uuid.UUID, query request.ListCollectorAppQuery) ([]response.CollectorApplicationResponse, int64, error)
	Approve(adminID, appID uuid.UUID) (*response.CollectorApplicationResponse, error)
	Reject(adminID, appID uuid.UUID, notes string) (*response.CollectorApplicationResponse, error)
}
