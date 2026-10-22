# Module 5 Challenge: Concurrent Systems

No new concepts here. This page is a one-page recap of Module 5 plus a checklist for designing concurrent code; the ten exercises each combine goroutines, channels, `select`, locks, atomics and `context` into a small component you'd find inside a real service.

## 1. Goroutines and WaitGroup

```go
var wg sync.WaitGroup
results := make([]int, len(inputs))   // one slot per goroutine: no lock needed
for i, in := range inputs {           // Go 1.22+: i and in are new each iteration
	wg.Go(func() {                    // Go 1.25: Add(1) + go + defer Done()
		results[i] = work(in)
	})
}
wg.Wait()                             // every goroutine has finished here
```

- `main` returning kills every goroutine. Something must **wait**: `wg.Wait()`, a channel receive, or a `Close` that waits.
- A panic in any goroutine crashes the program; `recover` only works in a `defer` **inside that goroutine**.
- `runtime.NumGoroutine()` counts live goroutines (useful for spotting leaks); `runtime.GOMAXPROCS(0)` reads how many run in parallel.

## 2. Channels

| Operation | open channel | closed channel | `nil` channel |
|---|---|---|---|
| `ch <- v` | blocks until received (or buffer has room) | **panic** | blocks forever |
| `v := <-ch` | blocks until a value arrives | zero value at once | blocks forever |
| `v, ok := <-ch` | `ok == true` | `ok == false` | blocks forever |
| `close(ch)` | wakes every receiver | **panic** | **panic** |
| `for v := range ch` | loops | ends after the buffered values | blocks forever |

- **The sender closes**, and only when nothing else will ever be sent. Receivers never close.
- Closing is a **broadcast**: every receiver wakes. Use `chan struct{}` for pure signals (`done`, `stop`, `release`).
- Unbuffered = hand-off (sender and receiver meet). Buffered `make(chan T, n)` = queue of `n`, also a **semaphore** or a **pool**.
- Directional types document ownership: `func Gen() <-chan int` (caller only receives), `func Sink(in chan<- int)`.

## 3. select and time

```go
select {
case v := <-in:          // whichever case is ready first;
case out <- x:           // if several are ready, one at random
case <-ctx.Done():
	return ctx.Err()
case <-time.After(d):    // timeout
default:                 // present: never block
}
```

| Tool | Use |
|---|---|
| `time.After(d)` | one-off timeout inside a `select` |
| `t := time.NewTimer(d)`; `t.Stop()`; `t.Reset(d)` | a timeout you cancel or restart (debounce) |
| `tk := time.NewTicker(d)`; `defer tk.Stop()` | periodic work (janitors, snapshots) |
| set a channel variable to `nil` | disable that `select` case |

## 4. sync and atomic

| Tool | Use it for |
|---|---|
| `sync.Mutex` + `defer mu.Unlock()` | several fields that must change together (an invariant) |
| `sync.RWMutex` (`RLock`) | many readers, rare writers |
| `sync.Once` / `OnceValue` / `OnceValues` | lazy init exactly once, safe from many goroutines |
| `atomic.Int64`, `atomic.Bool` | a single counter or flag |
| `atomic.Pointer[T]` | publish an **immutable** value; readers `Load` without locking |
| `sync.Map` | rarely: keys written once, read many times; otherwise map + mutex |

- Never copy a value that contains a mutex (`go vet` copylocks): use pointer receivers.
- Locking two things? Always in the **same global order** (e.g. by ID), or you can deadlock.
- Don't hold a lock while calling unknown code (a callback, `fn()`, a network call) or while blocking on a channel.

## 5. context

| Function | Gives you |
|---|---|
| `context.Background()` / `t.Context()` | the root in `main` / in a test |
| `WithCancel(parent)` | `ctx, cancel`: call `cancel()` to stop everything below |
| `WithTimeout(parent, d)` / `WithDeadline` | cancelled automatically; `ctx.Err()` = `DeadlineExceeded` |
| `WithCancelCause` + `context.Cause(ctx)` | say **why** it was cancelled |
| `WithValue(parent, key{}, v)` | request-scoped data only, unexported key type |
| `AfterFunc(ctx, f)` | run `f` once when `ctx` is done |
| `WithoutCancel(parent)` | keep values, drop cancellation (cleanup after the request ends) |

- `ctx` is the **first parameter**; always `defer cancel()`.
- Cancellation flows **down**: a child's deadline never outlives its parent's. Check `ctx.Err()` before expensive steps and in loops.
- Telling "my own timeout" from "the caller gave up": check the **parent** `ctx.Err()`.

