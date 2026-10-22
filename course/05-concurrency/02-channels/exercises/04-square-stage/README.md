# Square Stage

Build a two-step **pipeline**: `Generate` turns its arguments into a stream, and `Square` is a **stage** that reads one stream and writes another. Stages chain: `Square(Square(Generate(2)))` yields `16`.

```go
func Generate(nums ...int) <-chan int
func Square(in <-chan int) <-chan int
```

## Examples

| Pipeline | Values received, then the channel is closed |
|---|---|
| `Generate(1, 2, 3)` | `1 2 3` |
| `Generate()` | nothing |
| `Square(Generate(1, 2, 3))` | `1 4 9` |
| `Square(Generate(-3, 0))` | `9 0` |
| `Square(Square(Generate(1, 2, 3)))` | `1 16 81` |

## Rules
- Both functions **return immediately**; each starts one goroutine that does the work.
- `Generate` sends `nums` in order, then closes its channel.
- `Square` sends `v*v` for every `v` received from `in`, in order, **as each value arrives**. It doesn't wait for `in` to close before sending the first result.
- When `in` is closed, `Square` closes its output. It never closes `in` (it can't: `in` is receive-only).
- Never return a `nil` channel.
