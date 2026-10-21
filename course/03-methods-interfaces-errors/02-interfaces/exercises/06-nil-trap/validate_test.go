package validate_test

import (
	"fmt"
	"testing"

	validate "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/06-nil-trap"
)

// show prints a Problem without calling its methods, so typed nils are visible.
func show(p validate.Problem) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%T(%+v)", p, p)
}

// same reports whether got is the problem want describes (nil means "no problem").
func same(got validate.Problem, want *validate.FieldProblem) bool {
	if want == nil {
		return got == nil
	}
	fp, ok := got.(*validate.FieldProblem)
	return ok && fp != nil && *fp == *want
}

var (
	nameP  = &validate.FieldProblem{Field: "name", Reason: "required"}
	ageP   = &validate.FieldProblem{Field: "age", Reason: "out of range"}
	emailP = &validate.FieldProblem{Field: "email", Reason: "missing @"}
	good   = validate.User{Name: "Ann", Age: 30, Email: "ann@example.com"}
)

func TestProblemMethod(t *testing.T) {
	tests := []struct {
		p    *validate.FieldProblem
		want string
	}{
		{ageP, "age: out of range"},
		{nameP, "name: required"},
		{&validate.FieldProblem{}, ": "},
	}
	for _, tt := range tests {
		if got := tt.p.Problem(); got != tt.want {
			t.Errorf("(%+v).Problem() = %q, want %q", *tt.p, got, tt.want)
		}
	}
}

func TestChecks(t *testing.T) {
	tests := []struct {
		name  string
		check func(validate.User) validate.Problem
		u     validate.User
		want  *validate.FieldProblem
	}{
		{"CheckName", validate.CheckName, good, nil},
		{"CheckName", validate.CheckName, validate.User{Name: "Zoë"}, nil},
		{"CheckName", validate.CheckName, validate.User{}, nameP},
		{"CheckName", validate.CheckName, validate.User{Name: " \t "}, nameP},
		{"CheckAge", validate.CheckAge, good, nil},
		{"CheckAge", validate.CheckAge, validate.User{Age: 0}, nil},
		{"CheckAge", validate.CheckAge, validate.User{Age: 150}, nil},
		{"CheckAge", validate.CheckAge, validate.User{Age: 151}, ageP},
		{"CheckAge", validate.CheckAge, validate.User{Age: -1}, ageP},
		{"CheckEmail", validate.CheckEmail, good, nil},
		{"CheckEmail", validate.CheckEmail, validate.User{Email: "@"}, nil},
		{"CheckEmail", validate.CheckEmail, validate.User{Email: "ann.example.com"}, emailP},
		{"CheckEmail", validate.CheckEmail, validate.User{}, emailP},
	}
	for _, tt := range tests {
		got := tt.check(tt.u)
		if !same(got, tt.want) {
			want := "nil"
			if tt.want != nil {
				want = show(tt.want)
			}
			t.Errorf("%s(%+v) = %s, want %s", tt.name, tt.u, show(got), want)
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		u    validate.User
		want *validate.FieldProblem
	}{
		{good, nil},
		{validate.User{Name: "", Age: -1}, nameP},
		{validate.User{Name: "Bo", Age: 200}, ageP},
		{validate.User{Name: "Bo", Age: 20}, emailP},
		{validate.User{}, nameP},
	}
	for _, tt := range tests {
		got := validate.Validate(tt.u)
		if !same(got, tt.want) {
			want := "nil"
			if tt.want != nil {
				want = show(tt.want)
			}
			t.Errorf("Validate(%+v) = %s, want %s", tt.u, show(got), want)
		}
	}
}

func TestCollect(t *testing.T) {
	tests := []struct {
		u    validate.User
		want []*validate.FieldProblem
	}{
		{good, nil},
		{validate.User{Age: 200}, []*validate.FieldProblem{nameP, ageP, emailP}},
		{validate.User{Name: "Bo", Age: -5, Email: "bo@x"}, []*validate.FieldProblem{ageP}},
		{validate.User{Age: 5, Email: "x@y"}, []*validate.FieldProblem{nameP}},
		{validate.User{Name: "Bo"}, []*validate.FieldProblem{emailP}},
	}
	for _, tt := range tests {
		got := validate.Collect(tt.u)
		ok := len(got) == len(tt.want)
		for i := 0; ok && i < len(got); i++ {
			ok = same(got[i], tt.want[i])
		}
		if !ok {
			gotS := make([]string, len(got))
			for i, p := range got {
				gotS[i] = show(p)
			}
			wantS := make([]string, len(tt.want))
			for i, p := range tt.want {
				wantS[i] = show(p)
			}
			t.Errorf("Collect(%+v) = %v, want %v", tt.u, gotS, wantS)
		}
	}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		p    validate.Problem
		want bool
	}{
		{nil, true},
		{(*validate.FieldProblem)(nil), true},
		{&validate.FieldProblem{}, false},
		{nameP, false},
	}
	for _, tt := range tests {
		if got := validate.IsNil(tt.p); got != tt.want {
			t.Errorf("IsNil(%s) = %v, want %v", show(tt.p), got, tt.want)
		}
	}
}
