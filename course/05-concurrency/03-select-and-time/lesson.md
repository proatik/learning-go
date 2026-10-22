# select & Timers

## 1. `select`: wait on several channels at once

Lesson 2's `<-ch` waits on **one** channel. `select` waits on several and runs the **first case that can proceed**.

```go
a := make(chan string)
b := make(chan string)
go func() { time.Sleep(20 * time.Millisecond); a <- "from a" }()
go func() { time.Sleep(10 * time.Millisecond); b <- "from b" }()

for range 2 {
	select {
	case msg := <-a:
		fmt.Println(msg)
	case msg := <-b:
		fmt.Println(msg)
	}
}
// from b
// from a
```

- Each `case` is one send (`case ch <- v:`) or one receive (`case v := <-ch:`, `case v, ok := <-ch:`, `case <-ch:`).
- `select` **blocks** until at least one case can proceed, then runs exactly **one** case.
- If several are ready at once, it picks one **uniformly at random**. Never rely on case order:

```go
x := make(chan int, 1)
y := make(chan int, 1)
var counts [2]int
for range 1000 {
	x <- 0
	y <- 1
	select { // both are ready
	case i := <-x:
		counts[i]++
	case i := <-y:
		counts[i]++
	}
	select { // drain whichever is still full
	case <-x:
	case <-y:
	}
}
fmt.Println(counts) // roughly [500 500], different on every run (e.g. [503 497])
```

- `select {}` (no cases) blocks forever.

## 2. `default`: never block

A `default` case runs when **no other case is ready right now**. That turns any send or receive into a non-blocking "try".

```go
ch := make(chan int, 1)
for range 2 {
	select {
	case ch <- 1:
		fmt.Println("sent")
	default:
		fmt.Println("full")
	}
}
for range 2 {
	select {
	case v := <-ch:
		fmt.Println("got", v)
	default:
		fmt.Println("empty")
	}
}
// sent
// full
// got 1
// empty
```

| Use | Pattern |
|---|---|
| Try to send (drop if full) | `select { case ch <- v: default: }` |
| Try to receive (poll) | `select { case v := <-ch: …; default: … }` |
| Check "is this channel closed yet?" without waiting | `select { case <-done: return; default: }` |

- ⚠️ `default` inside a `for` loop with nothing else blocking is a **busy loop**: it spins at 100% CPU. Use `default` for one-off checks, not for waiting.

## 3. Timeouts with `time.After`

`time.After(d)` returns a `<-chan time.Time` that receives once, after `d`. As a `select` case, it puts a time limit on the other cases.

```go
func fetch() <-chan string {
	ch := make(chan string, 1) // buffered: the goroutine can finish even if nobody waits
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- "data"
	}()
	return ch
}

select {
case v := <-fetch():
	fmt.Println(v)
case <-time.After(10 * time.Millisecond):
	fmt.Println("timeout")
}
// timeout
```

- The **buffer of 1** matters: after a timeout nobody receives, and an unbuffered send would block that goroutine forever (a goroutine leak).
- Where you create the `time.After` decides what kind of timeout you get:

```go
// Overall deadline: one timer for the whole loop.
deadline := time.After(time.Second)
for {
	select {
	case v := <-results:
		handle(v)
	case <-deadline:
		return // 1s after the loop started, however busy it was
	}
}

// Idle timeout: a new timer on every iteration.
for {
	select {
	case v := <-results:
		handle(v)
	case <-time.After(time.Second):
		return // 1s with no value
	}
}
```

- Since Go 1.23 an un-fired `time.After` timer is garbage-collected normally, so `time.After` in a loop doesn't leak. (Old advice to avoid it is outdated.)
- For timeouts that cancel work across functions and goroutines, lesson 5's `context.WithTimeout` is the standard tool. `time.After` is for a local `select`.

## 4. `time.Timer`: a timeout you can stop and restart

`time.NewTimer(d)` returns a `*time.Timer` whose channel `C` receives once after `d`. Unlike `time.After`, you can cancel and reuse it.

```go
t := time.NewTimer(50 * time.Millisecond)
fmt.Println(t.Stop()) // true: it was running, now it won't fire
fmt.Println(t.Stop()) // false: already stopped
fmt.Println(t.Reset(10 * time.Millisecond)) // false: it wasn't running; now it is
<-t.C
fmt.Println("fired")
fmt.Println(t.Reset(10 * time.Millisecond)) // false: it had already fired
<-t.C
// true
// false
// false
// fired
// false
```

| Method | Does | Returns |
|---|---|---|
| `t.Stop()` | cancels the timer | `true` if it was running |
| `t.Reset(d)` | (re)starts it to fire after `d` from now | `true` if it was running |
| `time.AfterFunc(d, f)` | runs `f` in its own goroutine after `d`; `Stop` cancels it | the `*Timer` (its `C` is `nil`) |

