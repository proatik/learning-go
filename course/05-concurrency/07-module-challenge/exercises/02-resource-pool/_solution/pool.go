package pool

import (
	"context"
	"errors"
	"sync"
)

// ErrClosed is returned by Get once the pool is closed.
var ErrClosed = errors.New("pool: closed")

// Pool lends a fixed set of resources to many goroutines. Create one with New.
type Pool[T any] struct {
	idle      chan T        // resources not borrowed right now
	closed    chan struct{} // closed by Close
	closeOnce sync.Once
}

// New returns a pool holding all of resources, idle.
func New[T any](resources []T) *Pool[T] {
	p := &Pool[T]{
		idle:   make(chan T, len(resources)),
		closed: make(chan struct{}),
	}
	for _, r := range resources {
		p.idle <- r
	}
	return p
}

// Get borrows an idle resource, waiting until one is Put back, ctx is done
// (ctx.Err()) or the pool is closed (ErrClosed).
func (p *Pool[T]) Get(ctx context.Context) (T, error) {
	var zero T
	// select picks randomly among ready cases, so check closed first:
	// otherwise a closed pool with idle resources would still lend them.
	select {
	case <-p.closed:
		return zero, ErrClosed
	default:
	}
	select {
	case r := <-p.idle:
		return r, nil
	case <-ctx.Done():
		return zero, ctx.Err()
	case <-p.closed:
		return zero, ErrClosed
	}
}

// Put returns r to the pool. It does nothing after Close and panics if the
// pool is already full.
func (p *Pool[T]) Put(r T) {
	select {
	case <-p.closed:
		return
	default:
	}
	select {
	case p.idle <- r:
	default:
		panic("pool: Put into a full pool")
	}
}

// Idle returns the number of resources not currently borrowed.
func (p *Pool[T]) Idle() int {
	return len(p.idle)
}

// Close makes every waiting and future Get return ErrClosed. It is safe to
// call more than once.
func (p *Pool[T]) Close() {
	p.closeOnce.Do(func() { close(p.closed) })
}
