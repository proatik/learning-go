package api_test

import (
	"errors"
	"testing"
	"time"

	api "learngo/course/04-generics-packages-testing/04-packages/exercises/06-functional-options"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/06-functional-options/server"
)

func TestNewProduction(t *testing.T) {
	tests := []struct {
		name  string
		addr  string
		extra []server.Option
		want  string
	}{
		{"defaults", ":443", nil, "api@:443 timeout=10s maxConns=1000 middleware=[recover log]"},
		{"override conns", ":443", []server.Option{server.WithMaxConns(5)}, "api@:443 timeout=10s maxConns=5 middleware=[recover log]"},
		{"override name+timeout", ":8443", []server.Option{server.WithName("admin"), server.WithTimeout(2 * time.Second)}, "admin@:8443 timeout=2s maxConns=1000 middleware=[recover log]"},
		{"extra middleware", ":443", []server.Option{server.WithMiddleware("auth")}, "api@:443 timeout=10s maxConns=1000 middleware=[recover log auth]"},
	}
	for _, tt := range tests {
		s, err := api.NewProduction(tt.addr, tt.extra...)
		if err != nil || s == nil {
			t.Errorf("NewProduction(%q) [%s] = %v, %v, want a server, <nil>", tt.addr, tt.name, s, err)
			continue
		}
		if got := s.String(); got != tt.want {
			t.Errorf("NewProduction(%q) [%s].String() = %q, want %q", tt.addr, tt.name, got, tt.want)
		}
	}
}

func TestNewProductionErrors(t *testing.T) {
	if s, err := api.NewProduction(""); s != nil || !errors.Is(err, server.ErrNoAddr) {
		t.Errorf("NewProduction(%q) = %v, %v, want nil, %v", "", s, err, server.ErrNoAddr)
	}
	if s, err := api.NewProduction(":1", server.WithMaxConns(0)); s != nil || !errors.Is(err, server.ErrMaxConns) {
		t.Errorf("NewProduction(%q, WithMaxConns(0)) = %v, %v, want nil, %v", ":1", s, err, server.ErrMaxConns)
	}
}
