# Collect Errors

Run tasks that can fail concurrently, run **all** of them even when some fail, and return their errors joined **in input order**.

```go
func RunAll(tasks []func() error) error
```

## Examples

| Tasks (each sleeps 1s first) | Result |
|---|---|
| ok, ok, ok | `nil`, after **1s** |
| ok, `errA`, ok, `errB` | an error with message `"a failed\nb failed"`; `errors.Is(err, errA)` and `errors.Is(err, errB)` are `true` |
| `errB` (finishes last), `errA` (finishes first) | message `"b failed\na failed"`: input order, not finishing order |
| none (`nil` or empty) | `nil` |

## Rules
- Run every task in its own goroutine and wait for **all** of them; a failure doesn't stop the others.
- Store task `i`'s error in slot `i` of a slice made up front, then combine them with `errors.Join` after `Wait`.
- Return `nil` when no task failed (`errors.Join` already does this).
- The joined error must keep each original error reachable with `errors.Is`.
