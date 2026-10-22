# Batcher

Group a stream into batches: send a batch when it reaches `size` items **or** when `wait` has passed since its **first** item, whichever comes first. This is how log shippers and database writers cut down on round trips.

```go
func Batch[T any](in <-chan T, size int, wait time.Duration) <-chan []T
```

## Examples

All with `size = 3`, `wait = 1s`.

| Input (value@time) | `in` closed at | Output (batch@time) | Output closed at |
|---|---|---|---|
| `1 2 3 4`@0, `5`@500ms | 3s | `[1 2 3]`@0, `[4 5]`@1s | 3s |
| `1`@0, `2`@200ms | 400ms | `[1 2]`@400ms (flushed on close) | 400ms |
| `1`@0, `2`@100ms, `3`@200ms, `4`@1100ms | 5s | `[1 2 3]`@200ms, `[4]`@2100ms | 5s |
| nothing | 100ms | nothing | 100ms |

With `size = 1`: every item is sent as its own batch the moment it arrives.

## Rules
- Return the output channel at once; one goroutine does the work. The output is unbuffered.
- The `wait` clock for a batch starts when its **first** item arrives. After a batch is sent (full or timed out), the next batch starts empty and its clock starts at its own first item. A stale timer from the previous batch must not cut the next one short (stop it when a full batch is sent).
- When `in` is closed: send the partial batch (if not empty) **immediately**, then close the output.
- Never send an empty batch.
- Every batch is a **new** slice: the receiver may keep it, so don't reuse its backing array for the next batch.
- `size < 1` is treated as `1`.
- The tests never let an item arrive at the same instant a batch's `wait` runs out.
