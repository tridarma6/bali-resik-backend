package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type RewardUseCase interface {
	CreateReward(tenantID uuid.UUID, req *request.CreateRewardRequest) (*response.RewardResponse, error)
	ListRewards(tenantID uuid.UUID) ([]response.RewardResponse, error)
	UpdateReward(rewardID uuid.UUID, req *request.UpdateRewardRequest) (*response.RewardResponse, error)
	DeleteReward(rewardID uuid.UUID) error

	RedeemReward(tenantID, userID uuid.UUID, rewardID uuid.UUID) (*response.RewardTransactionResponse, error)
	GetTransactionHistory(tenantID, userID uuid.UUID, page, perPage int) ([]response.RewardTransactionResponse, int64, error)
	GetUserPoints(tenantID, userID uuid.UUID) (*response.UserPointsResponse, error)
}
