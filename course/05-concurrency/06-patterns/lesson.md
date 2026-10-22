# Concurrency Patterns

## 1. Why patterns

Lessons 1 to 5 gave you the building blocks: goroutines, channels, `select`, mutexes, `context`. Real services combine them in a few shapes that come up again and again. Learn the shapes and most concurrent code turns into "which pattern is this?".

| Pattern | Shape | Typical use |
|---|---|---|
| **Pipeline** | stage → stage → stage | stream processing, ETL, parsing then transforming |
| **Fan-out / fan-in** | one channel → many goroutines → one channel | spreading CPU or I/O work |
| **Worker pool** | fixed N goroutines read one job channel | queue consumers, bounded background work |
| **Semaphore** | buffered channel as a slot counter | "at most N requests to that API at once" |
| **Ordered results** | concurrent work, output in input order | batch calls where the caller needs `result[i]` for `input[i]` |
| **Error group** | N tasks, the first error cancels the rest | fan out N RPCs for one request |
| **Tee / pub-sub** | one value → every reader | events, cache invalidation, live updates |
| **Graceful shutdown** | stop intake → drain → exit, with a deadline | every long-running server |

Every pattern must answer the same four questions:

1. **Who closes each channel?**
2. **How does every goroutine end**, including when the consumer gives up?
3. **How much runs at once?**
4. **Where do errors go?**

## 2. Channel ownership

The goroutine that **creates** a channel **writes** to it and **closes** it. Everyone else only reads.

```go
// gen owns out: it makes it, is its only sender, and closes it.
func gen(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // runs however the goroutine ends
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out // receive-only: callers can't send or close
}
```

| Rule | Why |
|---|---|
| Return `<-chan T`, not `chan T` | the compiler stops readers from sending or closing |
| `defer close(out)` first thing in the goroutine | the channel is closed on every exit path |
| Never close a channel you only receive from | the owner may still send → `panic: send on closed channel` |
| Many senders? Close after **all** are done (`WaitGroup`, then `close`) | one sender closing early panics the others |
| Inputs belong to the caller | a stage closes its outputs, never its inputs |

## 3. Pipelines

A **pipeline** is a chain of stages. Each stage receives from upstream, does one thing, and sends downstream on a channel it owns.

```text
gen ──(1,2,3)──▶ square ──(1,4,9)──▶ consumer
 owns ch1         reads ch1, owns ch2   reads ch2
```

```go
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in { // ends when upstream closes in
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

ctx := context.Background()
for v := range square(ctx, gen(ctx, 1, 2, 3)) {
	fmt.Print(v, " ")
}
// 1 4 9
```

