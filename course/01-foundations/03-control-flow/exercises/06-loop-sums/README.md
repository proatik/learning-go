# Loop Sums

Implement:

```go
func SumTo(n int) int                     // 1 + 2 + … + n
func SumOfMultiples(limit, a, b int) int  // sum of k in [1, limit) divisible by a or b
func CountDivisors(n int) int             // how many d in [1, n] divide n
```

## Examples
| Call | Result | Why |
|---|---|---|
| `SumTo(4)` | `10` | 1+2+3+4 |
| `SumTo(0)` | `0` | nothing to add |
| `SumOfMultiples(10, 3, 5)` | `23` | 3+5+6+9 (10 itself is excluded) |
| `SumOfMultiples(20, 7, 7)` | `21` | 7+14 (counted once each) |
| `CountDivisors(12)` | `6` | 1, 2, 3, 4, 6, 12 |
| `CountDivisors(1)` | `1` | |

## Rules
- Any `n <= 0` or `limit <= 1` gives `0`. You may assume `a > 0` and `b > 0`.
- Use `for … := range …` over an integer in all three functions.
- In `SumOfMultiples`, use `continue` to skip numbers that are not multiples. A number divisible by both `a` and `b` is added **once**.
- Don't use the formula `n*(n+1)/2` — the point is to practise loops.
