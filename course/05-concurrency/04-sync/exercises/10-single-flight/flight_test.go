package flight

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var errDown = errors.New("database down")

// start runs f in a goroutine and returns a channel closed when f returns.
func start(f func()) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	return done
}

// finishes fails the test if done isn't closed within 2s.
func finishes(t *testing.T, done <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

// inFlight returns the number of entries in g.calls, and the dups of key's
// entry (-1 if there is none). ok is false if g.mu was locked.
func inFlight[V any](g *Group[V], key string) (n, dups int, ok bool) {
	if !g.mu.TryLock() {
		return 0, 0, false
	}
	defer g.mu.Unlock()
	dups = -1
	if c := g.calls[key]; c != nil {
		dups = c.dups
	}
	return len(g.calls), dups, true
}

func TestSequentialCallsAreNotCached(t *testing.T) {
	var g Group[int]
	var calls atomic.Int64
	fn := func() (int, error) { return int(calls.Add(1)) * 10, nil }
	for i := range 3 {
		var v int
		var shared bool
		var err error
		finishes(t, start(func() { v, shared, err = g.Do("k", fn) }), `Do("k", fn)`)
		if want := (i + 1) * 10; v != want || shared || err != nil {
			t.Errorf("Do(%q) call %d = (%d, %v, %v), want (%d, false, <nil>)", "k", i+1, v, shared, err, want)
		}
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("3 sequential Do(%q): fn calls = %d, want 3 (results are shared only while in flight)", "k", got)
	}
	if n, _, _ := inFlight(&g, "k"); n != 0 {
		t.Errorf("after every Do returned: len(g.calls) = %d, want 0", n)
	}
}

func TestErrorIsReturned(t *testing.T) {
	var g Group[string]
	var v string
	var shared bool
	var err error
	finishes(t, start(func() {
		v, shared, err = g.Do("user:1", func() (string, error) { return "", errDown })
	}), `Do("user:1", failing fn)`)
	if v != "" || shared || !errors.Is(err, errDown) {
		t.Errorf("Do(%q, failing fn) = (%q, %v, %v), want (\"\", false, %v)", "user:1", v, shared, err, errDown)
	}
	if n, _, _ := inFlight(&g, "user:1"); n != 0 {
		t.Errorf("after a failed Do: len(g.calls) = %d, want 0 (a failed call must be forgotten too)", n)
	}
}

type result struct {
	v      []string
	shared bool
	err    error
}

// concurrentDo starts n callers of g.Do(key, fn) while fn is blocked, waits
// until all of them have joined one call, then lets fn return.
func concurrentDo(t *testing.T, n int, val []string, fnErr error) (results []result, calls int64) {
	t.Helper()
	var g Group[[]string]
	var nCalls atomic.Int64
	release := make(chan struct{})
	fn := func() ([]string, error) {
		nCalls.Add(1)
		<-release
		return val, fnErr
	}
	results = make([]result, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			var r result
			r.v, r.shared, r.err = g.Do("k", fn)
			results[i] = r
		})
	}
	allDone := start(wg.Wait)
	defer func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}()

	deadline := time.After(2 * time.Second)
	lastLocked := false
	for joined := false; !joined; {
		select {
		case <-allDone:
			t.Fatalf("%d concurrent Do(%q): all returned while fn was still blocked (fn calls = %d), want them to wait for fn", n, "k", nCalls.Load())
		case <-deadline:
			if lastLocked {
				t.Fatalf("%d concurrent Do(%q): g.mu stayed locked while fn ran, want it unlocked during fn so others can join", n, "k")
			}
			_, dups, _ := inFlight(&g, "k")
			t.Fatalf("%d concurrent Do(%q) with fn blocked: g.calls[%q].dups = %d after 2s, want %d (fn calls = %d)", n, "k", "k", dups, n-1, nCalls.Load())
		case <-time.After(time.Millisecond):
		}
		_, dups, ok := inFlight(&g, "k")
		lastLocked = !ok
		joined = ok && dups == n-1
	}
	close(release)
	finishes(t, allDone, "the callers after fn returned")
	if c, _, _ := inFlight(&g, "k"); c != 0 {
		t.Errorf("after the shared call returned: len(g.calls) = %d, want 0", c)
	}
	return results, nCalls.Load()
}

