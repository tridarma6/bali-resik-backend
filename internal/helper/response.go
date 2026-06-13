package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func SuccessOK(ctx echo.Context, data interface{}, message string) error {
	return ctx.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func SuccessCreated(ctx echo.Context, data interface{}, message string) error {
	return ctx.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func SuccessPaginated(ctx echo.Context, data interface{}, meta *PaginationMeta, message string) error {
	return ctx.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta:    meta,
	})
}

func ErrorResponseJSON(ctx echo.Context, status int, code, message string) error {
	return ctx.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(ctx echo.Context, message string) error {
	return ErrorResponseJSON(ctx, http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(ctx echo.Context, message string) error {
	return ErrorResponseJSON(ctx, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(ctx echo.Context, message string) error {
	return ErrorResponseJSON(ctx, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(ctx echo.Context, message string) error {
	return ErrorResponseJSON(ctx, http.StatusNotFound, "NOT_FOUND", message)
}

func InternalError(ctx echo.Context, message string) error {
	return ErrorResponseJSON(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

func Conflict(ctx echo.Context, message string) error {
	return ErrorResponseJSON(ctx, http.StatusConflict, "CONFLICT", message)
}
