package ratelimit

import (
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestNewPanics(t *testing.T) {
	const want = "ratelimit: n and interval must be positive"
	cases := []struct {
		n        int
		interval time.Duration
	}{{0, time.Second}, {-1, time.Second}, {1, 0}, {1, -time.Second}}
	for _, c := range cases {
		func() {
			defer func() {
				if r := recover(); r != want {
					t.Errorf("New(%d, %v) panic = %v, want %q", c.n, c.interval, r, want)
				}
			}()
			l := New(c.n, c.interval)
			if l != nil && l.stop != nil {
				l.Stop()
			}
		}()
	}
}

// allowN calls l.Allow n times and returns the results.
func allowN(l *Limiter, n int) []bool {
	var got []bool
	for range n {
		got = append(got, l.Allow())
	}
	return got
}

func TestAllow(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := New(3, time.Second)
		defer l.Stop()
		start := time.Now()
		steps := []struct {
			sleep time.Duration
			want  []bool
		}{
			{0, []bool{true, true, true, false}},
			{time.Second, []bool{true, true, true, false}},               // 1s: refilled to 3
			{3 * time.Second, []bool{true, true, true, false}},           // 4s: unused windows don't pile up
			{500 * time.Millisecond, []bool{false, false, false, false}}, // 4.5s: nothing until 5s
			{500 * time.Millisecond, []bool{true, true, true, false}},    // 5s
		}
		for _, s := range steps {
			time.Sleep(s.sleep)
			synctest.Wait() // let the refill goroutine handle a tick due now
			if got := allowN(l, 4); !slices.Equal(got, s.want) {
				t.Fatalf("New(3, 1s), at %v: Allow() ×4 = %v, want %v", time.Since(start), got, s.want)
			}
		}
	})
}

func TestAllowPartialUse(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := New(3, time.Second)
		defer l.Stop()
		if !l.Allow() {
			t.Fatalf("New(3, 1s): first Allow() = false, want true")
		}
		time.Sleep(time.Second)
		synctest.Wait()
		want := []bool{true, true, true, false}
		if got := allowN(l, 4); !slices.Equal(got, want) {
			t.Errorf("New(3, 1s), 1 token used at 0s: at 1s Allow() ×4 = %v, want %v (refill to exactly 3)", got, want)
		}
	})
}

// waitAll runs f in a goroutine and fails if it hasn't finished within one
// (fake) minute.
func waitAll(t *testing.T, call string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(time.Minute):
		t.Fatalf("%s did not finish within a minute", call)
	}
}

func TestWaitSequential(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := New(3, time.Second)
		defer l.Stop()
		start := time.Now()
		var got []time.Duration
		waitAll(t, "New(3, 1s): Wait() ×7", func() {
			for range 7 {
				l.Wait()
				got = append(got, time.Since(start))
			}
		})
		s := time.Second
		want := []time.Duration{0, 0, 0, s, s, s, 2 * s}
		if !slices.Equal(got, want) {
			t.Errorf("New(3, 1s): Wait() ×7 returned at %v, want %v", got, want)
		}
	})
}

func TestWaitConcurrent(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := New(4, time.Second)
		defer l.Stop()
		start := time.Now()
		times := make(chan time.Duration, 10)
		waitAll(t, "New(4, 1s): 5 goroutines × 2 Wait()", func() {
			var wg sync.WaitGroup
			for range 5 {
				wg.Go(func() {
					for range 2 {
						l.Wait()
						times <- time.Since(start)
					}
				})
			}
			wg.Wait()
		})
		close(times)
		var got []time.Duration
		for d := range times {
			got = append(got, d)
		}
		slices.Sort(got)
		s := time.Second
		want := []time.Duration{0, 0, 0, 0, s, s, s, s, 2 * s, 2 * s}
		if !slices.Equal(got, want) {
			t.Errorf("New(4, 1s): 5 goroutines × 2 Wait() returned at %v, want %v", got, want)
		}
	})
}

func TestAllowManyGoroutines(t *testing.T) {
	l := New(10, time.Hour)
	defer l.Stop()
	counts := make(chan int, 50)
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			n := 0
			for range 100 {
				if l.Allow() {
					n++
				}
			}
			counts <- n
		})
	}
	wg.Wait()
	close(counts)
	total := 0
	for n := range counts {
		total += n
	}
	if total != 10 {
		t.Errorf("New(10, 1h): 50 goroutines × 100 Allow() = %d allowed, want 10", total)
	}
}

func TestStopEndsGoroutine(t *testing.T) {
	for _, n := range []int{1, 5} {
		synctest.Test(t, func(t *testing.T) {
			l := New(n, time.Second)
			time.Sleep(2500 * time.Millisecond)
			l.Stop()
			// If the refill goroutine is still running when this function
			// returns, synctest fails with "blocked goroutines remain".
		})
	}
}
