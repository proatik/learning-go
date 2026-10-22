package fetch

import (
	"errors"
	"fmt"
	"slices"
	"sync"
)

// Fetcher loads the value stored under a key, for example from another
// service.
type Fetcher interface {
	Fetch(key string) (string, error)
}

// ErrEmptyKey is reported for the empty key, which is never fetched.
var ErrEmptyKey = errors.New("fetch: empty key")

// Report is the outcome of FetchAll.
type Report struct {
	Values map[string]string // key -> value, for every successful fetch
	Errors map[string]error  // key -> error, for every failed key
	Failed []string          // the keys of Errors, sorted ascending
}

// result is what one goroutine produces for one key.
type result struct {
	value string
	err   error
}

// FetchAll fetches every distinct key concurrently and reports the values
// and errors.
func FetchAll(f Fetcher, keys []string) Report {
	unique := make([]string, 0, len(keys))
	seen := make(map[string]bool, len(keys))
	for _, k := range keys {
		if !seen[k] {
			seen[k] = true
			unique = append(unique, k)
		}
	}

	results := make([]result, len(unique))
	var wg sync.WaitGroup
	for i, key := range unique {
		if key == "" {
			results[i].err = ErrEmptyKey
			continue
		}
		wg.Go(func() { results[i] = fetchOne(f, key) })
	}
	wg.Wait()

	rep := Report{
		Values: make(map[string]string),
		Errors: make(map[string]error),
		Failed: []string{},
	}
	for i, key := range unique {
		if err := results[i].err; err != nil {
			rep.Errors[key] = err
			rep.Failed = append(rep.Failed, key)
		} else {
			rep.Values[key] = results[i].value
		}
	}
	slices.Sort(rep.Failed)
	return rep
}

// fetchOne fetches key, turning an error or a panic into a descriptive error.
func fetchOne(f Fetcher, key string) (res result) {
	defer func() {
		if r := recover(); r != nil {
			res = result{err: fmt.Errorf("fetch %q panicked: %v", key, r)}
		}
	}()
	v, err := f.Fetch(key)
	if err != nil {
		return result{err: fmt.Errorf("fetch %q: %w", key, err)}
	}
	return result{value: v}
}
