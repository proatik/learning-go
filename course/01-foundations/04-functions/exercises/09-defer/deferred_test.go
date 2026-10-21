package deferred

import (
	"strings"
	"testing"
)

func ExampleStack() {
	Stack(3)
	// Output:
	// push 1
	// push 2
	// push 3
	// pop 3
	// pop 2
	// pop 1
}

func ExampleStack_one() {
	Stack(1)
	// Output:
	// push 1
	// pop 1
}

func ExampleStack_twice() {
	Stack(2)
	Stack(0)
	Stack(2)
	// Output:
	// push 1
	// push 2
	// pop 2
	// pop 1
	// push 1
	// push 2
	// pop 2
	// pop 1
}

func TestParseAge(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"42", 42},
		{"0", 0},
		{"150", 150},
		{"7", 7},
		{"+30", 30},
	}
	for _, tt := range tests {
		got, err := ParseAge(tt.s)
		if err != nil {
			t.Errorf("ParseAge(%q) returned error %v, want nil", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseAge(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestParseAgeOutOfRange(t *testing.T) {
	for _, s := range []string{"200", "-1", "151", "-100"} {
		got, err := ParseAge(s)
		if err == nil {
			t.Errorf("ParseAge(%q): err = nil, want an error", s)
			continue
		}
		if got != 0 {
			t.Errorf("ParseAge(%q) = %d, want 0 alongside the error", s, got)
		}
		want := `ParseAge("` + s + `"): age out of range`
		if err.Error() != want {
			t.Errorf("ParseAge(%q) error = %q, want %q", s, err.Error(), want)
		}
	}
}

func TestParseAgeSyntax(t *testing.T) {
	for _, s := range []string{"abc", "", " 42", "4.5", "99999999999999999999"} {
		got, err := ParseAge(s)
		if err == nil {
			t.Errorf("ParseAge(%q): err = nil, want an error", s)
			continue
		}
		if got != 0 {
			t.Errorf("ParseAge(%q) = %d, want 0 alongside the error", s, got)
		}
		prefix := `ParseAge("` + s + `"): `
		msg := err.Error()
		if !strings.HasPrefix(msg, prefix) || len(msg) == len(prefix) {
			t.Errorf("ParseAge(%q) error = %q, want %q followed by the strconv error", s, msg, prefix)
		}
		if strings.Count(msg, "ParseAge(") != 1 {
			t.Errorf("ParseAge(%q) error = %q, want the prefix exactly once", s, msg)
		}
	}
}
