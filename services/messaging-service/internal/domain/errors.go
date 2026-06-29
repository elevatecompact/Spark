package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrValidation         = errors.New("validation error")
	ErrConvNotFound       = errors.New("conversation not found")
	ErrNotMember          = errors.New("not a conversation member")
	ErrNotAdmin           = errors.New("requires admin role")
	ErrMsgTooLong         = errors.New("message too long")
	ErrEditWindowExpired  = errors.New("edit window expired (1h)")
	ErrDuplicateReaction  = errors.New("already reacted with this emoji")
	ErrConversationLimit  = errors.New("conversation creation limit reached")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrConvNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrNotMember), errors.Is(err, ErrNotAdmin):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation), errors.Is(err, ErrMsgTooLong), errors.Is(err, ErrEditWindowExpired):
		return http.StatusBadRequest
	case errors.Is(err, ErrDuplicateReaction):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
