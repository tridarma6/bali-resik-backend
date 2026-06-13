package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type AuthUseCase interface {
	Register(req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
	RefreshToken(req *request.RefreshTokenRequest) (*response.AuthResponse, error)
	Logout(userID uuid.UUID, refreshToken string) error
}
