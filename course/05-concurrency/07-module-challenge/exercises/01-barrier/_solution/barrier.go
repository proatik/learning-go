package barrier

import "sync"

// Barrier makes n goroutines wait for each other, round after round.
// Create one with New.
type Barrier struct {
	mu      sync.Mutex
	n       int           // goroutines per round
	count   int           // arrived so far in this round
	round   int           // number of completed rounds
	release chan struct{} // closed to release the current round
}

// New returns a barrier for n goroutines. It panics if n < 1.
func New(n int) *Barrier {
	if n < 1 {
		panic("barrier: n must be at least 1")
	}
	return &Barrier{n: n, release: make(chan struct{})}
}

// Wait blocks until n goroutines have called Wait in the current round, then
// returns that round's number.
func (b *Barrier) Wait() int {
	b.mu.Lock()
	round, release := b.round, b.release
	b.count++
	if b.count == b.n {
		// Last to arrive: start a fresh round before waking this one, so
		// early leavers that call Wait again join the next round.
		b.count = 0
		b.round++
		b.release = make(chan struct{})
		b.mu.Unlock()
		close(release)
		return round
	}
	b.mu.Unlock()
	<-release
	return round
}
