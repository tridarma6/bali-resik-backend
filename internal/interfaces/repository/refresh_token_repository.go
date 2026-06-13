package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	RevokeByUserID(userID uuid.UUID) error
	RevokeByToken(token string) error
	DeleteExpired() error
}
