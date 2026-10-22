# Stream Results

A worker goroutine streams `Result` values (a value **or** an error) over a channel, and a consumer collects them, stopping at the first error. Size the buffer so the worker can **always finish**, even when the consumer stops reading early.

```go
type Result struct {
	Value int
	Err   error
}

func Run(inputs []int, f func(int) (int, error)) <-chan Result
func Collect(ch <-chan Result) ([]int, error)
```

## Examples

`double` returns `2*n`; `failOn3` returns `n` but fails with `errBad` for `3`.

| Call | Result |
|---|---|
| `Run([1 2 3], double)` | `{2 <nil>} {4 <nil>} {6 <nil>}`, then closed |
| `Run([1 3 5], failOn3)` | `{1 <nil>} {0 bad}`, then closed (`f` is never called for `5`) |
| `Run([], double)` | closed at once |
| `Collect(Run([1 2 3], double))` | `[2 4 6], nil` |
| `Collect(Run([1 3 5], failOn3))` | `[1], errBad` |
| `Collect(nil)` | `nil, nil` at once |

## Rules

**`Run`**
- Returns immediately. **One** goroutine calls `f` for each input **in order** and sends `Result{Value: v, Err: err}`.
- After sending the first result whose `Err != nil`, it stops: `f` is not called for the remaining inputs.
- The returned channel has capacity **exactly `len(inputs)`**, so the goroutine never blocks on a send and always exits, whether or not anyone reads. Then it closes the channel.

**`Collect`**
- Receives until the channel is closed, returning every `Value` in order and `nil`.
- At the first result with `Err != nil`, it stops reading **right away** and returns the values received **before** it, plus that error unchanged (callers use `errors.Is`).
- `Collect(nil)` returns `nil, nil` immediately: receiving from a nil channel would block forever.
