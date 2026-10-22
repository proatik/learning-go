# Sorted Map

A map that always lists its keys **in order**. It needs `<` on keys, so the key type is constrained by `cmp.Ordered`, while values stay `any`.

```go
type Entry[K, V any] struct {
	Key K
	Val V
}

type Map[K cmp.Ordered, V any] struct {
	keys []K     // always sorted, no duplicates
	vals map[K]V // key → value
}

func (m *Map[K, V]) Put(k K, v V)
func (m *Map[K, V]) Get(k K) (V, bool)
func (m *Map[K, V]) Delete(k K) bool
func (m *Map[K, V]) Len() int
func (m *Map[K, V]) Keys() []K
func (m *Map[K, V]) Entries() []Entry[K, V]
func (m *Map[K, V]) Range(lo, hi K) []Entry[K, V]
func (m *Map[K, V]) Min() (Entry[K, V], bool)
```

## Examples

| Code | Result |
|---|---|
| `var m Map[string, int]; m.Put("b", 2); m.Put("a", 1); m.Put("c", 3)` | `m.Keys()` = `[a b c]` |
| `m.Put("b", 20); m.Get("b")` | `20, true` (`Len()` still `3`) |
| `m.Get("zz")` | `0, false` |
| `m.Entries()` | `[{a 1} {b 20} {c 3}]` |
| `m.Range("b", "c")` | `[{b 20}]` (`lo` included, `hi` excluded) |
| `m.Min()` | `{a 1}, true` |
| `m.Delete("a")` / `m.Delete("a")` | `true` / `false` |

## Rules
- The zero `Map` is empty and usable: create the inner map on the first `Put`.
- `Put` on an existing key replaces the value without adding a duplicate key.
- `Get` on a missing key returns the zero value of `V` and `false`.
- `Keys`, `Entries` and `Range` return **new** slices in ascending key order (empty or nil when there's nothing to return).
- `Range(lo, hi)` returns the entries with `lo <= key < hi`; nothing if `lo >= hi`.
- `Min` on an empty map returns the zero `Entry` and `false`.
- Keep `keys` sorted as you go (binary search for the position); don't sort on every read.
