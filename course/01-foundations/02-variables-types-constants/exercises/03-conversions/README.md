# Type Conversions

Go never mixes `int` and `float64` automatically. You must convert with `T(v)`.

Implement:

```go
func BMI(weightKg float64, heightCm int) float64 // weight / (height in metres)²
func Truncate(f float64) int                     // drop the fractional part
func Pixels(inches float64, dpi int) int         // inches × dpi, fraction dropped
```

| Call | Result |
|---|---|
| `BMI(70, 175)` | `22.857142857142858` |
| `BMI(90, 200)` | `22.5` |
| `Truncate(3.99)` | `3` |
| `Truncate(-3.99)` | `-3` (toward zero, **not** down) |
| `Pixels(2.5, 96)` | `240` |
| `Pixels(1.1, 72)` | `79` (79.2 truncated) |

## Rules
- Heights are in **centimetres**: divide by `100` to get metres. Make sure that division happens in `float64`, not `int` (`175 / 100` is `1` in integer arithmetic!).
- Use conversions (`float64(x)`, `int(x)`) only. No `math` functions needed.
