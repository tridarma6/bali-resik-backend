package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
	"github.com/indim/bali-resik-backend/pkg/utils"
)

type AdminUseCaseImpl struct {
	tenantRepo repo.TenantRepository
	userRepo   repo.UserRepository
	roleRepo   repo.RoleRepository
	log        *logrus.Logger
}

func NewAdminUseCase(
	tenantRepo repo.TenantRepository,
	userRepo repo.UserRepository,
	roleRepo repo.RoleRepository,
	log *logrus.Logger,
) *AdminUseCaseImpl {
	return &AdminUseCaseImpl{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		log:        log,
	}
}

func (uc *AdminUseCaseImpl) CreateTenant(req *request.CreateTenantRequest) (*response.TenantResponse, error) {
	existing, _ := uc.tenantRepo.FindBySlug(req.Slug)
	if existing != nil {
		return nil, errors.NewDomainError("SLUG_EXISTS", "Tenant slug already exists", errors.ErrDuplicateEntry)
	}

	tenant := &models.Tenant{
		Name:       req.Name,
		Slug:       req.Slug,
		RegionType: req.RegionType,
		IsActive:   true,
	}

	if err := uc.tenantRepo.Create(tenant); err != nil {
		uc.log.WithError(err).Error("failed to create tenant")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create tenant", err)
	}

	return &response.TenantResponse{
		ID:         tenant.ID,
		Name:       tenant.Name,
		Slug:       tenant.Slug,
		RegionType: tenant.RegionType,
		IsActive:   tenant.IsActive,
		CreatedAt:  tenant.CreatedAt,
	}, nil
}

func (uc *AdminUseCaseImpl) ListTenants() ([]response.TenantResponse, error) {
	tenants, err := uc.tenantRepo.FindAll()
	if err != nil {
		uc.log.WithError(err).Error("failed to list tenants")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to list tenants", err)
	}

	result := make([]response.TenantResponse, 0, len(tenants))
	for _, t := range tenants {
		result = append(result, response.TenantResponse{
			ID:         t.ID,
			Name:       t.Name,
			Slug:       t.Slug,
			RegionType: t.RegionType,
			IsActive:   t.IsActive,
			CreatedAt:  t.CreatedAt,
		})
	}

	return result, nil
}

func (uc *AdminUseCaseImpl) CreateAdmin(req *request.CreateAdminRequest) (*response.UserResponse, error) {
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return nil, errors.NewDomainError("INVALID_TENANT_ID", "Invalid tenant ID", err)
	}

	tenant, err := uc.tenantRepo.FindByID(tenantID)
	if err != nil {
		return nil, errors.NewDomainError("TENANT_NOT_FOUND", "Tenant not found", err)
	}

	existingUser, _ := uc.userRepo.FindByEmailAndTenant(req.Email, tenant.ID)
	if existingUser != nil {
		return nil, errors.NewDomainError("EMAIL_EXISTS", "Email already registered in this tenant", errors.ErrEmailAlreadyUsed)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		uc.log.WithError(err).Error("failed to hash password")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to process request", err)
	}

	user := &models.User{
		TenantID:     tenant.ID,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(user); err != nil {
		uc.log.WithError(err).Error("failed to create admin user")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create user", err)
	}

	adminRole, err := uc.roleRepo.FindByName("admin_kabupaten")
	if err != nil {
		uc.log.WithError(err).Error("admin_kabupaten role not found")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Default role not found", err)
	}

	if err := uc.userRepo.AssignRole(user.ID, adminRole.ID, tenant.ID); err != nil {
		uc.log.WithError(err).Error("failed to assign admin role")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to assign role", err)
	}

	roles, _ := uc.userRepo.GetRoles(user.ID)
	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleNames = append(roleNames, r.Name)
	}

	return &response.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		TenantID: tenant.ID,
		Roles:    roleNames,
		IsActive: user.IsActive,
	}, nil
}
