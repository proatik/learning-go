package signup

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Form is a signup request.
type Form struct {
	Username string
	Email    string
	Password string
	Confirm  string
	Age      int
	Terms    bool
}

// Problems Validate can report (some of them wrapped with context).
var (
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameInvalid  = errors.New("username must be 3-16 letters, digits or _")
	ErrEmailInvalid     = errors.New("invalid email")
	ErrPasswordShort    = errors.New("password must be at least 8 characters")
	ErrPasswordNoDigit  = errors.New("password must contain a digit")
	ErrPasswordMismatch = errors.New("passwords do not match")
	ErrUnderage         = errors.New("must be at least 18")
	ErrTerms            = errors.New("terms must be accepted")
)

// Validate returns every problem with f joined into one error, or nil.
func Validate(f Form) error {
	var errs []error
	if f.Username == "" {
		errs = append(errs, ErrUsernameRequired)
	} else if !validUsername(f.Username) {
		errs = append(errs, ErrUsernameInvalid)
	}
	if !validEmail(f.Email) {
		errs = append(errs, fmt.Errorf("email %q: %w", f.Email, ErrEmailInvalid))
	}
	if utf8.RuneCountInString(f.Password) < 8 {
		errs = append(errs, ErrPasswordShort)
	}
	if !strings.ContainsAny(f.Password, "0123456789") {
		errs = append(errs, ErrPasswordNoDigit)
	}
	if f.Confirm != f.Password {
		errs = append(errs, ErrPasswordMismatch)
	}
	if f.Age < 18 {
		errs = append(errs, fmt.Errorf("age %d: %w", f.Age, ErrUnderage))
	}
	if !f.Terms {
		errs = append(errs, ErrTerms)
	}
	return errors.Join(errs...)
}

func validUsername(s string) bool {
	n := utf8.RuneCountInString(s)
	if n < 3 || n > 16 {
		return false
	}
	for _, r := range s {
		ok := r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
		if !ok {
			return false
		}
	}
	return true
}

func validEmail(s string) bool {
	user, domain, ok := strings.Cut(s, "@")
	return ok && user != "" && domain != "" &&
		!strings.Contains(domain, "@") && strings.Contains(domain, ".")
}

// Messages splits a joined error into the messages of its parts.
func Messages(err error) []string {
	if err == nil {
		return nil
	}
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return []string{err.Error()}
	}
	var msgs []string
	for _, e := range joined.Unwrap() {
		msgs = append(msgs, e.Error())
	}
	return msgs
}
