# Fan-In: Merge Channels

Combine any number of input channels into **one** output channel (fan-in).

```go
func Merge[T any](inputs ...<-chan T) <-chan T
```

```text
in1 ──┐
in2 ──┼──▶ Merge ──▶ out   (closed when in1, in2 and in3 are all closed)
in3 ──┘
```

## Examples

| Code | Output yields |
|---|---|
| `Merge(ch(1, 2, 3), ch(10, 20), ch(100))` | `1 2 3 10 20 100` in **some** interleaving, then closed |
| `Merge(ch(5))` | `5`, then closed |
| `Merge[int]()` | nothing: closed straight away |
| `Merge(nil, ch("a", "b"), nil)` | `"a" "b"`, then closed |

(`ch(…)` is a closed channel holding those values.)

## Rules
- Read **every input at the same time**, one goroutine per input. Reading one input to the end before starting the next is wrong: an input may stay open for a long time.
- Values from the same input keep their order. The interleaving between inputs is unspecified.
- `nil` inputs are skipped (receiving from a nil channel blocks forever).
- Close the output exactly once, **after** every input goroutine has finished. Merge owns the output, so it closes it; it never closes an input.
- Merge returns at once.
