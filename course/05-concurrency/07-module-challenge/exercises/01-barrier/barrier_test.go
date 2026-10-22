package barrier

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestNewPanics(t *testing.T) {
	for _, n := range []int{0, -1} {
		func() {
			defer func() {
				r := recover()
				if r != "barrier: n must be at least 1" {
					t.Errorf("New(%d) panic = %v, want %q", n, r, "barrier: n must be at least 1")
				}
			}()
			New(n)
		}()
	}
}

func TestOneGoroutine(t *testing.T) {
	b := New(1)
	for want := range 3 {
		if got := b.Wait(); got != want {
			t.Fatalf("New(1): Wait() call %d = %d, want %d", want+1, got, want)
		}
	}
}

func TestBlocksUntilAllArrive(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(3)
		rounds := make(chan int, 3)
		for range 2 {
			go func() { rounds <- b.Wait() }()
		}
		synctest.Wait()
		if got := len(rounds); got != 0 {
			t.Errorf("New(3), 2 goroutines call Wait: returned = %d, want 0 (they must block)", got)
		}
		go func() { rounds <- b.Wait() }()
		synctest.Wait()
		if got := len(rounds); got != 3 {
			t.Fatalf("New(3), third goroutine calls Wait: returned = %d, want 3", got)
		}
		for range 3 {
			if r := <-rounds; r != 0 {
				t.Errorf("New(3): first round Wait() = %d, want 0", r)
			}
		}
	})
}

func TestSecondRoundBlocks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(2)
		rounds := make(chan int, 4)
		go func() { rounds <- b.Wait() }()
		go func() { rounds <- b.Wait() }()
		synctest.Wait()
		if got := len(rounds); got != 2 {
			t.Fatalf("New(2), round 0 with 2 goroutines: returned = %d, want 2", got)
		}
		<-rounds
		<-rounds
		go func() { rounds <- b.Wait() }()
		synctest.Wait()
		if got := len(rounds); got != 0 {
			t.Errorf("New(2), 1 goroutine calls Wait in round 1: returned = %d, want 0 (it must block)", got)
		}
		go func() { rounds <- b.Wait() }()
		synctest.Wait()
		if got := len(rounds); got != 2 {
			t.Fatalf("New(2), second goroutine joins round 1: returned = %d, want 2", got)
		}
		for range 2 {
			if r := <-rounds; r != 1 {
				t.Errorf("New(2): second round Wait() = %d, want 1", r)
			}
		}
	})
}

func TestRoundsInLockstep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const n, rounds = 4, 50
		b := New(n)
		var arrived [rounds]atomic.Int32
		var mu sync.Mutex
		var problem string
		report := func(s string) {
			mu.Lock()
			defer mu.Unlock()
			if problem == "" {
				problem = s
			}
		}
		var wg sync.WaitGroup
		for id := range n {
			wg.Go(func() {
				for r := range rounds {
					time.Sleep(time.Duration(id) * time.Millisecond) // arrive at different times
					arrived[r].Add(1)
					got := b.Wait()
					if got != r {
						report(fmt.Sprintf("goroutine %d: Wait() in round %d = %d, want %d", id, r, got, r))
						return
					}
					if a := arrived[r].Load(); a != n {
						report(fmt.Sprintf("goroutine %d left round %d when arrived = %d, want %d", id, r, a, n))
						return
					}
				}
			})
		}
		wg.Wait()
		if problem != "" {
			t.Error(problem)
		}
	})
}

func TestManyRoundsConcurrent(t *testing.T) {
	const n, rounds = 10, 1000
	b := New(n)
	last := make([]int, n)
	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for id := range n {
			wg.Go(func() {
				for range rounds {
					last[id] = b.Wait()
				}
			})
		}
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%d goroutines x %d rounds did not finish within 2s (deadlock?)", n, rounds)
	}
	for id, got := range last {
		if got != rounds-1 {
			t.Fatalf("goroutine %d: last Wait() = %d, want %d", id, got, rounds-1)
		}
	}
}
