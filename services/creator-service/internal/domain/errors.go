package domain

import "errors"

var (
	ErrCreatorNotFound      = errors.New("creator not found")
	ErrCreatorAlreadyExists = errors.New("creator profile already exists")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrInvalidCategory      = errors.New("invalid category")
	ErrPortfolioNotFound    = errors.New("portfolio item not found")
	ErrScheduleConflict     = errors.New("schedule slot conflicts with existing slot")
	ErrSlotNotFound         = errors.New("schedule slot not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidInput         = errors.New("invalid input")
	ErrVerificationPending  = errors.New("verification already pending")
	ErrAlreadyVerified      = errors.New("creator already verified")
	ErrSelfFollow           = errors.New("cannot follow yourself")
)

func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, ErrCreatorNotFound):
		return 404
	case errors.Is(err, ErrCategoryNotFound):
		return 404
	case errors.Is(err, ErrPortfolioNotFound):
		return 404
	case errors.Is(err, ErrSlotNotFound):
		return 404
	case errors.Is(err, ErrCreatorAlreadyExists):
		return 409
	case errors.Is(err, ErrScheduleConflict):
		return 409
	case errors.Is(err, ErrVerificationPending):
		return 409
	case errors.Is(err, ErrAlreadyVerified):
		return 409
	case errors.Is(err, ErrInvalidCategory):
		return 400
	case errors.Is(err, ErrInvalidInput):
		return 400
	case errors.Is(err, ErrSelfFollow):
		return 400
	case errors.Is(err, ErrUnauthorized):
		return 401
	case errors.Is(err, ErrForbidden):
		return 403
	default:
		return 500
	}
}
