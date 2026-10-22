# Graceful Shutdown

A job server with a bounded queue and a fixed set of workers that can **shut down gracefully**: stop taking new work, finish everything already accepted, and give up waiting when the caller's deadline passes. This is the shape of `http.Server.Shutdown` and of every queue consumer in a microservice.

```go
var (
	ErrClosed    = errors.New("server: closed")
	ErrQueueFull = errors.New("server: queue full")
)

type Server struct {
	mu     sync.Mutex
	jobs   chan int      // the queue; closed by Shutdown
	closed bool          // set by Shutdown; guarded by mu
	done   chan struct{} // closed once every worker has exited
}

func New(workers, queue int, handle func(job int)) *Server
func (s *Server) Submit(job int) error
func (s *Server) Shutdown(ctx context.Context) error
```

```text
Submit ──▶ [ jobs queue ] ──▶ workers ──▶ handle(job)
Shutdown: closed = true ─▶ close(jobs) ─▶ workers drain the queue and exit ─▶ close(done)
          then wait for done OR ctx.Done(), whichever comes first
```

## Examples

| Code | Result |
|---|---|
| `New(2, 10, h)`, 6 jobs of 1s, `Shutdown(ctx with 10s timeout)` | `nil` after 3s; all 6 handled |
| same, `Shutdown(ctx with 1.5s timeout)` | `context.DeadlineExceeded` after 1.5s; 2 handled so far, the rest still finish |
| then `Shutdown(context.Background())` | `nil` at 3s, once the queue is drained |
| `Submit(7)` after `Shutdown` started | `ErrClosed` |
| `New(1, 1, h)`: one job running, one queued, `Submit(3)` | `ErrQueueFull` |
| `New(0, -1, h)` | 1 worker, no queue: `Submit` succeeds only while that worker is idle |
| `Shutdown` of an idle server | `nil` at once |

## Rules
- `New` starts `max(workers, 1)` workers and makes a queue with `max(queue, 0)` slots. Each worker ranges over the queue and calls `handle`.
- `Submit` **never blocks**. It returns `ErrClosed` once `Shutdown` has been called, `ErrQueueFull` if the job can't be handed over right away (a `select` with `default`), and `nil` otherwise.
- `Shutdown` marks the server closed and closes the queue **exactly once**, then waits for `done` or `ctx.Done()`. It returns `nil` once every accepted job has been handled, or `ctx.Err()` if the context ends first (the workers keep draining in the background). Calling it again just waits again.
- Every accepted job is handled exactly once, even when `Submit` and `Shutdown` race. Do the closed check and the send under the same mutex that `Shutdown` holds when it closes the queue: otherwise a `Submit` could send on a closed channel and panic.
