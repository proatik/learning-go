# Try Send, Try Receive

Send and receive **without ever blocking**, using `select` with a `default` case.

```go
func TrySend[T any](ch chan<- T, v T) bool
func TryRecv[T any](ch <-chan T) (v T, ok bool)
func DrainReady[T any](ch <-chan T) []T
```

## Examples

| Code | Result |
|---|---|
| `ch := make(chan int, 1); TrySend(ch, 5)` | `true` (buffer had room) |
| `TrySend(ch, 6)` (buffer now full) | `false`, returns at once |
| `TrySend(make(chan int), 1)` (no receiver waiting) | `false` |
| `TrySend(unbuffered, 1)` while a goroutine is blocked in `<-unbuffered` | `true` |
| `TryRecv(ch)` (holds `5`) | `5, true` |
| `TryRecv(ch)` (now empty) | `0, false` |
| `TryRecv` on a closed or `nil` channel | `0, false` |
| `DrainReady(ch)` with `1 2 3` buffered | `[1 2 3]`; the channel is now empty |
| `DrainReady(ch)` with nothing buffered | `nil` |

## Rules
- None of the three functions may block, not even for a moment's wait on another goroutine.
- `TrySend` returns `true` only if the value was handed over right now (buffer space, or a receiver already waiting). A `nil` channel returns `false`. The tests never pass a closed channel (sending on one panics).
- `TryRecv` returns `ok == true` only for a real value. Empty, closed and `nil` channels all give the zero value and `false`.
- `DrainReady` receives every value that is ready **now**, in order, and stops as soon as none is ready or the channel is closed. It returns `nil` if nothing was received.
