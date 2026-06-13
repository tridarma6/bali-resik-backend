package repository

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type RewardRepository interface {
	Create(reward *models.Reward) error
	FindByID(id uuid.UUID) (*models.Reward, error)
	FindByTenant(tenantID uuid.UUID) ([]models.Reward, error)
	Update(reward *models.Reward) error
	Delete(id uuid.UUID) error
}

type RewardTransactionRepository interface {
	Create(tx *models.RewardTransaction) error
	FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.RewardTransaction, int64, error)
	GetUserPoints(tenantID, userID uuid.UUID) (int, error)
}
