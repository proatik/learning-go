package profile_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	profile "learngo/course/03-methods-interfaces-errors/04-errors/exercises/02-validation-error"
)

func TestErrorMethod(t *testing.T) {
	tests := []struct {
		e    *profile.ValidationError
		want string
	}{
		{&profile.ValidationError{Field: "age", Reason: "must be at least 13"}, "age: must be at least 13"},
		{&profile.ValidationError{Field: "name", Reason: "must not be empty"}, "name: must not be empty"},
		{&profile.ValidationError{Field: "e-mail", Reason: "bad"}, "e-mail: bad"},
	}
	for _, tt := range tests {
		if got := tt.e.Error(); got != tt.want {
			t.Errorf("(%+v).Error() = %q, want %q", *tt.e, got, tt.want)
		}
	}
}

func TestValidateOK(t *testing.T) {
	tests := []profile.User{
		{Name: "Ann", Email: "ann@go.dev", Age: 30},
		{Name: "Bo", Email: "a@b", Age: 13},
		{Name: "Zoë Ünal", Email: "zoe@example.com", Age: 150},
		{Name: strings.Repeat("é", 20), Email: "x@y.z", Age: 40},
	}
	for _, u := range tests {
		if err := profile.Validate(u); err != nil {
			t.Errorf("Validate(%+v) error = %v (%T), want <nil>", u, err, err)
		}
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		u      profile.User
		field  string
		reason string
	}{
		{profile.User{Name: "", Email: "ann@go.dev", Age: 30}, "name", "must not be empty"},
		{profile.User{Name: "   ", Email: "ann@go.dev", Age: 30}, "name", "must not be empty"},
		{profile.User{Name: strings.Repeat("a", 21), Email: "ann@go.dev", Age: 30}, "name", "must be at most 20 characters"},
		{profile.User{Name: strings.Repeat("é", 21), Email: "ann@go.dev", Age: 30}, "name", "must be at most 20 characters"},
		{profile.User{Name: "Ann", Email: "ann.go.dev", Age: 30}, "email", "must look like user@host"},
		{profile.User{Name: "Ann", Email: "", Age: 30}, "email", "must look like user@host"},
		{profile.User{Name: "Ann", Email: "@go.dev", Age: 30}, "email", "must look like user@host"},
		{profile.User{Name: "Ann", Email: "ann@", Age: 30}, "email", "must look like user@host"},
		{profile.User{Name: "Ann", Email: "a@b@c", Age: 30}, "email", "must look like user@host"},
		{profile.User{Name: "Ann", Email: "ann@go.dev", Age: 12}, "age", "must be at least 13"},
		{profile.User{Name: "Ann", Email: "ann@go.dev", Age: -1}, "age", "must be at least 13"},
		{profile.User{Name: "Ann", Email: "ann@go.dev", Age: 151}, "age", "must be at most 150"},
		// several problems: the first one in field order wins
		{profile.User{Name: "", Email: "bad", Age: 1}, "name", "must not be empty"},
		{profile.User{Name: "Ann", Email: "bad", Age: 1}, "email", "must look like user@host"},
	}
	for _, tt := range tests {
		err := profile.Validate(tt.u)
		want := &profile.ValidationError{Field: tt.field, Reason: tt.reason}

		var ve *profile.ValidationError
		if !errors.As(err, &ve) {
			t.Errorf("Validate(%+v) error = %v (%T), want %v", tt.u, err, err, want)
			continue
		}
		if *ve != *want {
			t.Errorf("Validate(%+v) error = %+v, want %+v", tt.u, *ve, *want)
		}
		if got, wantMsg := err.Error(), tt.field+": "+tt.reason; got != wantMsg {
			t.Errorf("Validate(%+v) error = %q, want %q", tt.u, got, wantMsg)
		}
	}
}

func TestValidateThroughWrapping(t *testing.T) {
	u := profile.User{Name: "Ann", Email: "ann@go.dev", Age: 5}
	err := fmt.Errorf("signup: %w", profile.Validate(u))

	ve, ok := errors.AsType[*profile.ValidationError](err)
	if !ok {
		t.Fatalf("errors.AsType[*ValidationError](%v) ok = false, want true", err)
	}
	if ve.Field != "age" {
		t.Errorf("Validate(%+v) error field = %q, want %q", u, ve.Field, "age")
	}
	if got, want := err.Error(), "signup: age: must be at least 13"; got != want {
		t.Errorf("wrapped error = %q, want %q", got, want)
	}
}
