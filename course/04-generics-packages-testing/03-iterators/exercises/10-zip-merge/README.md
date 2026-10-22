# Zip and Merge

Walk two sequences side by side with `iter.Pull`, and always clean up.

```go
func Zip[A, B any](as iter.Seq[A], bs iter.Seq[B]) iter.Seq2[A, B] // pairs, until either ends
func Merge[V cmp.Ordered](a, b iter.Seq[V]) iter.Seq[V]            // two sorted sequences → one sorted
```

## Examples

| Code | Result |
|---|---|
| `for s, n := range Zip(slices.Values([]string{"x", "y", "z"}), slices.Values([]int{1, 2}))` | `(x, 1)`, `(y, 2)` |
| `Zip(names, naturals)` with `naturals` infinite | One pair per name |
| `slices.Collect(Merge(slices.Values([]int{1, 4, 9}), slices.Values([]int{2, 3, 10, 11})))` | `[1 2 3 4 9 10 11]` |
| `slices.Collect(Merge(slices.Values([]int{1, 2}), slices.Values([]int{2, 3})))` | `[1 2 2 3]` |
| `slices.Collect(Merge(empty, slices.Values([]int{5})))` | `[5]` |

## Rules
- `Zip` stops as soon as either sequence ends. It pulls at most one value from each sequence per pair.
- `Merge` assumes both inputs are sorted ascending and yields every value of both, in ascending order. On a tie, yield `a`'s value first. When one input ends, yield the rest of the other.
- `Merge` pulls lazily: a value is read only when it is needed for the next comparison, so infinite sorted inputs work.
- **Cleanup:** when your iterator returns, for **any** reason (an input ended, or the loop body stopped), every input sequence must have finished or been stopped, so its deferred cleanup has run. With `iter.Pull` that means `defer stop()`.
- Everything is lazy: building the iterator starts nothing. Stop as soon as `yield` returns `false`.