- **Go 1.23 semantics** (what you have): after `Stop` or `Reset` returns, `C` never delivers a stale value from before. You **don't** need the old "stop, then drain `<-t.C`" dance you'll see in older code.
- The classic use is an **inactivity timer**: one timer, `Reset` on every event.

```go
idle := time.NewTimer(d)
defer idle.Stop()
for {
	select {
	case ev := <-events:
		handle(ev)
		idle.Reset(d) // push the deadline back
	case <-idle.C:
		return // nothing for d
	}
}
```

```go
done := make(chan struct{})
time.AfterFunc(10*time.Millisecond, func() {
	fmt.Println("AfterFunc ran")
	close(done)
})
<-done
// AfterFunc ran
```

## 5. `time.Ticker`: something every interval

`time.NewTicker(d)` delivers on `C` every `d` until you call `Stop`.

```go
tk := time.NewTicker(10 * time.Millisecond)
defer tk.Stop()
start := time.Now()
for i := range 3 {
	<-tk.C
	fmt.Println(i+1, time.Since(start).Round(10*time.Millisecond))
}
// about 1 10ms, 2 20ms, 3 30ms (real time: the exact values depend on the machine)
```

- The first tick comes **after** one interval, not at time 0.
- `NewTicker(d)` and `Reset(d)` **panic** if `d <= 0`.
- A slow receiver doesn't get a backlog: the ticker adjusts the interval or **drops ticks**.
- `defer tk.Stop()` when you're done. It's no longer needed for garbage collection (Go 1.23), but a ticker keeps firing until it's stopped or unreachable.
- `time.Tick(d)` is a shortcut that returns only the channel (no way to stop it); it returns `nil` for `d <= 0`. Fine for a program that ticks forever, not for a function that returns.

| Need | Use |
|---|---|
| one deadline in a `select` | `time.After(d)` |
| a deadline you'll cancel or move | `time.NewTimer(d)` + `Stop` / `Reset` |
| run a function later | `time.AfterFunc(d, f)` |
| repeat every `d` | `time.NewTicker(d)` + `defer Stop()` |

## 6. Switching a case off with a `nil` channel

Lesson 2: a send or receive on a `nil` channel blocks **forever**. In a `select`, that means the case is simply **never chosen**. Setting a channel variable to `nil` switches its case off.

The key use: a closed channel is **always** ready (it returns the zero value forever), so a `select` loop must turn it off.

```go
func sum(a, b <-chan int) int {
	total := 0
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil // closed: never select this case again
				continue
			}
			total += v
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			total += v
		}
	}
	return total
}

fmt.Println(sum(gen(1, 2, 3), gen(10, 20))) // 36 (gen from lesson 2: sends its values, then closes)
```

The same trick switches a **timer** case on and off: keep `var fire <-chan time.Time` as `nil` while there's nothing to wait for, and set `fire = timer.C` when there is.

| Variable is | Its `select` case |
|---|---|
| a ready channel | can be chosen |
| a closed channel (receive) | always ready, gives zero values: turn it off! |
| `nil` | never chosen |

## 7. The `for` + `select` loop

Most concurrent code is a loop around a `select`, with one case for work and one for "stop".

```go
func worker(jobs <-chan int, done <-chan struct{}) {
	for {
		select {
		case j := <-jobs:
			fmt.Println("job", j)
		case <-done:
			return
		}
	}
}
```

- ⚠️ A plain `break` inside a `select` leaves only the **`select`**, not the loop. Use `return`, or a labeled `break`:

```go
loop:
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				break loop // leaves the for
			}
			fmt.Println(v)
		}
	}
```

- Prefer a small function with `return` over labels.

## 8. Testing time-based code with `testing/synctest`

Code with timeouts, tickers and backoff is hard to test with real time: tests are slow (a 5s timeout needs 5s) and flaky (a busy machine shifts timings). Go 1.25's `testing/synctest` fixes both.

`synctest.Test(t, f)` runs `f` in a **bubble**: an isolated group of goroutines with its own **fake clock**.

```go
// Countdown sends n, n-1, …, 1 on the returned channel, one per second,
// then closes it.
func Countdown(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for i := n; i > 0; i-- {
			<-tick.C
			ch <- i
		}
	}()
	return ch
}

func TestCountdown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now() // fake clock: 2000-01-01 00:00:00 UTC
		var got []string
		for v := range Countdown(3) {
			got = append(got, fmt.Sprintf("%d@%v", v, time.Since(start)))
		}
		want := []string{"3@1s", "2@2s", "1@3s"}
		if !slices.Equal(got, want) {
			t.Errorf("Countdown(3) = %v, want %v", got, want)
		}
	})
}
// --- PASS: TestCountdown (0.00s)
```

How the fake clock works:

