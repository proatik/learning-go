# Group, Partition, Pick the Best

Three helpers driven by a function you pass in: group by a key, split by a predicate, and pick the element with the largest key.

```go
func GroupBy[T any, K comparable](xs []T, key func(T) K) map[K][]T
func Partition[T any](xs []T, pred func(T) bool) (yes, no []T)
func MaxBy[T any, K cmp.Ordered](xs []T, key func(T) K) (T, bool)
```

## Examples

```go
type Order struct {
	ID     int
	City   string
	Amount int
}

orders := []Order{{1, "Oslo", 50}, {2, "Rome", 80}, {3, "Oslo", 20}}
```

| Call | Result |
|---|---|
| `GroupBy(orders, func(o Order) string { return o.City })` | `map[Oslo:[#1 #3] Rome:[#2]]` |
| `GroupBy([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })` | `map[false:[1 3] true:[2 4]]` |
| `Partition([]int{1, 2, 3, 4}, isEven)` | `[2 4], [1 3]` |
| `MaxBy(orders, func(o Order) int { return o.Amount })` | `{2 Rome 80}, true` |
| `MaxBy([]string{"go", "rust", "zig"}, func(s string) int { return len(s) })` | `"rust", true` |
| `MaxBy([]Order{}, amount)` | `Order{}, false` |

## Rules
- `GroupBy`: a non-nil map from each key to the elements with that key, **in their original order**. Empty input → empty map.
- `Partition`: `yes` holds the elements where `pred` is `true`, `no` the rest, both in original order. Every element ends up in exactly one of them.
- `MaxBy`:
  - the element with the **largest** `key`, and `true`
  - on a tie, the **first** such element wins
  - empty or nil → zero value and `false`
- Call `key` / `pred` exactly once per element. Never modify `xs`.
