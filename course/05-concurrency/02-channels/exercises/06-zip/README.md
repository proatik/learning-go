# Zip Two Channels

Pair up values from two channels in order, one from each per round, and stop as soon as **either** channel is closed. Use the comma-ok receive to tell a closed channel from a zero value.

```go
type Pair[X, Y any] struct {
	First  X
	Second Y
}

func Zip[X, Y any](a <-chan X, b <-chan Y) <-chan Pair[X, Y]
```

## Examples

| `a` receives | `b` receives | `Zip(a, b)` receives, then is closed |
|---|---|---|
| `1 2 3` | `"x" "y" "z"` | `{1 x} {2 y} {3 z}` |
| `1 2 3` | `"x"` | `{1 x}` |
| `1` | `"x" "y" "z"` | `{1 x}` |
| `0 0` | `"" ""` | `{0 } {0 }` (zero values are real values) |
| nothing | `"x"` | nothing |

## Rules
- `Zip` **returns immediately**; one goroutine does the work.
- Each round: receive from `a` **first**, then from `b`, then send the pair.
- If `a` is closed, stop **without** receiving from `b`. If `b` is closed, stop and drop the value already taken from `a`.
- Close the output when you stop. Don't read further values from either input.
- Both inputs are closed by their senders in the tests; values left in them stay there (the tests check exactly how many you read).
- Never return a `nil` channel.
