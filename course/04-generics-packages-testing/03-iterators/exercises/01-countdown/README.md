# Countdown

Write your first two iterators and respect the yield contract.

```go
func Countdown(n int) iter.Seq[int]        // n, n-1, …, 1
func Repeat[V any](v V, n int) iter.Seq[V] // v, n times
```

## Examples

| Code | Values |
|---|---|
| `slices.Collect(Countdown(3))` | `[3 2 1]` |
| `slices.Collect(Countdown(1))` | `[1]` |
| `slices.Collect(Countdown(0))` | `[]` (no values) |
| `slices.Collect(Repeat("go", 3))` | `[go go go]` |
| `slices.Collect(Repeat(7, -1))` | `[]` (no values) |
| `for v := range Countdown(1_000_000_000) { if v < 999_999_998 { break } }` | Stops after 3 values |

## Rules
- `n <= 0` produces an empty sequence. Never return a `nil` iterator.
- Stop as soon as `yield` returns `false`: the tests count how many times `yield` is called after a `break`.
- Building the iterator must do no work; values are produced only when someone ranges over it.
