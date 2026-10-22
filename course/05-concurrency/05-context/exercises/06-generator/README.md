# Leak-free Generators

Lesson 2's generators stopped on a done channel; these stop when a **context** is done. Every blocking send **and** receive must also watch `ctx.Done()`, or the goroutine leaks when the consumer walks away.

```go
func Count(ctx context.Context, start, step int) <-chan int
func Map(ctx context.Context, in <-chan int, f func(int) int) <-chan int
```

## Examples

| Code | Receives |
|---|---|
| `Count(ctx, 1, 1)` | `1, 2, 3, …` forever, until `ctx` is done |
| `Count(ctx, 10, -5)` | `10, 5, 0, -5, …` |
| `Map(ctx, Count(ctx, 1, 1), square)` | `1, 4, 9, 16, …` |
| `Map(ctx, in, double)` with `in` = `3, 4`, then closed | `6, 8`, then the channel is closed |
| `cancel()` while nobody is receiving | both channels are closed and both goroutines exit |
| `cancel()` while `Map` waits on an `in` that never sends | `Map`'s channel is closed and its goroutine exits |

## Rules
- Both return an **unbuffered** channel, fed by one goroutine they start. The goroutine closes its channel when it exits.
- `Count` sends `start`, `start+step`, `start+2*step`, … until `ctx` is done.
- `Map` sends `f(v)` for each `v` from `in`, in order. It stops when `in` is closed **or** `ctx` is done.
- After `ctx` is done the goroutines must exit and close their channels even if nobody ever receives again: the tests fail (or crash) on a leaked goroutine.
- Values already received before the cancel are fine; after it, at most the one value a goroutine was already offering may still get through.
