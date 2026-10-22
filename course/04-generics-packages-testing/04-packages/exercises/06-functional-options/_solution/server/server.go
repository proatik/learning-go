// Package server configures a network server with functional options.
package server

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// Defaults used when no option overrides them.
const (
	DefaultName     = "server"
	DefaultTimeout  = 30 * time.Second
	DefaultMaxConns = 100
)

// Validation errors returned (joined) by New.
var (
	ErrNoAddr   = errors.New("server: address is empty")
	ErrTimeout  = errors.New("server: timeout must be positive")
	ErrMaxConns = errors.New("server: maxConns must be at least 1")
)

// Server is a configured server. Create one with New.
type Server struct {
	addr       string
	name       string
	timeout    time.Duration
	maxConns   int
	middleware []string
}

// Option configures a Server; pass options to New.
type Option func(*Server)

// WithName sets the server's name.
func WithName(name string) Option {
	return func(s *Server) { s.name = name }
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

// WithMaxConns sets the maximum number of concurrent connections.
func WithMaxConns(n int) Option {
	return func(s *Server) { s.maxConns = n }
}

// WithMiddleware appends middleware names; repeated calls accumulate.
func WithMiddleware(names ...string) Option {
	return func(s *Server) { s.middleware = append(s.middleware, names...) }
}

// New returns a server listening on addr, configured by the defaults and
// then by opts in order. It reports every invalid setting with errors.Join.
func New(addr string, opts ...Option) (*Server, error) {
	s := &Server{
		addr:     addr,
		name:     DefaultName,
		timeout:  DefaultTimeout,
		maxConns: DefaultMaxConns,
	}
	for _, opt := range opts {
		opt(s)
	}

	var errs []error
	if s.addr == "" {
		errs = append(errs, ErrNoAddr)
	}
	if s.timeout <= 0 {
		errs = append(errs, ErrTimeout)
	}
	if s.maxConns < 1 {
		errs = append(errs, ErrMaxConns)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return s, nil
}

// Addr returns the listen address.
func (s *Server) Addr() string { return s.addr }

// Name returns the server's name.
func (s *Server) Name() string { return s.name }

// Timeout returns the request timeout.
func (s *Server) Timeout() time.Duration { return s.timeout }

// MaxConns returns the connection limit.
func (s *Server) MaxConns() int { return s.maxConns }

// Middleware returns a copy of the middleware names, in order.
func (s *Server) Middleware() []string { return slices.Clone(s.middleware) }

// String describes the server as
// "<name>@<addr> timeout=<timeout> maxConns=<n> middleware=<list>".
func (s *Server) String() string {
	return fmt.Sprintf("%s@%s timeout=%v maxConns=%d middleware=%v",
		s.name, s.addr, s.timeout, s.maxConns, s.middleware)
}
