package gauge

import "sync/atomic"

// StoreMax sets m to v if v is greater than m's current value, and reports
// whether it did. It is safe to call from many goroutines at once.
func StoreMax(m *atomic.Int64, v int64) bool {
	for {
		old := m.Load()
		if v <= old {
			return false
		}
		if m.CompareAndSwap(old, v) {
			return true
		}
		// Another goroutine changed m between Load and CompareAndSwap: retry.
	}
}

// Gauge tracks how many things are in progress right now, and the highest
// that number has ever been. It is safe for concurrent use without a mutex.
// The zero value is a gauge at 0 with a peak of 0.
type Gauge struct {
	current atomic.Int64
	peak    atomic.Int64
}

// Inc adds 1 to the current value, updates the peak, and returns the new
// current value.
func (g *Gauge) Inc() int64 {
	n := g.current.Add(1)
	StoreMax(&g.peak, n)
	return n
}

// Dec subtracts 1 from the current value and returns the new current value.
func (g *Gauge) Dec() int64 {
	return g.current.Add(-1)
}

// Current returns the current value.
func (g *Gauge) Current() int64 {
	return g.current.Load()
}

// Peak returns the highest current value seen so far.
func (g *Gauge) Peak() int64 {
	return g.peak.Load()
}
