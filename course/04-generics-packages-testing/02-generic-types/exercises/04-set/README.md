# Generic Set

A `Set[T comparable]` built on a map, with the usual set algebra. A method can't ask for more than `comparable`, so sorting is a **package function** with a tighter constraint.

```go
type Set[T comparable] map[T]struct{}

func New[T comparable](vals ...T) Set[T]

func (s Set[T]) Add(vals ...T)
func (s Set[T]) Remove(v T) bool
func (s Set[T]) Has(v T) bool
func (s Set[T]) Len() int
func (s Set[T]) Union(o Set[T]) Set[T]
func (s Set[T]) Intersect(o Set[T]) Set[T]
func (s Set[T]) Difference(o Set[T]) Set[T]
func (s Set[T]) SubsetOf(o Set[T]) bool
func (s Set[T]) Equal(o Set[T]) bool

func Sorted[T cmp.Ordered](s Set[T]) []T
```

## Examples

| Code | Result |
|---|---|
| `New(3, 1, 3, 2).Len()` | `3` (duplicates collapse) |
| `New("a").Has("a")`, `New("a").Has("b")` | `true`, `false` |
| `s := New(1, 2); s.Remove(2)` / `s.Remove(9)` | `true` / `false` |
| `Sorted(New(1, 2, 3).Union(New(3, 4)))` | `[1 2 3 4]` |
| `Sorted(New(1, 2, 3).Intersect(New(2, 3, 4)))` | `[2 3]` |
| `Sorted(New(1, 2, 3).Difference(New(2)))` | `[1 3]` |
| `New(1, 2).SubsetOf(New(1, 2, 3))` | `true` |
| `New(1, 2).Equal(New(2, 1))` | `true` |
| `Sorted(New("pear", "Apple", "fig"))` | `["Apple" "fig" "pear"]` |

## Rules
- `New` with no arguments returns an empty, **non-nil** set ready for `Add`.
- `Remove` reports whether `v` was present.
- `Union`, `Intersect` and `Difference` return a **new non-nil** set and never modify `s` or `o`.
- A **nil** `Set` is a valid empty set for every method except `Add` and `Remove` (reading a nil map is fine; writing panics).
- The empty set is a subset of every set. `Equal` means same elements.
- `Sorted` returns the elements in ascending order (empty or nil for an empty set).
