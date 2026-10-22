# Split Evens and Odds

Route each value from one channel to one of **two** output channels, and close **both** when the input ends.

```go
func Split(in <-chan int) (evens, odds <-chan int)
```

## Examples

| `in` receives, then is closed | `evens` receives | `odds` receives |
|---|---|---|
| `1 2 3 4 5` | `2 4` | `1 3 5` |
| `-3 -2 0 7` | `-2 0` | `-3 7` |
| `2 4` | `2 4` | nothing |
| nothing | nothing | nothing |

## Rules
- `Split` **returns immediately**; one goroutine reads `in` and sends each value to `evens` or `odds`.
- Each output keeps the input order.
- `0` and negative even numbers are even. Careful: in Go, `-3 % 2` is `-1`, not `1`.
- When `in` is closed, close **both** outputs, even one that never received a value.
- The outputs are unbuffered, so a send waits for its reader. Callers must read **both** outputs at the same time (the tests use two goroutines); a caller that only reads `evens` would block `Split` on the first odd value.
- Never return a `nil` channel.
