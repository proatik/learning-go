# Filter, Map and Reduce

Write the three classic generic sequence functions: two lazy adapters and one consumer.

```go
func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V]
func Map[In, Out any](seq iter.Seq[In], f func(In) Out) iter.Seq[Out]
func Reduce[V, A any](seq iter.Seq[V], init A, f func(A, V) A) A
```

## Examples

| Code | Result |
|---|---|
| `slices.Collect(Filter(slices.Values([]int{1, 2, 3, 4}), isEven))` | `[2 4]` |
| `slices.Collect(Map(slices.Values([]string{"go", "hi"}), strings.ToUpper))` | `[GO HI]` |
| `slices.Collect(Map(slices.Values([]string{"a", "bcd"}), func(s string) int { return len(s) }))` | `[1 3]` |
| `Reduce(slices.Values([]int{1, 2, 3}), 0, func(sum, v int) int { return sum + v })` | `6` |
| `Reduce(slices.Values([]int{}), 10, add)` | `10` |
| `Reduce(Map(Filter(nums, isEven), square), 0, add)` | sum of the even squares |

## Rules
- `Filter` yields the values for which `keep` returns `true`, in order.
- `Map` yields `f(v)` for each `v`, in order. `In` and `Out` may be different types.
- **Lazy:** calling `Filter` or `Map` must not call `seq`, `keep` or `f`. Work happens only while someone ranges over the result.
- `keep` and `f` are called **exactly once** per source value the consumer reaches, never again for the same value.
- A `break` in the consumer stops the source: no more source values are produced.
- `Reduce` starts with `acc = init`, sets `acc = f(acc, v)` for every value, and returns `acc` (so an empty sequence returns `init`).
