package erru

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Common errors
var (
	ErrBadRequest          = New(http.StatusBadRequest, "Bad Request")
	ErrUnauthorized        = New(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "Forbidden")
	ErrNotFound            = New(http.StatusNotFound, "Not Found")
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal Server Error")
)
