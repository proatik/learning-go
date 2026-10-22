// Package api builds the production HTTP server configuration.
package api

import (
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/06-functional-options/server"
)

// NewProduction returns a server named "api" with a 10s timeout, 1000
// connections and the "recover" and "log" middleware. Options in extra are
// applied last, so they override those settings.
func NewProduction(addr string, extra ...server.Option) (*server.Server, error) {
	opts := []server.Option{
		server.WithName("api"),
		server.WithTimeout(10 * time.Second),
		server.WithMaxConns(1000),
		server.WithMiddleware("recover", "log"),
	}
	return server.New(addr, append(opts, extra...)...)
}
