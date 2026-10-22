# Channels

## 1. What a channel is

Lesson 1's goroutines shared results through a pre-sized slice and a `WaitGroup`. A **channel** is the other tool: a typed pipe that one goroutine **sends** values into and another **receives** them from.

```go
ch := make(chan string) // a channel of string values

go func() {
	ch <- "hello" // send: put a value into the channel
}()

msg := <-ch      // receive: take a value out
fmt.Println(msg) // hello
fmt.Printf("%T\n", ch) // chan string
```

| Syntax | Meaning |
|---|---|
| `chan T` | the type: a channel carrying `T` values |
| `make(chan T)` | an **unbuffered** channel |
| `make(chan T, n)` | a **buffered** channel holding up to `n` values |
| `ch <- v` | send `v` (a statement) |
| `v := <-ch` | receive a value (an expression; `<-ch` alone discards it) |
| `v, ok := <-ch` | receive, with `ok == false` once the channel is closed and empty (§4) |
| `close(ch)` | no more values will be sent |
| `len(ch)`, `cap(ch)` | values waiting in the buffer, buffer size |

- A channel is a **reference** value, like a map: copying it or passing it to a function shares the same channel.
- The zero value is `nil`: `var ch chan int` is a nil channel. Always create channels with `make` (§7 explains what `nil` does).
- `chan struct{}` carries no data (`struct{}` takes zero bytes). Use it when only the **event** matters: "done", "start", "stop".

## 2. Unbuffered channels synchronise

On an unbuffered channel, a send **waits until a receiver takes the value**, and a receive waits until a sender offers one. The two goroutines meet at the channel, so it's a hand-off and a synchronisation point at once.

```go
var result int
done := make(chan struct{})

go func() {
	result = 42
	done <- struct{}{} // blocks until main is ready to receive
}()

<-done              // blocks until the goroutine sends
fmt.Println(result) // 42
```

- Everything the goroutine did **before** the send is visible to main **after** the receive. That's why reading `result` here is safe without any lock.
- This replaces a `WaitGroup` when you have exactly one goroutine to wait for, and it can carry the result too: `res := make(chan int); go func() { res <- compute() }(); v := <-res`.

## 3. Buffered channels

A buffered channel is a fixed-size FIFO queue. A send only blocks when the buffer is **full**; a receive only blocks when it's **empty**.

```go
q := make(chan int, 3)
q <- 1 // doesn't block: there's room
q <- 2
fmt.Println(len(q), cap(q)) // 2 3
fmt.Println(<-q, <-q)       // 1 2 (first in, first out)
fmt.Println(len(q))         // 0
```

| | Unbuffered `make(chan T)` | Buffered `make(chan T, n)` |
|---|---|---|
| Send blocks until | a receiver takes the value | there's a free slot |
| Receive blocks until | a sender offers a value | a value is in the buffer |
| Guarantees | the receiver **has** the value when the send returns | only that the value is **queued** |
| Typical use | hand-off, signalling, step-by-step coordination | absorbing bursts, decoupling a fast producer from a slow consumer, letting a sender finish without a reader |

- `make(chan T, 0)` is the same as `make(chan T)`.
- Don't add a buffer "to fix a deadlock". If the program only works with a buffer of exactly `n`, the design depends on `n` and breaks when the input grows. Size a buffer deliberately (e.g. one slot per result, §10).
- `len(ch)` is a snapshot: another goroutine can change it the moment after you read it. It's only reliable when one goroutine owns the channel.

## 4. close, range and comma-ok

`close(ch)` means "no more values". Values already in the buffer are still delivered; after that, every receive returns **immediately** with the zero value.

```go
ch := make(chan int, 2)
ch <- 7
close(ch)

v, ok := <-ch
fmt.Println(v, ok) // 7 true  (a real value, sent before close)
v, ok = <-ch
fmt.Println(v, ok) // 0 false (closed and empty)
```

`for v := range ch` receives until the channel is **closed and drained**. It's the normal way to consume a channel:

```go
ch := make(chan int)
go func() {
	for i := range 3 {
		ch <- i * 10
	}
	close(ch) // ends the receiver's range loop
}()

for v := range ch {
	fmt.Println(v) // 0, then 10, then 20
}
```

| Receive form | On a closed, empty channel |
|---|---|
| `v := <-ch` | returns the zero value at once, so you can't tell it from a real `0` |
| `v, ok := <-ch` | `ok == false`: use it when you read one value at a time |
| `for v := range ch` | the loop ends |

- ⚠️ A plain `for { v := <-ch; … }` never stops on a closed channel: it spins forever on zero values. Use `range`, or check `ok`.
- If nobody closes the channel, `range` waits forever (§9).

## 5. Who closes: only the sender

