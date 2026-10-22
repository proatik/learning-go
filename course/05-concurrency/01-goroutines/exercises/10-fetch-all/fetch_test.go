package fetch_test

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	fetch "learngo/course/05-concurrency/01-goroutines/exercises/10-fetch-all"
)

var ErrNotFound = errors.New("not found")

// fakeFetcher serves fixed data. Each call sleeps delay (the tests run it
// on synctest's fake clock) and is counted per key.
type fakeFetcher struct {
	delay time.Duration

	mu    sync.Mutex
	calls map[string]int
}

func (f *fakeFetcher) Fetch(key string) (string, error) {
	f.mu.Lock()
	if f.calls == nil {
		f.calls = make(map[string]int)
	}
	f.calls[key]++
	f.mu.Unlock()

	time.Sleep(f.delay)
	switch key {
	case "a":
		return "apple", nil
	case "b":
		return "banana", nil
	case "c":
		return "", nil // an empty value is still a success
	case "x", "y":
		return "", ErrNotFound
	case "p":
		panic("boom")
	case "q":
		panic(fmt.Errorf("wrapped %w", ErrNotFound))
	default:
		return "value-" + key, nil
	}
}

func (f *fakeFetcher) callCount() map[string]int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return maps.Clone(f.calls)
}

// errText renders an error map deterministically for messages.
func errText(m map[string]error) map[string]string {
	out := make(map[string]string, len(m))
	for k, err := range m {
		out[k] = fmt.Sprint(err)
	}
	return out
}

// checkReport compares rep with the expected values and error messages.
func checkReport(t *testing.T, keys []string, rep fetch.Report, values, errs map[string]string) {
	t.Helper()
	if rep.Values == nil || rep.Errors == nil {
		t.Fatalf("FetchAll(%q) Values = %v, Errors = %v, want non-nil maps", keys, rep.Values, rep.Errors)
	}
	if !maps.Equal(rep.Values, values) {
		t.Errorf("FetchAll(%q).Values = %v, want %v", keys, rep.Values, values)
	}
	if got := errText(rep.Errors); !maps.Equal(got, errs) {
		t.Errorf("FetchAll(%q).Errors = %q, want %q", keys, got, errs)
	}
	wantFailed := slices.Sorted(maps.Keys(errs))
	if wantFailed == nil {
		wantFailed = []string{}
	}
	if !slices.Equal(rep.Failed, wantFailed) {
		t.Errorf("FetchAll(%q).Failed = %q, want %q", keys, rep.Failed, wantFailed)
	}
}

func TestFetchAll(t *testing.T) {
	tests := []struct {
		keys   []string
		values map[string]string
		errs   map[string]string
		calls  map[string]int
	}{
		{
			keys:   []string{"b", "a"},
			values: map[string]string{"a": "apple", "b": "banana"},
			errs:   map[string]string{},
			calls:  map[string]int{"a": 1, "b": 1},
		},
		{
			keys:   []string{"x", "a", "x"},
			values: map[string]string{"a": "apple"},
			errs:   map[string]string{"x": `fetch "x": not found`},
			calls:  map[string]int{"a": 1, "x": 1},
		},
		{
			keys:   []string{"p", "", "a"},
			values: map[string]string{"a": "apple"},
			errs:   map[string]string{"": "fetch: empty key", "p": `fetch "p" panicked: boom`},
			calls:  map[string]int{"a": 1, "p": 1},
		},
		{
			keys:   []string{"y", "x", "c", "q", "c", "y", "", ""},
			values: map[string]string{"c": ""},
			errs: map[string]string{
				"":  "fetch: empty key",
				"q": `fetch "q" panicked: wrapped not found`,
				"x": `fetch "x": not found`,
				"y": `fetch "y": not found`,
			},
			calls: map[string]int{"c": 1, "q": 1, "x": 1, "y": 1},
		},
		{
			keys:   []string{"héllo", "世界"},
			values: map[string]string{"héllo": "value-héllo", "世界": "value-世界"},
			errs:   map[string]string{},
			calls:  map[string]int{"héllo": 1, "世界": 1},
		},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			f := &fakeFetcher{delay: time.Second}
			start := time.Now()
			rep := fetch.FetchAll(f, tt.keys)
			elapsed := time.Since(start)
			checkReport(t, tt.keys, rep, tt.values, tt.errs)
			if got := f.callCount(); !maps.Equal(got, tt.calls) {
				t.Errorf("FetchAll(%q): Fetch calls per key = %v, want %v", tt.keys, got, tt.calls)
			}
			if elapsed != time.Second {
				t.Errorf("FetchAll(%q) with 1s fetches took %v, want 1s", tt.keys, elapsed)
			}
		})
	}
}

func TestErrorsUnwrap(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		keys := []string{"x", "", "a"}
		rep := fetch.FetchAll(&fakeFetcher{delay: time.Second}, keys)
		if !errors.Is(rep.Errors["x"], ErrNotFound) {
			t.Errorf("FetchAll(%q): errors.Is(Errors[\"x\"], ErrNotFound) = false, want true", keys)
		}
		if rep.Errors[""] != fetch.ErrEmptyKey {
			t.Errorf("FetchAll(%q).Errors[\"\"] = %v, want ErrEmptyKey itself", keys, rep.Errors[""])
		}
	})
}

func TestEmpty(t *testing.T) {
	for _, keys := range [][]string{nil, {}} {
		synctest.Test(t, func(t *testing.T) {
			f := &fakeFetcher{delay: time.Second}
			rep := fetch.FetchAll(f, keys)
			checkReport(t, keys, rep, map[string]string{}, map[string]string{})
			if got := f.callCount(); len(got) != 0 {
				t.Errorf("FetchAll(%q): Fetch calls = %v, want none", keys, got)
			}
		})
	}
}

func TestOnlyEmptyKeys(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		keys := []string{"", "", ""}
		f := &fakeFetcher{delay: time.Second}
		start := time.Now()
		rep := fetch.FetchAll(f, keys)
		checkReport(t, keys, rep, map[string]string{}, map[string]string{"": "fetch: empty key"})
		if got := f.callCount(); len(got) != 0 {
			t.Errorf("FetchAll(%q): Fetch calls = %v, want none", keys, got)
		}
		if elapsed := time.Since(start); elapsed != 0 {
			t.Errorf("FetchAll(%q) took %v, want 0s (nothing to fetch)", keys, elapsed)
		}
	})
}

func TestManyKeys(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var keys []string
		values := make(map[string]string)
		errs := make(map[string]string)
		for i := range 300 {
			k := fmt.Sprintf("k%03d", i)
			keys = append(keys, k, k) // every key twice
			values[k] = "value-" + k
		}
		keys = append(keys, "x", "p")
		errs["x"] = `fetch "x": not found`
		errs["p"] = `fetch "p" panicked: boom`
		f := &fakeFetcher{delay: 2 * time.Second}
		start := time.Now()
		rep := fetch.FetchAll(f, keys)
		elapsed := time.Since(start)
		checkReport(t, []string{"k000", "k000", "…", "k299", "k299", "x", "p"}, rep, values, errs)
		total := 0
		for _, n := range f.callCount() {
			total += n
		}
		if total != 302 {
			t.Errorf("FetchAll(302 distinct keys, each twice): Fetch calls = %d, want 302", total)
		}
		if elapsed != 2*time.Second {
			t.Errorf("FetchAll(302 keys) with 2s fetches took %v, want 2s", elapsed)
		}
	})
}
