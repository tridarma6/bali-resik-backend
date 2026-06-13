package usecase

import (
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type AdminUseCase interface {
	CreateTenant(req *request.CreateTenantRequest) (*response.TenantResponse, error)
	ListTenants() ([]response.TenantResponse, error)
	CreateAdmin(req *request.CreateAdminRequest) (*response.UserResponse, error)
}
