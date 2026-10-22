# AfterFunc Bridges

Use `context.AfterFunc` to connect contexts to things that don't understand them: a blocking `Read` that has no `ctx` parameter, and a **second** context.

```go
var ErrClosed = errors.New("reader closed")

// Reader is a blocking source that knows nothing about contexts.
type Reader interface {
	Read() (string, error) // blocks until a message arrives; returns ErrClosed once Close is called
	Close() error
}

func ReadCtx(ctx context.Context, r Reader) (string, error)
func Merge(a, b context.Context) (context.Context, context.CancelFunc)
```

## Examples

| Code | Result |
|---|---|
| `ReadCtx(ctx, r)`, a message `"hi"` arrives after 1s | `"hi", nil` after 1s; `r` is **not** closed, not even when `ctx` is canceled later |
| `ReadCtx(ctx, r)`, no message, `cancel()` after 2s | `"", context.Canceled` after 2s; `r.Close()` was called once |
| `ReadCtx(ctx, r)` with a 500ms timeout, no message | `"", context.DeadlineExceeded` after 500ms |
| `ReadCtx(ctx, r)` with `ctx` already canceled | `"", context.Canceled`, even if a message was waiting |
| `ReadCtx(ctx, r)`, `Read` fails with `ErrBroken` | `"", ErrBroken`; `r` not closed |
| `m, _ := Merge(a, b)`, then `a` canceled with cause `X` | `m.Err()` = `Canceled`, `context.Cause(m)` = `X` |
| `m, _ := Merge(a, b)`, then `b` canceled with cause `Y` | `m.Err()` = `Canceled`, `context.Cause(m)` = `Y` |
| `m, _ := Merge(a, b)`, `b` times out after 2s | after 2s: `m.Err()` = `Canceled`, `context.Cause(m)` = `DeadlineExceeded` |
| `m, cancel := Merge(a, b)`; `cancel()` | `m.Err()` and `context.Cause(m)` = `Canceled`; `a` and `b` untouched |

## Rules
- **ReadCtx:** register `stop := context.AfterFunc(ctx, …)` that closes `r`, then call `r.Read()` once, then call `stop()`.
  - `stop()` returns `false` → the context won (it was canceled and `r` is being closed): return `""` and `ctx.Err()`, whatever `Read` returned.
  - `stop()` returns `true` → the read finished first: return exactly what `Read` returned (`""` with its error on failure).
  - Never call `Close` yourself; only the `AfterFunc` callback may.
- **Merge:** the result is a child of `a` (it has `a`'s values and deadline) that is **also** canceled when `b` is done, with `context.Cause(b)` as its cause.
  - Build it with `context.WithCancelCause(a)` and `context.AfterFunc(b, …)`.
  - The returned `CancelFunc` stops the `AfterFunc` registration and cancels the result with cause `context.Canceled`.