## 6. Patterns

| Pattern | Shape | Who closes what |
|---|---|---|
| Generator | `func Gen(ctx) <-chan T` | the generator closes its output |
| Pipeline | stage `in <-chan A` → `out <-chan B` | each stage closes its own `out` when `in` ends |
| Fan-out / fan-in | n workers read one channel; merge their outputs | merger closes after a `WaitGroup` of the workers |
| Worker pool | `jobs` channel + n workers | the producer closes `jobs`; workers exit on `range` end |
| Semaphore | `sem := make(chan struct{}, n)`; `sem <- struct{}{}` … `<-sem` | nobody (it's never closed) |
| errgroup | first error cancels a shared `ctx` | `Wait` returns after all goroutines end |
| Ordered results | write to `results[i]` | — |
| Pub/sub | per-subscriber channel | the bus closes subscriber channels under its lock |
| Graceful shutdown | `stop` signal → drain → `done` | owner closes `stop`, worker closes `done` |

## 7. Testing concurrent code

| Need | Technique |
|---|---|
| Time (timeouts, tickers, retries) | `synctest.Test(t, func(t *testing.T) {...})`: fake clock, `time.Sleep(time.Hour)` returns at once |
| "Everyone is now blocked" | `synctest.Wait()`, then check state or `select { case v := <-ch: … default: t.Fatal(…) }` |
| Proves it runs concurrently | n tasks of `time.Sleep(time.Second)` take exactly `1s` of fake time |
| Proves a limit | count `running` with `atomic.Int64`, keep the max, assert `max == limit` |
| Proves thread safety | hammer: 50 goroutines × 1000 ops, then assert exact totals and invariants |
| Never hang | bound every wait outside a bubble: `select { case <-done: case <-time.After(2*time.Second): t.Fatal(…) }` |
| Leaks | a synctest bubble can't end with blocked goroutines: it panics |

## Approaching a concurrent problem — checklist

- **Draw it first.** Boxes for goroutines, arrows for channels, a cylinder for each piece of shared state.
- **One owner per channel.** Who sends? Who **closes** (exactly one sender, after the last send)? Who receives, and how do they learn it's over (`range` end, `done`, `ctx`)?
- **Every goroutine has an exit.** For each `go`, name the event that makes it return, and who waits for it (`WaitGroup`, `done` channel, `Close`).
- **Every blocking operation has an escape.** A send or receive that could wait forever needs a `select` with `ctx.Done()` or `stop`.
- **List the shared state and its guard.** Each field is either owned by one goroutine, behind a mutex (held briefly, never across a callback), atomic, or immutable after publishing.
- **Cancellation flows down.** One `ctx` in, derived contexts for each request, `cancel` always deferred, parent checked before retrying.
- **Results don't depend on scheduling.** Index by position, sort at the end, never rely on goroutine order or map order.
- **Test it with fake time.** Write the synctest test first: exact elapsed times, exact max concurrency, `synctest.Wait()` before checking state.

## Common mistakes

| Mistake | Fix |
|---|---|
| Closing a channel from the receiver, or twice | only the single sender closes; use `sync.Once` if several paths may |
| `wg.Add` inside the goroutine | `wg.Go(f)`, or `Add` before `go` |
| Goroutine stuck on a send nobody receives | buffer it, or `select` on `ctx.Done()` too |
| Mutex held while calling `fn()` or sending on a channel | copy what you need, unlock, then call |
| Reading a shared field without the lock "just to check" | lock it, or make it atomic |
| `time.Sleep` in tests to "wait until it's done" | `synctest.Wait()` or a channel |
| `select` assumed to prefer the first case | it picks randomly: check the priority case with a separate non-blocking `select` first |
| Retrying when the caller cancelled | check the parent `ctx.Err() == nil` before every retry |
| Forgetting `defer cancel()` / `defer tk.Stop()` | add it on the next line, every time |
| A background goroutine with no way to stop it | give the type a `Close` that signals `stop` and waits for `done` |

## Key takeaways

- Concurrency is about **ownership**: of each channel (who closes it), each goroutine (who waits for it) and each piece of state (what guards it).
- Channels coordinate and hand work around; mutexes protect invariants; atomics handle single values and immutable snapshots.
- `context` carries cancellation and deadlines down the call tree; every blocking call should be able to give up.
- The classic patterns (pool, semaphore, pipeline, fan-in, errgroup, pub/sub, graceful shutdown) are combinations of the same few tools.
- Deterministic tests are possible: `testing/synctest` gives exact time, plus hammer tests and max-concurrency counters for everything else.
