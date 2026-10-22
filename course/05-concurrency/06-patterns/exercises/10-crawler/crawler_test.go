package crawler

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

var errNotFound = errors.New("not found")

// fakeFetcher serves a fixed link graph. Each fetch takes delay[url] (or
// defaultDelay) of fake time, gives up when ctx ends, and is counted.
type fakeFetcher struct {
	graph        map[string][]string
	delay        map[string]time.Duration
	defaultDelay time.Duration

	mu            sync.Mutex
	calls         map[string]int
	running, peak atomic.Int64
}

func (f *fakeFetcher) Fetch(ctx context.Context, url string) ([]string, error) {
	f.mu.Lock()
	if f.calls == nil {
		f.calls = make(map[string]int)
	}
	f.calls[url]++
	f.mu.Unlock()

	now := f.running.Add(1)
	defer f.running.Add(-1)
	for p := f.peak.Load(); now > p && !f.peak.CompareAndSwap(p, now); p = f.peak.Load() {
	}

	d, ok := f.delay[url]
	if !ok {
		d = f.defaultDelay
	}
	select {
	case <-time.After(d):
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	links, ok := f.graph[url]
	if !ok {
		return nil, errNotFound
	}
	return links, nil
}

// checkOnce fails if any URL was fetched more than once, and returns how many
// distinct URLs were fetched.
func (f *fakeFetcher) checkOnce(t *testing.T, name string) int {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	for url, n := range f.calls {
		if n > 1 {
			t.Errorf("%s fetched %q %d times, want once", name, url, n)
		}
	}
	return len(f.calls)
}

var site = map[string][]string{
	"home": {"a", "b"},
	"a":    {"b", "c", "a"},
	"b":    {"home", "d"},
	"c":    {"e"},
	"d":    {},
	"e":    {"home", "c"},
}

func TestCrawlDepths(t *testing.T) {
	tests := []struct {
		depth int
		want  []string
	}{
		{0, []string{"home"}},
		{1, []string{"a", "b", "home"}},
		{2, []string{"a", "b", "c", "d", "home"}},
		{3, []string{"a", "b", "c", "d", "e", "home"}},
		{10, []string{"a", "b", "c", "d", "e", "home"}},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			f := &fakeFetcher{graph: site, defaultDelay: time.Second}
			name := fmt.Sprintf("Crawl(home, depth %d, limit 3)", tt.depth)
			got, err := Crawl(t.Context(), f, "home", tt.depth, 3)
			if err != nil {
				t.Errorf("%s error = %v, want <nil>", name, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("%s = %q, want %q", name, got, tt.want)
			}
			if n := f.checkOnce(t, name); n != len(tt.want) {
				t.Errorf("%s fetched %d distinct URLs, want %d", name, n, len(tt.want))
			}
		})
	}
}

func TestCrawlNegativeDepth(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeFetcher{graph: site}
		got, err := Crawl(t.Context(), f, "home", -1, 3)
		if len(got) != 0 || err != nil {
			t.Errorf("Crawl(home, depth -1) = (%q, %v), want ([], <nil>)", got, err)
		}
		if n := f.checkOnce(t, "Crawl(depth -1)"); n != 0 {
			t.Errorf("Crawl(home, depth -1) fetched %d URLs, want 0", n)
		}
	})
}

// TestCrawlShortestPath: "c" is 2 links away through the slow page "b", but
// 3 links away through the fast pages "a" → "x". With depth 3, "d" (behind
// "c") is within reach and must be fetched.
func TestCrawlShortestPath(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeFetcher{
			graph: map[string][]string{
				"home": {"a", "b"},
				"a":    {"x"},
				"x":    {"c"},
				"b":    {"c"},
				"c":    {"d"},
				"d":    {"z"},
			},
			delay:        map[string]time.Duration{"b": 5 * time.Second},
			defaultDelay: time.Second,
		}
		got, err := Crawl(t.Context(), f, "home", 3, 4)
		if err != nil {
			t.Errorf("Crawl error = %v, want <nil>", err)
		}
		if want := []string{"a", "b", "c", "d", "home", "x"}; !slices.Equal(got, want) {
			t.Errorf("Crawl(home, depth 3) with a slow b = %q, want %q (reach each URL by its shortest path)", got, want)
		}
		f.checkOnce(t, "Crawl(home, depth 3)")
	})
}

func TestCrawlErrors(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeFetcher{
			graph: map[string][]string{
				"home": {"a", "zz-missing", "m-missing"},
				"a":    {"m-missing", "home"},
			},
			defaultDelay: time.Second,
		}
		got, err := Crawl(t.Context(), f, "home", 5, 2)
		if want := []string{"a", "home"}; !slices.Equal(got, want) {
			t.Errorf("Crawl with broken links = %q, want %q (failed fetches are not visited)", got, want)
		}
		if !errors.Is(err, errNotFound) {
			t.Errorf("Crawl with broken links error = %v, want one wrapping %v", err, errNotFound)
		}
		want := "fetch m-missing: not found\nfetch zz-missing: not found"
		if err == nil || err.Error() != want {
			t.Errorf("Crawl with broken links error = %q, want %q", fmt.Sprint(err), want)
		}
		f.checkOnce(t, "Crawl with broken links")
	})
}

