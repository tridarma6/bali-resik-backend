package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type UserUseCase interface {
	GetProfile(tenantID, userID uuid.UUID) (*response.UserResponse, error)
	UpdateProfile(tenantID, userID uuid.UUID, req *request.UpdateProfileRequest) (*response.UserResponse, error)
	ChangePassword(tenantID, userID uuid.UUID, req *request.ChangePasswordRequest) error
	UpdateAvatar(tenantID, userID uuid.UUID, avatarURL string) (*response.UserResponse, error)
}
