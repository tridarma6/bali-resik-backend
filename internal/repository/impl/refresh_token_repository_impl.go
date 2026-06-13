package impl

import (
	"time"

	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"gorm.io/gorm"
)

type RefreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{db: db}
}

func (r *RefreshTokenRepositoryImpl) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepositoryImpl) FindByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("token = ? AND revoked_at IS NULL", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepositoryImpl) RevokeByUserID(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}

func (r *RefreshTokenRepositoryImpl) RevokeByToken(token string) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked_at", &now).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}
