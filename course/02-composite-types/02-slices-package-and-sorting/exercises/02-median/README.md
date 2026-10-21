# Median and Spread

Compute statistics **without** disturbing the caller's data.

```go
func Median(nums []float64) (float64, error)
func Spread(nums []int) (int, error)
```

- `Median`: the middle value of the sorted numbers. With an **even** count, the average of the two middle values.
- `Spread`: largest value minus smallest value.
- Both return `0` and a non-nil error for an empty (or nil) slice.

| Call | Result |
|---|---|
| `Median([]float64{3, 1, 2})` | `2, nil` |
| `Median([]float64{4, 1, 3, 2})` | `2.5, nil` |
| `Median([]float64{7})` | `7, nil` |
| `Median(nil)` | `0, <error>` |
| `Spread([]int{4, -2, 9})` | `11, nil` |
| `Spread([]int{5})` | `0, nil` |
| `Spread([]int{})` | `0, <error>` |

## Rules
- **Do not modify `nums`**: after the call, the caller's slice must be in its original order. Sort a `slices.Clone`.
- `Spread` uses `slices.Min` and `slices.Max` (remember: they panic on an empty slice, so check first).
