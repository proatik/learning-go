package fetch

import "errors"

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

// FetchAll fetches every distinct key concurrently and reports the values
// and errors.
func FetchAll(f Fetcher, keys []string) Report {
	// TODO
	return Report{}
}
