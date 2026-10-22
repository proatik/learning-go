package fetch

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

var errNotFound = errors.New("404 not found")

// step is what one Get attempt does: wait delay (or until ctx is done), then
// return body or err.
type step struct {
	delay time.Duration
	body  string
	err   error
}

// fakeClient follows a plan per URL; the last step repeats. URLs without a
// plan take 1s and return "body of <url>".
type fakeClient struct {
	mu                  sync.Mutex
	plan                map[string][]step
	calls               map[string]int
	running, maxRunning atomic.Int64
}

func newClient(plan map[string][]step) *fakeClient {
	return &fakeClient{plan: plan, calls: map[string]int{}}
}

func (f *fakeClient) Get(ctx context.Context, url string) ([]byte, error) {
	f.mu.Lock()
	n := f.calls[url]
	f.calls[url]++
	steps := f.plan[url]
	f.mu.Unlock()
	st := step{delay: time.Second, body: "body of " + url}
	if len(steps) > 0 {
		st = steps[min(n, len(steps)-1)]
	}
	cur := f.running.Add(1)
	defer f.running.Add(-1)
	for {
		m := f.maxRunning.Load()
		if cur <= m || f.maxRunning.CompareAndSwap(m, cur) {
			break
		}
	}
	timer := time.NewTimer(st.delay)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if st.err != nil {
		return nil, st.err
	}
	return []byte(st.body), nil
}

func (f *fakeClient) totalCalls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	total := 0
	for _, n := range f.calls {
		total += n
	}
	return total
}

func wantLen(t *testing.T, urls []string, got []Result) {
	t.Helper()
	if len(got) != len(urls) {
		t.Fatalf("Download(%q) returned %d results, want %d", urls, len(got), len(urls))
	}
	for i, r := range got {
		if r.URL != urls[i] {
			t.Fatalf("Download(%q): results[%d].URL = %q, want %q (input order)", urls, i, r.URL, urls[i])
		}
	}
}

func TestDownloadInOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newClient(map[string][]step{
			"slow": {{delay: 3 * time.Second, body: "S"}},
			"mid":  {{delay: 2 * time.Second, body: "M"}},
			"fast": {{delay: 0, body: "F"}},
		})
		urls := []string{"slow", "mid", "fast", "fast"}
		got := Download(context.Background(), c, urls, 4, time.Minute)
		wantLen(t, urls, got)
		for i, want := range []string{"S", "M", "F", "F"} {
			r := got[i]
			if string(r.Body) != want || r.Err != nil || r.Attempts != 1 {
				t.Errorf("results[%d] = {Body: %q, Err: %v, Attempts: %d}, want {Body: %q, Err: <nil>, Attempts: 1}", i, r.Body, r.Err, r.Attempts, want)
			}
		}
		if got := Download(context.Background(), c, nil, 4, time.Minute); len(got) != 0 {
			t.Errorf("Download(nil) = %v, want an empty result", got)
		}
	})
}

func TestRetries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		wrapped := fmt.Errorf("503 unavailable: %w", ErrTemporary)
		c := newClient(map[string][]step{
			"flaky":     {{err: ErrTemporary}, {body: "ok"}},
			"wrapped":   {{err: wrapped}, {body: "ok"}},
			"down":      {{err: ErrTemporary}, {err: wrapped}},
			"missing":   {{err: errNotFound}},
			"thenfails": {{err: ErrTemporary}, {err: errNotFound}},
		})
		urls := []string{"flaky", "wrapped", "down", "missing", "thenfails"}
		got := Download(context.Background(), c, urls, 2, time.Minute)
		wantLen(t, urls, got)
		tests := []struct {
			body     string
			err      error
			attempts int
		}{
			{"ok", nil, 2},
			{"ok", nil, 2},
			{"", wrapped, 2},
			{"", errNotFound, 1},
			{"", errNotFound, 2},
		}
		for i, tt := range tests {
			r := got[i]
			errOK := (tt.err == nil && r.Err == nil) || (tt.err != nil && errors.Is(r.Err, tt.err))
			if string(r.Body) != tt.body || !errOK || r.Attempts != tt.attempts {
				t.Errorf("%q: result = {Body: %q, Err: %v, Attempts: %d}, want {Body: %q, Err: %v, Attempts: %d}",
					urls[i], r.Body, r.Err, r.Attempts, tt.body, tt.err, tt.attempts)
			}
		}
	})
}

func TestPerRequestTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newClient(map[string][]step{
			"hangs":    {{delay: 5 * time.Second, body: "late"}},
			"slowOnce": {{delay: 5 * time.Second, body: "late"}, {delay: 500 * time.Millisecond, body: "ok"}},
		})
		urls := []string{"hangs", "slowOnce"}
		start := time.Now()
		got := Download(context.Background(), c, urls, 3, 2*time.Second)
		wantLen(t, urls, got)
		if r := got[0]; !errors.Is(r.Err, context.DeadlineExceeded) || r.Attempts != 2 {
			t.Errorf(`"hangs" (5s, timeout 2s): result = {Err: %v, Attempts: %d}, want {Err: %v, Attempts: 2}`, r.Err, r.Attempts, context.DeadlineExceeded)
		}
		if r := got[1]; string(r.Body) != "ok" || r.Err != nil || r.Attempts != 2 {
			t.Errorf(`"slowOnce": result = {Body: %q, Err: %v, Attempts: %d}, want {Body: "ok", Err: <nil>, Attempts: 2}`, r.Body, r.Err, r.Attempts)
		}
		if el := time.Since(start); el != 4*time.Second {
			t.Errorf("Download with a 2s timeout and one retry took %v, want 4s", el)
		}
	})
}

func TestLimit(t *testing.T) {
	for _, tt := range []struct {
		limit, wantMax int
		want           time.Duration
	}{{3, 3, 4 * time.Second}, {10, 10, time.Second}, {0, 1, 10 * time.Second}, {-2, 1, 10 * time.Second}} {
		synctest.Test(t, func(t *testing.T) {
			c := newClient(nil)
			urls := make([]string, 10)
			for i := range urls {
				urls[i] = fmt.Sprintf("u%d", i)
			}
			start := time.Now()
			got := Download(context.Background(), c, urls, tt.limit, time.Minute)
			wantLen(t, urls, got)
			if m := c.maxRunning.Load(); m != int64(tt.wantMax) {
				t.Errorf("Download(10 URLs, limit %d): max Get calls at once = %d, want %d", tt.limit, m, tt.wantMax)
			}
			if el := time.Since(start); el != tt.want {
				t.Errorf("Download(10 URLs of 1s, limit %d) took %v, want %v", tt.limit, el, tt.want)
			}
			for i, r := range got {
				if want := "body of " + urls[i]; string(r.Body) != want {
					t.Fatalf("results[%d].Body = %q, want %q", i, r.Body, want)
				}
			}
		})
	}
}

func TestCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newClient(nil) // every URL takes 1s
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		time.AfterFunc(1500*time.Millisecond, cancel)
		urls := []string{"a", "b", "c", "d", "e", "f"}
		start := time.Now()
		got := Download(ctx, c, urls, 2, time.Minute)
		wantLen(t, urls, got)
		if el := time.Since(start); el != 1500*time.Millisecond {
			t.Errorf("Download cancelled at 1.5s returned after %v, want 1.5s", el)
		}
		var done, interrupted, notStarted int
		for _, r := range got {
			switch {
			case r.Err == nil && r.Attempts == 1:
				done++
			case errors.Is(r.Err, context.Canceled) && r.Attempts == 1:
				interrupted++
			case errors.Is(r.Err, context.Canceled) && r.Attempts == 0:
				notStarted++
			default:
				t.Errorf("%q: unexpected result {Err: %v, Attempts: %d}", r.URL, r.Err, r.Attempts)
			}
		}
		if done != 2 || interrupted != 2 || notStarted != 2 {
			t.Errorf("6 URLs of 1s, limit 2, cancelled at 1.5s: (done, interrupted, not started) = (%d, %d, %d), want (2, 2, 2)", done, interrupted, notStarted)
		}
		if n := c.totalCalls(); n != 4 {
			t.Errorf("Get calls = %d, want 4 (no retry after cancel, no Get for URLs that never started)", n)
		}
	})
}

func TestParentDeadlineIsNotRetried(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newClient(map[string][]step{"slow": {{delay: 5 * time.Second}}})
		ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		defer cancel()
		start := time.Now()
		got := Download(ctx, c, []string{"slow"}, 1, time.Minute)
		wantLen(t, []string{"slow"}, got)
		if r := got[0]; !errors.Is(r.Err, context.DeadlineExceeded) || r.Attempts != 1 {
			t.Errorf("ctx deadline 1.5s, Get takes 5s: result = {Err: %v, Attempts: %d}, want {Err: %v, Attempts: 1} (the parent is dead: no retry)", r.Err, r.Attempts, context.DeadlineExceeded)
		}
		if el := time.Since(start); el != 1500*time.Millisecond {
			t.Errorf("Download returned after %v, want 1.5s", el)
		}
	})
}
