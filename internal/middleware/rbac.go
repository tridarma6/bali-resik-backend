package middleware

import (
	"github.com/labstack/echo/v4"
)

type RBACMiddleware struct {
	allowedRoles []string
}

func NewRBACMiddleware(allowedRoles ...string) *RBACMiddleware {
	return &RBACMiddleware{
		allowedRoles: allowedRoles,
	}
}

func (m *RBACMiddleware) RequireRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		role := GetUserRole(ctx)
		if role == "" {
			return echo.ErrForbidden
		}

		for _, allowed := range m.allowedRoles {
			if role == allowed {
				return next(ctx)
			}
		}

		return echo.ErrForbidden
	}
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	mw := NewRBACMiddleware(roles...)
	return mw.RequireRole
}
