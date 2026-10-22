# Batch a Stream

Group a stream of values into slices of `size` (the last one may be shorter), e.g. to insert rows into a database 100 at a time instead of one by one.

```go
func Batch(in <-chan string, size int) <-chan []string
```

## Examples

| `in` receives, then is closed | `size` | `Batch` receives, then is closed |
|---|---|---|
| `a b c d e` | `2` | `[a b] [c d] [e]` |
| `a b c d` | `2` | `[a b] [c d]` |
| `a b` | `5` | `[a b]` |
| `a b c` | `1` | `[a] [b] [c]` |
| nothing | `3` | nothing (no empty batch) |
| anything | `0` | panics: `batch: size must be at least 1` |

## Rules
- If `size < 1`, `Batch` panics with exactly `"batch: size must be at least 1"`, **in the caller's goroutine**, before starting anything.
- Otherwise it **returns immediately**; one goroutine reads `in`.
- Send a batch as soon as it holds `size` values. When `in` is closed, send the remaining values as a final, shorter batch (only if there are any), then close the output.
- Never send an empty batch. Values keep their input order.
- ⚠️ Every batch must be a **new slice**. The receiver keeps the batches it got; if you reuse one backing array (`buf = buf[:0]`), your next appends overwrite batches already sent.
- Never return a `nil` channel.
