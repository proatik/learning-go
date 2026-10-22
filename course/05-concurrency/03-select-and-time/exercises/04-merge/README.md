# Merge Two Channels

Merge two input channels into one output channel with a **single** goroutine and one `select`, using the **nil-channel trick** to switch off an input once it's closed.

```go
func Merge[T any](a, b <-chan T) <-chan T
```

## Examples

| Inputs | Output (then closed) |
|---|---|
| `a`: `1 2 3`, `b`: `10 20` | the 5 values, `a`'s in order `1 2 3`, `b`'s in order `10 20`, interleaved in any way |
| `a` closed at once, `b`: `5 6` | `5 6` (no zero values from the closed `a`) |
| `a` stays open, `b`: `1 2 3` then closed | `1 2 3` arrive **while** `a` is still open |
| `Merge(nil, b)`, `b`: `7` | `7` |
| `Merge[int](nil, nil)` | closed right away, no values |

## Rules
- Return the output channel at once; one goroutine does all the forwarding. The output is unbuffered.
- Forward every value from both inputs, keeping each input's own order. Don't wait for one input to finish before reading the other.
- A closed input must stop being selected: receiving from a closed channel succeeds immediately with the zero value, forever. Set that input variable to `nil` so its case can never fire again.
- A `nil` input counts as already finished.
- Close the output once **both** inputs are finished, and then let the goroutine exit (the tests check for a leaked goroutine).
