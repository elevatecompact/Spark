package validators

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,30}$`)
	displayNameRegex = regexp.MustCompile(`^[^\p{Cc}\p{Cf}]{1,50}$`)
	htmlTagRegex    = regexp.MustCompile(`<[^>]*>`)
	scriptRegex     = regexp.MustCompile(`(?i)<script[\s>]`)
	eventRegex      = regexp.MustCompile(`(?i)\son\w+\s*=`)

	minPasswordLength = 8
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	msgs := make([]string, len(ve))
	for i, e := range ve {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

func (ve ValidationErrors) AsError() error {
	if len(ve) == 0 {
		return nil
	}
	return ve
}

func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}

	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return &ValidationError{Field: "username", Message: "username is required"}
	}

	if !usernameRegex.MatchString(username) {
		return &ValidationError{
			Field:   "username",
			Message: "username must be 3-30 characters and contain only letters, numbers, dots, underscores, and hyphens",
		}
	}

	if strings.HasPrefix(username, ".") || strings.HasSuffix(username, ".") {
		return &ValidationError{Field: "username", Message: "username cannot start or end with a dot"}
	}

	if strings.HasPrefix(username, "-") || strings.HasSuffix(username, "-") {
		return &ValidationError{Field: "username", Message: "username cannot start or end with a hyphen"}
	}

	return nil
}

func ValidatePassword(password string) error {
	var errs ValidationErrors

	if password == "" {
		errs = append(errs, ValidationError{Field: "password", Message: "password is required"})
		return errs.AsError()
	}

	if len(password) < minPasswordLength {
		errs = append(errs, ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("password must be at least %d characters", minPasswordLength),
		})
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errs = append(errs, ValidationError{Field: "password", Message: "password must contain at least one uppercase letter"})
	}
	if !hasLower {
		errs = append(errs, ValidationError{Field: "password", Message: "password must contain at least one lowercase letter"})
	}
	if !hasDigit {
		errs = append(errs, ValidationError{Field: "password", Message: "password must contain at least one digit"})
	}
	if !hasSpecial {
		errs = append(errs, ValidationError{Field: "password", Message: "password must contain at least one special character"})
	}

	return errs.AsError()
}

func ValidateDisplayName(name string) error {
	if name == "" {
		return &ValidationError{Field: "display_name", Message: "display name is required"}
	}

	if !displayNameRegex.MatchString(name) {
		return &ValidationError{
			Field:   "display_name",
			Message: "display name must be between 1 and 50 characters and contain no control characters",
		}
	}

	return nil
}

func SanitizeHTML(input string) string {
	s := scriptRegex.ReplaceAllString(input, "&lt;script")
	s = eventRegex.ReplaceAllString(s, " blocked-event")
	s = htmlTagRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "javascript:", "")
	s = strings.ReplaceAll(s, "vbscript:", "")
	return s
}
