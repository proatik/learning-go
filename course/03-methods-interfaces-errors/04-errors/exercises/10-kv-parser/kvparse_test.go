package kvparse

import (
	"errors"
	"maps"
	"testing"
)

func TestParseErrorMethods(t *testing.T) {
	e := &ParseError{Line: 2, Col: 5, Err: ErrMissingEquals}
	if got, want := e.Error(), "line 2, col 5: missing '='"; got != want {
		t.Errorf("(%+v).Error() = %q, want %q", *e, got, want)
	}
	if got := e.Unwrap(); got != ErrMissingEquals {
		t.Errorf("(%+v).Unwrap() = %v, want %v", *e, got, ErrMissingEquals)
	}
	if !errors.Is(e, ErrMissingEquals) {
		t.Errorf("errors.Is(%+v, ErrMissingEquals) = false, want true", *e)
	}
}

func TestParseOK(t *testing.T) {
	tests := []struct {
		in   string
		want map[string]string
	}{
		{"", map[string]string{}},
		{"a=1;b=2", map[string]string{"a": "1", "b": "2"}},
		{" host = go.dev ; port=443 ;", map[string]string{"host": "go.dev", "port": "443"}},
		{"a=1\n\nb=x=y", map[string]string{"a": "1", "b": "x=y"}},
		{"empty=", map[string]string{"empty": ""}},
		{";;\n  \n\t;", map[string]string{}},
		{"name=Zoë;city=東京", map[string]string{"name": "Zoë", "city": "東京"}},
		{"\tk.v-1_X\t=\t spaced value \t", map[string]string{"k.v-1_X": "spaced value"}},
		{"a=1;\nb=2;\n", map[string]string{"a": "1", "b": "2"}},
	}
	for _, tt := range tests {
		got, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q) error = %v, want <nil>", tt.in, err)
			continue
		}
		if got == nil || !maps.Equal(got, tt.want) {
			t.Errorf("Parse(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		in        string
		line, col int
		want      error
	}{
		{"oops", 1, 1, ErrMissingEquals},
		{"a=1;oops", 1, 5, ErrMissingEquals},
		{"a=1;  oops  ", 1, 7, ErrMissingEquals},
		{"name=Zoë;x", 1, 10, ErrMissingEquals},
		{"a=東京;b", 1, 6, ErrMissingEquals},
		{"\tx", 1, 2, ErrMissingEquals},
		{"a=1\nb=2\n\nnope", 4, 1, ErrMissingEquals},
		{"=1", 1, 1, ErrEmptyKey},
		{"a=1\n  =2", 2, 3, ErrEmptyKey},
		{"x=1; \t = 2", 1, 8, ErrEmptyKey},
		{"a=1\nnö=2", 2, 2, ErrInvalidKey},
		{"my key=1", 1, 3, ErrInvalidKey},
		{"a=1;  b$c = 2", 1, 8, ErrInvalidKey},
		{"é=1", 1, 1, ErrInvalidKey},
		{"a=1\nb=2;a=3", 2, 5, ErrDuplicateKey},
		{"a=1;  a=2", 1, 7, ErrDuplicateKey},
		{"a=1;a = 2;oops", 1, 5, ErrDuplicateKey}, // the first problem wins
		{"k=v\n\n\n  k = again", 4, 3, ErrDuplicateKey},
	}
	for _, tt := range tests {
		got, err := Parse(tt.in)
		want := &ParseError{Line: tt.line, Col: tt.col, Err: tt.want}
		pe, ok := errors.AsType[*ParseError](err)
		if !ok {
			t.Errorf("Parse(%q) error = %v, want %v", tt.in, err, want)
			continue
		}
		if *pe != *want {
			t.Errorf("Parse(%q) error = %+v, want %+v", tt.in, *pe, *want)
		}
		if !errors.Is(err, tt.want) {
			t.Errorf("Parse(%q) error = %v, want %v", tt.in, err, tt.want)
		}
		if err.Error() != want.Error() {
			t.Errorf("Parse(%q) error = %q, want %q", tt.in, err.Error(), want.Error())
		}
		if got != nil {
			t.Errorf("Parse(%q) = %v, want nil map with an error", tt.in, got)
		}
	}
}
