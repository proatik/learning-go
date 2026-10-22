# Windows and Dedupe

Write an iterator that walks a slice in overlapping windows, and an adapter that remembers state between values.

```go
func Windows[T any](s []T, k int) iter.Seq[[]T]       // every run of k neighbours
func Dedupe[V comparable](seq iter.Seq[V]) iter.Seq[V] // drop consecutive duplicates
```

## Examples

| Code | Result |
|---|---|
| `slices.Collect(Windows([]int{1, 2, 3, 4}, 2))` | `[[1 2] [2 3] [3 4]]` |
| `slices.Collect(Windows([]int{1, 2, 3}, 3))` | `[[1 2 3]]` |
| `slices.Collect(Windows([]int{1, 2}, 3))` | `[]` |
| `slices.Collect(Windows([]int{1, 2}, 0))` | `[]` |
| `slices.Collect(Dedupe(slices.Values([]int{1, 1, 2, 2, 2, 1, 3, 3})))` | `[1 2 1 3]` |
| `slices.Collect(Dedupe(slices.Values([]string{"", "", "a"})))` | `["" a]` |

## Rules
- `Windows` yields `s[0:k]`, `s[1:k+1]`, …, `s[len(s)-k:]`, in order.
- `k < 1` or `k > len(s)` yields nothing (no panic).
- Each window is a sub-slice of `s` (no copying), with its **capacity clipped to `k`**, so an `append` to a window can't overwrite `s`. Use a full slice expression `s[i : i+k : i+k]`.
- `Dedupe` yields a value only when it differs from the **previous source value**. Non-neighbouring repeats are kept (`1 2 1`). The very first value is always yielded, even if it is a zero value.
- `Dedupe` is lazy, pulls only what it needs, and stops as soon as `yield` returns `false`.
