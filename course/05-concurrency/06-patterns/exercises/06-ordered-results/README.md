# Ordered Concurrent Results

Process a **stream** concurrently, yet emit the results in input order, each one as soon as it (and every earlier one) is ready.

```go
func Process[T, R any](ctx context.Context, workers int, in <-chan T, f func(T) R) <-chan R
```

```text
in: a b c d ─▶ dispatcher ──▶ queue of promises [pa pb pc] ──▶ emitter ──▶ out: f(a) f(b) f(c) f(d)
                    │                              ▲
                    └─ go: pa <- f(a) ─────────────┘   (one goroutine per item, at most `workers` at once)
```

## Examples

| Code | Result |
|---|---|
| `Process(ctx, 2, ch("go", "héllo"), strings.ToUpper)` | yields `"GO"`, `"HÉLLO"`, then closed |
| 6 items sleeping 6s, 5s, …, 1s with 6 workers | yields in **input** order, all at 6s |
| 5 items of 1s with 2 workers | results arrive at 1s, 1s, 2s, 2s, 3s |
| 7 items of 1s with 3 workers | at most 3 calls at once; closed after 3s |
| `cancel()` while `in` is open or the reader has stopped | the output closes, every goroutine exits |

(`ch(…)` is a closed channel holding those values.)

## Rules
- Up to `max(workers, 1)` calls of `f` run at the same time.
- The `i`-th result sent is `f` of the `i`-th value received. Send each result as soon as it **and all earlier ones** are done: never wait for the whole input.
- **Bound how far you run ahead:** if the reader stops reading, no more than `3 × workers` calls of `f` may be started. (A one-slot promise channel per item, kept in a queue channel with capacity `workers`, gives you this for free.)
- Close the output after the last result, once `in` is closed.
- When `ctx` is done, close the output promptly and let every goroutine end. Every receive and send must also `select` on `ctx.Done()`. A worker whose result nobody will read must not block forever: give each promise channel a buffer of 1.
- Never close `in`.
