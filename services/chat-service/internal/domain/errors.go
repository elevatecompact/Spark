package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrForbidden         = errors.New("forbidden")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrValidation        = errors.New("validation error")
	ErrRateLimited       = errors.New("rate limited")
	ErrInternalServer    = errors.New("internal server error")
	ErrRoomNotFound      = errors.New("room not found")
	ErrMessageTooLong    = errors.New("message too long")
	ErrUserMuted         = errors.New("user is muted in this room")
	ErrUserBanned        = errors.New("user is banned from this room")
	ErrSlowMode          = errors.New("slow mode enabled, please wait")
	ErrDuplicateEntry    = errors.New("duplicate entry")
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
	case errors.Is(err, ErrRoomNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrMessageTooLong):
		return http.StatusBadRequest
	case errors.Is(err, ErrUserMuted):
		return http.StatusForbidden
	case errors.Is(err, ErrUserBanned):
		return http.StatusForbidden
	case errors.Is(err, ErrSlowMode):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrDuplicateEntry):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
