# Bounded Parallel Map

Apply `f` to every element concurrently, but never run more than `limit` calls at once. A **buffered channel** is the semaphore: its capacity is the number of free slots.

```go
func Map[T, R any](limit int, in []T, f func(T) R) []R
```

```text
sem := make(chan struct{}, limit)
sem <- struct{}{}   acquire a slot (blocks while limit calls are running)
<-sem               release it
```

## Examples

| Code | Result |
|---|---|
| `Map(2, []int{3, -1, 0}, strconv.Itoa)` | `["3" "-1" "0"]` |
| `Map(3, 10 items that sleep 1s, …)` | at most 3 run at once; returns after 4s |
| `Map(50, 10 items that sleep 1s, …)` | all 10 run at once; returns after 1s |
| `Map(0, 3 items that sleep 1s, …)` | limit 1: one at a time, 3s |
| `Map(3, []int{}, strconv.Itoa)` | an empty slice |

## Rules
- At most `max(limit, 1)` calls of `f` run at the same time, and up to that many **do** run at once.
- **Acquire before `go`:** take a slot, then start the goroutine, which releases the slot when `f` returns. This also caps the number of goroutines, not just the number of running calls (100 items with limit 3 never have 100 goroutines waiting).
- `result[i]` is `f(in[i])`, whatever order the calls finish in. Each goroutine writes its own index of a pre-sized slice, so no mutex is needed.
- `Map` returns only after every call has finished.