| Operation | Open | Closed | nil |
|---|---|---|---|
| `ch <- v` | sends (or blocks) | ❌ `panic: send on closed channel` | blocks forever |
| `<-ch` | receives (or blocks) | zero value, `ok == false`, never blocks | blocks forever |
| `close(ch)` | closes | ❌ `panic: close of closed channel` | ❌ `panic: close of nil channel` |

- The **sender** closes, because only the sender knows there's nothing more to send. A receiver that closes can make the sender panic.
- **Several senders?** None of them may close on its own. Wait for all of them (`WaitGroup`), then close once from one place:

```go
out := make(chan int)
var wg sync.WaitGroup
for i := range 3 {
	wg.Go(func() { out <- i })
}
go func() {
	wg.Wait()
	close(out) // after the last send
}()
```

- Closing is a **signal**, not cleanup. A channel nobody uses any more is garbage-collected whether it's closed or not. Close only when a receiver needs to know the stream ended (it uses `range`, `ok`, or waits on it as a done signal).
- `defer close(out)` at the top of the sending goroutine closes it on every return path.

## 6. Directional channel types

A channel type can be restricted to one direction. The compiler then enforces who sends and who receives.

| Type | Can | Can't |
|---|---|---|
| `chan T` | send, receive, close | — |
| `chan<- T` (send-only) | send, close | receive |
| `<-chan T` (receive-only) | receive | send, close |

```go
func produce(out chan<- int) { // can only send
	out <- 1
	close(out)
}

func consume(in <-chan int) { // can only receive
	for v := range in {
		fmt.Println("got", v)
	}
}

ch := make(chan int, 1)
produce(ch) // chan int converts implicitly to chan<- int
consume(ch) // ... and to <-chan int; prints: got 1
```

| You write | Result |
|---|---|
| `r <- 1` with `r <-chan int` | ❌ `invalid operation: cannot send to receive-only channel <-chan int r` |
| `<-s` with `s chan<- int` | ❌ `invalid operation: cannot receive from send-only channel chan<- int s` |
| `close(r)` with `r <-chan int` | ❌ `invalid operation: cannot close receive-only channel r` |
| `var back chan int = r` | ❌ `cannot use r (variable of type <-chan int) as chan int value`: a direction can be removed, never added back |

- The arrow shows the data flow: `chan<-` points **into** the channel, `<-chan` points **out of** it.
- Put directions in **function signatures** (parameters and results). Inside the function that owns the channel, `make` a bidirectional one.

## 7. nil channels block forever

```go
var ch chan int // nil: never made
<-ch            // fatal error: all goroutines are asleep - deadlock!
                // goroutine 1 [chan receive (nil chan)]:
```

- Sending to or receiving from a nil channel blocks **forever**; it never panics. `close` of a nil channel panics.
- The usual cause is a forgotten `make`: a struct field `ch chan int` or a function that returns `nil` by mistake.
- It's also a tool: lesson 3 sets a channel variable to `nil` on purpose to switch off one case of a `select`.

## 8. Generators: functions that return `<-chan T`

A **generator** starts a goroutine that produces values and returns the receive end. The caller can only read; the goroutine owns the channel and closes it.

```go
// Count sends 1, 2, …, n and then closes the channel.
func Count(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // the owner closes
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out // chan int → <-chan int, implicitly
}

for v := range Count(3) {
	fmt.Print(v, " ") // 1 2 3
}
```

A **stage** reads one channel and writes another, so stages chain into a **pipeline**:

```go
// Double sends 2*v for every v from in, and closes when in is closed.
func Double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * 2
		}
	}()
	return out
}

for v := range Double(Count(3)) {
	fmt.Print(v, " ") // 2 4 6
}
```

- The function **returns at once**; the work happens in the goroutine while the caller reads. Values stream through one at a time.
- Closing propagates: `Count` closes → `Double`'s range ends → `Double` closes → the caller's range ends.
- Lesson 6 builds on this (fan-out, fan-in, worker pools).

## 9. "all goroutines are asleep - deadlock!"

When **every** goroutine is blocked on a channel (or a lock), nothing can ever wake them, and the runtime stops the program:

```go
ch := make(chan int)
go func() { ch <- 1 }() // forgot close(ch)

for v := range ch {
	fmt.Println(v) // prints 1, then:
}
// fatal error: all goroutines are asleep - deadlock!
// goroutine 1 [chan receive]:
```

| Cause | Fix |
|---|---|
| unbuffered send in `main` with no other goroutine to receive: `ch := make(chan int); ch <- 1` | receive in another goroutine, or send from one |
| `range` over a channel nobody closes | the sender closes when it's done |
| receiving from a nil channel | `make` it |
| `wg.Wait()` before reading a channel the workers are blocked sending to | read in one goroutine, `Wait` + `close` in another (§5) |

