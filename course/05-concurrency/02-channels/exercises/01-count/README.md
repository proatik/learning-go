# Count Generator

Write two **generators**: functions that start a goroutine, return a receive-only channel, and close it when there's nothing left to send.

```go
func Count(n int) <-chan int
func Runes(s string) <-chan rune
```

## Examples

| Call | Values received, then the channel is closed |
|---|---|
| `Count(3)` | `1 2 3` |
| `Count(1)` | `1` |
| `Count(0)` | nothing |
| `Count(-5)` | nothing |
| `Runes("go")` | `'g' 'o'` |
| `Runes("héllo")` | `'h' 'é' 'l' 'l' 'o'` |
| `Runes("")` | nothing |

## Rules
- Both functions **return immediately**; a goroutine sends the values.
- `Count` sends `1, 2, …, n` in order. For `n <= 0` it sends nothing.
- `Runes` sends each rune (not byte) of `s` in order.
- Both always **close** the channel after the last value, so `for v := range Count(n)` ends.
- Never return a `nil` channel.
