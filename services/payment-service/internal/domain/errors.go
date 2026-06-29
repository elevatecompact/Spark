package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrForbidden            = errors.New("forbidden")
	ErrValidation           = errors.New("validation error")
	ErrIntentNotFound       = errors.New("payment intent not found")
	ErrMethodNotFound       = errors.New("payment method not found")
	ErrProcessorDisabled    = errors.New("payment processor is disabled")
	ErrWebhookNotFound      = errors.New("webhook event not found")
	ErrRefundFailed         = errors.New("refund failed")
	ErrConfirmFailed        = errors.New("confirmation failed")
	ErrCancelFailed         = errors.New("cancellation failed")
	ErrDuplicateIdempotency = errors.New("duplicate idempotency key")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrProcessorError       = errors.New("processor error")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrIntentNotFound), errors.Is(err, ErrMethodNotFound), errors.Is(err, ErrWebhookNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrInvalidAmount), errors.Is(err, ErrProcessorDisabled):
		return http.StatusBadRequest
	case errors.Is(err, ErrDuplicateIdempotency):
		return http.StatusConflict
	case errors.Is(err, ErrProcessorError), errors.Is(err, ErrRefundFailed), errors.Is(err, ErrConfirmFailed), errors.Is(err, ErrCancelFailed):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
