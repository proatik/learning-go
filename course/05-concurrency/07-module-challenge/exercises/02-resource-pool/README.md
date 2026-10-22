# Resource Pool

Share a fixed set of expensive resources (database connections, browsers, buffers) between many goroutines. `Get` borrows one, **waiting** if all are in use, until one is returned with `Put`, the context ends, or the pool is closed.

```go
var ErrClosed = errors.New("pool: closed")

type Pool[T any] struct {
	idle      chan T        // resources not borrowed right now
	closed    chan struct{} // closed by Close
	closeOnce sync.Once
}

func New[T any](resources []T) *Pool[T]

func (p *Pool[T]) Get(ctx context.Context) (T, error)
func (p *Pool[T]) Put(r T)
func (p *Pool[T]) Idle() int
func (p *Pool[T]) Close()
```

## Examples

| Code (`p := New([]string{"a", "b"})`) | Result |
|---|---|
| `p.Idle()` | `2` |
| `p.Get(ctx)`, `p.Get(ctx)` | `"a", nil` and `"b", nil` (in some order); `Idle()` = `0` |
| third `p.Get(ctx)` | blocks until a goroutine calls `p.Put("a")`, then returns `"a", nil` |
| `p.Get(ctx)` with a 1s timeout, nothing returned | after 1s: `"", context.DeadlineExceeded` |
| `p.Close()`, then `p.Get(ctx)` | `"", ErrClosed` at once |
| `p.Close()` while a `Get` is waiting | that `Get` returns `"", ErrClosed` |
| `p.Put(x)` after `Close` | ignored |

## Rules
- `New` makes a pool holding all of `resources`, idle; the pool's capacity is `len(resources)`.
- `Get` returns an idle resource at once if there is one. Otherwise it blocks until a resource is `Put` back, `ctx` is done (return the zero value and `ctx.Err()`), or the pool is closed (zero value and `ErrClosed`).
- After `Close`, `Get` always returns `ErrClosed`, even if resources are idle.
- `Put` never blocks. It returns `r` to the pool; after `Close` it does nothing. A `Put` into a pool that is already full is a bug in the caller: panic with exactly `"pool: Put into a full pool"`.
- `Idle` returns the number of idle resources. `Close` can be called many times.
- Safe for concurrent use.
