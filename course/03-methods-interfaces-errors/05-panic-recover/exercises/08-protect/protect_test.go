package middleware_test

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	mw "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/08-protect"
)

// catch runs f and returns the value it panicked with (nil if it did not).
func catch(f func()) (r any) {
	defer func() { r = recover() }()
	f()
	return nil
}

func first(s string) string { return s[:1] }

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func boom(string) string { panic("boom") }

func TestPanicErrorMethods(t *testing.T) {
	tests := []struct {
		e      *mw.PanicError
		msg    string
		unwrap error
	}{
		{&mw.PanicError{Input: "x", Value: "bad"}, `panic handling "x": bad`, nil},
		{&mw.PanicError{Input: "", Value: 42}, `panic handling "": 42`, nil},
		{&mw.PanicError{Input: "a\tb", Value: io.EOF}, `panic handling "a\tb": EOF`, io.EOF},
		{&mw.PanicError{Input: "ü", Value: mw.ErrAbort}, `panic handling "ü": abort handler`, mw.ErrAbort},
	}
	for _, tt := range tests {
		if got := tt.e.Error(); got != tt.msg {
			t.Errorf("(%+v).Error() = %q, want %q", *tt.e, got, tt.msg)
		}
		if got := tt.e.Unwrap(); got != tt.unwrap {
			t.Errorf("(%+v).Unwrap() = %v, want %v", *tt.e, got, tt.unwrap)
		}
	}
}

func TestProtectNoPanic(t *testing.T) {
	tests := []struct {
		h    func(string) string
		in   string
		want string
	}{
		{strings.ToUpper, "go", "GO"},
		{strings.ToUpper, "", ""},
		{first, "hello", "h"},
		{reverse, "héllo", "olléh"},
	}
	for _, tt := range tests {
		var got string
		var err error
		if r := catch(func() { got, err = mw.Protect(tt.h)(tt.in) }); r != nil {
			t.Errorf("Protect(h)(%q) panicked: %v", tt.in, r)
			continue
		}
		if got != tt.want || err != nil {
			t.Errorf("Protect(h)(%q) = %q, %v; want %q, <nil>", tt.in, got, err, tt.want)
		}
	}
}

func TestProtectPanics(t *testing.T) {
	sentinel := errors.New("db down")
	tests := []struct {
		name  string
		h     func(string) string
		in    string
		msg   string
		value any // expected PanicError.Value, compared with == when comparable
	}{
		{"string panic", boom, "x", `panic handling "x": boom`, "boom"},
		{"int panic", func(string) string { panic(7) }, "q", `panic handling "q": 7`, 7},
		{"error panic", func(string) string { panic(sentinel) }, "id=3", `panic handling "id=3": db down`, sentinel},
		{"run-time panic", first, "", `panic handling "": runtime error: slice bounds out of range [:1] with length 0`, nil},
	}
	for _, tt := range tests {
		var got string
		var err error
		if r := catch(func() { got, err = mw.Protect(tt.h)(tt.in) }); r != nil {
			t.Errorf("%s: Protect(h)(%q) panicked: %v", tt.name, tt.in, r)
			continue
		}
		var pe *mw.PanicError
		if !errors.As(err, &pe) {
			t.Errorf("%s: Protect(h)(%q) error = %v, want a *PanicError", tt.name, tt.in, err)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("%s: Protect(h)(%q) error = %q, want %q", tt.name, tt.in, err.Error(), tt.msg)
		}
		if pe.Input != tt.in {
			t.Errorf("%s: Protect(h)(%q) PanicError.Input = %q, want %q", tt.name, tt.in, pe.Input, tt.in)
		}
		if tt.value != nil && pe.Value != tt.value {
			t.Errorf("%s: Protect(h)(%q) PanicError.Value = %v, want %v", tt.name, tt.in, pe.Value, tt.value)
		}
		if got != "" {
			t.Errorf("%s: Protect(h)(%q) = %q, want \"\"", tt.name, tt.in, got)
		}
	}
}

func TestProtectUnwrapChain(t *testing.T) {
	var err error
	if r := catch(func() { _, err = mw.Protect(func(string) string { panic(io.EOF) })("x") }); r != nil {
		t.Fatalf("Protect(h)(\"x\") panicked: %v", r)
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("Protect(panic(io.EOF))(\"x\") error = %v, want errors.Is(err, io.EOF)", err)
	}
	if r := catch(func() { _, err = mw.Protect(first)("") }); r != nil {
		t.Fatalf("Protect(first)(\"\") panicked: %v", r)
	}
	var re runtime.Error
	if !errors.As(err, &re) {
		t.Errorf("Protect(first)(\"\") error = %v, want errors.As(err, runtime.Error)", err)
	}
}

func TestProtectRepanicsAbort(t *testing.T) {
	wrapped := fmt.Errorf("client gone: %w", mw.ErrAbort)
	for _, v := range []error{mw.ErrAbort, wrapped} {
		called := false
		h := func(string) string { called = true; panic(v) }
		r := catch(func() { _, _ = mw.Protect(h)("x") })
		if !called {
			t.Errorf("Protect(h)(\"x\") did not call h")
			continue
		}
		if r != v {
			t.Errorf("Protect(panic(%v))(\"x\") panicked with %v, want a re-panic with the same value %v", v, r, v)
		}
	}
}

func TestPipeline(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		steps []func(string) string
		want  string
		msg   string // "" means no error
	}{
		{"two steps", "go", []func(string) string{strings.ToUpper, reverse}, "OG", ""},
		{"no steps", "same", nil, "same", ""},
		{"panic in step 1", "x", []func(string) string{strings.ToUpper, boom, reverse}, "", `step 1: panic handling "X": boom`},
		{"four steps", "ab", []func(string) string{strings.TrimSpace, strings.ToLower, first, first}, "a", ""},
		{"empty input reaches first", " ", []func(string) string{strings.TrimSpace, first}, "", `step 1: panic handling "": runtime error: slice bounds out of range [:1] with length 0`},
	}
	for _, tt := range tests {
		var got string
		var err error
		if r := catch(func() { got, err = mw.Pipeline(tt.in, tt.steps) }); r != nil {
			t.Errorf("%s: Pipeline(%q) panicked: %v", tt.name, tt.in, r)
			continue
		}
		if tt.msg == "" {
			if got != tt.want || err != nil {
				t.Errorf("%s: Pipeline(%q) = %q, %v; want %q, <nil>", tt.name, tt.in, got, err, tt.want)
			}
			continue
		}
		if err == nil || err.Error() != tt.msg || got != "" {
			t.Errorf("%s: Pipeline(%q) = %q, %v; want \"\", %s", tt.name, tt.in, got, err, tt.msg)
		}
		var pe *mw.PanicError
		if err != nil && !errors.As(err, &pe) {
			t.Errorf("%s: Pipeline(%q) error = %v, want it to wrap a *PanicError", tt.name, tt.in, err)
		}
	}
}

func TestPipelineAbortPassesThrough(t *testing.T) {
	ran := false
	steps := []func(string) string{
		strings.ToUpper,
		func(string) string { panic(mw.ErrAbort) },
		func(s string) string { ran = true; return s },
	}
	r := catch(func() { _, _ = mw.Pipeline("x", steps) })
	if r != mw.ErrAbort {
		t.Errorf("Pipeline with an aborting step panicked with %v, want %v", r, mw.ErrAbort)
	}
	if ran {
		t.Errorf("Pipeline ran the step after the abort")
	}
}
