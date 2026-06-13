package router

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/auth"
	"github.com/indim/bali-resik-backend/internal/middleware"
)

type Router struct {
	e          *echo.Echo
	log        *logrus.Logger
	jwtService auth.JWTService
	authMW     *middleware.AuthMiddleware
}

func New(
	e *echo.Echo,
	log *logrus.Logger,
	jwtService auth.JWTService,
) *Router {
	return &Router{
		e:          e,
		log:        log,
		jwtService: jwtService,
		authMW:     middleware.NewAuthMiddleware(jwtService, log),
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

	r.SetupAuthRoutes()
	r.SetupAdminRoutes()
}

func (r *Router) SetupAuthRoutes() {
	api := r.e.Group("/api/v1")

	api.POST("/auth/login", func(c echo.Context) error {
		return c.JSON(501, map[string]string{"message": "Not implemented yet"})
	})
	api.POST("/auth/register", func(c echo.Context) error {
		return c.JSON(501, map[string]string{"message": "Not implemented yet"})
	})
	api.POST("/auth/refresh", func(c echo.Context) error {
		return c.JSON(501, map[string]string{"message": "Not implemented yet"})
	})
	api.POST("/auth/logout", func(c echo.Context) error {
		return c.JSON(501, map[string]string{"message": "Not implemented yet"})
	})
}

func (r *Router) SetupAdminRoutes() {
	admin := r.e.Group("/api/v1/admin")
	admin.Use(r.authMW.Authenticate)
	admin.Use(middleware.RequireRoles("super_admin"))

	admin.GET("/tenants", func(c echo.Context) error {
		return c.JSON(501, map[string]string{"message": "Not implemented yet"})
	})
}
