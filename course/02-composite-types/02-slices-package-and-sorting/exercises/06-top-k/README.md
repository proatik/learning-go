# Top-K and K-th Smallest

Leaderboard-style questions answered by sorting a copy.

```go
func TopK(nums []int, k int) []int
func KthSmallest(nums []int, k int) (int, error)
func KthLargestDistinct(nums []int, k int) (int, error)
```

| Function | Returns |
|---|---|
| `TopK` | the `k` largest numbers, **largest first**. `k > len(nums)`: all of them. `k <= 0`: empty slice |
| `KthSmallest` | the `k`-th smallest number, counting from **1** (duplicates count separately) |
| `KthLargestDistinct` | the `k`-th largest **distinct** value, counting from 1 |

`KthSmallest` and `KthLargestDistinct` return `0` and an error when there is no such element (`k < 1` or `k` too big).

| Call | Result |
|---|---|
| `TopK([]int{5, 1, 9, 3, 9}, 3)` | `[9 9 5]` |
| `TopK([]int{2, 1}, 10)` | `[2 1]` |
| `TopK([]int{2, 1}, 0)` | `[]` |
| `KthSmallest([]int{7, 2, 9, 2}, 1)` | `2, nil` |
| `KthSmallest([]int{7, 2, 9, 2}, 2)` | `2, nil` |
| `KthSmallest([]int{7, 2, 9, 2}, 3)` | `7, nil` |
| `KthSmallest([]int{7, 2, 9, 2}, 5)` | `0, <error>` |
| `KthLargestDistinct([]int{7, 2, 9, 2, 9}, 2)` | `7, nil` |
| `KthLargestDistinct([]int{4, 4, 4}, 2)` | `0, <error>` |

## Rules
- Never modify `nums`.
- `TopK` returns a slice that does **not** share memory with `nums`.
- `KthLargestDistinct`: sort, `slices.Compact`, then index from the end (or `slices.Reverse`).
