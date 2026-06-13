package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/indim/bali-resik-backend/internal/auth"
	"github.com/indim/bali-resik-backend/internal/config"
	"github.com/indim/bali-resik-backend/internal/database"
	"github.com/indim/bali-resik-backend/internal/logger"
	"github.com/indim/bali-resik-backend/internal/router"
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

	e := echo.New()

	jwtService := auth.NewJWTService(&cfg.JWT)

	r := router.New(e, log, jwtService)
	r.Setup()

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.WithField("address", addr).Info("server starting")

	if err := e.Start(addr); err != nil {
		log.WithError(err).Fatal("server failed to start")
	}
}