- **Closing flows downstream.** `gen` closes `ch1`, so `square`'s `range` ends, so it closes `ch2`, so the consumer's `range` ends. No extra signals needed.
- **Cancellation flows to every stage.** Pass the same `ctx` to all of them and put `case <-ctx.Done()` next to **every send** (and every receive that isn't a `range`).
- Stages **stream**: each value moves on as soon as it is ready. A stage that collects its whole input before emitting kills the benefit.
- Unbuffered channels between stages are the default. A small buffer smooths out uneven stage speeds but doesn't fix a slow stage.

## 4. Goroutine leaks

A goroutine that can never finish is a **leak**: its stack and everything it references stay in memory for good. In a long-running service, leaks pile up with every request.

The classic leak: the consumer stops early, and the producer is stuck on a send nobody will receive.

```go
func leakyGen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n // blocks forever once the reader is gone
		}
	}()
	return out
}

before := runtime.NumGoroutine()
first := <-leakyGen(1, 2, 3) // read one value, walk away
time.Sleep(10 * time.Millisecond)
fmt.Println(first, runtime.NumGoroutine()-before) // 1 1   ← one goroutine stuck for good

before = runtime.NumGoroutine()
ctx, cancel := context.WithCancel(context.Background())
first = <-gen(ctx, 1, 2, 3) // the ctx-aware version from §2
cancel()                    // "I'm done": the producer's select returns
time.Sleep(10 * time.Millisecond)
fmt.Println(first, runtime.NumGoroutine()-before) // 1 0
```

| Leak cause | Fix |
|---|---|
| send with no receiver left | `select { case out <- v: case <-ctx.Done(): return }` |
| receive from a channel nobody closes | the owner closes it; or `select` on `ctx.Done()` too |
| worker `range`s over a jobs channel that is never closed | whoever owns `jobs` closes it once no more jobs will come |
| result channel nobody reads after a timeout | buffer it (`make(chan R, 1)`), so the late sender doesn't block |

- The consumer that stops early must **say so**: call `cancel()` (usually `defer cancel()`).
- In tests, `synctest.Test` catches leaks for you: the bubble fails if goroutines are still blocked when the test function returns.

## 5. Fan-out and fan-in

**Fan-out:** several goroutines receive from the **same** channel; each value goes to exactly one of them. **Fan-in (merge):** several channels are combined into one.

```text
            ┌─▶ worker ─┐
in ─────────┼─▶ worker ─┼──▶ out       fan-out on the left, fan-in on the right
            └─▶ worker ─┘
```

```go
// merge copies every input to one output, closed after the last input ends.
func merge[T any](inputs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, in := range inputs {
		wg.Go(func() {
			for v := range in {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()  // every copier is done…
		close(out) // …so now it's safe to close
	}()
	return out
}

var got []int
for v := range merge(gen(ctx, 1, 2), gen(ctx, 10, 20), gen(ctx, 100)) {
	got = append(got, v)
}
slices.Sort(got) // arrival order varies from run to run
fmt.Println(got) // [1 2 10 20 100]
```

- One goroutine **per input**: reading the inputs one after another would stall on a slow input.
- Order **within** one input is kept. Order **across** inputs is not.
- The closer is a separate goroutine, so `merge` returns at once.
- With `ctx`, the copier's send becomes a `select` with `ctx.Done()`, just like a pipeline stage.

## 6. Worker pools

A **worker pool** is fan-out with a **fixed** number of long-lived goroutines. It caps concurrency and reuses goroutines across many jobs.

```go
jobs := make(chan string)
results := make(chan string)

var wg sync.WaitGroup
for range 3 { // 3 workers, however many jobs there are
	wg.Go(func() {
		for j := range jobs { // exits when jobs is closed and empty
			results <- strings.ToUpper(j)
		}
	})
}
go func() { // the producer owns jobs
	defer close(jobs)
	for _, j := range []string{"a", "b", "c", "d"} {
		jobs <- j
	}
}()
go func() { // results closes once every worker has exited
	wg.Wait()
	close(results)
}()

var rs []string
for r := range results {
	rs = append(rs, r)
}
slices.Sort(rs)
fmt.Println(rs) // [A B C D]
```

| Who | Owns / closes |
|---|---|
| producer | `jobs`: closes it after the last job |
| workers | nothing: they only read `jobs` and send to `results` |
| closer goroutine | `results`: closes it after `wg.Wait()` |

- Pool size: about `runtime.GOMAXPROCS(0)` for CPU-bound work; for I/O-bound work, whatever the downstream system can handle.
- Results arrive in **completion** order. Carry an index or ID in the result if the caller needs to match them up.

## 7. Bounding parallelism with a semaphore

When you'd rather start a goroutine per item but must cap how many run at once, use a **buffered channel as a semaphore**. Its capacity is the number of slots.

```text
sem := make(chan struct{}, 2)       [ _ _ ]   2 free slots
sem <- struct{}{}   acquire          [ x _ ]   blocks when both slots are taken
<-sem               release          [ _ _ ]
```

```go
urls := []string{"a", "bb", "ccc", "dddd", "eeeee"}
lens := make([]int, len(urls)) // pre-sized: each goroutine writes its own index
sem := make(chan struct{}, 2)  // at most 2 at once

var wg sync.WaitGroup
for i, u := range urls {
	sem <- struct{}{} // acquire BEFORE starting the goroutine
	wg.Go(func() {
		defer func() { <-sem }() // release when done
		lens[i] = len(u)
	})
}
wg.Wait()
fmt.Println(lens) // [1 2 3 4 5]
```

| Worker pool | Semaphore |
|---|---|
| N goroutines, many jobs each | one goroutine per item, N at a time |
| needs a jobs channel and a close | just a buffered channel |
| best for an endless stream (queue consumer) | best for a known batch (a slice) |

- **Acquire before `go`.** Acquiring inside the goroutine still caps the running work, but 10 000 items would start 10 000 goroutines that all sit waiting.
- `struct{}` takes no memory: the channel carries only "a slot is taken".
- With a context: `select { case sem <- struct{}{}: case <-ctx.Done(): stop launching }`.

## 8. Keeping results in order

Concurrent work finishes in any order. Two ways to give the caller results in input order:

| Input | Technique |
|---|---|
| a **slice** (size known) | write `out[i]` into a pre-sized slice (§7). Distinct indexes need no mutex |
| a **stream** (channel) | a queue of per-item **promise** channels, read in order |

The streaming version: the dispatcher gives each item a one-slot channel, the "promise" of its result, and queues the promises in input order. The emitter waits for each promise in turn.

```text
in: a b c ─▶ dispatcher ── go: pa <- f(a) ... ─┐
                 │                              ▼
                 └──▶ queue [pa pb pc] ──▶ emitter: out <- <-pa; out <- <-pb; ...
```

- Each promise has a buffer of **1**, so a worker can always finish its send, even after the emitter has gone because of a cancel.
- The queue's capacity bounds how far ahead of a slow reader the dispatcher can run.
- Head-of-line blocking is the price: one slow item holds back every later result, even finished ones.

## 9. Error propagation: the first error cancels the rest

Fanning out N calls for one request, you usually want: run them all at once, and if **one** fails, cancel the others and return that error. All the pieces are from earlier lessons:

```go
ctx, cancel := context.WithCancelCause(context.Background())
var (
	once     sync.Once
	firstErr error
	wg       sync.WaitGroup
)
fail := func(err error) {
	once.Do(func() { // keep only the first error
		firstErr = err
		cancel(err) // every task watching ctx stops
	})
}

wg.Go(func() {
	select {
	case <-time.After(time.Second):
		fmt.Println("slow task finished")
	case <-ctx.Done():
		fmt.Println("slow task stopped:", context.Cause(ctx))
	}
})
wg.Go(func() { fail(errors.New("db down")) })

wg.Wait()
cancel(nil) // always release the context
fmt.Println("first error:", firstErr)
// slow task stopped: db down
// first error: db down
```

| Strategy | Use when |
|---|---|
| **first error wins, cancel the rest** | the result is useless if any part fails (one page built from 3 services) |
| **collect every error** (`errors.Join`, lesson 1's error slice) | the parts are independent (a batch import: report every failed row) |
| **per-item result + error** (`struct{ V T; Err error }`) | the caller decides per item (a crawler that skips broken pages) |

- Tasks must **watch `ctx`**, otherwise cancelling does nothing.
- Still `wg.Wait()` for every task, even after a failure, so none outlives the call.
- The widely used `golang.org/x/sync/errgroup` package packages exactly this. It isn't in the standard library, so exercise 7 has you write your own.

## 10. Tee and pub/sub: every reader gets every value

Fan-out gives each value to **one** reader. Sometimes every reader needs **every** value.

```text
tee:                         pub/sub (broker):
        ┌──▶ out1                       ┌──▶ [buf] subscriber 1
in ─────┤                    Publish ───┼──▶ [buf] subscriber 2
        └──▶ out2                       └──▶ [full] subscriber 3 → policy?
```

- **Tee** (fixed two readers): deliver each value to both outputs before reading the next. Use a `select` over both sends and set a channel variable to `nil` once it has its copy. A nil channel's case never fires (lesson 3), so this sends to whichever reader is ready first.
- **Pub/sub** (readers come and go): a broker keeps a map of subscriber channels under a mutex. The key design question is the **slow-subscriber policy**:

| Policy | Code | Trade-off |
|---|---|---|
| block | `ch <- v` | one stuck subscriber stalls everyone |
| drop | `select { case ch <- v: default: }` | publisher never blocks; slow readers miss values |
| disconnect | drop, then unsubscribe and close it | readers must resubscribe and catch up |

- The broker **owns** every subscriber channel, so the broker closes them (on unsubscribe or `Close`), exactly once.
- Send while holding the lock (fine, because the drop policy's sends never block). Otherwise an unsubscribe can close a channel in the middle of a send.

## 11. Graceful shutdown

A server must stop without losing accepted work: **stop intake → drain → exit**, bounded by a deadline.

```text
Shutdown(ctx):
  1. mark closed        new Submit calls get ErrClosed
  2. close(jobs)        workers finish the queue, their range ends
  3. wait for workers   ┐ whichever comes first:
     or ctx.Done()      ┘ nil, or ctx.Err() (workers keep draining)
```

```go
func (s *Server) Submit(job int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed { // checked under the same lock Shutdown holds to close
		return ErrClosed
	}
	select {
	case s.jobs <- job:
		return nil
	default: // never block while holding the lock
		return ErrQueueFull
	}
}
```

- The "check closed, then send" pair must be **atomic** with respect to `close(jobs)`. Otherwise a `Submit` can slip in between and hit `panic: send on closed channel`.
- Turn "all workers exited" into a channel (`done`, closed after `wg.Wait()`) so it can sit in a `select` with `ctx.Done()`.
- This is how `http.Server.Shutdown(ctx)` behaves: stop listening, wait for active requests, return `ctx.Err()` if the deadline passes first.

## 12. Choosing a pattern

| You need… | Reach for |
|---|---|
| a sequence of transformations on a stream | pipeline, ctx in every stage |
| N goroutines feeding one consumer | fan-in: `WaitGroup`, then `close` |
| to process a never-ending queue with N workers | worker pool |
| to process a slice with at most N at once | semaphore, results by index |
| results of a stream in input order | promise queue |
| all-or-nothing concurrent calls | first error cancels (errgroup) |
| every reader to see every value | tee (fixed readers) / pub-sub (dynamic readers) |
| to stop a long-running component cleanly | graceful shutdown with a deadline |

## Common mistakes

| Mistake | Fix |
|---|---|
| A consumer closes the channel it reads from | Only the owner (the sender) closes. Signal "stop" with `ctx` instead |
| One of several senders closes the shared output | `WaitGroup` over all senders, then one `close` |
| Bare `out <- v` in a stage | `select` with `case <-ctx.Done(): return`: otherwise an early-exiting consumer leaks the stage |
| Forgetting `cancel()` when the consumer stops early | `ctx, cancel := context.WithCancel(…)`; `defer cancel()` |
| Workers `range` over a jobs channel that is never closed | The producer closes `jobs` after its last job |
| Merging by reading inputs one after another | One goroutine per input |
| `go func() { sem <- struct{}{}; … }()` for 10 000 items | Acquire **before** `go`, so at most N goroutines exist |
| Results by append from many goroutines | Write to `out[i]` in a pre-sized slice (or use a mutex) |
| Tee sends `out1 <- v` then `out2 <- v` | `select` over both sends, nil-ing each one once it's done |
| Blocking send to a subscriber while holding the broker's lock | Non-blocking `select` with `default` (drop) |
| `close(jobs)` racing with `Submit`'s send | Check `closed` and send under the same mutex `Shutdown` uses |
| Waiting forever in `Shutdown` | `select` on a `done` channel and `ctx.Done()` |
| Tasks that ignore `ctx` in an error group | Pass `ctx` down and check it, or cancelling is useless |

## Key takeaways

- The creator of a channel sends on it and closes it. Return `<-chan T`. Several senders: `WaitGroup`, then one `close`.
- Closing flows **downstream** (`range` ends, the stage closes its own output); cancellation flows **to every stage** through one shared `ctx`.
- Every blocking send or receive needs an exit: `select` with `ctx.Done()`, or a buffer of 1 for a result nobody may read. Otherwise you leak a goroutine.
- Bound parallelism with a worker pool (stream of jobs) or a buffered-channel semaphore (batch). Acquire before `go`.
- Keep order with indexes into a pre-sized slice, or a queue of promise channels for streams.
- Errors: first error cancels the rest (`WithCancelCause` + `sync.Once`), collect them all (`errors.Join`), or return a per-item result. Pick on purpose.
- Broadcasting (tee, pub/sub) needs a slow-reader policy; graceful shutdown is stop intake → drain → wait, bounded by a deadline.
