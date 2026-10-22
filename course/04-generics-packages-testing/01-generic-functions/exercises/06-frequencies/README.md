# Frequencies and Sorted Keys

Count values, list map keys in a stable order, and find the most common value, for any key type.

```go
func Frequencies[T comparable](xs []T) map[T]int
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K
func MostCommon[T cmp.Ordered](xs []T) (value T, count int, ok bool)
```

## Examples

| Call | Result |
|---|---|
| `Frequencies([]string{"a", "b", "a"})` | `map[a:2 b:1]` |
| `Frequencies([]int{})` | `map[]` (empty, **not nil**) |
| `SortedKeys(map[string]int{"b": 1, "a": 2, "c": 0})` | `[a b c]` |
| `SortedKeys(map[int]bool{})` | `[]` (length 0) |
| `MostCommon([]string{"go", "rust", "go"})` | `"go", 2, true` |
| `MostCommon([]int{3, 1, 3, 1, 2})` | `1, 2, true` (tie: smallest wins) |
| `MostCommon([]int{})` | `0, 0, false` |

## Rules
- `Frequencies`: a new map from each distinct value to how many times it appears. Always non-nil.
- `SortedKeys`: the keys of `m` in **ascending** order. Map iteration order is random, so sort the result. Works for any value type `V`.
- `MostCommon`:
  - the value with the highest count, that count, and `true`
  - on a tie, the **smallest** value wins
  - empty or nil → zero value, `0`, `false`
  - Build it from `Frequencies` and `SortedKeys`.
- `slices.Sort` is allowed.
