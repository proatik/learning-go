package fetch

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		e    HTTPError
		want string
	}{
		{HTTPError{"/users/7", 404}, "/users/7: status 404"},
		{HTTPError{"https://go.dev/x", 503}, "https://go.dev/x: status 503"},
		{HTTPError{"", 500}, ": status 500"},
	}
	for _, tt := range tests {
		e := tt.e
		if got := e.Error(); got != tt.want {
			t.Errorf("(%+v).Error() = %q, want %q", e, got, tt.want)
		}
	}
}

func TestIsCategories(t *testing.T) {
	tests := []struct {
		code                      int
		client, server, retryable bool
	}{
		{400, true, false, false},
		{404, true, false, false},
		{429, true, false, true},
		{499, true, false, false},
		{500, false, true, false},
		{501, false, true, false},
		{502, false, true, true},
		{503, false, true, true},
		{504, false, true, true},
		{599, false, true, false},
		{399, false, false, false},
		{600, false, false, false},
	}
	for _, tt := range tests {
		e := &HTTPError{URL: "/a", Code: tt.code}
		wrapped := fmt.Errorf("load: %w", e)
		for _, c := range []struct {
			target error
			want   bool
		}{{ErrClient, tt.client}, {ErrServer, tt.server}, {ErrRetryable, tt.retryable}, {ErrInvalidStatus, false}} {
			if got := e.Is(c.target); got != c.want {
				t.Errorf("(%v).Is(%v) = %t, want %t", e, c.target, got, c.want)
			}
			if got := errors.Is(wrapped, c.target); got != c.want {
				t.Errorf("errors.Is(%v, %v) = %t, want %t", wrapped, c.target, got, c.want)
			}
		}
	}
}

func TestIsSameCode(t *testing.T) {
	e := fmt.Errorf("ctx: %w", &HTTPError{URL: "/a", Code: 404})
	tests := []struct {
		target error
		want   bool
	}{
		{&HTTPError{URL: "/other", Code: 404}, true},
		{&HTTPError{Code: 404}, true},
		{&HTTPError{URL: "/a", Code: 500}, false},
		{(*HTTPError)(nil), false},
		{errors.New("client error"), false}, // same text as ErrClient, different value
		{nil, false},
	}
	for _, tt := range tests {
		if got := errors.Is(e, tt.target); got != tt.want {
			t.Errorf("errors.Is(%v, %#v) = %t, want %t", e, tt.target, got, tt.want)
		}
	}
}

func TestCheck(t *testing.T) {
	for _, code := range []int{100, 200, 204, 301, 399} {
		if err := Check("/a", code); err != nil {
			t.Errorf("Check(%q, %d) error = %v, want <nil>", "/a", code, err)
		}
	}
	for _, code := range []int{400, 404, 500, 599} {
		err := Check("/a", code)
		he, ok := errors.AsType[*HTTPError](err)
		if !ok {
			t.Errorf("Check(%q, %d) error = %v, want *HTTPError", "/a", code, err)
			continue
		}
		if *he != (HTTPError{URL: "/a", Code: code}) {
			t.Errorf("Check(%q, %d) error = %+v, want {URL:/a Code:%d}", "/a", code, *he, code)
		}
	}
	for _, code := range []int{0, 42, 99, 600, -404} {
		err := Check("/a", code)
		if !errors.Is(err, ErrInvalidStatus) {
			t.Errorf("Check(%q, %d) error = %v, want %v", "/a", code, err, ErrInvalidStatus)
			continue
		}
		if want := fmt.Sprintf("/a: status %d: invalid status code", code); err.Error() != want {
			t.Errorf("Check(%q, %d) error = %q, want %q", "/a", code, err.Error(), want)
		}
		if errors.Is(err, ErrClient) || errors.Is(err, ErrServer) {
			t.Errorf("Check(%q, %d) error = %v, should not match a category", "/a", code, err)
		}
	}
}

func TestAdvice(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{nil, "ok"},
		{Check("/a", 200), "ok"},
		{Check("/a", 429), "retry later"},
		{Check("/a", 503), "retry later"},
		{Check("/a", 404), "fix the request"},
		{fmt.Errorf("x: %w", Check("/a", 400)), "fix the request"},
		{Check("/a", 500), "server problem"},
		{fmt.Errorf("y: %w", fmt.Errorf("x: %w", Check("/a", 501))), "server problem"},
		{Check("/a", 42), "unknown error"},
		{errors.New("boom"), "unknown error"},
		{fmt.Errorf("x: %v", Check("/a", 404)), "unknown error"}, // %v cuts the chain
	}
	for _, tt := range tests {
		if got := Advice(tt.err); got != tt.want {
			t.Errorf("Advice(%v) = %q, want %q", tt.err, got, tt.want)
		}
	}
}
