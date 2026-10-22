package ratelimit

import "time"

// Limiter allows at most n events per interval. Create one with New and
// call Stop when done.
type Limiter struct {
	tokens chan struct{} // buffer of n: one value = one allowed event
	stop   chan struct{} // closed by Stop
}

// New returns a full limiter that refills to n tokens every interval. It
// panics if n < 1 or interval <= 0.
func New(n int, interval time.Duration) *Limiter {
	// TODO
	return &Limiter{}
}

// Allow takes a token if one is free and reports whether it did. It never
// blocks.
func (l *Limiter) Allow() bool {
	// TODO
	return false
}

// Wait blocks until a token is free, then takes it.
func (l *Limiter) Wait() {
	// TODO
}

// Stop stops the refill goroutine. Call it once.
func (l *Limiter) Stop() {
	// TODO
}
