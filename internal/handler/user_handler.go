package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/helper"
	"github.com/indim/bali-resik-backend/internal/middleware"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type UserHandler struct {
	userUseCase ucase.UserUseCase
	log         *logrus.Logger
}

func NewUserHandler(userUseCase ucase.UserUseCase, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		log:         log,
	}
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	resp, err := h.userUseCase.GetProfile(tenantID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessOK(c, resp, "Profile retrieved")
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var req request.UpdateProfileRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.userUseCase.UpdateProfile(tenantID, userID, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessOK(c, resp, "Profile updated")
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var req request.ChangePasswordRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	if err := h.userUseCase.ChangePassword(tenantID, userID, &req); err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessOK(c, nil, "Password changed")
}

func (h *UserHandler) UploadAvatar(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return helper.BadRequest(c, "Avatar image file is required")
	}

	if file.Size > 5*1024*1024 {
		return helper.BadRequest(c, "Avatar image must be less than 5MB")
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return helper.BadRequest(c, "Only JPG, PNG, and WebP images are allowed")
	}

	filename := fmt.Sprintf("%s_%d%s", userID.String(), time.Now().UnixNano(), ext)
	uploadDir := "uploads/avatars"
	os.MkdirAll(uploadDir, 0755)

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		h.log.WithError(err).Error("failed to create avatar file")
		return helper.InternalError(c, "Failed to upload avatar")
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return helper.InternalError(c, "Failed to read uploaded file")
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		h.log.WithError(err).Error("failed to save avatar file")
		return helper.InternalError(c, "Failed to save avatar")
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	resp, err := h.userUseCase.UpdateAvatar(tenantID, userID, avatarURL)
	if err != nil {
		return h.handleError(c, err)
	}

	return helper.SuccessOK(c, resp, "Avatar updated")
}

func (h *UserHandler) handleError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled user error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "FORBIDDEN":
		return helper.Forbidden(c, domainErr.Message)
	case "INVALID_CREDENTIALS":
		return helper.Unauthorized(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}
