package middleware

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	ContextKeyTenantID  = "tenant_id"
	ContextKeyUserID    = "user_id"
	ContextKeyUserRole  = "user_role"
	ContextKeyUserEmail = "user_email"

	SuperAdminRole = "super_admin"
)

var (
	ErrTenantRequired   = errors.New("tenant context is required")
	ErrTenantMismatch   = errors.New("tenant mismatch")
	ErrAccessDenied     = errors.New("access denied")
)

func GetTenantID(ctx echo.Context) (uuid.UUID, error) {
	tenantID, ok := ctx.Get(ContextKeyTenantID).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrTenantRequired
	}
	return tenantID, nil
}

func GetUserID(ctx echo.Context) (uuid.UUID, error) {
	userID, ok := ctx.Get(ContextKeyUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user context is required")
	}
	return userID, nil
}

func GetUserRole(ctx echo.Context) string {
	role, _ := ctx.Get(ContextKeyUserRole).(string)
	return role
}

func IsSuperAdmin(ctx echo.Context) bool {
	return GetUserRole(ctx) == SuperAdminRole
}
