package router

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/auth"
	"github.com/indim/bali-resik-backend/internal/handler"
	"github.com/indim/bali-resik-backend/internal/middleware"
)

type Router struct {
	e                 *echo.Echo
	log               *logrus.Logger
	jwtService        auth.JWTService
	authHandler       *handler.AuthHandler
	adminHandler      *handler.AdminHandler
	pickupHandler     *handler.PickupHandler
	reportHandler     *handler.WasteReportHandler
	rewardHandler     *handler.RewardHandler
	educationHandler  *handler.EducationHandler
	notifHandler      *handler.NotificationHandler
	analyticsHandler  *handler.AnalyticsHandler

	authMW        *middleware.AuthMiddleware
}

func New(
	e *echo.Echo,
	log *logrus.Logger,
	jwtService auth.JWTService,
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	pickupHandler *handler.PickupHandler,
	reportHandler *handler.WasteReportHandler,
	rewardHandler *handler.RewardHandler,
	educationHandler *handler.EducationHandler,
	notifHandler *handler.NotificationHandler,
	analyticsHandler *handler.AnalyticsHandler,
) *Router {
	return &Router{
		e:                e,
		log:              log,
		jwtService:       jwtService,
		authHandler:      authHandler,
		adminHandler:     adminHandler,
		pickupHandler:    pickupHandler,
		reportHandler:    reportHandler,
		rewardHandler:    rewardHandler,
		educationHandler: educationHandler,
		notifHandler:     notifHandler,
		analyticsHandler: analyticsHandler,
		authMW:           middleware.NewAuthMiddleware(jwtService, log),
	}
}

func (r *Router) Setup() {
	r.e.Use(echomw.RequestID())
	r.e.Use(echomw.Secure())
	r.e.Use(middleware.CORSConfig())
	r.e.Use(echomw.Recover())
	r.e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))

	r.e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status": "ok",
		})
	})

	r.setupAuthRoutes()
	r.setupAdminRoutes()
	r.setupPickupRoutes()
	r.setupWasteReportRoutes()
	r.setupRewardRoutes()
	r.setupEducationRoutes()
	r.setupNotificationRoutes()
	r.setupAnalyticsRoutes()
}

func (r *Router) setupAuthRoutes() {
	api := r.e.Group("/api/v1")

	api.POST("/auth/register", r.authHandler.Register)
	api.POST("/auth/login", r.authHandler.Login)
	api.POST("/auth/refresh", r.authHandler.RefreshToken)

	authGroup := api.Group("")
	authGroup.Use(r.authMW.Authenticate)
	authGroup.POST("/auth/logout", r.authHandler.Logout)
}

func (r *Router) setupAdminRoutes() {
	admin := r.e.Group("/api/v1/admin")
	admin.Use(r.authMW.Authenticate)
	admin.Use(middleware.RequireRoles("super_admin"))

	admin.POST("/tenants", r.adminHandler.CreateTenant)
	admin.GET("/tenants", r.adminHandler.ListTenants)
	admin.POST("/admins", r.adminHandler.CreateAdmin)
}

func (r *Router) setupPickupRoutes() {
	pickups := r.e.Group("/api/v1/pickups")
	pickups.Use(r.authMW.Authenticate)

	pickups.POST("", r.pickupHandler.Create, middleware.RequireRoles("citizen"))
	pickups.GET("/mine", r.pickupHandler.ListMine, middleware.RequireRoles("citizen"))
	pickups.GET("/assigned", r.pickupHandler.ListAssigned, middleware.RequireRoles("collector"))
	pickups.GET("", r.pickupHandler.List, middleware.RequireRoles("admin_kabupaten", "super_admin"))
	pickups.GET("/:id", r.pickupHandler.GetByID)
	pickups.PUT("/:id/assign", r.pickupHandler.AssignCollector, middleware.RequireRoles("admin_kabupaten", "super_admin"))
	pickups.PUT("/:id/status", r.pickupHandler.UpdateStatus, middleware.RequireRoles("collector", "admin_kabupaten"))
	pickups.DELETE("/:id", r.pickupHandler.Cancel)
}

func (r *Router) setupWasteReportRoutes() {
	reports := r.e.Group("/api/v1/reports")
	reports.Use(r.authMW.Authenticate)

	reports.POST("", r.reportHandler.Create, middleware.RequireRoles("citizen"))
	reports.GET("/mine", r.reportHandler.ListMine, middleware.RequireRoles("citizen"))
	reports.GET("/nearby", r.reportHandler.FindNearby)
	reports.GET("", r.reportHandler.List, middleware.RequireRoles("admin_kabupaten", "super_admin"))
	reports.GET("/:id", r.reportHandler.GetByID)
	reports.PUT("/:id/status", r.reportHandler.UpdateStatus, middleware.RequireRoles("admin_kabupaten"))
	reports.POST("/:id/images", r.reportHandler.UploadImage)
}

func (r *Router) setupRewardRoutes() {
	rewards := r.e.Group("/api/v1/rewards")
	rewards.Use(r.authMW.Authenticate)

	rewards.GET("", r.rewardHandler.List)
	rewards.GET("/points", r.rewardHandler.GetPoints)
	rewards.GET("/transactions", r.rewardHandler.GetTransactions)
	rewards.POST("/redeem", r.rewardHandler.Redeem, middleware.RequireRoles("citizen"))

	adminRewards := r.e.Group("/api/v1/rewards")
	adminRewards.Use(r.authMW.Authenticate)
	adminRewards.Use(middleware.RequireRoles("admin_kabupaten", "super_admin"))
	adminRewards.POST("", r.rewardHandler.Create)
	adminRewards.PUT("/:id", r.rewardHandler.Update)
	adminRewards.DELETE("/:id", r.rewardHandler.Delete)
}

func (r *Router) setupEducationRoutes() {
	education := r.e.Group("/api/v1/education")
	education.Use(r.authMW.Authenticate)

	education.GET("", r.educationHandler.List)
	education.GET("/:id", r.educationHandler.GetByID)

	adminEducation := r.e.Group("/api/v1/education")
	adminEducation.Use(r.authMW.Authenticate)
	adminEducation.Use(middleware.RequireRoles("admin_kabupaten", "super_admin"))
	adminEducation.POST("", r.educationHandler.Create)
	adminEducation.PUT("/:id", r.educationHandler.Update)
	adminEducation.DELETE("/:id", r.educationHandler.Delete)
}

func (r *Router) setupNotificationRoutes() {
	notifs := r.e.Group("/api/v1/notifications")
	notifs.Use(r.authMW.Authenticate)

	notifs.GET("", r.notifHandler.List)
	notifs.GET("/unread-count", r.notifHandler.GetUnreadCount)
	notifs.PUT("/:id/read", r.notifHandler.MarkAsRead)
	notifs.PUT("/read-all", r.notifHandler.MarkAllAsRead)
	notifs.DELETE("/:id", r.notifHandler.Delete)
}

func (r *Router) setupAnalyticsRoutes() {
	analytics := r.e.Group("/api/v1/analytics")
	analytics.Use(r.authMW.Authenticate)

	analytics.GET("/dashboard", r.analyticsHandler.GetDashboard, middleware.RequireRoles("admin_kabupaten", "super_admin"))

	r.e.GET("/api/v1/admin/regional-stats", r.analyticsHandler.GetRegionalStats,
		r.authMW.Authenticate,
		middleware.RequireRoles("super_admin"),
	)
}
