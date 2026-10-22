package stats

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSequential(t *testing.T) {
	var s Stats
	if got := s.Snapshot(); got != (Counts{}) {
		t.Errorf("Stats{}.Snapshot() = %+v, want %+v", got, Counts{})
	}
	s.Hit()
	s.Hit()
	s.Hit()
	s.Miss()
	if got, want := s.Snapshot(), (Counts{Hits: 3, Misses: 1}); got != want {
		t.Errorf("3 Hit(), 1 Miss(): Snapshot() = %+v, want %+v", got, want)
	}
	if got, want := s.Reset(), (Counts{Hits: 3, Misses: 1}); got != want {
		t.Errorf("Reset() = %+v, want %+v", got, want)
	}
	if got := s.Snapshot(); got != (Counts{}) {
		t.Errorf("after Reset(): Snapshot() = %+v, want %+v", got, Counts{})
	}
	if got := s.Reset(); got != (Counts{}) {
		t.Errorf("second Reset() = %+v, want %+v", got, Counts{})
	}
	s.Miss()
	if got, want := s.Snapshot(), (Counts{Misses: 1}); got != want {
		t.Errorf("after Reset(), Miss(): Snapshot() = %+v, want %+v", got, want)
	}
}

func TestCounts(t *testing.T) {
	tests := []struct {
		in    Counts
		total int64
		rate  float64
	}{
		{Counts{}, 0, 0},
		{Counts{Hits: 3, Misses: 1}, 4, 0.75},
		{Counts{Hits: 5}, 5, 1},
		{Counts{Misses: 7}, 7, 0},
		{Counts{Hits: 1, Misses: 2}, 3, 1.0 / 3},
	}
	for _, tt := range tests {
		if got := tt.in.Total(); got != tt.total {
			t.Errorf("%+v.Total() = %d, want %d", tt.in, got, tt.total)
		}
		if got := tt.in.HitRate(); got != tt.rate {
			t.Errorf("%+v.HitRate() = %v, want %v", tt.in, got, tt.rate)
		}
	}
}

func TestConcurrent(t *testing.T) {
	var s Stats
	var wg sync.WaitGroup
	for g := range 50 {
		wg.Go(func() {
			for i := range 1000 {
				if (g+i)%4 == 0 {
					s.Miss()
				} else {
					s.Hit()
				}
				_ = s.Snapshot()
			}
		})
	}
	wg.Wait()
	if got, want := s.Snapshot(), (Counts{Hits: 37500, Misses: 12500}); got != want {
		t.Errorf("50 goroutines × 1000 Hit()/Miss(): Snapshot() = %+v, want %+v", got, want)
	}
}

func TestResetLosesNothing(t *testing.T) {
	var s Stats
	var hits, misses atomic.Int64 // the sum of every Reset() result
	stop := make(chan struct{})
	var resetter sync.WaitGroup
	resetter.Go(func() {
		for {
			c := s.Reset()
			hits.Add(c.Hits)
			misses.Add(c.Misses)
			select {
			case <-stop:
				return
			default:
			}
		}
	})
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			for range 2000 {
				s.Hit()
				s.Miss()
			}
		})
	}
	wg.Wait()
	close(stop)
	resetter.Wait()
	last := s.Reset()
	got := Counts{Hits: hits.Load() + last.Hits, Misses: misses.Load() + last.Misses}
	if want := (Counts{Hits: 40000, Misses: 40000}); got != want {
		t.Errorf("sum of all Reset() results during 40000 Hit() + 40000 Miss() = %+v, want %+v", got, want)
	}
}
