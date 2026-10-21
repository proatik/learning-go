package validate

import "strings"

// Problem describes one validation failure.
type Problem interface {
	Problem() string
}

// FieldProblem is a Problem with a single field.
type FieldProblem struct {
	Field, Reason string
}

var _ Problem = (*FieldProblem)(nil)

// Problem returns "<Field>: <Reason>".
func (p *FieldProblem) Problem() string {
	return p.Field + ": " + p.Reason
}

// User is the data being validated.
type User struct {
	Name  string
	Age   int
	Email string
}

// CheckName requires a non-blank name.
func CheckName(u User) Problem {
	if strings.TrimSpace(u.Name) == "" {
		return &FieldProblem{Field: "name", Reason: "required"}
	}
	return nil // a nil interface, not a nil *FieldProblem
}

// CheckAge requires 0 <= Age <= 150.
func CheckAge(u User) Problem {
	if u.Age < 0 || u.Age > 150 {
		return &FieldProblem{Field: "age", Reason: "out of range"}
	}
	return nil
}

// CheckEmail requires an "@" in the email.
func CheckEmail(u User) Problem {
	if !strings.Contains(u.Email, "@") {
		return &FieldProblem{Field: "email", Reason: "missing @"}
	}
	return nil
}

// checks lists the validators in the order they run.
var checks = []func(User) Problem{CheckName, CheckAge, CheckEmail}

// Validate returns the first problem found (name, age, email), or nil.
func Validate(u User) Problem {
	for _, check := range checks {
		if p := check(u); p != nil {
			return p
		}
	}
	return nil
}

// Collect returns every problem found, in check order.
func Collect(u User) []Problem {
	var out []Problem
	for _, check := range checks {
		if p := check(u); p != nil {
			out = append(out, p)
		}
	}
	return out
}

// IsNil reports whether p is nil or holds a nil *FieldProblem.
func IsNil(p Problem) bool {
	if p == nil {
		return true
	}
	fp, ok := p.(*FieldProblem)
	return ok && fp == nil
}
