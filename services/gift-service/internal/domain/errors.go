package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrForbidden           = errors.New("forbidden")
	ErrValidation          = errors.New("validation error")
	ErrGiftItemNotFound    = errors.New("gift item not found")
	ErrGiftNotFound        = errors.New("gift not found")
	ErrGiftCardNotFound    = errors.New("gift card not found")
	ErrGiftCardExpired     = errors.New("gift card expired")
	ErrGiftCardRedeemed    = errors.New("gift card already redeemed")
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrCampaignInactive    = errors.New("campaign is not active")
	ErrCampaignBudgetExhausted = errors.New("campaign budget exhausted")
	ErrGiftSendingDisabled = errors.New("gift sending is disabled")
	ErrGiftCardsDisabled   = errors.New("gift cards are disabled")
	ErrCampaignMatchingDisabled = errors.New("campaign matching is disabled")
	ErrAmountTooSmall      = errors.New("amount below minimum")
	ErrAmountTooLarge      = errors.New("amount above maximum")
	ErrNotOwner            = errors.New("not owner")
	ErrGiftNotCompleted    = errors.New("gift is not in completed status")
	ErrBatchLimitExceeded  = errors.New("batch limit exceeded (max 50)")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrGiftItemNotFound), errors.Is(err, ErrGiftNotFound), errors.Is(err, ErrGiftCardNotFound), errors.Is(err, ErrCampaignNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrNotOwner):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrAmountTooSmall), errors.Is(err, ErrAmountTooLarge), errors.Is(err, ErrCampaignInactive), errors.Is(err, ErrBatchLimitExceeded):
		return http.StatusBadRequest
	case errors.Is(err, ErrGiftCardExpired), errors.Is(err, ErrGiftCardRedeemed), errors.Is(err, ErrCampaignBudgetExhausted):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
