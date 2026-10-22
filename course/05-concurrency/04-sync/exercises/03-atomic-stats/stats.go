package stats

import "sync/atomic"

// Stats counts cache hits and misses. It is safe for concurrent use without
// a mutex. The zero value is ready to use.
type Stats struct {
	hits   atomic.Int64
	misses atomic.Int64
}

// Counts is a plain copy of the counters at one moment.
type Counts struct {
	Hits   int64
	Misses int64
}

// Hit records one cache hit.
func (s *Stats) Hit() {
	// TODO
}

// Miss records one cache miss.
func (s *Stats) Miss() {
	// TODO
}

// Snapshot returns the current counts.
func (s *Stats) Snapshot() Counts {
	// TODO
	return Counts{}
}

// Reset sets both counters to 0 and returns the counts they held.
func (s *Stats) Reset() Counts {
	// TODO
	return Counts{}
}

// Total returns Hits + Misses.
func (c Counts) Total() int64 {
	// TODO
	return 0
}

// HitRate returns Hits / Total as a fraction from 0 to 1, or 0 if Total is 0.
func (c Counts) HitRate() float64 {
	// TODO
	return 0
}
