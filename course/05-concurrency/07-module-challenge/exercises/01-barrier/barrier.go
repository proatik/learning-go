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
	// TODO
	return &Barrier{}
}

// Wait blocks until n goroutines have called Wait in the current round, then
// returns that round's number.
func (b *Barrier) Wait() int {
	// TODO
	return 0
}
