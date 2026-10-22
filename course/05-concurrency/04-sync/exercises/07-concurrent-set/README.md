# Concurrent Set

A generic set that any number of goroutines can share. `Add` and `Remove` report whether they changed anything, and that answer must be right even under contention.

```go
type Set[T comparable] struct {
	mu    sync.RWMutex   // guards items
	items map[T]struct{} // created on first Add
}

func (s *Set[T]) Add(v T) bool
func (s *Set[T]) Remove(v T) bool
func (s *Set[T]) Has(v T) bool
func (s *Set[T]) Len() int
func (s *Set[T]) Values() []T
```

## Examples

| Code | Result |
|---|---|
| `var s Set[string]; s.Has("go")` | `false` (the zero value is an empty set) |
| `s.Add("go")` / `s.Add("go")` | `true` / `false` (already there) |
| `s.Remove("go")` / `s.Remove("go")` | `true` / `false` |
| `s.Add("a"); s.Add("b"); s.Len()` | `2` |
| `s.Values()` | `[a b]` in any order |
| 50 goroutines each `Add` the same 200 values | `Add` returns `true` exactly 200 times in total |

## Rules
- The zero value works: create the map in `Add` (reads, `delete` and `len` on a nil map are fine).
- `Add` returns `true` only for the call that actually inserted `v`; `Remove` only for the call that actually deleted it. Check and change under the **same** `Lock`.
- Read-only methods (`Has`, `Len`, `Values`) use `RLock`; `Add` and `Remove` use `Lock`. The tests hold `s.mu.RLock()` and check that readers still get through while `Add` waits.
- `Values` returns a **new** slice (any order; empty or nil for an empty set).
