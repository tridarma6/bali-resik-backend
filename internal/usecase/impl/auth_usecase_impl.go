package impl

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/auth"
	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
	"github.com/indim/bali-resik-backend/pkg/utils"
)

type AuthUseCaseImpl struct {
	userRepo         repo.UserRepository
	tenantRepo       repo.TenantRepository
	roleRepo         repo.RoleRepository
	refreshTokenRepo repo.RefreshTokenRepository
	jwtService       auth.JWTService
	log              *logrus.Logger
}

func NewAuthUseCase(
	userRepo repo.UserRepository,
	tenantRepo repo.TenantRepository,
	roleRepo repo.RoleRepository,
	refreshTokenRepo repo.RefreshTokenRepository,
	jwtService auth.JWTService,
	log *logrus.Logger,
) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		userRepo:         userRepo,
		tenantRepo:       tenantRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		log:              log,
	}
}

func (uc *AuthUseCaseImpl) Register(req *request.RegisterRequest) (*response.AuthResponse, error) {
	tenant, err := uc.tenantRepo.FindBySlug(req.TenantSlug)
	if err != nil {
		return nil, errors.NewDomainError("TENANT_NOT_FOUND", "Tenant not found", err)
	}

	existingUser, err := uc.userRepo.FindByEmailAndTenant(req.Email, tenant.ID)
	if err == nil && existingUser != nil {
		return nil, errors.NewDomainError("EMAIL_EXISTS", "Email already registered in this region", errors.ErrEmailAlreadyUsed)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		uc.log.WithError(err).Error("failed to hash password")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to process registration", err)
	}

	user := &models.User{
		TenantID:     tenant.ID,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(user); err != nil {
		uc.log.WithError(err).Error("failed to create user")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create user", err)
	}

	citizenRole, err := uc.roleRepo.FindByName("citizen")
	if err != nil {
		uc.log.WithError(err).Error("citizen role not found")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Default role not found", err)
	}

	if err := uc.userRepo.AssignRole(user.ID, citizenRole.ID, tenant.ID); err != nil {
		uc.log.WithError(err).Error("failed to assign citizen role")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to assign role", err)
	}

	return uc.generateAuthResponse(user, tenant.ID)
}

func (uc *AuthUseCaseImpl) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.NewDomainError("INVALID_CREDENTIALS", "Invalid email or password", errors.ErrInvalidCreds)
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.NewDomainError("INVALID_CREDENTIALS", "Invalid email or password", errors.ErrInvalidCreds)
	}

	if !user.IsActive {
		return nil, errors.NewDomainError("ACCOUNT_DISABLED", "Account has been disabled", nil)
	}

	return uc.generateAuthResponse(user, user.TenantID)
}

func (uc *AuthUseCaseImpl) RefreshToken(req *request.RefreshTokenRequest) (*response.AuthResponse, error) {
	storedToken, err := uc.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewDomainError("INVALID_TOKEN", "Invalid or expired refresh token", err)
	}

	if storedToken.IsExpired() {
		return nil, errors.NewDomainError("TOKEN_EXPIRED", "Refresh token has expired", errors.ErrTokenExpired)
	}

	if err := uc.refreshTokenRepo.RevokeByToken(storedToken.Token); err != nil {
		uc.log.WithError(err).Error("failed to revoke old refresh token")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to refresh token", err)
	}

	user, err := uc.userRepo.FindByID(storedToken.UserID)
	if err != nil {
		return nil, errors.NewDomainError("USER_NOT_FOUND", "User not found", err)
	}

	return uc.generateAuthResponse(user, user.TenantID)
}

func (uc *AuthUseCaseImpl) Logout(userID uuid.UUID, refreshToken string) error {
	if err := uc.refreshTokenRepo.RevokeByUserID(userID); err != nil {
		uc.log.WithError(err).Error("failed to revoke refresh tokens on logout")
		return errors.NewDomainError("INTERNAL_ERROR", "Failed to logout", err)
	}

	return nil
}

func (uc *AuthUseCaseImpl) generateAuthResponse(user *models.User, tenantID uuid.UUID) (*response.AuthResponse, error) {
	roles, err := uc.userRepo.GetRoles(user.ID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get user roles")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to get user roles", err)
	}

	roleNames := make([]string, 0, len(roles))
	primaryRole := "citizen"
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
		if role.Name == "super_admin" || role.Name == "admin_kabupaten" {
			primaryRole = role.Name
		}
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, tenantID, primaryRole, user.Email)
	if err != nil {
		uc.log.WithError(err).Error("failed to generate access token")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to generate token", err)
	}

	refreshTokenStr, err := generateRandomToken()
	if err != nil {
		uc.log.WithError(err).Error("failed to generate refresh token")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to generate refresh token", err)
	}

	refreshTokenTTL := time.Duration(10080) * time.Minute

	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}

	if err := uc.refreshTokenRepo.Create(rt); err != nil {
		uc.log.WithError(err).Error("failed to store refresh token")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to store refresh token", err)
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    15 * 60,
		User: response.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			TenantID: tenantID,
			Roles:    roleNames,
			IsActive: user.IsActive,
		},
	}, nil
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
