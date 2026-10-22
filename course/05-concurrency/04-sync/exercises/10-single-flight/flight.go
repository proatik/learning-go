package flight

import "sync"

// call is one in-flight (or finished) call of fn for a key.
type call[V any] struct {
	done chan struct{} // closed when fn has returned and val, err are set
	val  V
	err  error
	dups int // how many callers joined this call instead of starting their own
}

// Group runs at most one fn per key at a time and shares its result with
// every caller that asks for the same key meanwhile. The zero value is ready
// to use.
type Group[V any] struct {
	mu    sync.Mutex          // guards calls
	calls map[string]*call[V] // in-flight calls by key
}

// Do calls fn and returns its results, unless a call for key is already in
// flight: then it waits for that call and returns its results instead.
// shared reports whether the results were given to more than one caller.
func (g *Group[V]) Do(key string, fn func() (V, error)) (v V, shared bool, err error) {
	// TODO
	var zero V
	return zero, false, nil
}
