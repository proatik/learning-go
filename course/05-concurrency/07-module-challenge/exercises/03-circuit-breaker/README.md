# Circuit Breaker

When a downstream service is failing, hammering it with more requests makes things worse. A **circuit breaker** wraps calls to it: after `threshold` consecutive failures it **opens** and rejects calls instantly for a `cooldown`, then lets **one** trial call through to see whether the service has recovered.

```go
type State int

const (
	Closed   State = iota // calls go through
	Open                  // calls are rejected with ErrOpen
	HalfOpen              // one trial call is allowed
)

func (s State) String() string // "closed", "open", "half-open" (given)

var ErrOpen = errors.New("breaker: open")

type Breaker struct {
	mu        sync.Mutex
	threshold int
	cooldown  time.Duration
	state     State
	failures  int       // consecutive failures while Closed
	openedAt  time.Time // when the breaker last opened
	probing   bool      // a HalfOpen trial call is running
}

func New(threshold int, cooldown time.Duration) *Breaker

func (b *Breaker) Call(fn func() error) error
func (b *Breaker) State() State
```

## Examples

| Code (`b := New(3, 10*time.Second)`) | Result |
|---|---|
| `b.Call(ok)` | `nil`; `State()` = `closed` |
| `b.Call(fail)` twice, then `b.Call(ok)` | the success resets the count: still `closed` |
| `b.Call(fail)` three times in a row | each returns `fail`'s error; now `open` |
| `b.Call(ok)` 5s later | `ErrOpen`, and `ok` is **not** called |
| `b.Call(ok)` 10s after opening | `ok` runs (the trial call); returns `nil`; now `closed` |
| a trial call that fails | returns its error; `open` again, cooldown restarts |
| another `Call` while the trial call is running | `ErrOpen` at once |

## Rules
- `New` starts `Closed`. `threshold < 1` counts as `1`.
- **Closed:** run `fn`. An error increments the consecutive-failure count; reaching `threshold` switches to `Open` and records the time. A success resets the count to `0`.
- **Open:** if less than `cooldown` has passed since opening, return `ErrOpen` without calling `fn`. Once `cooldown` has passed (`>=`), switch to `HalfOpen` and make **this** call the trial call.
- **HalfOpen:** only one trial call at a time; every other `Call` returns `ErrOpen` while it runs. Trial success → `Closed` (count `0`). Trial failure → `Open` again, with a fresh `openedAt`.
- A non-trial call that started while `Closed` but finishes after the breaker left `Closed` doesn't change anything.
- `Call` returns `fn`'s error unchanged. Use `time.Now()` for time.
- `State` returns the current state. The `Open` → `HalfOpen` switch happens inside `Call`, so `State` still reports `open` after the cooldown until the next `Call`.
- Safe for concurrent use. **Don't hold the mutex while `fn` runs**: calls in the `Closed` state must run concurrently.
