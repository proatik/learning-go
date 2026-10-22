# First Success Wins

Ask several replicas for the same answer at once, return the **first** successful reply, and cancel the others so they stop wasting work.

```go
var ErrNoFetchers = errors.New("first: no fetchers")

func First(ctx context.Context, fetchers ...func(context.Context) (string, error)) (string, error)
```

## Examples

| Fetchers | Returns | After |
|---|---|---|
| `a` = "A" after 3s, `b` = "B" after 1s, `c` = "C" after 2s | `"B", nil`; `a` and `c` see their ctx canceled at 1s | 1s |
| `a` fails after 100ms, `b` = "B" after 2s, `c` = "C" after 3s | `"B", nil` | 2s |
| all three fail (after 3s, 1s, 2s) | `"", errors.Join(errA, errB, errC)` | 3s |
| none | `"", ErrNoFetchers` | 0s |
| `ctx` times out after 1s, fetchers need 5s | `"", err` with `errors.Is(err, context.DeadlineExceeded)` | 1s |

## Rules
- Start every fetcher in its **own goroutine** at once, all with one shared child of `ctx` (`context.WithCancel`).
- Return as soon as one fetcher succeeds (`err == nil`), and **cancel** the child so the others stop. Their late results are ignored.
- If every fetcher fails, return `""` and `errors.Join` of all the errors **in argument order** (not finishing order).
- Fetchers are expected to return soon after their ctx is canceled. Don't leave any goroutine stuck: a loser that finishes after you've returned must not block forever on a send (a buffered channel with room for every result fixes that).
- The child must be canceled before `First` returns, on every path.
