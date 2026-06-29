package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("forbidden")
	ErrValidation         = errors.New("validation error")
	ErrNotifNotFound      = errors.New("notification not found")
	ErrDeviceNotFound     = errors.New("device not found")
	ErrTemplateNotFound   = errors.New("template not found")
	ErrChannelDisabled    = errors.New("notification channel is disabled")
	ErrRateLimited        = errors.New("rate limit exceeded")
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

func NewDomainErrorMsg(err error, msg string, httpCode int) *DomainError {
	return &DomainError{Err: err, Message: msg, HTTPCode: httpCode}
}

func HTTPStatusFromError(err error) int {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrNotifNotFound), errors.Is(err, ErrDeviceNotFound), errors.Is(err, ErrTemplateNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrChannelDisabled):
		return http.StatusBadRequest
	case errors.Is(err, ErrRateLimited):
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
