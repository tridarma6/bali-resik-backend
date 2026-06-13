package middleware

import (
	"bytes"
	"io"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var auditSkipPaths = map[string]bool{
	"/health":              true,
	"/api/v1/auth/login":   true,
	"/api/v1/auth/refresh": true,
	"/api/v1/auth/register": true,
}

var mutateMethods = map[string]bool{
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

type AuditMiddleware struct {
	logger *logrus.Logger
}

func NewAuditMiddleware(logger *logrus.Logger) *AuditMiddleware {
	return &AuditMiddleware{logger: logger}
}

func (m *AuditMiddleware) AuditLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if auditSkipPaths[c.Path()] {
			return next(c)
		}

		if !mutateMethods[c.Request().Method] {
			return next(c)
		}

		bodyBytes, _ := io.ReadAll(c.Request().Body)
		c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

		err := next(c)

		tenantID, tenantErr := GetTenantID(c)
		userID, userErr := GetUserID(c)
		userRole := GetUserRole(c)

		if tenantErr == nil && userErr == nil {
			m.logger.WithFields(logrus.Fields{
				"audit":       true,
				"tenant_id":   tenantID.String(),
				"user_id":     userID.String(),
				"role":        userRole,
				"action":      c.Request().Method + " " + c.Request().URL.Path,
				"entity_type": extractAuditEntityType(c.Path()),
				"entity_id":   extractAuditEntityID(c),
				"status":      c.Response().Status,
				"ip":          c.RealIP(),
				"user_agent":  c.Request().UserAgent(),
				"body":        string(bodyBytes),
			}).Info("audit_log")
		}

		return err
	}
}

func extractAuditEntityType(path string) string {
	parts := splitPath(path)
	for _, part := range parts {
		if part != "" && part != "api" && part != "v1" && part != "admin" && part[0] != ':' {
			return part
		}
	}
	return "unknown"
}

func extractAuditEntityID(c echo.Context) string {
	id := c.Param("id")
	if id != "" {
		if _, err := uuid.Parse(id); err == nil {
			return id
		}
	}
	return ""
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, ch := range path {
		if ch == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
