package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type PickupUseCase interface {
	CreatePickup(tenantID, userID uuid.UUID, req *request.CreatePickupRequest) (*response.PickupResponse, error)
	GetPickup(tenantID, pickupID uuid.UUID) (*response.PickupResponse, error)
	ListPickups(tenantID uuid.UUID, query request.ListPickupQuery) ([]response.PickupResponse, int64, error)
	ListMyPickups(tenantID, userID uuid.UUID, query request.ListPickupQuery) ([]response.PickupResponse, int64, error)
	ListCollectorPickups(tenantID, collectorID uuid.UUID, query request.ListPickupQuery) ([]response.PickupResponse, int64, error)
	AssignCollector(tenantID, pickupID, collectorID uuid.UUID) (*response.PickupResponse, error)
	UpdatePickupStatus(tenantID, pickupID uuid.UUID, status string) (*response.PickupResponse, error)
	CancelPickup(tenantID, pickupID, userID uuid.UUID) error
}
