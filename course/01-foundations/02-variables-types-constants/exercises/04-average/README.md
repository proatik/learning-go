# Average Without the Integer-Division Trap

Implement:

```go
func IntAverage(a, b, c int) int     // integer average, truncated toward zero
func Average(a, b, c int) float64    // exact average
func Percent(part, total int) float64 // part as a percentage of total
```

| Call | Result |
|---|---|
| `IntAverage(1, 2, 2)` | `1` |
| `IntAverage(-1, -2, -2)` | `-1` (toward zero, not `-2`) |
| `Average(1, 2, 2)` | `1.6666666666666667` |
| `Average(-1, 0, 0)` | `-0.3333333333333333` |
| `Percent(1, 3)` | `33.33333333333333` |
| `Percent(7, 4)` | `175` |

## Rules
- `float64((a + b + c) / 3)` is **wrong**: the division already happened in `int`. Convert **before** you divide.
- You may assume `total != 0`.
