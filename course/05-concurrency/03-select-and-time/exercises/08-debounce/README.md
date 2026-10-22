# Debounce

Pass on only the **latest** value once the input has been quiet for `d`. A search box does this: it waits until you stop typing before it sends a query.

```go
func Debounce[T any](in <-chan T, d time.Duration) <-chan T
```

## Examples

All with `d = 200ms`.

| Input (value@time) | `in` closed at | Output (value@time) | Output closed at |
|---|---|---|---|
| `1`@0, `2`@50ms, `3`@100ms, `4`@400ms, `5`@450ms | 1s | `3`@300ms, `5`@650ms | 1s |
| `1`@0, `2`@300ms, `3`@600ms | 2s | `1`@200ms, `2`@500ms, `3`@800ms | 2s |
| `1`@0 | 100ms | `1`@100ms (flushed on close) | 100ms |
| nothing | 500ms | nothing | 500ms |

## Rules
- Return the output channel at once; one goroutine does the work. The output is unbuffered.
- Every input value replaces the pending one and restarts the quiet period. When `d` passes with no new input, send the pending value.
- When `in` is closed: send the pending value (if any) **immediately**, then close the output.
- Never send the same pending value twice, and send nothing when nothing is pending.
- Tip: keep the timer's channel in a variable that is `nil` while nothing is pending, so that `select` case is switched off.
- The tests never send a value at the same instant a quiet period ends.
