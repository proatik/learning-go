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
	g.mu.Lock()
	if c, ok := g.calls[key]; ok {
		c.dups++
		g.mu.Unlock()
		<-c.done // close(c.done) happens after val and err are set
		return c.val, true, c.err
	}
	if g.calls == nil {
		g.calls = make(map[string]*call[V])
	}
	c := &call[V]{done: make(chan struct{})}
	g.calls[key] = c
	g.mu.Unlock()

	c.val, c.err = fn() // without the lock: other keys must not wait for us

	g.mu.Lock()
	delete(g.calls, key) // from now on, a new Do for key starts a new call
	shared = c.dups > 0
	g.mu.Unlock()
	close(c.done)
	return c.val, shared, c.err
}
