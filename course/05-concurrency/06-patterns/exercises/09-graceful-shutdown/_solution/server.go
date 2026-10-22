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
	s := &Server{
		jobs: make(chan int, max(queue, 0)),
		done: make(chan struct{}),
	}
	var wg sync.WaitGroup
	for range max(workers, 1) {
		wg.Go(func() {
			for j := range s.jobs {
				handle(j)
			}
		})
	}
	go func() {
		wg.Wait()
		close(s.done)
	}()
	return s
}

// Submit queues job without blocking. It returns ErrClosed once Shutdown has
// been called and ErrQueueFull if no worker or queue slot is free.
func (s *Server) Submit(job int) error {
	// The lock makes "check closed, then send" atomic with respect to
	// Shutdown's close(s.jobs), so a send can never hit a closed channel.
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	select {
	case s.jobs <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// Shutdown stops accepting jobs and waits until every accepted job has been
// handled and the workers have exited. If ctx ends first it returns
// ctx.Err(); the workers keep draining the queue in the background.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	if !s.closed {
		s.closed = true
		close(s.jobs) // workers finish what is queued, then their range ends
	}
	s.mu.Unlock()

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
