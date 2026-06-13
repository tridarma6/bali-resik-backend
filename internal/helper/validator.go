package helper

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func ValidateRequest(ctx echo.Context, req interface{}) error {
	if err := ctx.Bind(req); err != nil {
		return ErrorResponseJSON(ctx, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
	}

	if err := validate.Struct(req); err != nil {
		var errMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errMessages = append(errMessages, formatValidationError(err))
		}
		return ErrorResponseJSON(ctx, http.StatusBadRequest, "VALIDATION_ERROR", strings.Join(errMessages, "; "))
	}

	return nil
}

func formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + err.Param() + " characters"
	case "max":
		return field + " must be at most " + err.Param() + " characters"
	case "uuid":
		return field + " must be a valid UUID"
	default:
		return field + " is invalid"
	}
}
