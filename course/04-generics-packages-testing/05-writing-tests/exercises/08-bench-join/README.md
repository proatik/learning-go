# A benchmark

Write `TestJoinInts`, which catches every broken version of `JoinInts`, and `BenchmarkJoinInts`, which measures it with `b.Loop()`.

```go
// JoinInts formats nums in decimal, separated by sep.
func JoinInts(nums []int, sep string) string
```

`JoinInts` is already written (in the read-only checker). You write the test and the benchmark in `joiner_test.go`.

## Examples

| Call | Result |
|---|---|
| `JoinInts(nil, ",")` | `""` |
| `JoinInts([]int{7}, ",")` | `"7"` |
| `JoinInts([]int{1, 2, 3}, ", ")` | `"1, 2, 3"` |
| `JoinInts([]int{-1, 0, 42}, "/")` | `"-1/0/42"` |

## Rules

- `TestJoinInts`: table-driven, reported as `"JoinInts(%v, %q) = %q, want %q"`. It must pass on the correct code and fail for every planted bug.
- `BenchmarkJoinInts`:
  - build the input (for example 100 numbers) **before** the loop
  - call `b.ReportAllocs()`
  - call `JoinInts` inside `for b.Loop() { ... }`
- **The checker** runs the benchmark once (`-bench=. -benchtime=1x`). It must succeed, call `JoinInts`, use `b.Loop()`, and print an `allocs/op` column.
