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

// Problem returns "<Field>: <Reason>".
func (p *FieldProblem) Problem() string {
	// TODO
	return ""
}

// User is the data being validated.
type User struct {
	Name  string
	Age   int
	Email string
}

// CheckName requires a non-blank name.
// TODO: this returns a non-nil Problem even for a valid name. Why?
func CheckName(u User) Problem {
	var p *FieldProblem
	if strings.TrimSpace(u.Name) == "" {
		p = &FieldProblem{Field: "name", Reason: "required"}
	}
	return p
}

// CheckAge requires 0 <= Age <= 150.
// TODO: same bug as CheckName.
func CheckAge(u User) Problem {
	var p *FieldProblem
	if u.Age < 0 || u.Age > 150 {
		p = &FieldProblem{Field: "age", Reason: "out of range"}
	}
	return p
}

// CheckEmail requires an "@" in the email.
func CheckEmail(u User) Problem {
	// TODO
	return nil
}

// Validate returns the first problem found (name, age, email), or nil.
func Validate(u User) Problem {
	// TODO
	return nil
}

// Collect returns every problem found, in check order.
func Collect(u User) []Problem {
	// TODO
	return nil
}

// IsNil reports whether p is nil or holds a nil *FieldProblem.
func IsNil(p Problem) bool {
	// TODO
	return false
}