- The goroutine list after the message shows **where** each one is stuck: `[chan send]`, `[chan receive]`, `[chan receive (nil chan)]`.
- ⚠️ The runtime only detects it when **all** goroutines are stuck. If anything else is still running (an HTTP server, a timer, a test framework), a stuck goroutine just waits silently. That's a **goroutine leak** (§10), and it's the more common bug.

## 10. Goroutine leaks

A goroutine blocked forever on a channel is never freed. Stop reading a generator early and its goroutine stays stuck on its next send:

```go
before := runtime.NumGoroutine()
for v := range Count(1000) {
	if v == 2 {
		break // nobody will ever receive 3
	}
}
time.Sleep(10 * time.Millisecond)
fmt.Println(runtime.NumGoroutine() - before) // 1: the Count goroutine is stuck on out <- 3
```

Two ways to make sure a producer can always finish:

| Fix | When |
|---|---|
| **Buffer one slot per value**: `make(chan Result, len(jobs))` | the number of values is known up front; the producer never blocks, even if the reader stops early |
| **A done channel** (§11) | the stream is long or endless; the consumer tells the producer to stop |

## 11. The done channel: stopping a producer

Closing a channel wakes **every** receiver at once, so `close(done)` is a broadcast "stop" signal. The producer waits on two things at once with **`select`**, which blocks until one of its cases can proceed and then runs that one (lesson 3 covers it fully):

```go
// Naturals sends 1, 2, 3, … until done is closed.
func Naturals(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; ; i++ {
			select {
			case out <- i: // the consumer took i
			case <-done: // the consumer wants us to stop
				return
			}
		}
	}()
	return out
}

done := make(chan struct{})
for v := range Naturals(done) {
	fmt.Print(v, " ") // 1 2 3
	if v == 3 {
		break
	}
}
close(done) // the producer returns; NumGoroutine is back to where it started
```

- `done` is `<-chan struct{}`: the producer only listens. The **consumer** creates and closes it, because it's the one that knows when to stop. (It's the consumer's *signal* channel; `out` is still closed by its sender.)
- Use `close(done)`, not `done <- struct{}{}`: a send wakes one goroutine, a close wakes them all, and it works even if a producer isn't listening at that moment.
- Put **every** blocking send in the producer inside such a `select`, or it can still get stuck on the one that isn't.
- After `close(done)` both cases may be ready at once, and `select` picks randomly, so the producer may send a value or two more before it notices. Consumers that stop early just ignore them.
- `context.Context` (lesson 5) is this exact pattern with cancellation reasons and deadlines added; its `Done()` method returns a `<-chan struct{}`.

## Common mistakes

| Mistake | Fix |
|---|---|
| Unbuffered `ch <- v` in `main` with no receiver → deadlock | Send from a goroutine, or receive in one |
| `range` over a channel that's never closed | The sender `defer close(out)`s when it finishes |
| The receiver closes the channel, the sender panics `send on closed channel` | Only the sender closes; with several senders, `wg.Wait()` then `close` once |
| Closing twice (`close of closed channel`) | One owner closes, in one place (`defer close(out)`) |
| `for { v := <-ch }` spins on zero values after close | `for v := range ch`, or `v, ok := <-ch` and stop when `!ok` |
| `var ch chan int` without `make` → blocks forever | `ch := make(chan int)` |
| Adding a buffer to "fix" a deadlock | Find the missing receiver / close; size buffers on purpose |
| Breaking out of a generator's range → leaked producer | A done channel, or a buffer big enough for every value |
| Returning `chan T` from a generator | Return `<-chan T` so callers can't send or close |
| Calling `close` on a `<-chan T` parameter | Won't compile; the owner of the bidirectional channel closes it |

## Key takeaways

- `make(chan T)` is unbuffered: send and receive meet, so it **synchronises**, and everything before the send is visible after the receive. `make(chan T, n)` is a FIFO queue that blocks only when full or empty.
- `close` means "no more values". Receivers drain the buffer, then get the zero value with `ok == false`; `range` stops.
- **Only the sender closes.** Send on a closed channel, or close twice, panics. Closing isn't needed for garbage collection.
- Directional types (`chan<- T`, `<-chan T`) document and enforce who sends and who receives; `chan T` converts to either implicitly, never back.
- A nil channel blocks forever on send and receive.
- Generators return `<-chan T`, own the goroutine and `defer close(out)`; stages chain into pipelines.
- `fatal error: all goroutines are asleep - deadlock!` means every goroutine is blocked. If something else is running you get a silent **leak** instead.
- Make sure every producer can finish: a buffer sized to the output, or a **done channel** that the consumer closes and the producer `select`s on.
