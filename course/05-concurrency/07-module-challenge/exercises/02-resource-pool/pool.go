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
	// TODO
	return &Pool[T]{}
}

// Get borrows an idle resource, waiting until one is Put back, ctx is done
// (ctx.Err()) or the pool is closed (ErrClosed).
func (p *Pool[T]) Get(ctx context.Context) (T, error) {
	// TODO
	var zero T
	return zero, nil
}

// Put returns r to the pool. It does nothing after Close and panics if the
// pool is already full.
func (p *Pool[T]) Put(r T) {
	// TODO
}

// Idle returns the number of resources not currently borrowed.
func (p *Pool[T]) Idle() int {
	// TODO
	return 0
}

// Close makes every waiting and future Get return ErrClosed. It is safe to
// call more than once.
func (p *Pool[T]) Close() {
	// TODO
}
