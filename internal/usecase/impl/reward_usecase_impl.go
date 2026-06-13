package impl

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/constants"
	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
)

type RewardUseCaseImpl struct {
	rewardRepo repo.RewardRepository
	txRepo     repo.RewardTransactionRepository
	log        *logrus.Logger
}

func NewRewardUseCase(
	rewardRepo repo.RewardRepository,
	txRepo repo.RewardTransactionRepository,
	log *logrus.Logger,
) *RewardUseCaseImpl {
	return &RewardUseCaseImpl{
		rewardRepo: rewardRepo,
		txRepo:     txRepo,
		log:        log,
	}
}

func (uc *RewardUseCaseImpl) CreateReward(tenantID uuid.UUID, req *request.CreateRewardRequest) (*response.RewardResponse, error) {
	reward := &models.Reward{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		PointsCost:  req.PointsCost,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := uc.rewardRepo.Create(reward); err != nil {
		uc.log.WithError(err).Error("failed to create reward")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create reward", err)
	}

	return uc.toRewardResponse(reward), nil
}

func (uc *RewardUseCaseImpl) ListRewards(tenantID uuid.UUID) ([]response.RewardResponse, error) {
	rewards, err := uc.rewardRepo.FindByTenant(tenantID)
	if err != nil {
		uc.log.WithError(err).Error("failed to list rewards")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to list rewards", err)
	}

	result := make([]response.RewardResponse, 0, len(rewards))
	for _, r := range rewards {
		result = append(result, *uc.toRewardResponse(&r))
	}

	return result, nil
}

func (uc *RewardUseCaseImpl) UpdateReward(rewardID uuid.UUID, req *request.UpdateRewardRequest) (*response.RewardResponse, error) {
	reward, err := uc.rewardRepo.FindByID(rewardID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Reward not found", errors.ErrNotFound)
	}

	if req.Name != "" {
		reward.Name = req.Name
	}
	if req.Description != "" {
		reward.Description = req.Description
	}
	if req.PointsCost > 0 {
		reward.PointsCost = req.PointsCost
	}
	if req.Stock != nil {
		reward.Stock = req.Stock
	}
	if req.ImageURL != "" {
		reward.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		reward.IsActive = *req.IsActive
	}

	if err := uc.rewardRepo.Update(reward); err != nil {
		uc.log.WithError(err).Error("failed to update reward")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to update reward", err)
	}

	return uc.toRewardResponse(reward), nil
}

func (uc *RewardUseCaseImpl) DeleteReward(rewardID uuid.UUID) error {
	if err := uc.rewardRepo.Delete(rewardID); err != nil {
		return errors.NewDomainError("NOT_FOUND", "Reward not found", errors.ErrNotFound)
	}
	return nil
}

func (uc *RewardUseCaseImpl) RedeemReward(tenantID, userID uuid.UUID, rewardID uuid.UUID) (*response.RewardTransactionResponse, error) {
	reward, err := uc.rewardRepo.FindByID(rewardID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Reward not found", errors.ErrNotFound)
	}

	if !reward.IsActive {
		return nil, errors.NewDomainError("REWARD_INACTIVE", "This reward is no longer available", nil)
	}

	if reward.Stock != nil && *reward.Stock <= 0 {
		return nil, errors.NewDomainError("OUT_OF_STOCK", "This reward is out of stock", errors.ErrInsufficientStock)
	}

	currentPoints, err := uc.txRepo.GetUserPoints(tenantID, userID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get user points")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to process redemption", err)
	}

	if currentPoints < reward.PointsCost {
		return nil, errors.NewDomainError("INSUFFICIENT_POINTS",
			"Insufficient points. You need "+fmt.Sprintf("%d", reward.PointsCost)+" points",
			errors.ErrInsufficientPoints)
	}

	tx := &models.RewardTransaction{
		TenantID:    tenantID,
		UserID:      userID,
		Points:      reward.PointsCost,
		Type:        constants.RewardTypeRedeem,
		ReferenceID: &reward.ID,
		Description: "Redeemed: " + reward.Name,
	}

	if err := uc.txRepo.Create(tx); err != nil {
		uc.log.WithError(err).Error("failed to create redemption transaction")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to process redemption", err)
	}

	if reward.Stock != nil {
		newStock := *reward.Stock - 1
		reward.Stock = &newStock
		uc.rewardRepo.Update(reward)
	}

	return uc.toTxResponse(tx), nil
}

func (uc *RewardUseCaseImpl) GetTransactionHistory(tenantID, userID uuid.UUID, page, perPage int) ([]response.RewardTransactionResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	transactions, total, err := uc.txRepo.FindByUser(tenantID, userID, page, perPage)
	if err != nil {
		uc.log.WithError(err).Error("failed to get transaction history")
		return nil, 0, errors.NewDomainError("INTERNAL_ERROR", "Failed to get transaction history", err)
	}

	result := make([]response.RewardTransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		result = append(result, *uc.toTxResponse(&t))
	}

	return result, total, nil
}

func (uc *RewardUseCaseImpl) GetUserPoints(tenantID, userID uuid.UUID) (*response.UserPointsResponse, error) {
	points, err := uc.txRepo.GetUserPoints(tenantID, userID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get user points")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to get user points", err)
	}

	return &response.UserPointsResponse{TotalPoints: points}, nil
}

func (uc *RewardUseCaseImpl) toRewardResponse(r *models.Reward) *response.RewardResponse {
	return &response.RewardResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		PointsCost:  r.PointsCost,
		Stock:       r.Stock,
		ImageURL:    r.ImageURL,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
	}
}

func (uc *RewardUseCaseImpl) toTxResponse(t *models.RewardTransaction) *response.RewardTransactionResponse {
	return &response.RewardTransactionResponse{
		ID:          t.ID,
		Points:      t.Points,
		Type:        t.Type,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
	}
}
