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
	return &Breaker{threshold: max(threshold, 1), cooldown: cooldown}
}

// Call runs fn unless the breaker is open, and returns fn's error. It
// returns ErrOpen without calling fn when the breaker rejects the call.
func (b *Breaker) Call(fn func() error) error {
	trial, ok := b.allow()
	if !ok {
		return ErrOpen
	}
	err := fn() // without the lock: closed calls run concurrently
	b.record(trial, err)
	return err
}

// allow decides whether a call may run and whether it is the trial call.
func (b *Breaker) allow() (trial, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case Closed:
		return false, true
	case Open:
		if time.Since(b.openedAt) < b.cooldown {
			return false, false
		}
		b.state = HalfOpen
	}
	if b.probing {
		return false, false
	}
	b.probing = true
	return true, true
}

// record updates the state with the outcome of a call.
func (b *Breaker) record(trial bool, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if trial {
		b.probing = false
		if err != nil {
			b.state, b.openedAt = Open, time.Now()
			return
		}
		b.state, b.failures = Closed, 0
		return
	}
	if b.state != Closed {
		return // a slow call that started before the breaker opened
	}
	if err == nil {
		b.failures = 0
		return
	}
	b.failures++
	if b.failures >= b.threshold {
		b.state, b.openedAt, b.failures = Open, time.Now(), 0
	}
}

// State returns the breaker's current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
