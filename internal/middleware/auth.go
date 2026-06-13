package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/auth"
)

type AuthMiddleware struct {
	jwtService auth.JWTService
	log        *logrus.Logger
}

func NewAuthMiddleware(jwtService auth.JWTService, log *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		log:        log,
	}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.ErrUnauthorized
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			m.log.WithError(err).Warn("invalid JWT token")
			return echo.ErrUnauthorized
		}

		ctx.Set(ContextKeyUserID, claims.UserID)
		ctx.Set(ContextKeyTenantID, claims.TenantID)
		ctx.Set(ContextKeyUserRole, claims.Role)
		ctx.Set(ContextKeyUserEmail, claims.Email)

		return next(ctx)
	}
}
