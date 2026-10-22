# With Timeout

Run a function in a goroutine and return its result, or give up with `ErrTimeout` once `d` has passed.

```go
var ErrTimeout = errors.New("timeout")

func Do[T any](d time.Duration, f func() T) (T, error)
```

## Examples

| Call | Result | Returns after |
|---|---|---|
| `Do(2*time.Second, f)`, `f` takes 1s and returns `42` | `42, nil` | 1s |
| `Do(2*time.Second, f)`, `f` takes 3s | `0, ErrTimeout` | 2s |
| `Do(time.Second, f)`, `f` returns `"hi"` at once | `"hi", nil` | 0s |
| `Do(0, f)` or `Do(-time.Second, f)` | `zero, ErrTimeout`, `f` is never called | 0s |

## Rules
- Call `f` exactly once, in a **new goroutine**, and wait for whichever comes first: its result or `d` elapsing (`select` + `time.After`).
- On timeout return the zero `T` and exactly `ErrTimeout` (the tests use `errors.Is`).
- `d <= 0` returns `ErrTimeout` immediately without calling `f`.
- `Do` doesn't stop `f` on timeout (it can't), but the goroutine must still be able to **finish** later: give the result channel a buffer of 1 so its send never blocks. The tests run in a `synctest` bubble, which fails with `blocked goroutines remain` if that goroutine is stuck.
