package signup

import "errors"

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
	// TODO
	return nil
}

// Messages splits a joined error into the messages of its parts.
func Messages(err error) []string {
	// TODO
	return nil
}
