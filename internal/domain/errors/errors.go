package errors

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrTenantMismatch   = errors.New("tenant mismatch")
	ErrEmailAlreadyUsed = errors.New("email already used")
	ErrInvalidCreds     = errors.New("invalid credentials")
	ErrTokenExpired     = errors.New("token expired")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInsufficientPoints = errors.New("insufficient points")
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
