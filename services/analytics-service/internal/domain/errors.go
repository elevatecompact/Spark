package domain

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrForbidden           = errors.New("forbidden")
	ErrValidation          = errors.New("validation error")
	ErrDashboardNotFound   = errors.New("dashboard not found")
	ErrReportNotFound      = errors.New("report not found")
	ErrFunnelNotFound      = errors.New("funnel not found")
	ErrTemplateNotFound    = errors.New("template not found")
	ErrInvalidQuery        = errors.New("invalid query")
	ErrExportFailed        = errors.New("export failed")
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
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrDashboardNotFound), errors.Is(err, ErrReportNotFound), errors.Is(err, ErrFunnelNotFound), errors.Is(err, ErrTemplateNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrInvalidQuery):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
