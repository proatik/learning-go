# Statistics for Any Number

Compute mean, variance, median and a summary for **any** numeric slice: `[]int`, `[]uint8`, `[]float32`, `[]Celsius`... Results are always `float64`.

```go
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Summary struct {
	Count                         int
	Min, Max, Mean, Median, StdDev float64
}

func Mean[T Number](xs []T) (float64, bool)
func Variance[T Number](xs []T) (float64, bool)
func Median[T Number](xs []T) (float64, bool)
func Summarize[T Number](xs []T) (Summary, bool)
```

## Examples

| Call | Result |
|---|---|
| `Mean([]int{1, 2, 3, 4})` | `2.5, true` |
| `Mean([]uint8{200, 200})` | `200, true` (not `72`: see Rules) |
| `Variance([]int{2, 4, 4, 4, 5, 5, 7, 9})` | `4, true` |
| `Variance([]uint{1, 3})` | `1, true` |
| `Median([]int{3, 1, 2})` | `2, true` |
| `Median([]int{4, 1, 3, 2})` | `2.5, true` |
| `Median([]Celsius{21.5})` | `21.5, true` |
| `Mean([]int{})` (also `Variance`, `Median`) | `0, false` |
| `Summarize([]int{2, 4, 4, 4, 5, 5, 7, 9})` | `{Count:8 Min:2 Max:9 Mean:5 Median:4.5 StdDev:2}, true` |

## Rules
- Empty or nil input → `0` (or `Summary{}`) and `false`. Otherwise `true`.
- **Convert each element to `float64` before doing arithmetic.** Adding in `T` overflows small types (`uint8(200) + uint8(200)` is `144`), and `x - mean` in an unsigned `T` wraps around.
- `Mean`: sum / count.
- `Variance`: the **population** variance, the mean of `(x - mean)²`.
- `Median`: the middle value of the **sorted** values. For an even count, the mean of the two middle values.
  - **Don't modify `xs`**: sort a copy (`slices.Clone` + `slices.Sort`; every `Number` type is ordered, so `slices.Sort` accepts it).
- `Summarize`: `Count = len(xs)`, `Min` / `Max` as `float64`, and `StdDev = math.Sqrt(variance)`. Reuse the other three functions.