func TestConcurrentCallersShareOneCall(t *testing.T) {
	want := []string{"ann", "bob"}
	results, calls := concurrentDo(t, 10, want, nil)
	if calls != 1 {
		t.Errorf("10 concurrent Do(%q): fn calls = %d, want 1", "k", calls)
	}
	for i, r := range results {
		if !slices.Equal(r.v, want) || !r.shared || r.err != nil {
			t.Errorf("10 concurrent Do(%q): caller %d got (%v, %v, %v), want (%v, true, <nil>)", "k", i, r.v, r.shared, r.err, want)
		}
	}
}

func TestConcurrentCallersShareTheError(t *testing.T) {
	results, calls := concurrentDo(t, 5, nil, errDown)
	if calls != 1 {
		t.Errorf("5 concurrent Do(%q), failing fn: fn calls = %d, want 1", "k", calls)
	}
	for i, r := range results {
		if r.v != nil || !r.shared || !errors.Is(r.err, errDown) {
			t.Errorf("5 concurrent Do(%q), failing fn: caller %d got (%v, %v, %v), want ([], true, %v)", "k", i, r.v, r.shared, r.err, errDown)
		}
	}
}

func TestDifferentKeysRunInParallel(t *testing.T) {
	var g Group[string]
	aStarted, bStarted := make(chan struct{}), make(chan struct{})
	// each fn waits until the other one has started: only possible if the
	// two keys really run at the same time.
	both := func(mine, other chan struct{}, v string) func() (string, error) {
		return func() (string, error) {
			close(mine)
			select {
			case <-other:
				return v, nil
			case <-time.After(2 * time.Second):
				return "", errors.New("timed out waiting for the other key's fn")
			}
		}
	}
	var va, vb string
	var ea, eb error
	doneA := start(func() { va, _, ea = g.Do("a", both(aStarted, bStarted, "A")) })
	doneB := start(func() { vb, _, eb = g.Do("b", both(bStarted, aStarted, "B")) })
	finishes(t, doneA, `Do("a")`)
	finishes(t, doneB, `Do("b")`)
	if va != "A" || ea != nil {
		t.Errorf("Do(%q) during Do(%q) = (%q, %v), want (%q, <nil>): different keys must not wait for each other", "a", "b", va, ea, "A")
	}
	if vb != "B" || eb != nil {
		t.Errorf("Do(%q) during Do(%q) = (%q, %v), want (%q, <nil>): different keys must not wait for each other", "b", "a", vb, eb, "B")
	}
}

func TestNewCallAfterSharedCall(t *testing.T) {
	var g Group[int]
	var calls atomic.Int64
	release := make(chan struct{})
	fn := func() (int, error) {
		calls.Add(1)
		<-release
		return 1, nil
	}
	first := start(func() { g.Do("k", fn) })
	second := make(chan struct{})
	deadline := time.After(2 * time.Second)
	for {
		if _, dups, ok := inFlight(&g, "k"); ok && dups >= 0 {
			break
		}
		select {
		case <-first:
			t.Fatalf("Do(%q) returned while fn was blocked, want it to wait for fn", "k")
		case <-deadline:
			close(release)
			t.Fatalf("Do(%q): no entry in g.calls within 2s while fn ran, want g.calls[%q] set", "k", "k")
		case <-time.After(time.Millisecond):
		}
	}
	close(release)
	finishes(t, first, `first Do("k")`)
	go func() {
		defer close(second)
		v, shared, err := g.Do("k", func() (int, error) { calls.Add(1); return 2, nil })
		if v != 2 || shared || err != nil {
			t.Errorf("Do(%q) after the first call finished = (%d, %v, %v), want (2, false, <nil>)", "k", v, shared, err)
		}
	}()
	finishes(t, second, `second Do("k")`)
	if got := calls.Load(); got != 2 {
		t.Errorf("Do(%q) twice, one after the other: fn calls = %d, want 2", "k", got)
	}
}
