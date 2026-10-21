# Running Prefix Sums

Build a new slice of running totals, then use it to answer range-sum questions in O(1):

```go
func PrefixSums(s []int) []int
func RangeSum(prefix []int, lo, hi int) int
```

`PrefixSums(s)[i]` is `s[0] + s[1] + ... + s[i]`. `RangeSum(PrefixSums(s), lo, hi)` is the sum of `s[lo:hi]` (half-open, like a slice expression).

| Call | Result |
|---|---|
| `PrefixSums([]int{1, 2, 3, 4})` | `[1 3 6 10]` |
| `PrefixSums([]int{5, -5, 5})` | `[5 0 5]` |
| `PrefixSums(nil)` | empty slice (`len 0`) |
| `RangeSum([]int{1, 3, 6, 10}, 1, 3)` | `5` (= 2 + 3) |
| `RangeSum([]int{1, 3, 6, 10}, 0, 4)` | `10` |
| `RangeSum([]int{1, 3, 6, 10}, 2, 2)` | `0` (empty range) |

## Rules
- `PrefixSums` must **not** modify `s`. Allocate the result once with `make([]int, len(s))`.
- `RangeSum` must not loop: use at most two lookups into `prefix`.
- You may assume `0 <= lo <= hi <= len(prefix)`.
