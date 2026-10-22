package server

import (
	"context"
	"errors"
	"sync"
)

// Errors returned by Submit.
var (
	ErrClosed    = errors.New("server: closed")
	ErrQueueFull = errors.New("server: queue full")
)

// Server runs submitted jobs on a fixed set of workers. Create one with New.
type Server struct {
	mu     sync.Mutex
	jobs   chan int      // the queue; closed by Shutdown
	closed bool          // set by Shutdown; guarded by mu
	done   chan struct{} // closed once every worker has exited
}

// New starts workers goroutines (at least 1) that call handle for each
// submitted job. Up to queue jobs (at least 0) can wait for a free worker.
func New(workers, queue int, handle func(job int)) *Server {
	// TODO
	return &Server{}
}

// Submit queues job without blocking. It returns ErrClosed once Shutdown has
// been called and ErrQueueFull if no worker or queue slot is free.
func (s *Server) Submit(job int) error {
	// TODO
	return nil
}

// Shutdown stops accepting jobs and waits until every accepted job has been
// handled and the workers have exited. If ctx ends first it returns
// ctx.Err(); the workers keep draining the queue in the background.
func (s *Server) Shutdown(ctx context.Context) error {
	// TODO
	return nil
}
