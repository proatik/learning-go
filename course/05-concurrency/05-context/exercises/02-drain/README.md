# Drain Until Canceled

Write the standard worker loop: take jobs from a channel and handle them until the channel is closed **or** the context is done.

```go
func Drain(ctx context.Context, in <-chan string, handle func(string)) (int, error)
```

## Examples

| Situation | Returns |
|---|---|
| `in` holds `"a"`, `"b"`, `"c"`, then is closed | `3, nil` (handled `a`, `b`, `c` in order) |
| `in` is closed and empty | `0, nil` |
| `ctx` already canceled, `in` holds 5 jobs | `0, context.Canceled` (nothing handled) |
| `handle` cancels `ctx` while handling the 2nd of 5 jobs | `2, context.Canceled` |
| `in` stays open and empty; `ctx` times out after 2s | `0, context.DeadlineExceeded` after 2s |
| `in` is `nil`; `cancel()` after 1s | `0, context.Canceled` after 1s |

## Rules
- Call `handle` once per job, in the order received. Return how many jobs were handled.
- When `in` is closed, return the count and `nil`.
- When `ctx` is done, return the count and `ctx.Err()`, even while waiting for a job that never comes.
- **Cancellation wins:** if `ctx` is already done at the start of an iteration, stop without taking another job, even if jobs are waiting. (`select` alone picks randomly between ready cases.)
- A `nil` `in` never delivers anything, so only `ctx` can end the loop.
