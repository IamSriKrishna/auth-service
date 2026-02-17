package utils

import "net/http"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

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
