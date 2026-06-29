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
	ErrInternal           = errors.New("internal error")
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrWalletFrozen       = errors.New("wallet is frozen")
	ErrWalletClosed       = errors.New("wallet is closed")
	ErrNegativeAmount     = errors.New("amount must be positive")
	ErrDuplicateIdempotency = errors.New("duplicate idempotency key")
	ErrBalanceExceeded    = errors.New("balance exceeds maximum")
	ErrPayoutMinimum      = errors.New("amount below minimum payout")
	ErrTransactionFailed  = errors.New("transaction failed")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrWalletNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation), errors.Is(err, ErrNegativeAmount), errors.Is(err, ErrPayoutMinimum):
		return http.StatusBadRequest
	case errors.Is(err, ErrInsufficientFunds):
		return http.StatusPaymentRequired
	case errors.Is(err, ErrWalletFrozen), errors.Is(err, ErrWalletClosed):
		return http.StatusForbidden
	case errors.Is(err, ErrDuplicateIdempotency):
		return http.StatusConflict
	case errors.Is(err, ErrBalanceExceeded):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
