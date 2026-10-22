# Parallel Sum

Split a slice into `k` chunks, sum each chunk in **its own goroutine**, and combine the partial sums.

```go
// sumChunk adds up one chunk. It's given; call it once per chunk.
// (The tests replace it with a slow version to check the chunks run concurrently.)
var sumChunk = func(part []int) int { … }

func Sum(nums []int, k int) (total int, partials []int)
```

## Examples

| Call | `partials` | `total` |
|---|---|---|
| `Sum([]int{1, 2, 3, 4, 5, 6}, 3)` | `[3 7 11]` (chunks `[1 2] [3 4] [5 6]`) | `21` |
| `Sum([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3)` | `[6 15 34]` (chunks `[1 2 3] [4 5 6] [7 8 9 10]`) | `55` |
| `Sum([]int{5, -5, 7}, 1)` | `[7]` | `7` |
| `Sum([]int{1, 2}, 10)` | `[1 2]` (`k` lowered to `len(nums)`) | `3` |
| `Sum([]int{4, 4}, 0)` / `Sum([]int{4, 4}, -2)` | `[8]` (`k` raised to `1`) | `8` |
| `Sum(nil, 3)` / `Sum([]int{}, 3)` | empty | `0` |

## Rules
- Clamp `k`: if `k > len(nums)`, use `len(nums)`; then if `k < 1`, use `1`. An empty `nums` gives `0` and an empty `partials` (no goroutines, no `sumChunk` calls).
- Chunk `i` (for `i` from `0` to `k-1`) is `nums[i*n/k : (i+1)*n/k]` with `n = len(nums)`. Chunks are never empty.
- Start one goroutine per chunk; it stores `sumChunk(chunk)` in `partials[i]`. `partials` has length `k`, in chunk order.
- `total` is the sum of `partials`, computed after every goroutine has finished.
