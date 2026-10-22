# Stoppable Stream

An endless generator and a `Take` stage that can both be **stopped** through a done channel, so that no goroutine is ever left behind. `Take(done, From(done, 1), 3)` yields `1 2 3`; closing `done` afterwards shuts `From` down.

```go
func From(done <-chan struct{}, start int) <-chan int
func Take(done <-chan struct{}, in <-chan int, n int) <-chan int
```

## Examples

| Code | Values received, then the channel is closed |
|---|---|
| `Take(done, From(done, 1), 3)` | `1 2 3` |
| `Take(done, From(done, -2), 4)` | `-2 -1 0 1` |
| `Take(done, in, 3)` with `in` holding `7 8 9 10` | `7 8 9`; `10` is left unread in `in` |
| `Take(done, in, 10)` with `in` holding `7 8`, then closed | `7 8` |
| `Take(done, in, 0)` | nothing; `in` is not read at all |
| `From(done, 1)`, read `1 2`, then `close(done)` | maybe a value or two more, then closed |

## Rules

**Both functions**
- Return immediately; each starts **one** goroutine and closes its output when that goroutine returns.
- Stop as soon as `done` is closed, **wherever** the goroutine is blocked: waiting to send, or (for `Take`) waiting to receive. Use `select` with a `case <-done:` around **every** blocking channel operation.
- After `close(done)`, a reader that is still reading may get a few more values, but the output must be closed promptly and the goroutine must exit, even if nobody reads at all.

**`From`**
- Sends `start, start+1, start+2, …` forever, until `done` is closed.

**`Take`**
- Forwards the first `n` values from `in`, in order, then closes its output and returns, **without** receiving any further value from `in`.
- If `in` is closed first, close the output early. For `n <= 0`, close at once without reading `in`.
- `Take` never closes `in` (it doesn't own it); `From` keeps running until the caller closes `done`.
