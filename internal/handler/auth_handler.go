package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/helper"
	"github.com/indim/bali-resik-backend/internal/middleware"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type AuthHandler struct {
	authUseCase ucase.AuthUseCase
	log         *logrus.Logger
}

func NewAuthHandler(authUseCase ucase.AuthUseCase, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		log:         log,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req request.RegisterRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.authUseCase.Register(&req)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Registration successful")
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.authUseCase.Login(&req)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return helper.SuccessOK(c, resp, "Login successful")
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req request.RefreshTokenRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.authUseCase.RefreshToken(&req)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return helper.SuccessOK(c, resp, "Token refreshed")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User not authenticated")
	}

	refreshToken := c.Request().Header.Get("X-Refresh-Token")

	if err := h.authUseCase.Logout(userID, refreshToken); err != nil {
		return h.handleAuthError(c, err)
	}

	return helper.SuccessOK(c, nil, "Logout successful")
}

func (h *AuthHandler) handleAuthError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled auth error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "INVALID_CREDENTIALS":
		return helper.Unauthorized(c, domainErr.Message)
	case "TENANT_NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "EMAIL_EXISTS":
		return helper.Conflict(c, domainErr.Message)
	case "INVALID_TOKEN", "TOKEN_EXPIRED":
		return helper.Unauthorized(c, domainErr.Message)
	case "ACCOUNT_DISABLED":
		return helper.Forbidden(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
