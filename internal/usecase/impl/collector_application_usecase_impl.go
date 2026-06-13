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

type CollectorApplicationUseCaseImpl struct {
	appRepo repo.CollectorApplicationRepository
	userRepo repo.UserRepository
	roleRepo repo.RoleRepository
	log      *logrus.Logger
}

func NewCollectorApplicationUseCase(
	appRepo repo.CollectorApplicationRepository,
	userRepo repo.UserRepository,
	roleRepo repo.RoleRepository,
	log *logrus.Logger,
) *CollectorApplicationUseCaseImpl {
	return &CollectorApplicationUseCaseImpl{
		appRepo:  appRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
		log:      log,
	}
}

func (uc *CollectorApplicationUseCaseImpl) Submit(tenantID, userID uuid.UUID) (*response.CollectorApplicationResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "User not found", err)
	}

	roles, err := uc.userRepo.GetRoles(user.ID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get user roles")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to verify user roles", err)
	}
	for _, r := range roles {
		if r.Name == constants.RoleCollector {
			return nil, errors.NewDomainError("ALREADY_COLLECTOR", "User is already a collector", errors.ErrDuplicateEntry)
		}
	}

	existing, _ := uc.appRepo.FindPendingByUser(tenantID, userID)
	if existing != nil {
		return nil, errors.NewDomainError("DUPLICATE_APPLICATION", "You already have a pending application", errors.ErrDuplicateEntry)
	}

	app := &models.CollectorApplication{
		TenantID: tenantID,
		UserID:   userID,
		Status:   constants.AppStatusPending,
	}

	if err := uc.appRepo.Create(app); err != nil {
		uc.log.WithError(err).Error("failed to create collector application")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to submit application", err)
	}

	return response.ToCollectorApplicationResponse(app), nil
}

func (uc *CollectorApplicationUseCaseImpl) ListMyApplications(tenantID, userID uuid.UUID, query request.ListCollectorAppQuery) ([]response.CollectorApplicationResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	apps, total, err := uc.appRepo.FindByUser(tenantID, userID, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list collector applications")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list applications", err)
	}

	result := make([]response.CollectorApplicationResponse, 0, len(apps))
	for _, app := range apps {
		result = append(result, *response.ToCollectorApplicationResponse(&app))
	}

	return result, total, nil
}

func (uc *CollectorApplicationUseCaseImpl) ListAll(tenantID uuid.UUID, query request.ListCollectorAppQuery) ([]response.CollectorApplicationResponse, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 || query.PerPage > 100 {
		query.PerPage = 20
	}

	apps, total, err := uc.appRepo.FindByTenant(tenantID, query.Status, query.Page, query.PerPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to list collector applications")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to list applications", err)
	}

	result := make([]response.CollectorApplicationResponse, 0, len(apps))
	for _, app := range apps {
		result = append(result, *response.ToCollectorApplicationResponse(&app))
	}

	return result, total, nil
}

func (uc *CollectorApplicationUseCaseImpl) Approve(adminID, appID uuid.UUID) (*response.CollectorApplicationResponse, error) {
	app, err := uc.appRepo.FindByID(appID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Application not found", err)
	}

	if app.Status != constants.AppStatusPending {
		return nil, errors.NewDomainError("INVALID_STATUS", "Application is not in pending status", errors.ErrInvalidInput)
	}

	collectorRole, err := uc.roleRepo.FindByName(constants.RoleCollector)
	if err != nil {
		uc.log.WithError(err).Error("collector role not found")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Default role not found", err)
	}

	if err := uc.userRepo.AssignRole(app.UserID, collectorRole.ID, app.TenantID); err != nil {
		uc.log.WithError(err).Error("failed to assign collector role")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to assign collector role", err)
	}

	now := time.Now()
	app.Status = constants.AppStatusApproved
	app.ReviewedBy = &adminID
	app.ReviewedAt = &now

	if err := uc.appRepo.Update(app); err != nil {
		uc.log.WithError(err).Error("failed to update application")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to approve application", err)
	}

	updated, err := uc.appRepo.FindByID(app.ID)
	if err != nil {
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve updated application", err)
	}

	return response.ToCollectorApplicationResponse(updated), nil
}

func (uc *CollectorApplicationUseCaseImpl) Reject(adminID, appID uuid.UUID, notes string) (*response.CollectorApplicationResponse, error) {
	app, err := uc.appRepo.FindByID(appID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Application not found", err)
	}

	if app.Status != constants.AppStatusPending {
		return nil, errors.NewDomainError("INVALID_STATUS", "Application is not in pending status", errors.ErrInvalidInput)
	}

	now := time.Now()
	app.Status = constants.AppStatusRejected
	app.AdminNotes = notes
	app.ReviewedBy = &adminID
	app.ReviewedAt = &now

	if err := uc.appRepo.Update(app); err != nil {
		uc.log.WithError(err).Error("failed to update application")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to reject application", err)
	}

	return response.ToCollectorApplicationResponse(app), nil
}