func TestCrawlLimit(t *testing.T) {
	tests := []struct {
		leaves, limit int
		wantPeak      int64
		wantElapsed   time.Duration
	}{
		{20, 4, 4, 6 * time.Second}, // home, then 5 rounds of 4
		{3, 0, 1, 4 * time.Second},  // limit < 1 means 1
		{3, -2, 1, 4 * time.Second},
		{20, 100, 20, 2 * time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			graph := map[string][]string{}
			want := []string{"home"}
			for i := range tt.leaves {
				leaf := fmt.Sprintf("leaf%02d", i)
				graph["home"] = append(graph["home"], leaf)
				graph[leaf] = []string{"home"}
				want = append(want, leaf)
			}
			slices.Sort(want)
			f := &fakeFetcher{graph: graph, defaultDelay: time.Second}
			name := fmt.Sprintf("Crawl(%d leaves of 1s, limit %d)", tt.leaves, tt.limit)
			start := time.Now()
			got, err := Crawl(t.Context(), f, "home", 1, tt.limit)
			elapsed := time.Since(start)
			if err != nil || !slices.Equal(got, want) {
				t.Errorf("%s = (%q, %v), want (%q, <nil>)", name, got, err, want)
			}
			if p := f.peak.Load(); p != tt.wantPeak {
				t.Errorf("%s max concurrent fetches = %d, want %d", name, p, tt.wantPeak)
			}
			if elapsed != tt.wantElapsed {
				t.Errorf("%s took %v, want %v", name, elapsed, tt.wantElapsed)
			}
		})
	}
}

// TestCrawlTimeout: a chain of 1s pages with a 2.5s deadline. home and a
// finish, the fetch of b is cut short, c is never started.
func TestCrawlTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeFetcher{
			graph:        map[string][]string{"home": {"a"}, "a": {"b"}, "b": {"c"}, "c": {"d"}, "d": {}},
			defaultDelay: time.Second,
		}
		ctx, cancel := context.WithTimeout(t.Context(), 2500*time.Millisecond)
		defer cancel()
		start := time.Now()
		got, err := Crawl(ctx, f, "home", 10, 2)
		if d := time.Since(start); d != 2500*time.Millisecond {
			t.Errorf("Crawl with a 2.5s timeout returned after %v, want 2.5s (pass ctx to Fetch)", d)
		}
		if want := []string{"a", "home"}; !slices.Equal(got, want) {
			t.Errorf("Crawl with a 2.5s timeout = %q, want %q", got, want)
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Crawl with a 2.5s timeout error = %v, want one wrapping %v", err, context.DeadlineExceeded)
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.calls["c"] != 0 {
			t.Errorf("Crawl fetched %q after the deadline, want no new fetches once ctx is done", "c")
		}
	})
}

func TestCrawlCancelledContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		f := &fakeFetcher{graph: site}
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		got, err := Crawl(ctx, f, "home", 3, 3)
		if len(got) != 0 || !errors.Is(err, context.Canceled) {
			t.Errorf("Crawl(cancelled ctx) = (%q, %v), want ([], %v)", got, err, context.Canceled)
		}
		if n := f.checkOnce(t, "Crawl(cancelled ctx)"); n != 0 {
			t.Errorf("Crawl(cancelled ctx) fetched %d URLs, want 0", n)
		}
	})
}

// TestCrawlBigGraph compares the crawl of a 300-page graph with a simple
// sequential breadth-first search.
func TestCrawlBigGraph(t *testing.T) {
	graph := map[string][]string{}
	seed := uint32(42)
	page := func(i int) string { return fmt.Sprintf("p%03d", i) }
	for i := range 300 {
		for range 4 {
			seed = seed*1664525 + 1013904223
			graph[page(i)] = append(graph[page(i)], page(int(seed>>16)%300))
		}
	}
	const depth = 4
	dist := map[string]int{"p000": 0}
	queue := []string{"p000"}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		if dist[u] == depth {
			continue
		}
		for _, l := range graph[u] {
			if _, ok := dist[l]; !ok {
				dist[l] = dist[u] + 1
				queue = append(queue, l)
			}
		}
	}
	var want []string
	for u := range dist {
		want = append(want, u)
	}
	slices.Sort(want)

	synctest.Test(t, func(t *testing.T) {
		f := &fakeFetcher{graph: graph, defaultDelay: time.Millisecond}
		got, err := Crawl(t.Context(), f, "p000", depth, 8)
		if err != nil {
			t.Errorf("Crawl(300-page graph) error = %v, want <nil>", err)
		}
		if !slices.Equal(got, want) {
			t.Errorf("Crawl(300-page graph, depth %d) visited %d pages, want %d", depth, len(got), len(want))
		}
		f.checkOnce(t, "Crawl(300-page graph)")
		if p := f.peak.Load(); p > 8 {
			t.Errorf("Crawl(limit 8) max concurrent fetches = %d, want at most 8", p)
		}
	})
}
