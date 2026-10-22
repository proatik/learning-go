# Parallel Merge Sort

Sort a slice with a **parallel merge sort**: split it in half, sort both halves at the same time, then merge them. Splitting stops at a **depth limit**, below which each piece is sorted sequentially, so the number of goroutines stays bounded.

```go
func SortFunc[T any](s []T, cmp func(a, b T) int, depth int)
func Sort[T cmp.Ordered](s []T, depth int)
```

## Examples

| Call | Result |
|---|---|
| `Sort([]int{5, 2, 8, 1, 9, 3}, 1)` | `[1 2 3 5 8 9]`: halves `[5 2 8]` and `[1 9 3]` sorted concurrently, then merged |
| `Sort([]string{"b", "c", "a"}, 0)` | `[a b c]`: depth 0 is fully sequential |
| `SortFunc(people, byAge, 3)` | sorted by `byAge`, with up to 8 pieces sorted at once |
| `Sort([]int{}, 4)` / `Sort[int](nil, 4)` | nothing to do |
| `Sort([]int{2, 1}, 10)` | `[1 2]`: depth larger than needed is fine |

## Rules
- Sort `s` **in place**, ascending by `cmp` (negative: `a` first, `0`: equal, positive: `b` first). The result must equal `slices.SortFunc(s, cmp)`'s; stability is not required.
- `depth <= 0` or `len(s) < 2`: sort sequentially with `slices.SortFunc` and return.
- Otherwise: split at `mid := len(s) / 2`, sort `s[:mid]` and `s[mid:]` **concurrently** with `depth-1`, wait for both, then **merge** them (through a temporary slice) back into `s`.
- So depth `d` sorts up to `2^d` pieces at the same time. Never more: the tests watch how many `cmp` calls run at once.
- Every goroutine you start has finished when `SortFunc` returns.
- `Sort` is `SortFunc` with `cmp.Compare`.
