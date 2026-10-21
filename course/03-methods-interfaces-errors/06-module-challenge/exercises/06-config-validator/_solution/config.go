package config

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

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
func (e *FieldError) Error() string { return e.Field + ": " + e.Err.Error() }

// Unwrap exposes the problem to errors.Is.
func (e *FieldError) Unwrap() error { return e.Err }

// Config is a service configuration.
type Config struct {
	Name     string
	Port     int
	Env      string // "dev", "staging" or "prod"
	Replicas int
	Tags     []string
}

var envs = []string{"dev", "staging", "prod"}

// Validate checks every field and returns all problems joined with
// errors.Join, or nil if the config is valid.
func (c Config) Validate() error {
	var errs []error
	add := func(field string, err error) {
		errs = append(errs, &FieldError{Field: field, Err: err})
	}

	if strings.TrimSpace(c.Name) == "" {
		add("name", ErrRequired)
	}
	if c.Port < 1 || c.Port > 65535 {
		add("port", ErrOutOfRange)
	}
	switch {
	case c.Env == "":
		add("env", ErrRequired)
	case !slices.Contains(envs, c.Env):
		add("env", ErrInvalid)
	}
	switch {
	case c.Replicas < 1 || c.Replicas > 10:
		add("replicas", ErrOutOfRange)
	case c.Env == "prod" && c.Replicas < 2:
		add("replicas", ErrInvalid)
	}
	seen := map[string]bool{}
	for i, tag := range c.Tags {
		field := "tags[" + strconv.Itoa(i) + "]"
		switch {
		case strings.TrimSpace(tag) == "":
			add(field, ErrRequired)
		case seen[tag]:
			add(field, ErrInvalid)
		}
		seen[tag] = true
	}
	return errors.Join(errs...) // nil when errs is empty
}

// Fields returns the Field of every *FieldError in err's tree, depth-first
// in order. It follows both Unwrap() error and Unwrap() []error.
func Fields(err error) []string {
	var out []string
	var walk func(error)
	walk = func(err error) {
		if err == nil {
			return
		}
		if fe, ok := err.(*FieldError); ok {
			out = append(out, fe.Field)
		}
		switch u := err.(type) {
		case interface{ Unwrap() []error }:
			for _, e := range u.Unwrap() {
				walk(e)
			}
		case interface{ Unwrap() error }:
			walk(u.Unwrap())
		}
	}
	walk(err)
	return out
}
