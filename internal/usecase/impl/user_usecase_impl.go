package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
	pkg "github.com/indim/bali-resik-backend/pkg/utils"
)

type UserUseCaseImpl struct {
	userRepo repo.UserRepository
	log      *logrus.Logger
}

func NewUserUseCase(userRepo repo.UserRepository, log *logrus.Logger) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepo: userRepo,
		log:      log,
	}
}

func (uc *UserUseCaseImpl) GetProfile(tenantID, userID uuid.UUID) (*response.UserResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		uc.log.WithError(err).Error("user not found")
		return nil, errors.NewDomainError("NOT_FOUND", "User not found", err)
	}

	if user.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied", nil)
	}

	roles, err := uc.userRepo.GetRoles(user.ID)
	if err != nil {
		uc.log.WithError(err).Warn("failed to get user roles")
		roles = nil
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	resp := &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		TenantID:  tenantID,
		Roles:     roleNames,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return resp, nil
}

func (uc *UserUseCaseImpl) UpdateProfile(tenantID, userID uuid.UUID, req *request.UpdateProfileRequest) (*response.UserResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		uc.log.WithError(err).Error("user not found")
		return nil, errors.NewDomainError("NOT_FOUND", "User not found", err)
	}

	if user.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied", nil)
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := uc.userRepo.Update(user); err != nil {
		uc.log.WithError(err).Error("failed to update user")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to update profile", err)
	}

	roles, err := uc.userRepo.GetRoles(user.ID)
	if err != nil {
		uc.log.WithError(err).Warn("failed to get user roles")
		roles = nil
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	resp := &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		TenantID:  tenantID,
		Roles:     roleNames,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return resp, nil
}

func (uc *UserUseCaseImpl) ChangePassword(tenantID, userID uuid.UUID, req *request.ChangePasswordRequest) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		uc.log.WithError(err).Error("user not found")
		return errors.NewDomainError("NOT_FOUND", "User not found", err)
	}

	if user.TenantID != tenantID {
		return errors.NewDomainError("FORBIDDEN", "Access denied", nil)
	}

	if !pkg.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return errors.NewDomainError("INVALID_CREDENTIALS", "Current password is incorrect", nil)
	}

	hashedPassword, err := pkg.HashPassword(req.NewPassword)
	if err != nil {
		uc.log.WithError(err).Error("failed to hash password")
		return errors.NewDomainError("INTERNAL_ERROR", "Failed to change password", err)
	}

	user.PasswordHash = hashedPassword
	if err := uc.userRepo.Update(user); err != nil {
		uc.log.WithError(err).Error("failed to update password")
		return errors.NewDomainError("INTERNAL_ERROR", "Failed to change password", err)
	}

	return nil
}

func (uc *UserUseCaseImpl) UpdateAvatar(tenantID, userID uuid.UUID, avatarURL string) (*response.UserResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		uc.log.WithError(err).Error("user not found")
		return nil, errors.NewDomainError("NOT_FOUND", "User not found", err)
	}

	if user.TenantID != tenantID {
		return nil, errors.NewDomainError("FORBIDDEN", "Access denied", nil)
	}

	user.AvatarURL = avatarURL
	if err := uc.userRepo.Update(user); err != nil {
		uc.log.WithError(err).Error("failed to update avatar")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to update avatar", err)
	}

	roles, err := uc.userRepo.GetRoles(user.ID)
	if err != nil {
		uc.log.WithError(err).Warn("failed to get user roles")
		roles = nil
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	resp := &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		TenantID:  tenantID,
		Roles:     roleNames,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return resp, nil
}
