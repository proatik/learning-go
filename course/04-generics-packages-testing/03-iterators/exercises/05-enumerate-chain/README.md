# Enumerate, Chain, Keys and Values

Move between `iter.Seq` and `iter.Seq2`, and join several sequences into one.

```go
func Enumerate[V any](seq iter.Seq[V]) iter.Seq2[int, V]  // (0, v0), (1, v1), …
func Chain[V any](seqs ...iter.Seq[V]) iter.Seq[V]        // every value of seqs[0], then seqs[1], …
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K]      // the first value of each pair
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V]    // the second value of each pair
```

## Examples

| Code | Result |
|---|---|
| `for i, s := range Enumerate(slices.Values([]string{"a", "b"}))` | `0 a`, `1 b` |
| `slices.Collect(Chain(slices.Values([]int{1, 2}), slices.Values([]int{3})))` | `[1 2 3]` |
| `slices.Collect(Chain[int]())` | `[]` |
| `slices.Collect(Keys(slices.All([]string{"x", "y"})))` | `[0 1]` |
| `slices.Collect(Values(slices.All([]string{"x", "y"})))` | `[x y]` |
| `slices.Collect(Values(Enumerate(seq)))` | the same values as `seq` |

## Rules
- `Enumerate` counts from `0` in the order the values arrive.
- `Chain` runs the sequences one after another. A `nil` sequence in `seqs` is skipped. With no sequences the result is empty.
- `Chain` starts a sequence only when the previous one has finished: after a `break` in the first sequence's values, the later sequences are **never called**.
- Everything is lazy: building the iterator calls nothing.
- Stop as soon as `yield` returns `false`. Careful with the nested loops in `Chain`: `break` only leaves the inner one.
