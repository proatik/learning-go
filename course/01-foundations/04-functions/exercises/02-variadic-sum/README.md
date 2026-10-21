# Variadic Sum and Mean

Implement two variadic functions:

```go
func Sum(nums ...int) int
func Mean(nums ...int) (float64, error)
```

| Call | Result |
|---|---|
| `Sum()` | `0` |
| `Sum(4)` | `4` |
| `Sum(1, 2, 3)` | `6` |
| `Sum(xs...)` with `xs := []int{10, -5}` | `5` |
| `Mean(1, 2, 3, 4)` | `2.5, nil` |
| `Mean(-3, 3)` | `0, nil` |
| `Mean()` | `0, <non-nil error>` |

## Rules
- `Sum` loops over `nums` with `range`.
- `Mean` must **reuse** `Sum` by spreading its slice: `Sum(nums...)`.
- `Mean` of no numbers is an error (you can't divide by zero); return `0` with it.
- Convert to `float64` **before** dividing, or `7 / 2` becomes `3`.
