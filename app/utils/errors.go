package utils

import "net/http"

// HTTPError represents an error with an associated HTTP status code
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// Common error constructors
func NewNotFoundError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NewBadRequestError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewInternalServerError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
