package config

import "errors"

// Problems a field can have. FieldError wraps one of them.
var (
	ErrRequired   = errors.New("required")
	ErrOutOfRange = errors.New("out of range")
	ErrInvalid    = errors.New("invalid")
)

// FieldError reports one problem with one field.
type FieldError struct {
	Field string
	Err   error
}

// Error returns "FIELD: PROBLEM".
func (e *FieldError) Error() string {
	// TODO
	return ""
}

// TODO: add an Unwrap method so errors.Is can see Err.

// Config is a service configuration.
type Config struct {
	Name     string
	Port     int
	Env      string // "dev", "staging" or "prod"
	Replicas int
	Tags     []string
}

// Validate checks every field and returns all problems joined with
// errors.Join, or nil if the config is valid.
func (c Config) Validate() error {
	// TODO
	return nil
}

// Fields returns the Field of every *FieldError in err's tree, depth-first
// in order. It follows both Unwrap() error and Unwrap() []error.
func Fields(err error) []string {
	// TODO
	return nil
}
