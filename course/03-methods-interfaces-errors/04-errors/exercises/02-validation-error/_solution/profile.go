package profile

import (
	"strings"
	"unicode/utf8"
)

// User is a sign-up request.
type User struct {
	Name  string
	Email string
	Age   int
}

// ValidationError reports which field of a User is invalid and why.
type ValidationError struct {
	Field  string
	Reason string
}

// Error returns "<field>: <reason>".
func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Reason
}

// Validate returns the first problem with u as a *ValidationError, or nil.
func Validate(u User) error {
	name := strings.TrimSpace(u.Name)
	if name == "" {
		return &ValidationError{Field: "name", Reason: "must not be empty"}
	}
	if utf8.RuneCountInString(u.Name) > 20 {
		return &ValidationError{Field: "name", Reason: "must be at most 20 characters"}
	}
	user, host, ok := strings.Cut(u.Email, "@")
	if !ok || user == "" || host == "" || strings.Contains(host, "@") {
		return &ValidationError{Field: "email", Reason: "must look like user@host"}
	}
	if u.Age < 13 {
		return &ValidationError{Field: "age", Reason: "must be at least 13"}
	}
	if u.Age > 150 {
		return &ValidationError{Field: "age", Reason: "must be at most 150"}
	}
	return nil
}
