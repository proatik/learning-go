# Your Own errgroup

Write a small version of the widely used `errgroup` package: run tasks concurrently, wait for all of them, return the **first** error, and cancel the shared context the moment a task fails.

```go
type Group struct {
	wg     sync.WaitGroup
	once   sync.Once               // guards err: only the first error is kept
	err    error                   // the first non-nil error returned by a task
	cancel context.CancelCauseFunc // nil for the zero Group
}

func WithContext(ctx context.Context) (*Group, context.Context)
func (g *Group) Go(f func() error)
func (g *Group) Wait() error
```

```text
g, ctx := WithContext(parent)
g.Go(fetchUser)   ──┐
g.Go(fetchOrders) ──┼── first error ──▶ cancel(ctx, err) ──▶ the other tasks see ctx.Done() and stop
g.Go(fetchPrefs)  ──┘
err := g.Wait()   ◀── waits for all three, returns that first error
```

## Examples

| Code | Result |
|---|---|
| 5 tasks of 1s returning `nil`, then `Wait()` | `nil`, after 1s (they run at the same time) |
| tasks fail at 2s (`errA`) and 1s (`errB`), a third succeeds at 3s | `Wait()` = `errB`, returned at 3s |
| one task fails at 1s, another waits on `ctx.Done()` | the second stops at 1s; `Wait()` = the first task's error |
| after that failure | `ctx.Err()` = `context.Canceled`, `context.Cause(ctx)` = the error |
| no failures, after `Wait()` | `ctx.Err()` = `context.Canceled`, `context.Cause(ctx)` = `context.Canceled` |
| `var g Group; g.Go(f); g.Wait()` | works: the zero Group just has no context to cancel |

## Rules
- `WithContext` derives the group's context with `context.WithCancelCause` and keeps the cancel function.
- `Go` starts `f` in a new goroutine at once.
- Only the **first** non-nil error is kept (`sync.Once`). Return it unchanged, so `errors.Is` still works on it. Then cancel the context with that error as the cause.
- `Wait` waits for **every** task, even after a failure, then cancels the context (with a `nil` cause if nothing failed) and returns the first error or `nil`.
- No failure means no cancel before `Wait`: `ctx.Err()` stays `nil` until then.
- The zero `Group` works: `cancel` is `nil`, so check before calling it.
