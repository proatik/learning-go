# Drain and Sum

Consume a channel until its sender closes it: `Sum` adds up every value, `Drain` collects every value of any type into a slice.

```go
func Sum(ch <-chan int) int
func Drain[T any](ch <-chan T) []T
```

## Examples

| Channel receives, then is closed | Call | Result |
|---|---|---|
| `1 2 3` | `Sum(ch)` | `6` |
| `-4 4 10` | `Sum(ch)` | `10` |
| nothing | `Sum(ch)` | `0` |
| `"a" "b" "a"` | `Drain(ch)` | `["a" "b" "a"]` |
| `0 0` | `Drain(ch)` | `[0 0]` (zeros are real values) |
| nothing | `Drain(ch)` | empty (`nil` or `[]T{}`) |

## Rules
- Read **every** value until the channel is closed, then return. Values keep their order in `Drain`.
- The sender may still be producing when you start: block on each receive, don't give up when the channel is momentarily empty (`len(ch) == 0` tells you nothing here).
- A real `0` is not "closed". Stop only when the channel reports it is closed.
- You only receive: the parameter is `<-chan T`, so you can't (and mustn't) close it.
