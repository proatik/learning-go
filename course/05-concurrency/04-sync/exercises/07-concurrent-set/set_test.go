package syncset

import (
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// run calls f in a goroutine and fails the test if it doesn't return within
// 2s, so a self-deadlock fails the test instead of hanging it.
func run(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

func TestZeroValue(t *testing.T) {
	var s Set[string]
	run(t, "methods on a zero Set", func() {
		if s.Has("a") {
			t.Errorf("Set{}.Has(%q) = true, want false", "a")
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Set{}.Len() = %d, want 0", got)
		}
		if s.Remove("a") {
			t.Errorf("Set{}.Remove(%q) = true, want false", "a")
		}
		if got := s.Values(); len(got) != 0 {
			t.Errorf("Set{}.Values() = %v, want []", got)
		}
	})
}

func TestSequential(t *testing.T) {
	var s Set[string]
	run(t, "Add/Has/Remove/Len", func() {
		steps := []struct {
			call string
			got  bool
			want bool
		}{
			{`Add("go")`, s.Add("go"), true},
			{`Add("go")`, s.Add("go"), false},
			{`Add("")`, s.Add(""), true},
			{`Add("日本")`, s.Add("日本"), true},
			{`Has("go")`, s.Has("go"), true},
			{`Has("")`, s.Has(""), true},
			{`Has("rust")`, s.Has("rust"), false},
			{`Remove("go")`, s.Remove("go"), true},
			{`Remove("go")`, s.Remove("go"), false},
			{`Has("go")`, s.Has("go"), false},
			{`Add("go")`, s.Add("go"), true},
		}
		for _, st := range steps {
			if st.got != st.want {
				t.Errorf("%s = %v, want %v", st.call, st.got, st.want)
			}
		}
		if got := s.Len(); got != 3 {
			t.Errorf("Len() = %d, want 3", got)
		}
		got := s.Values()
		slices.Sort(got)
		if want := []string{"", "go", "日本"}; !slices.Equal(got, want) {
			t.Errorf("Values() (sorted) = %q, want %q", got, want)
		}
	})
}

func TestOtherTypes(t *testing.T) {
	type point struct{ X, Y int }
	var s Set[point]
	run(t, "Set[point]", func() {
		s.Add(point{1, 2})
		s.Add(point{1, 2})
		s.Add(point{0, 0})
		if got := s.Len(); got != 2 {
			t.Errorf("Set[point]: Len() = %d, want 2", got)
		}
		if !s.Has(point{0, 0}) {
			t.Errorf("Set[point]: Has({0 0}) = false, want true")
		}
	})
}

func TestReadersShareTheLock(t *testing.T) {
	s := &Set[int]{items: map[int]struct{}{7: {}}}
	s.mu.RLock()
	var has bool
	run(t, "Has(7) while another reader holds s.mu.RLock()", func() { has = s.Has(7) })
	run(t, "Len() while another reader holds s.mu.RLock()", func() { s.Len() })
	run(t, "Values() while another reader holds s.mu.RLock()", func() { s.Values() })
	if !has {
		t.Errorf("Has(7) = false, want true")
	}
	added := make(chan struct{})
	go func() {
		defer close(added)
		s.Add(8)
	}()
	time.Sleep(20 * time.Millisecond)
	select {
	case <-added:
		t.Errorf("Add(8) while a reader holds s.mu.RLock() = returned, want it to wait (use s.mu.Lock())")
	default:
	}
	s.mu.RUnlock()
	run(t, "Add(8) after the reader unlocked", func() { <-added })
}

func TestConcurrentAddIsExact(t *testing.T) {
	const goroutines, values = 50, 200
	var s Set[int]
	var added, removed atomic.Int64
	run(t, "50 goroutines adding the same 200 values", func() {
		var wg sync.WaitGroup
		for g := range goroutines {
			wg.Go(func() {
				for i := range values {
					v := (i + g*7) % values // every goroutine, a different order
					if s.Add(v) {
						added.Add(1)
					}
					s.Has(v)
				}
			})
		}
		wg.Wait()
	})
	if got := added.Load(); got != values {
		t.Errorf("%d goroutines Add the same %d values: Add returned true %d times, want %d", goroutines, values, got, values)
	}
	if got := s.Len(); got != values {
		t.Errorf("Len() = %d, want %d", got, values)
	}
	run(t, "50 goroutines removing the same 200 values", func() {
		var wg sync.WaitGroup
		for g := range goroutines {
			wg.Go(func() {
				for i := range values {
					if s.Remove((i + g*13) % values) {
						removed.Add(1)
					}
					_ = s.Values()
				}
			})
		}
		wg.Wait()
	})
	if got := removed.Load(); got != values {
		t.Errorf("%d goroutines Remove the same %d values: Remove returned true %d times, want %d", goroutines, values, got, values)
	}
	if got := s.Len(); got != 0 {
		t.Errorf("Len() after removing everything = %d, want 0", got)
	}
}
