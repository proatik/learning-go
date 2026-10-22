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
	s.hits.Add(1)
}

// Miss records one cache miss.
func (s *Stats) Miss() {
	s.misses.Add(1)
}

// Snapshot returns the current counts.
func (s *Stats) Snapshot() Counts {
	return Counts{Hits: s.hits.Load(), Misses: s.misses.Load()}
}

// Reset sets both counters to 0 and returns the counts they held.
func (s *Stats) Reset() Counts {
	return Counts{Hits: s.hits.Swap(0), Misses: s.misses.Swap(0)}
}

// Total returns Hits + Misses.
func (c Counts) Total() int64 {
	return c.Hits + c.Misses
}

// HitRate returns Hits / Total as a fraction from 0 to 1, or 0 if Total is 0.
func (c Counts) HitRate() float64 {
	if c.Total() == 0 {
		return 0
	}
	return float64(c.Hits) / float64(c.Total())
}
