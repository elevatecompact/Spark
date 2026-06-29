package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound        = errors.New("resource not found")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation error")
	ErrNoEmbedding     = errors.New("embedding not found")
	ErrNoModel         = errors.New("no active model")
	ErrFeedNotAvail    = errors.New("feed not available")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrNoEmbedding), errors.Is(err, ErrNoModel):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrFeedNotAvail):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
