# Counters with Closures

Each function below returns a **closure** that keeps its own private state between calls.

```go
func Counter() func() int
func CountBy(start, step int) func() int
func Accumulator(initial int) func(int) int
```

| Code | Calls return |
|---|---|
| `c := Counter()`; `c()`, `c()`, `c()` | `1`, `2`, `3` |
| `c := CountBy(10, 5)`; `c()`, `c()`, `c()` | `10`, `15`, `20` |
| `c := CountBy(3, -1)`; `c()`, `c()`, `c()` | `3`, `2`, `1` |
| `acc := Accumulator(100)`; `acc(10)`, `acc(-30)`, `acc(0)` | `110`, `80`, `80` |

## Rules
- Every call to `Counter`, `CountBy` or `Accumulator` creates **independent** state: two counters never affect each other.
- No package-level variables: the state lives in a local variable captured by the closure.
- `CountBy` returns `start` on the **first** call.
- `Accumulator` adds its argument to a running total and returns the new total.
