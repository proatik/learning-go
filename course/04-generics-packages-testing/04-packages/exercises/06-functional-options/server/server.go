// Package server configures a network server with functional options.
package server

import (
	"errors"
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
	// TODO: return a func(*Server) that sets s.name
	return func(*Server) {}
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) Option {
	// TODO
	return func(*Server) {}
}

// WithMaxConns sets the maximum number of concurrent connections.
func WithMaxConns(n int) Option {
	// TODO
	return func(*Server) {}
}

// WithMiddleware appends middleware names; repeated calls accumulate.
func WithMiddleware(names ...string) Option {
	// TODO
	return func(*Server) {}
}

// New returns a server listening on addr, configured by the defaults and
// then by opts in order. It reports every invalid setting with errors.Join.
func New(addr string, opts ...Option) (*Server, error) {
	// TODO: defaults → apply options → validate
	return nil, nil
}

// Addr returns the listen address.
func (s *Server) Addr() string {
	// TODO
	return ""
}

// Name returns the server's name.
func (s *Server) Name() string {
	// TODO
	return ""
}

// Timeout returns the request timeout.
func (s *Server) Timeout() time.Duration {
	// TODO
	return 0
}

// MaxConns returns the connection limit.
func (s *Server) MaxConns() int {
	// TODO
	return 0
}

// Middleware returns a copy of the middleware names, in order.
func (s *Server) Middleware() []string {
	// TODO
	return nil
}

// String describes the server as
// "<name>@<addr> timeout=<timeout> maxConns=<n> middleware=<list>".
func (s *Server) String() string {
	// TODO
	return ""
}
