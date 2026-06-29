package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrForbidden           = errors.New("forbidden")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrValidation          = errors.New("validation error")
	ErrPlanNotFound        = errors.New("plan not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrPlanInactive        = errors.New("plan is inactive")
	ErrAlreadySubscribed   = errors.New("already subscribed to this plan")
	ErrMaxSubscriptions    = errors.New("maximum active subscriptions reached")
	ErrSubscriptionActive  = errors.New("subscription is already active")
	ErrInvoiceNotFound     = errors.New("invoice not found")
	ErrNotOwner            = errors.New("not the subscription owner")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrPlanNotFound), errors.Is(err, ErrSubscriptionNotFound), errors.Is(err, ErrInvoiceNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrNotOwner):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation), errors.Is(err, ErrPlanInactive), errors.Is(err, ErrAlreadySubscribed), errors.Is(err, ErrMaxSubscriptions):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
