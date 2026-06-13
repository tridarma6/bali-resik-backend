package impl

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmailAndTenant(email string, tenantID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ? AND tenant_id = ?", email, tenantID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByTenant(tenantID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepositoryImpl) AssignRole(userID, roleID, tenantID uuid.UUID) error {
	ur := models.UserRole{
		UserID:   userID,
		RoleID:   roleID,
		TenantID: tenantID,
	}
	return r.db.Create(&ur).Error
}

func (r *UserRepositoryImpl) GetRoles(userID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Raw(`
		SELECT r.* FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`, userID).Scan(&roles).Error
	return roles, err
}
