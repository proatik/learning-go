package breaker

import (
	"errors"
	"sync"
	"time"
)

// State is the position of a circuit breaker.
type State int

const (
	Closed   State = iota // calls go through
	Open                  // calls are rejected with ErrOpen
	HalfOpen              // one trial call is allowed
)

// String returns "closed", "open" or "half-open".
func (s State) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	}
	return "unknown"
}

// ErrOpen is returned by Call when the breaker rejects the call.
var ErrOpen = errors.New("breaker: open")

// Breaker stops calling a failing function until a cooldown has passed.
// Create one with New.
type Breaker struct {
	mu        sync.Mutex
	threshold int
	cooldown  time.Duration
	state     State
	failures  int       // consecutive failures while Closed
	openedAt  time.Time // when the breaker last opened
	probing   bool      // a HalfOpen trial call is running
}

// New returns a closed breaker that opens after threshold consecutive
// failures and allows a trial call after cooldown.
func New(threshold int, cooldown time.Duration) *Breaker {
	// TODO
	return &Breaker{}
}

// Call runs fn unless the breaker is open, and returns fn's error. It
// returns ErrOpen without calling fn when the breaker rejects the call.
func (b *Breaker) Call(fn func() error) error {
	// TODO
	return nil
}

// State returns the breaker's current state.
func (b *Breaker) State() State {
	// TODO
	return Closed
}
