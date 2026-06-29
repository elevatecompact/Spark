package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrValidation   = errors.New("validation error")
	ErrForbidden    = errors.New("forbidden")
	ErrSoldOut      = errors.New("ticket tier sold out")
	ErrAlreadyRSVPd = errors.New("already registered")
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
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrSoldOut), errors.Is(err, ErrAlreadyRSVPd):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
