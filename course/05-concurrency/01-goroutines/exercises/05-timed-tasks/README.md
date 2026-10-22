# Timed Tasks

Run tasks concurrently and measure two things: how long **each** task took, and how long the whole batch took.

```go
func Timed(tasks []func()) (total time.Duration, each []time.Duration)
```

## Examples

| Tasks | `each` | `total` |
|---|---|---|
| sleep 1s, sleep 3s, sleep 2s | `[1s 3s 2s]` | `3s` (the longest, not the sum `6s`) |
| 5 tasks that sleep 1s | `[1s 1s 1s 1s 1s]` | `1s` |
| sleep 2s, return at once | `[2s 0s]` | `2s` |
| none (`nil` or empty) | empty | `0s` |

## Rules
- Run every task in its own goroutine, all at the same time.
- `each[i]` is the time task `i` itself took: `time.Now()` just before calling it and `time.Since` just after, **inside** its goroutine. `each` has `len(tasks)` elements, in input order.
- `total` is measured around the whole batch: from before the first goroutine starts until `Wait` returns.
- Return only after every task has finished.
- The tests use a fake clock (lesson 3), so these durations are exact.
