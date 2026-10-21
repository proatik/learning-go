package checks

import (
	"errors"
	"testing"
)

// Failures are printed with %#v: it shows the dynamic type, so a typed nil
// appears as (*checks.CheckError)(nil).

func TestCheckHost(t *testing.T) {
	for _, h := range []string{"go.dev", "localhost", "127.0.0.1", "日本.jp"} {
		if err := CheckHost(h); err != nil {
			t.Errorf("CheckHost(%q) = %#v, want <nil>", h, err)
		}
	}
	for _, h := range []string{"", "   ", "\t\n"} {
		ce, ok := errors.AsType[*CheckError](CheckHost(h))
		if !ok || ce == nil {
			t.Errorf("CheckHost(%q) error = %#v, want &CheckError{Name: \"host\", ...}", h, ce)
			continue
		}
		if ce.Name != "host" || ce.Problem != "must not be empty" {
			t.Errorf("CheckHost(%q) error = %+v, want {Name:host Problem:must not be empty}", h, *ce)
		}
	}
}

func TestCheckPort(t *testing.T) {
	for _, p := range []int{1, 80, 8080, 65535} {
		if err := CheckPort(p); err != nil {
			t.Errorf("CheckPort(%d) = %#v, want <nil>", p, err)
		}
	}
	for _, p := range []int{0, -1, 65536, 1 << 20} {
		ce, ok := errors.AsType[*CheckError](CheckPort(p))
		if !ok || ce == nil {
			t.Errorf("CheckPort(%d) error = %#v, want &CheckError{Name: \"port\", ...}", p, ce)
			continue
		}
		if got, want := ce.Error(), "port: must be 1-65535"; got != want {
			t.Errorf("CheckPort(%d) error = %q, want %q", p, got, want)
		}
	}
}

func TestCheckAll(t *testing.T) {
	if err := CheckAll("go.dev", 443); err != nil {
		t.Errorf("CheckAll(%q, %d) = %#v, want <nil>", "go.dev", 443, err)
	}

	tests := []struct {
		host string
		port int
		msg  string
	}{
		{"", 443, "host: must not be empty"},
		{"go.dev", 0, "port: must be 1-65535"},
		{"", 0, "host: must not be empty\nport: must be 1-65535"},
	}
	for _, tt := range tests {
		err := CheckAll(tt.host, tt.port)
		if err == nil {
			t.Errorf("CheckAll(%q, %d) = <nil>, want %q", tt.host, tt.port, tt.msg)
			continue
		}
		joined, ok := err.(interface{ Unwrap() []error })
		if !ok {
			t.Errorf("CheckAll(%q, %d) = %#v, want an errors.Join result", tt.host, tt.port, err)
			continue
		}
		for _, e := range joined.Unwrap() {
			if ce, ok := e.(*CheckError); ok && ce == nil {
				t.Errorf("CheckAll(%q, %d) contains a typed nil %#v", tt.host, tt.port, e)
			}
		}
		if t.Failed() {
			continue
		}
		if got := err.Error(); got != tt.msg {
			t.Errorf("CheckAll(%q, %d) error = %q, want %q", tt.host, tt.port, got, tt.msg)
		}
	}
}

func TestFirstFailure(t *testing.T) {
	e1 := &CheckError{Name: "a", Problem: "x"}
	e2 := &CheckError{Name: "b", Problem: "y"}
	tests := []struct {
		in   []*CheckError
		want *CheckError
	}{
		{[]*CheckError{nil, e1, nil}, e1},
		{[]*CheckError{e2, e1}, e2},
		{[]*CheckError{nil, nil, e2}, e2},
	}
	for _, tt := range tests {
		got := FirstFailure(tt.in)
		if ce, ok := got.(*CheckError); !ok || ce != tt.want {
			t.Errorf("FirstFailure(%v) = %#v, want %#v", tt.in, got, tt.want)
		}
	}
	for _, in := range [][]*CheckError{nil, {}, {nil}, {nil, nil}} {
		if got := FirstFailure(in); got != nil {
			t.Errorf("FirstFailure(%v) = %#v, want <nil>", in, got)
		}
	}
}
