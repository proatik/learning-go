package ages

import (
	"errors"
	"testing"
)

func TestParseAgeValid(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"42", 42},
		{"0", 0},
		{"150", 150},
		{"  7 ", 7},
		{"\t18\n", 18},
		{"+5", 5},
		{"007", 7},
	}
	for _, tt := range tests {
		got, err := ParseAge(tt.in)
		if err != nil {
			t.Errorf("ParseAge(%q) error = %v, want <nil>", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseAge(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestParseAgeErrors(t *testing.T) {
	sentinels := []error{ErrEmpty, ErrNotNumber, ErrOutOfRange}
	tests := []struct {
		in      string
		want    error
		wantMsg string
	}{
		{"", ErrEmpty, "empty age"},
		{"   ", ErrEmpty, "empty age"},
		{"abc", ErrNotNumber, `parse age "abc": not a number`},
		{" 4 2 ", ErrNotNumber, `parse age "4 2": not a number`},
		{"12.5", ErrNotNumber, `parse age "12.5": not a number`},
		{"ⅷ", ErrNotNumber, `parse age "ⅷ": not a number`},
		{"151", ErrOutOfRange, "age 151: out of range"},
		{"-1", ErrOutOfRange, "age -1: out of range"},
		{" 1000 ", ErrOutOfRange, "age 1000: out of range"},
	}
	for _, tt := range tests {
		got, err := ParseAge(tt.in)
		if !errors.Is(err, tt.want) {
			t.Errorf("ParseAge(%q) error = %v, want %v", tt.in, err, tt.want)
			continue
		}
		for _, s := range sentinels {
			if s != tt.want && errors.Is(err, s) {
				t.Errorf("ParseAge(%q) error = %v, should not match %v", tt.in, err, s)
			}
		}
		if err.Error() != tt.wantMsg {
			t.Errorf("ParseAge(%q) error = %q, want %q", tt.in, err.Error(), tt.wantMsg)
		}
		if got != 0 {
			t.Errorf("ParseAge(%q) = %d, want 0 with an error", tt.in, got)
		}
	}
}

func TestParseAgeEmptyIsUnwrapped(t *testing.T) {
	if _, err := ParseAge(""); err != ErrEmpty {
		t.Errorf("ParseAge(%q) error = %v, want ErrEmpty itself (not wrapped)", "", err)
	}
}
