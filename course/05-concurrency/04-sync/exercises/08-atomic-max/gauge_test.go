package gauge

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStoreMax(t *testing.T) {
	tests := []struct {
		start, v  int64
		want      bool
		wantValue int64
	}{
		{0, 5, true, 5},
		{5, 3, false, 5},
		{5, 5, false, 5},
		{5, 6, true, 6},
		{0, -1, false, 0},
		{-10, -3, true, -3},
		{math.MinInt64, math.MinInt64 + 1, true, math.MinInt64 + 1},
		{math.MaxInt64 - 1, math.MaxInt64, true, math.MaxInt64},
		{math.MaxInt64, 0, false, math.MaxInt64},
	}
	for _, tt := range tests {
		var m atomic.Int64
		m.Store(tt.start)
		got := StoreMax(&m, tt.v)
		if got != tt.want {
			t.Errorf("StoreMax(%d, %d) = %v, want %v", tt.start, tt.v, got, tt.want)
		}
		if m.Load() != tt.wantValue {
			t.Errorf("StoreMax(%d, %d): value = %d, want %d", tt.start, tt.v, m.Load(), tt.wantValue)
		}
	}
}

func TestStoreMaxConcurrent(t *testing.T) {
	for round := range 5 {
		var m atomic.Int64
		m.Store(-1)
		stop := make(chan struct{})
		var decreased atomic.Bool
		var watcher sync.WaitGroup
		watcher.Go(func() { // the max must never go down
			prev := int64(-1)
			for {
				cur := m.Load()
				if cur < prev {
					decreased.Store(true)
				}
				prev = cur
				select {
				case <-stop:
					return
				default:
				}
			}
		})
		var wins atomic.Int64
		var wg sync.WaitGroup
		for g := range 50 {
			wg.Go(func() {
				for i := range 1000 {
					if StoreMax(&m, int64(i*50+g)) {
						wins.Add(1)
					}
				}
			})
		}
		wg.Wait()
		close(stop)
		watcher.Wait()
		if got := m.Load(); got != 49999 {
			t.Fatalf("round %d: 50 goroutines StoreMax values 0..49999: value = %d, want 49999", round, got)
		}
		if decreased.Load() {
			t.Fatalf("round %d: the stored max went down during concurrent StoreMax calls, want it to only go up", round)
		}
		if w := wins.Load(); w < 1 || w > 50000 {
			t.Fatalf("round %d: StoreMax returned true %d times, want between 1 and 50000", round, w)
		}
	}
}

func TestGaugeSequential(t *testing.T) {
	var g Gauge
	steps := []struct {
		call      string
		got, want int64
		peak      int64
	}{
		{"Inc()", g.Inc(), 1, 1},
		{"Inc()", g.Inc(), 2, 2},
		{"Inc()", g.Inc(), 3, 3},
		{"Dec()", g.Dec(), 2, 3},
		{"Inc()", g.Inc(), 3, 3},
		{"Dec()", g.Dec(), 2, 3},
		{"Dec()", g.Dec(), 1, 3},
		{"Dec()", g.Dec(), 0, 3},
		{"Dec()", g.Dec(), -1, 3},
		{"Inc()", g.Inc(), 0, 3},
	}
	var g2 Gauge
	for i, s := range steps {
		if s.got != s.want {
			t.Errorf("step %d: %s = %d, want %d", i+1, s.call, s.got, s.want)
		}
		// replay on g2 to check Current and Peak after each step
		if s.call == "Inc()" {
			g2.Inc()
		} else {
			g2.Dec()
		}
		if c := g2.Current(); c != s.want {
			t.Errorf("step %d: after %s, Current() = %d, want %d", i+1, s.call, c, s.want)
		}
		if p := g2.Peak(); p != s.peak {
			t.Errorf("step %d: after %s, Peak() = %d, want %d", i+1, s.call, p, s.peak)
		}
	}
	var z Gauge
	if z.Current() != 0 || z.Peak() != 0 {
		t.Errorf("Gauge{}: Current(), Peak() = %d, %d, want 0, 0", z.Current(), z.Peak())
	}
}

// overlap makes n goroutines call Inc, wait until all n are "in progress",
// then call Dec. It returns the values Inc returned.
func overlap(t *testing.T, g *Gauge, n int) []int64 {
	t.Helper()
	got := make([]int64, n)
	var ready, finished sync.WaitGroup
	ready.Add(n)
	release := make(chan struct{})
	for i := range n {
		finished.Go(func() {
			got[i] = g.Inc()
			ready.Done()
			<-release
			g.Dec()
		})
	}
	allReady := make(chan struct{})
	go func() {
		ready.Wait()
		close(allReady)
	}()
	select {
	case <-allReady:
	case <-time.After(2 * time.Second):
		t.Fatalf("%d goroutines did not all call Inc() within 2s", n)
	}
	if c := g.Current(); c != int64(n) {
		t.Errorf("%d goroutines in progress: Current() = %d, want %d", n, c, n)
	}
	close(release)
	finished.Wait()
	return got
}

func TestGaugePeakUnderLoad(t *testing.T) {
	for round := range 20 {
		var g Gauge
		got := overlap(t, &g, 50)
		var seen [51]bool
		for _, v := range got {
			if v < 1 || v > 50 || seen[v] {
				t.Fatalf("round %d: 50 overlapping Inc() returned %v, want each of 1..50 exactly once", round, got)
			}
			seen[v] = true
		}
		if p := g.Peak(); p != 50 {
			t.Fatalf("round %d: 50 overlapping Inc()/Dec(): Peak() = %d, want 50", round, p)
		}
		if c := g.Current(); c != 0 {
			t.Fatalf("round %d: after every Dec(): Current() = %d, want 0", round, c)
		}
		overlap(t, &g, 10)
		if p := g.Peak(); p != 50 {
			t.Fatalf("round %d: a later burst of 10: Peak() = %d, want 50 (the peak never goes down)", round, p)
		}
	}
}
