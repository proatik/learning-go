# Safe Divide

Let Go's run-time panic happen, then **recover** it and turn it into an `error`.

```go
func SafeDivide(a, b int) (q int, err error)
func Average(nums []int) (avg int, err error)
```

| Call | Result |
|---|---|
| `SafeDivide(7, 2)` | `3, nil` |
| `SafeDivide(-7, 2)` | `-3, nil` |
| `SafeDivide(5, 0)` | `0`, error `safe divide: runtime error: integer divide by zero` |
| `SafeDivide(math.MinInt, -1)` | `math.MinInt, nil` (wraps, no panic) |
| `Average([]int{2, 4, 9})` | `5, nil` |
| `Average(nil)` | `0`, error `average: runtime error: integer divide by zero` |

## Rules
- Do **not** check `b == 0` or `len(nums) == 0`: compute `a / b` (and `sum / len(nums)`) directly and handle the panic in a deferred function.
- Error text: `safe divide: <panic value>` and `average: <panic value>` (use `%v`).
- On error, the numeric result is `0`.
- Integer division truncates toward zero, as in plain Go.
