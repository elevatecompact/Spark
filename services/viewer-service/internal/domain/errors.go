package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound              = errors.New("resource not found")
	ErrForbidden             = errors.New("forbidden")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrValidation            = errors.New("validation error")
	ErrRateLimited           = errors.New("rate limited")
	ErrInternalServer        = errors.New("internal server error")
	ErrContentNotFound       = errors.New("content not found")
	ErrDuplicateEntry        = errors.New("duplicate entry")
	ErrMaxBookmarksReached   = errors.New("maximum bookmarks reached")
	ErrMaxWatchLaterReached  = errors.New("maximum watch later items reached")
	ErrAlreadyRated          = errors.New("already rated")
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

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(err error, httpCode int) *DomainError {
	return &DomainError{Err: err, HTTPCode: httpCode}
}

func NewDomainErrorMsg(err error, msg string, httpCode int) *DomainError {
	return &DomainError{Err: err, Message: msg, HTTPCode: httpCode}
}

func HTTPStatusFromError(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, ErrRateLimited):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrContentNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrDuplicateEntry):
		return http.StatusConflict
	case errors.Is(err, ErrMaxBookmarksReached):
		return http.StatusConflict
	case errors.Is(err, ErrMaxWatchLaterReached):
		return http.StatusConflict
	case errors.Is(err, ErrAlreadyRated):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
