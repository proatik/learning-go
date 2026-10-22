package server_test

import (
	"errors"
	"testing"
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/06-functional-options/server"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		addr string
		opts []server.Option
		want string
	}{
		{"defaults", ":8080", nil, "server@:8080 timeout=30s maxConns=100 middleware=[]"},
		{"name+timeout", ":80", []server.Option{server.WithName("web"), server.WithTimeout(5 * time.Second)}, "web@:80 timeout=5s maxConns=100 middleware=[]"},
		{"max conns", "localhost:9000", []server.Option{server.WithMaxConns(1)}, "server@localhost:9000 timeout=30s maxConns=1 middleware=[]"},
		{"middleware accumulates", ":80", []server.Option{server.WithMiddleware("log"), server.WithMiddleware("auth", "gzip")}, "server@:80 timeout=30s maxConns=100 middleware=[log auth gzip]"},
		{"later wins", ":80", []server.Option{server.WithTimeout(time.Second), server.WithTimeout(90 * time.Second)}, "server@:80 timeout=1m30s maxConns=100 middleware=[]"},
		{"empty middleware call", ":80", []server.Option{server.WithMiddleware()}, "server@:80 timeout=30s maxConns=100 middleware=[]"},
	}
	for _, tt := range tests {
		s, err := server.New(tt.addr, tt.opts...)
		if err != nil || s == nil {
			t.Errorf("New(%q) [%s] = %v, %v, want a server, <nil>", tt.addr, tt.name, s, err)
			continue
		}
		if got := s.String(); got != tt.want {
			t.Errorf("New(%q) [%s].String() = %q, want %q", tt.addr, tt.name, got, tt.want)
		}
	}
}

func TestGetters(t *testing.T) {
	s, err := server.New(":1", server.WithName("x"), server.WithTimeout(time.Minute), server.WithMaxConns(7), server.WithMiddleware("a", "b"))
	if err != nil || s == nil {
		t.Fatalf("New(%q, …) = %v, %v, want a server, <nil>", ":1", s, err)
	}
	if got := s.Addr(); got != ":1" {
		t.Errorf("Addr() = %q, want %q", got, ":1")
	}
	if got := s.Name(); got != "x" {
		t.Errorf("Name() = %q, want %q", got, "x")
	}
	if got := s.Timeout(); got != time.Minute {
		t.Errorf("Timeout() = %v, want %v", got, time.Minute)
	}
	if got := s.MaxConns(); got != 7 {
		t.Errorf("MaxConns() = %d, want %d", got, 7)
	}
	mw := s.Middleware()
	if len(mw) != 2 || mw[0] != "a" || mw[1] != "b" {
		t.Fatalf("Middleware() = %q, want %q", mw, []string{"a", "b"})
	}
	mw[0] = "changed"
	if got := s.Middleware()[0]; got != "a" {
		t.Errorf("Middleware()[0] after editing the returned slice = %q, want %q", got, "a")
	}
}

func TestDefaultsAreConstants(t *testing.T) {
	s, err := server.New(":1")
	if err != nil || s == nil {
		t.Fatalf("New(%q) = %v, %v, want a server, <nil>", ":1", s, err)
	}
	if s.Name() != server.DefaultName || s.Timeout() != server.DefaultTimeout || s.MaxConns() != server.DefaultMaxConns {
		t.Errorf("New(%q) = %q, %v, %d, want %q, %v, %d", ":1", s.Name(), s.Timeout(), s.MaxConns(),
			server.DefaultName, server.DefaultTimeout, server.DefaultMaxConns)
	}
	if got := s.Middleware(); len(got) != 0 {
		t.Errorf("default Middleware() = %q, want []", got)
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		opts    []server.Option
		want    []error
		notWant []error
	}{
		{"no addr", "", nil, []error{server.ErrNoAddr}, []error{server.ErrTimeout, server.ErrMaxConns}},
		{"zero timeout", ":1", []server.Option{server.WithTimeout(0)}, []error{server.ErrTimeout}, []error{server.ErrNoAddr, server.ErrMaxConns}},
		{"negative timeout", ":1", []server.Option{server.WithTimeout(-time.Second)}, []error{server.ErrTimeout}, nil},
		{"zero conns", ":1", []server.Option{server.WithMaxConns(0)}, []error{server.ErrMaxConns}, []error{server.ErrNoAddr, server.ErrTimeout}},
		{"everything", "", []server.Option{server.WithTimeout(0), server.WithMaxConns(-1)}, []error{server.ErrNoAddr, server.ErrTimeout, server.ErrMaxConns}, nil},
	}
	for _, tt := range tests {
		s, err := server.New(tt.addr, tt.opts...)
		if s != nil {
			t.Errorf("New(%q) [%s] server = %v, want nil", tt.addr, tt.name, s)
		}
		for _, w := range tt.want {
			if !errors.Is(err, w) {
				t.Errorf("New(%q) [%s] error = %v, want %v", tt.addr, tt.name, err, w)
			}
		}
		for _, nw := range tt.notWant {
			if errors.Is(err, nw) {
				t.Errorf("New(%q) [%s] error = %v, want no %v", tt.addr, tt.name, err, nw)
			}
		}
	}
}

func TestValidationOrder(t *testing.T) {
	_, err := server.New("", server.WithTimeout(0), server.WithMaxConns(0))
	want := "server: address is empty\nserver: timeout must be positive\nserver: maxConns must be at least 1"
	if err == nil || err.Error() != want {
		t.Errorf("New(%q, WithTimeout(0), WithMaxConns(0)) error = %v, want %q", "", err, want)
	}
}
