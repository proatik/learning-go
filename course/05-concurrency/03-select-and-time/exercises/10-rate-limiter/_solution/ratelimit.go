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
	if n < 1 || interval <= 0 {
		panic("ratelimit: n and interval must be positive")
	}
	l := &Limiter{
		tokens: make(chan struct{}, n),
		stop:   make(chan struct{}),
	}
	l.fill()
	go l.refill(interval)
	return l
}

// refill tops the bucket up on every tick until Stop is called.
func (l *Limiter) refill(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			l.fill()
		case <-l.stop:
			return
		}
	}
}

// fill tops the bucket up to n tokens. It counts the missing tokens first:
// "send until full" would hand a token straight to every waiting goroutine
// and let more than n through.
func (l *Limiter) fill() {
	for range cap(l.tokens) - len(l.tokens) {
		select {
		case l.tokens <- struct{}{}:
		default:
			return
		}
	}
}

// Allow takes a token if one is free and reports whether it did. It never
// blocks.
func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until a token is free, then takes it.
func (l *Limiter) Wait() {
	<-l.tokens
}

// Stop stops the refill goroutine. Call it once.
func (l *Limiter) Stop() {
	close(l.stop)
}
