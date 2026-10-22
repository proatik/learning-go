package group

import (
	"context"
	"sync"
)

// Group runs tasks in goroutines and collects the first error. The zero
// Group is valid and does not cancel anything.
type Group struct {
	wg     sync.WaitGroup
	once   sync.Once               // guards err: only the first error is kept
	err    error                   // the first non-nil error returned by a task
	cancel context.CancelCauseFunc // nil for the zero Group
}

// WithContext returns a new Group and a context derived from ctx. The context
// is cancelled, with the first error as its cause, as soon as a task fails,
// or when Wait returns, whichever happens first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	// TODO
	return &Group{}, ctx
}

// Go runs f in a new goroutine. If f returns an error and it is the group's
// first, it is recorded and the group's context is cancelled.
func (g *Group) Go(f func() error) {
	// TODO
}

// Wait blocks until every function started with Go has returned, cancels the
// group's context, and returns the first error (or nil).
func (g *Group) Wait() error {
	// TODO
	return nil
}
