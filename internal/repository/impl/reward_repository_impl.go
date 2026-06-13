package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/helper"
	"gorm.io/gorm"
)

type RewardRepositoryImpl struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) *RewardRepositoryImpl {
	return &RewardRepositoryImpl{db: db}
}

func (r *RewardRepositoryImpl) Create(reward *models.Reward) error {
	return r.db.Create(reward).Error
}

func (r *RewardRepositoryImpl) FindByID(id uuid.UUID) (*models.Reward, error) {
	var reward models.Reward
	err := r.db.First(&reward, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

func (r *RewardRepositoryImpl) FindByTenant(tenantID uuid.UUID) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("points_cost ASC").
		Find(&rewards).Error
	return rewards, err
}

func (r *RewardRepositoryImpl) Update(reward *models.Reward) error {
	return r.db.Save(reward).Error
}

func (r *RewardRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Reward{}, "id = ?", id).Error
}

type RewardTransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewRewardTransactionRepository(db *gorm.DB) *RewardTransactionRepositoryImpl {
	return &RewardTransactionRepositoryImpl{db: db}
}

func (r *RewardTransactionRepositoryImpl) Create(tx *models.RewardTransaction) error {
	return r.db.Create(tx).Error
}

func (r *RewardTransactionRepositoryImpl) FindByUser(tenantID, userID uuid.UUID, page, perPage int) ([]models.RewardTransaction, int64, error) {
	var transactions []models.RewardTransaction
	query := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	param := helper.PaginationParam{Page: page, PerPage: perPage}
	paginatedQuery, meta := helper.Paginate(query, param)

	err := paginatedQuery.Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, meta.Total, nil
}

func (r *RewardTransactionRepositoryImpl) GetUserPoints(tenantID, userID uuid.UUID) (int, error) {
	type result struct {
		Total int
	}
	var res result
	err := r.db.Model(&models.RewardTransaction{}).
		Select("COALESCE(SUM(CASE WHEN type = 'earn' THEN points ELSE 0 END) - SUM(CASE WHEN type = 'redeem' THEN points ELSE 0 END), 0) as total").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Scan(&res).Error

	return res.Total, err
}
