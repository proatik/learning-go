package zeros_test

import (
	"errors"
	"strconv"
	"testing"

	zr "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/03-zero-values"
)

type Point struct{ X, Y int }

// mustPanic runs f and returns the value it panicked with (nil if it didn't).
func mustPanic(f func()) (r any) {
	defer func() { r = recover() }()
	f()
	return nil
}

func TestMustOK(t *testing.T) {
	if got := zr.Must(strconv.Atoi("42")); got != 42 {
		t.Errorf("Must(strconv.Atoi(%q)) = %d, want %d", "42", got, 42)
	}
	if got := zr.Must("hello", nil); got != "hello" {
		t.Errorf("Must(%q, nil) = %q, want %q", "hello", got, "hello")
	}
	if got := zr.Must(0, nil); got != 0 {
		t.Errorf("Must(0, nil) = %d, want 0", got)
	}
	if got := zr.Must(Point{1, 2}, nil); got != (Point{1, 2}) {
		t.Errorf("Must(%v, nil) = %v, want %v", Point{1, 2}, got, Point{1, 2})
	}
}

func TestMustPanics(t *testing.T) {
	r := mustPanic(func() { zr.Must(strconv.Atoi("x")) })
	err, ok := r.(error)
	if !ok {
		t.Fatalf("Must(strconv.Atoi(%q)) panic value = %v, want the error", "x", r)
	}
	if _, ok := errors.AsType[*strconv.NumError](err); !ok {
		t.Errorf("Must(strconv.Atoi(%q)) panic value = %v, want a *strconv.NumError", "x", err)
	}

	sentinel := errors.New("boom")
	r = mustPanic(func() { zr.Must("ignored", sentinel) })
	if r != sentinel {
		t.Errorf("Must(%q, boom) panic value = %v, want %v", "ignored", r, sentinel)
	}
}

func TestFirstOr(t *testing.T) {
	tests := []struct {
		in       []string
		fallback string
		want     string
	}{
		{[]string{"a", "b"}, "none", "a"},
		{[]string{""}, "none", ""},
		{[]string{}, "none", "none"},
		{nil, "none", "none"},
		{nil, "", ""},
	}
	for _, tt := range tests {
		if got := zr.FirstOr(tt.in, tt.fallback); got != tt.want {
			t.Errorf("FirstOr(%q, %q) = %q, want %q", tt.in, tt.fallback, got, tt.want)
		}
	}
	if got := zr.FirstOr([]int{0, 5}, 9); got != 0 {
		t.Errorf("FirstOr(%v, %d) = %d, want %d", []int{0, 5}, 9, got, 0)
	}
	if got := zr.FirstOr[float64](nil, 1.5); got != 1.5 {
		t.Errorf("FirstOr(nil, %v) = %v, want %v", 1.5, got, 1.5)
	}
}

func TestIsZero(t *testing.T) {
	var nilPtr *int
	checks := []struct {
		name string
		got  bool
		want bool
	}{
		{"IsZero(0)", zr.IsZero(0), true},
		{"IsZero(7)", zr.IsZero(7), false},
		{"IsZero(-1)", zr.IsZero(-1), false},
		{`IsZero("")`, zr.IsZero(""), true},
		{`IsZero(" ")`, zr.IsZero(" "), false},
		{"IsZero(0.0)", zr.IsZero(0.0), true},
		{"IsZero(false)", zr.IsZero(false), true},
		{"IsZero(true)", zr.IsZero(true), false},
		{"IsZero(Point{})", zr.IsZero(Point{}), true},
		{"IsZero(Point{X: 1})", zr.IsZero(Point{X: 1}), false},
		{"IsZero((*int)(nil))", zr.IsZero(nilPtr), true},
		{"IsZero(new(int))", zr.IsZero(new(int)), false},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %t, want %t", c.name, c.got, c.want)
		}
	}
}

func TestFirstNonZero(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{[]string{"", "", "guest", "admin"}, "guest"},
		{[]string{"admin", "guest"}, "admin"},
		{[]string{"", ""}, ""},
		{[]string{}, ""},
		{nil, ""},
	}
	for _, tt := range tests {
		if got := zr.FirstNonZero(tt.in...); got != tt.want {
			t.Errorf("FirstNonZero(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
	if got := zr.FirstNonZero(0, 0, 8080, 80); got != 8080 {
		t.Errorf("FirstNonZero(0, 0, 8080, 80) = %d, want %d", got, 8080)
	}
	if got := zr.FirstNonZero(0, 0); got != 0 {
		t.Errorf("FirstNonZero(0, 0) = %d, want %d", got, 0)
	}
	if got := zr.FirstNonZero[int](); got != 0 {
		t.Errorf("FirstNonZero() = %d, want %d", got, 0)
	}
	if got := zr.FirstNonZero(Point{}, Point{0, 3}); got != (Point{0, 3}) {
		t.Errorf("FirstNonZero(%v, %v) = %v, want %v", Point{}, Point{0, 3}, got, Point{0, 3})
	}
}