- Inside the bubble, `time.Now`, `Sleep`, `After`, `NewTimer`, `NewTicker`, `AfterFunc` all use the bubble's clock. It starts at **2000-01-01 00:00:00 UTC**.
- Time **only moves when every goroutine in the bubble is durably blocked**. Then it jumps straight to the next timer. So "3 seconds" takes no real time, and `time.Since` is **exact**: `3s`, never `3.0021s`.
- That makes timing assertions deterministic: "returns after exactly 250ms", "retries at 0, 100ms, 300ms", "N tasks of 1s each, run concurrently, take 1s total".

**`synctest.Wait()`** blocks until every **other** goroutine in the bubble is durably blocked. Use it to check state "once everything has settled", without sleeping:

```go
func TestAfterFunc(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		fired := false
		time.AfterFunc(time.Second, func() { fired = true })

		time.Sleep(999 * time.Millisecond)
		synctest.Wait() // everything else in the bubble is now blocked
		if fired {
			t.Fatal("fired before 1s")
		}

		time.Sleep(time.Millisecond)
		synctest.Wait() // the AfterFunc goroutine has run and exited
		if !fired {
			t.Fatal("not fired at 1s")
		}
	})
}
// --- PASS: TestAfterFunc (0.00s)
```

What counts as **durably blocked** (only another bubble goroutine can wake it):

| Durably blocked ✅ | Not durably blocked ❌ (time won't advance) |
|---|---|
| send / receive on a channel **created in the bubble** | a channel created **outside** the bubble |
| a `select` whose cases are all bubble channels | `sync.Mutex` / `RWMutex` lock |
| `time.Sleep` | network or file I/O, system calls |
| `sync.WaitGroup.Wait` (with `Add`/`Go` called in the bubble), `sync.Cond.Wait` | |

**The pitfall: a bubble must end clean.** `synctest.Test` waits for every goroutine in the bubble to exit. If they can't, it **panics**, and the panic kills the whole test binary:

| Situation | Result |
|---|---|
| `f` returns while a goroutine is still blocked (a leak) | `panic: deadlock: main bubble goroutine has exited but blocked goroutines remain` |
| everything is blocked and no timer can wake anyone | `panic: deadlock: all goroutines in bubble are blocked` |
| a timer or `AfterFunc` that hasn't fired yet when `f` returns | fine: no goroutine is waiting on it |

```go
synctest.Test(t, func(t *testing.T) {
	ch := make(chan int)
	go func() { ch <- 1 }() // nobody ever receives
})
// panic: deadlock: main bubble goroutine has exited but blocked goroutines remain
```

So in a bubble:

- Let background goroutines **finish** before `f` returns: `time.Sleep(time.Hour)` (instant) lets sleeping senders complete; buffered channels let senders exit without a receiver; call `Stop` / close a done channel for long-running loops.
- This is also a feature: a leaked goroutine **can't** go unnoticed.
- Don't call `t.Run`, `t.Parallel` or `t.Deadline` inside the bubble. Use one `synctest.Test` per table case instead.
- Operating on a bubble's channel or timer from **outside** the bubble panics.

## Common mistakes

| Mistake | Fix |
|---|---|
| Relying on `select` case order | Several ready cases are chosen at random; design for any order |
| `for { select { …; default: } }` to wait | Busy loop at 100% CPU; drop the `default` so `select` blocks |
| `break` inside `select` to leave the loop | `return`, or a labeled `break loop` |
| Unbuffered result channel + timeout | The goroutine leaks on timeout; `make(chan T, 1)` |
| `time.After` inside the loop when you meant one overall deadline | Create it once, before the loop |
| Keeping a closed channel in a `select` loop | It's always ready with zero values; set it to `nil` |
| Old "stop then drain `<-t.C`" before `Reset` | Not needed since Go 1.23; just `Reset` |
| `time.NewTicker(0)` | Panics; validate `d > 0` first |
| Forgetting `tk.Stop()` in a function that returns | `defer tk.Stop()` right after `NewTicker` |
| Testing time code with real `time.Sleep` | Slow and flaky; use `synctest.Test` and its fake clock |
| Leaving goroutines blocked when a bubble returns | Let them finish (sleep, buffer, stop); otherwise the test binary panics |

## Key takeaways

- `select` waits on several channel operations and runs one ready case, chosen at random if several are ready.
- `default` makes a `select` non-blocking: try-send, try-receive, "is it closed yet?".
- `time.After` adds a timeout case; give result channels a buffer of 1 so goroutines can still finish.
- `time.Timer` can be `Stop`ped and `Reset` (no draining since Go 1.23); `time.Ticker` repeats until `Stop`.
- A `nil` channel's case is never chosen: use it to switch off closed inputs and idle timers.
- `for` + `select` is the shape of most concurrent loops; leave it with `return`, not `break`.
- `testing/synctest` gives tests a fake clock: time-based tests run instantly and exactly. `synctest.Wait` waits until everything has settled. A bubble must end with no blocked goroutines.
