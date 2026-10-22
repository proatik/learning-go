# WaitGroup.Go

Call `f(0)`, `f(1)`, …, `f(n-1)`, each in its own goroutine started with Go 1.25's `wg.Go`, wait for all of them, and return how many you started.

```go
func Times(n int, f func(i int)) int
```

## Examples

| Call | Result |
|---|---|
| `Times(4, f)` where `f` sleeps 1s | `4`, after **1s**; `f` got `0`, `1`, `2`, `3` once each |
| `Times(1, f)` | `1`; `f(0)` ran |
| `Times(0, f)` | `0`; `f` is never called |
| `Times(-3, f)` | `0`; `f` is never called |

## Rules
- Use `var wg sync.WaitGroup` and `wg.Go(func() { … })`; no `Add` / `Done`.
- Every `i` from `0` to `n-1` is passed to `f` **exactly once**, each call in its own goroutine.
- `Times` returns only after every call to `f` has returned.
- The result is the number of goroutines started: `n`, or `0` when `n <= 0`.
