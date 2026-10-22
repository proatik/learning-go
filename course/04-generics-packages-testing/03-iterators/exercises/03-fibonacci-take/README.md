# Fibonacci and Take

Write an **infinite** sequence and the generic adapter that makes it safe to consume.

```go
func Fibonacci() iter.Seq[int]                       // 0, 1, 1, 2, 3, 5, 8, … forever
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] // the first n values of seq
```

## Examples

| Code | Result |
|---|---|
| `slices.Collect(Take(Fibonacci(), 8))` | `[0 1 1 2 3 5 8 13]` |
| `slices.Collect(Take(slices.Values([]string{"a", "b"}), 5))` | `[a b]` (the source ran out first) |
| `slices.Collect(Take(Fibonacci(), 0))` | `[]` |
| `for v := range Fibonacci() { if v > 50 { break } }` | Stops at `55` |

## Rules
- `Fibonacci` starts `0, 1` and each next value is the sum of the previous two. It never ends by itself (ignore `int` overflow).
- `Take` yields at most `n` values. `n <= 0` yields nothing.
- `Take` must **pull no more than it uses**: after the `n`-th value it returns without asking `seq` for another one. With `n <= 0` it never calls `seq`.
- `Take(seq, n)` only builds the iterator; `seq` runs when someone ranges over the result.
- Both stop as soon as `yield` returns `false`.
