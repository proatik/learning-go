# Ordered Map

Build a generic map that **remembers insertion order**, with O(1) `Get`/`Set`/`Delete` and iterators that walk it forwards and backwards.

```go
type entry[K comparable, V any] struct {
	key        K
	value      V
	prev, next *entry[K, V]
}

type OrderedMap[K comparable, V any] struct {
	index map[K]*entry[K, V]
	head  *entry[K, V] // oldest
	tail  *entry[K, V] // newest
}

func New[K comparable, V any]() *OrderedMap[K, V]
func Collect[K comparable, V any](seq iter.Seq2[K, V]) *OrderedMap[K, V]

func (m *OrderedMap[K, V]) Set(key K, value V)
func (m *OrderedMap[K, V]) Get(key K) (V, bool)
func (m *OrderedMap[K, V]) Delete(key K) bool
func (m *OrderedMap[K, V]) Len() int
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V]
func (m *OrderedMap[K, V]) Keys() iter.Seq[K]
func (m *OrderedMap[K, V]) Backward() iter.Seq2[K, V]
```

## Examples

| Code | Result |
|---|---|
| `m := New[string, int](); m.Set("b", 2); m.Set("a", 1)` | `All()` yields `b 2`, `a 1` |
| `m.Set("b", 20)` | `All()` yields `b 20`, `a 1` — position kept |
| `m.Delete("b"); m.Set("b", 5)` | `All()` yields `a 1`, `b 5` — re-added at the end |
| `slices.Collect(m.Keys())` | `[a b]` |
| `m.Backward()` | yields `b 5`, `a 1` |
| `Collect(slices.All([]string{"x", "y"}))` | a map `0→x, 1→y` |

## Rules
- **New** creates the `index` map. A missing key in `Get` returns the zero `V` and `false`.
- **Set** on a new key appends it at the end; on an existing key it only updates the value.
- **Delete** unlinks the entry in O(1) and reports whether the key was there.
- **All**, **Keys** and **Backward** must honour the yield contract (stop as soon as `yield` returns `false`).
- Deleting the **current** key inside a `for k := range m.All()` loop must not skip or repeat any other key.
- **Collect** keeps each key's **first** position and **last** value.
