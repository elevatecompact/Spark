package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrValidation     = errors.New("validation error")
	ErrScanFailed     = errors.New("scan failed")
)

type DomainError struct {
	Err      error
	Message  string
	HTTPCode int
}

func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *DomainError) Unwrap() error { return e.Err }

func NewDomainError(err error, httpCode int) *DomainError {
	return &DomainError{Err: err, HTTPCode: httpCode}
}

func HTTPStatusFromError(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrValidation), errors.Is(err, ErrScanFailed):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
