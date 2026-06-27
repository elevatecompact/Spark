package domain

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailTaken           = errors.New("email already taken")
	ErrUsernameTaken        = errors.New("username already taken")
	ErrSessionExpired       = errors.New("session expired")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("token expired")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidClientID      = errors.New("invalid client id")
	ErrInvalidClientSecret  = errors.New("invalid client secret")
	ErrInvalidGrantType     = errors.New("invalid grant type")
	ErrInvalidRedirectURI   = errors.New("invalid redirect uri")
	ErrInvalidScope         = errors.New("invalid scope")
	ErrAuthorizationCodeUsed = errors.New("authorization code already used")
	ErrPasswordTooWeak      = errors.New("password too weak")
	ErrUserSuspended        = errors.New("user is suspended")
	ErrUserBanned           = errors.New("user is banned")
	ErrForbidden            = errors.New("forbidden")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrRateLimited          = errors.New("rate limited")
	ErrPasskeyNotFound      = errors.New("passkey not found")
	ErrPasskeyRegistration  = errors.New("passkey registration failed")
	ErrPasskeyAuthentication = errors.New("passkey authentication failed")
	ErrInternalServer       = errors.New("internal server error")
	ErrValidation           = errors.New("validation error")
)

type DomainError struct {
	Err       error
	Message   string
	HTTPCode  int
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
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, ErrEmailTaken):
		return http.StatusConflict
	case errors.Is(err, ErrUsernameTaken):
		return http.StatusConflict
	case errors.Is(err, ErrSessionExpired):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrExpiredToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidRefreshToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidClientID):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidClientSecret):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidGrantType):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidRedirectURI):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidScope):
		return http.StatusBadRequest
	case errors.Is(err, ErrAuthorizationCodeUsed):
		return http.StatusBadRequest
	case errors.Is(err, ErrPasswordTooWeak):
		return http.StatusBadRequest
	case errors.Is(err, ErrUserSuspended):
		return http.StatusForbidden
	case errors.Is(err, ErrUserBanned):
		return http.StatusForbidden
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrRateLimited):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrPasskeyNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrPasskeyRegistration):
		return http.StatusBadRequest
	case errors.Is(err, ErrPasskeyAuthentication):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
