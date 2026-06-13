package impl

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/constants"
	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
)

type PickupUseCaseImpl struct {
	pickupRepo repo.PickupRepository
	userRepo   repo.UserRepository
	log        *logrus.Logger
}

func NewPickupUseCase(
	pickupRepo repo.PickupRepository,
	userRepo repo.UserRepository,
	log *logrus.Logger,
) *PickupUseCaseImpl {
	return &PickupUseCaseImpl{
		pickupRepo: pickupRepo,
		userRepo:   userRepo,
		log:        log,
	}
}

func (uc *PickupUseCaseImpl) CreatePickup(tenantID, userID uuid.UUID, req *request.CreatePickupRequest) (*response.PickupResponse, error) {
	var scheduledDate *time.Time
	if req.ScheduledDate != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledDate)
		if err != nil {
			return nil, errors.NewDomainError("INVALID_DATE", "Invalid scheduled date format, use RFC3339", err)
		}
		scheduledDate = &t
	}

	pickup := &models.PickupRequest{
		TenantID:      tenantID,
		UserID:        userID,
		WasteType:     req.WasteType,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Address:       req.Address,
		Status:        constants.PickupStatusPending,
		ScheduledDate: scheduledDate,
		Notes:         req.Notes,
	}

	if err := uc.pickupRepo.Create(pickup); err != nil {
		uc.log.WithError(err).Error("failed to create pickup request")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create pickup request", err)
	}

	return uc.toPickupResponse(pickup), nil
}

func (uc *PickupUseCaseImpl) GetPickup(tenantID, pickupID uuid.UUID) (*response.PickupResponse, error) {
	pickup, err := uc.pickupRepo.FindByID(pickupID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Pickup request not found", errors.ErrNotFound)
	}

	if pickup.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied to this pickup request", errors.ErrForbidden)
	}

	return uc.toPickupResponse(pickup), nil
}

func (uc *PickupUseCaseImpl) ListPickups(tenantID uuid.UUID, query request.ListPickupQuery) ([]response.PickupResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	pickups, total, err := uc.pickupRepo.FindByTenant(tenantID, query.Status, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list pickups")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list pickups", err)
	}

	result := make([]response.PickupResponse, 0, len(pickups))
	for _, p := range pickups {
		result = append(result, *uc.toPickupResponse(&p))
	}

	return result, total, nil
}

func (uc *PickupUseCaseImpl) ListMyPickups(tenantID, userID uuid.UUID, query request.ListPickupQuery) ([]response.PickupResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	pickups, total, err := uc.pickupRepo.FindByUser(tenantID, userID, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list user pickups")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list pickups", err)
	}

	result := make([]response.PickupResponse, 0, len(pickups))
	for _, p := range pickups {
		result = append(result, *uc.toPickupResponse(&p))
	}

	return result, total, nil
}

func (uc *PickupUseCaseImpl) ListCollectorPickups(tenantID, collectorID uuid.UUID, query request.ListPickupQuery) ([]response.PickupResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	pickups, total, err := uc.pickupRepo.FindByCollector(tenantID, collectorID, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list collector pickups")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list pickups", err)
	}

	result := make([]response.PickupResponse, 0, len(pickups))
	for _, p := range pickups {
		result = append(result, *uc.toPickupResponse(&p))
	}

	return result, total, nil
}

func (uc *PickupUseCaseImpl) AssignCollector(tenantID, pickupID, collectorID uuid.UUID) (*response.PickupResponse, error) {
	pickup, err := uc.pickupRepo.FindByID(pickupID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Pickup request not found", errors.ErrNotFound)
	}

	if pickup.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied", errors.ErrForbidden)
	}

	if pickup.Status != constants.PickupStatusPending {
		return nil, errors.NewDomainError("INVALID_STATUS", "Can only assign collector to pending pickups", nil)
	}

	collector, err := uc.userRepo.FindByID(collectorID)
	if err != nil {
		return nil, errors.NewDomainError("COLLECTOR_NOT_FOUND", "Collector not found", err)
	}

	if collector.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Collector must be from the same region", errors.ErrForbidden)
	}

	if err := uc.pickupRepo.AssignCollector(pickupID, collectorID); err != nil {
		uc.log.WithError(err).Error("failed to assign collector")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to assign collector", err)
	}

	pickup, err = uc.pickupRepo.FindByID(pickupID)
	if err != nil {
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve updated pickup", err)
	}

	return uc.toPickupResponse(pickup), nil
}

func (uc *PickupUseCaseImpl) UpdatePickupStatus(tenantID, pickupID uuid.UUID, status string) (*response.PickupResponse, error) {
	pickup, err := uc.pickupRepo.FindByID(pickupID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Pickup request not found", errors.ErrNotFound)
	}

	if pickup.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied", errors.ErrForbidden)
	}

	if err := uc.validateStatusTransition(pickup.Status, status); err != nil {
		return nil, err
	}

	if err := uc.pickupRepo.UpdateStatus(pickupID, status); err != nil {
		uc.log.WithError(err).Error("failed to update pickup status")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to update status", err)
	}

	pickup.Status = status
	return uc.toPickupResponse(pickup), nil
}

func (uc *PickupUseCaseImpl) CancelPickup(tenantID, pickupID, userID uuid.UUID) error {
	pickup, err := uc.pickupRepo.FindByID(pickupID)
	if err != nil {
		return errors.NewDomainError("NOT_FOUND", "Pickup request not found", errors.ErrNotFound)
	}

	if pickup.TenantID != tenantID {
		return errors.NewDomainError("FORBIDDEN", "Access denied", errors.ErrForbidden)
	}

	if pickup.UserID != userID {
		return errors.NewDomainError("FORBIDDEN", "You can only cancel your own pickup requests", errors.ErrForbidden)
	}

	if pickup.Status != constants.PickupStatusPending {
		return errors.NewDomainError("INVALID_STATUS", "Can only cancel pending pickup requests", nil)
	}

	return uc.pickupRepo.CancelPickup(pickupID)
}

func (uc *PickupUseCaseImpl) validateStatusTransition(current, next string) error {
	transitions := map[string][]string{
		constants.PickupStatusPending:    {constants.PickupStatusAssigned, constants.PickupStatusCancelled},
		constants.PickupStatusAssigned:   {constants.PickupStatusInProgress, constants.PickupStatusCancelled},
		constants.PickupStatusInProgress: {constants.PickupStatusCompleted},
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

func (uc *PickupUseCaseImpl) toPickupResponse(p *models.PickupRequest) *response.PickupResponse {
	resp := &response.PickupResponse{
		ID:            p.ID,
		WasteType:     p.WasteType,
		Latitude:      p.Latitude,
		Longitude:     p.Longitude,
		Address:       p.Address,
		Status:        p.Status,
		ScheduledDate: p.ScheduledDate,
		Notes:         p.Notes,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.User != nil {
		resp.User = &response.UserRef{
			ID:       p.User.ID,
			FullName: p.User.FullName,
			Email:    p.User.Email,
			Phone:    p.User.Phone,
		}
	}

	if p.Collector != nil {
		resp.Collector = &response.UserRef{
			ID:       p.Collector.ID,
			FullName: p.Collector.FullName,
			Email:    p.Collector.Email,
			Phone:    p.Collector.Phone,
		}
	}

	return resp
}
