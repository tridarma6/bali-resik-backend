package seed

import (
	"github.com/sirupsen/logrus"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/pkg/utils"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, log *logrus.Logger) {
	log.Info("seeding database...")

	seedRoles(db)
	seedTenants(db)
	seedSuperAdmin(db)

	log.Info("database seeding completed")
}

func seedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "super_admin", Description: "Super Administrator with access to all tenants"},
		{Name: "admin_kabupaten", Description: "Regional administrator for a specific Kabupaten/Kota"},
		{Name: "citizen", Description: "Regular citizen user"},
		{Name: "collector", Description: "Waste collection partner"},
	}

	for _, role := range roles {
		var existing models.Role
		result := db.Where("name = ?", role.Name).First(&existing)
		if result.Error != nil {
			db.Create(&role)
		}
	}
}

func seedTenants(db *gorm.DB) {
	tenants := []models.Tenant{
		{Name: "Kota Denpasar", Slug: "kota-denpasar", RegionType: "kota"},
		{Name: "Kabupaten Badung", Slug: "kabupaten-badung", RegionType: "kabupaten"},
		{Name: "Kabupaten Gianyar", Slug: "kabupaten-gianyar", RegionType: "kabupaten"},
		{Name: "Kabupaten Tabanan", Slug: "kabupaten-tabanan", RegionType: "kabupaten"},
		{Name: "Kabupaten Buleleng", Slug: "kabupaten-buleleng", RegionType: "kabupaten"},
		{Name: "Kabupaten Jembrana", Slug: "kabupaten-jembrana", RegionType: "kabupaten"},
		{Name: "Kabupaten Karangasem", Slug: "kabupaten-karangasem", RegionType: "kabupaten"},
		{Name: "Kabupaten Bangli", Slug: "kabupaten-bangli", RegionType: "kabupaten"},
		{Name: "Kabupaten Klungkung", Slug: "kabupaten-klungkung", RegionType: "kabupaten"},
	}

	for _, tenant := range tenants {
		var existing models.Tenant
		result := db.Where("slug = ?", tenant.Slug).First(&existing)
		if result.Error != nil {
			db.Create(&tenant)
		}
	}
}

func seedSuperAdmin(db *gorm.DB) {
	var existing models.User
	result := db.Where("email = ?", "superadmin@baliresik.go.id").First(&existing)
	if result.Error == nil {
		return
	}

	hashedPassword, err := utils.HashPassword("SuperAdmin123!")
	if err != nil {
		return
	}

	var kotaDenpasar models.Tenant
	db.Where("slug = ?", "kota-denpasar").First(&kotaDenpasar)

	admin := models.User{
		TenantID:     kotaDenpasar.ID,
		Email:        "superadmin@baliresik.go.id",
		PasswordHash: hashedPassword,
		FullName:     "Super Admin",
		IsActive:     true,
	}

	admin.ID = utils.NewUUID()

	if err := db.Create(&admin).Error; err != nil {
		return
	}

	var superAdminRole models.Role
	db.Where("name = ?", "super_admin").First(&superAdminRole)

	db.Create(&models.UserRole{
		UserID:   admin.ID,
		RoleID:   superAdminRole.ID,
		TenantID: kotaDenpasar.ID,
	})
}
