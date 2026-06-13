package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/indim/bali-resik-backend/internal/auth"
	"github.com/indim/bali-resik-backend/internal/config"
	"github.com/indim/bali-resik-backend/internal/database"
	"github.com/indim/bali-resik-backend/internal/database/seed"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/handler"
	"github.com/indim/bali-resik-backend/internal/logger"
	repoimpl "github.com/indim/bali-resik-backend/internal/repository/impl"
	"github.com/indim/bali-resik-backend/internal/router"
	ucaseimpl "github.com/indim/bali-resik-backend/internal/usecase/impl"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger(&cfg.Logger)
	log.Info("starting bali-resik-backend server")

	db, err := database.NewPostgresDB(&cfg.Database, log)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	runMigrations(db, log)
	seed.Run(db, log)

	jwtService := auth.NewJWTService(&cfg.JWT)

	tenantRepo := repoimpl.NewTenantRepository(db)
	userRepo := repoimpl.NewUserRepository(db)
	roleRepo := repoimpl.NewRoleRepository(db)
	refreshTokenRepo := repoimpl.NewRefreshTokenRepository(db)
	pickupRepo := repoimpl.NewPickupRepository(db)
	reportRepo := repoimpl.NewWasteReportRepository(db)
	rewardRepo := repoimpl.NewRewardRepository(db)
	txRepo := repoimpl.NewRewardTransactionRepository(db)
	eduRepo := repoimpl.NewEducationRepository(db)

	authUseCase := ucaseimpl.NewAuthUseCase(userRepo, tenantRepo, roleRepo, refreshTokenRepo, jwtService, log)
	adminUseCase := ucaseimpl.NewAdminUseCase(tenantRepo, userRepo, roleRepo, log)
	pickupUseCase := ucaseimpl.NewPickupUseCase(pickupRepo, userRepo, log)
	reportUseCase := ucaseimpl.NewWasteReportUseCase(reportRepo, userRepo, log)
	rewardUseCase := ucaseimpl.NewRewardUseCase(rewardRepo, txRepo, log)
	eduUseCase := ucaseimpl.NewEducationUseCase(eduRepo, log)

	authHandler := handler.NewAuthHandler(authUseCase, log)
	adminHandler := handler.NewAdminHandler(adminUseCase, log)
	pickupHandler := handler.NewPickupHandler(pickupUseCase, log)
	reportHandler := handler.NewWasteReportHandler(reportUseCase, log)
	rewardHandler := handler.NewRewardHandler(rewardUseCase, log)
	educationHandler := handler.NewEducationHandler(eduUseCase, log)

	e := echo.New()

	r := router.New(e, log, jwtService, authHandler, adminHandler, pickupHandler, reportHandler, rewardHandler, educationHandler)
	r.Setup()

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.WithField("address", addr).Info("server starting")

	if err := e.Start(addr); err != nil {
		log.WithError(err).Fatal("server failed to start")
	}
}

func runMigrations(db *gorm.DB, log *logrus.Logger) {
	log.Info("running database migrations")

	if err := db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.RefreshToken{},
		&models.PickupRequest{},
		&models.WasteReport{},
		&models.ReportImage{},
		&models.Reward{},
		&models.RewardTransaction{},
		&models.EducationalContent{},
	); err != nil {
		log.WithError(err).Fatal("failed to run database migrations")
	}

	log.Info("database migrations completed")
}