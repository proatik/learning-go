# Goroutines

## 1. The `go` statement

A **goroutine** is a function running concurrently with the rest of the program. Put `go` in front of any function or method call to start one.

```go
func say(s string) { fmt.Println(s) }

go say("hello")            // a named function
go func() {                // a function literal; note the () that calls it
	fmt.Println("world")
}()
```

- `go` returns **immediately**. The caller doesn't wait, and the call's return values are thrown away (`v := go f()` is a syntax error).
- The arguments are evaluated **in the caller**, right away: `go say(name)` copies `name` now, even if it changes later.
- Goroutines are cheap: a few KB of stack that grows as needed. Starting 10,000 of them is normal.
- The runtime schedules them onto OS threads. You never pick a thread, and you can't predict **which goroutine runs first**.

| Term | Meaning |
|---|---|
| **concurrency** | several tasks in progress at once (structure) |
| **parallelism** | several tasks executing at the same instant on different CPUs (execution) |

Goroutines give you concurrency; you get parallelism for free when there are several CPUs (§9).

## 2. When `main` returns, everything stops

The program ends when `main` returns. Running goroutines are **killed**, not waited for.

```go
func main() {
	go say("hello from a goroutine")
	fmt.Println("main done")
}
// Output (usually): main done
// The goroutine is killed before it gets to print.
```

- Adding `time.Sleep` "to let it finish" is a guess, not a fix: on a busy machine the sleep may be too short. Wait **explicitly** instead (§3).
- The same holds in tests: a test function that returns doesn't wait for goroutines it started.

## 3. `sync.WaitGroup`: `Add`, `Done`, `Wait`

A `WaitGroup` is a counter of unfinished tasks.

| Method | Does |
|---|---|
| `wg.Add(n)` | counter += n; call it **before** `go` |
| `wg.Done()` | counter -= 1; call it when the task finishes (usually `defer wg.Done()`) |
| `wg.Wait()` | blocks until the counter is 0 |

```go
var wg sync.WaitGroup
for _, name := range []string{"ann", "bob", "cy"} {
	wg.Add(1)
	go func() {
		defer wg.Done() // runs even if the body returns early
		fmt.Println("hi", name)
	}()
}
wg.Wait()
fmt.Println("all done")
// Output: the three "hi" lines in any order, then "all done" (always last).
```

- The zero value is ready to use. **Never copy** a `WaitGroup` after first use: pass `*sync.WaitGroup`, or better, keep it local (lesson 4 shows `go vet`'s copylocks check).
- `Add` **inside** the goroutine is a bug: `Wait` may run first, see 0 and return early. `go vet` reports `WaitGroup.Add called from inside new goroutine`.
- More `Done` than `Add` → `panic: sync: negative WaitGroup counter`.
- More `Add` than `Done` → `Wait` blocks forever. If nothing else can run, the runtime stops with `fatal error: all goroutines are asleep - deadlock!`.

## 4. `wg.Go` (Go 1.25)

`wg.Go(f)` does the `Add(1)`, the `go` and the `defer Done()` for you. Prefer it in new code.

```go
var wg sync.WaitGroup
for _, name := range []string{"ann", "bob", "cy"} {
	wg.Go(func() {
		fmt.Println("hi", name)
	})
}
wg.Wait()
```

| Old style | Go 1.25+ |
|---|---|
| `wg.Add(1); go func() { defer wg.Done(); work() }()` | `wg.Go(work)` or `wg.Go(func() { work() })` |

- `f` takes no arguments and returns nothing: `func()`. Pass data in by capturing variables (§5).
- The docs say "The function f must not panic": a panic still crashes the program (§8).

## 5. Loop variables are per iteration (Go 1.22)

Since Go 1.22, each iteration of a `for` loop gets **its own** copy of the loop variable, so a closure started in the loop sees that iteration's value.

```go
var wg sync.WaitGroup
for i := range 3 {
	wg.Go(func() { fmt.Println(i) }) // each goroutine has its own i
}
wg.Wait()
// Output: 0, 1 and 2, in any order.
```

- Before Go 1.22 all three closures shared one `i` and often printed `3 3 3`. Old code works around it with `i := i` or `go func(i int) { … }(i)`; you don't need that any more (the module's `go.mod` says `go 1.27`).
- This applies to `for i := range n`, `for i, v := range s` and three-clause loops `for i := 0; i < n; i++`.
- Variables declared **outside** the loop are still shared by every goroutine.

## 6. Passing results out

Goroutines can't return values to the caller. The simplest safe pattern: **pre-size a slice and let goroutine `i` write only index `i`**.

```go
nums := []int{1, 2, 3, 4, 5}
squares := make([]int, len(nums)) // allocated before any goroutine starts

var wg sync.WaitGroup
for i, n := range nums {
	wg.Go(func() {
		squares[i] = n * n // each goroutine owns one element
	})
}
wg.Wait() // after this, every write is visible
fmt.Println(squares) // [1 4 9 16 25]
```

| Shared write | Safe? | Why |
|---|---|---|
| `results[i] = v`, each goroutine a **different** `i`, slice made before | ✅ | different elements are different variables; nothing else is touched |
| `results = append(results, v)` from many goroutines | ❌ | every goroutine reads and writes the **same** slice header; results get lost |
| `m[key] = v` on one map from many goroutines | ❌ | a map is one shared data structure; the runtime usually stops with `fatal error: concurrent map writes` |
| `count++` on a shared `int` | ❌ | read-modify-write: two goroutines can read the same old value |

- Reading the results is only safe **after** `wg.Wait()`: `Wait` returning guarantees every goroutine's writes are visible.
- Results come out **in input order** for free, whatever order the goroutines finish in.
- Need one map? Give each goroutine its **own** local map, store it at its index, and merge them after `Wait`.
- These ❌ rows are **data races**. Lesson 4 covers how to detect them and how to share safely with mutexes and atomics; channels (lesson 2) are the other way to pass results.

## 7. Errors from goroutines

Same pattern: an `[]error` with one slot per task, combined after `Wait`.

```go
tasks := []func() error{
	func() error { return nil },
	func() error { return errors.New("disk full") },
	func() error { return errors.New("timeout") },
}

errs := make([]error, len(tasks))
var wg sync.WaitGroup
for i, task := range tasks {
	wg.Go(func() { errs[i] = task() })
}
wg.Wait()

err := errors.Join(errs...) // nil if every element is nil; nils are skipped
fmt.Println(err)
// Output:
// disk full
// timeout
```

- The order is the **input** order, so the message is deterministic.
- `errors.Is` / `errors.As` see through `errors.Join` (Module 3).
- This runs **every** task even after one fails. Stopping the others early needs cancellation: lessons 5 (`context`) and 6 (errgroup).

## 8. Panics inside goroutines

A panic that isn't recovered **in the goroutine where it happened** crashes the whole program, whatever `main` does.

```go
defer func() { recover() }() // ❌ useless: it's in main, the panic is in another goroutine
go func() { panic("boom") }()
// panic: boom
//
// goroutine 7 [running]:
// ...
```

- `recover` only works in a deferred call **in the same goroutine**. To survive a panicking task, recover inside each goroutine and turn the panic into an error:

```go
func safely(i int, task func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("task %d panicked: %v", i, r)
		}
	}()
	task()
	return nil
}

errs := make([]error, 2)
var wg sync.WaitGroup
for i, task := range []func(){func() {}, func() { panic("boom") }} {
	wg.Go(func() { errs[i] = safely(i, task) })
}
wg.Wait()
fmt.Println(errs) // [<nil> task 1 panicked: boom]
```

- A **named result** (`err error`) is what lets the deferred function change the return value (Module 3).
- The `recover` is inside `safely`, which runs inside the goroutine, so `wg.Go`'s "f must not panic" rule is kept.
- Runtime errors recover too: calling a nil `func()` gives `runtime error: invalid memory address or nil pointer dereference`; `panic(nil)` gives a `*runtime.PanicNilError` ("panic called with nil argument").
- Recover at a **boundary** (one task, one request), report it, and keep going. Don't sprinkle `recover` everywhere.

## 9. How many goroutines, how many CPUs

| Function | Returns |
|---|---|
| `runtime.NumGoroutine()` | goroutines that currently exist (including `main`) |
| `runtime.NumCPU()` | logical CPUs usable by the process |
| `runtime.GOMAXPROCS(0)` | how many goroutines may **execute** at the same instant; `0` only reads it |

```go
fmt.Println(runtime.NumGoroutine()) // 1: just main
fmt.Println(runtime.NumCPU(), runtime.GOMAXPROCS(0))
// Output on an 8-CPU machine: 8 8 (machine-dependent)
```

- `GOMAXPROCS` defaults to the CPU count (in containers, the CPU limit). You almost never change it; `GOMAXPROCS=1 go run .` is handy to show that code still works with **no** parallelism.
- With `GOMAXPROCS=1` goroutines are still **concurrent**: while one sleeps or waits, another runs. 100 goroutines that each `time.Sleep(time.Second)` take about 1s, not 100s.
- `NumGoroutine` growing without bound means a **goroutine leak**: goroutines stuck waiting forever. Every goroutine you start needs a way to finish.
- CPU-bound work gains up to about `GOMAXPROCS`× from goroutines; waiting-bound work (network, disk, sleeps) gains much more.

## 10. How the checkers test concurrent code

The exercises' tests use `testing/synctest`, a fake clock where `time.Sleep(time.Second)` returns instantly but the clock still moves exactly 1s, so "3 one-second tasks finish in exactly 1s" is a deterministic check. Lesson 3 teaches it.

## Common mistakes

| Mistake | Fix |
|---|---|
| `main` (or a test) returns before its goroutines finish | `wg.Wait()` before returning |
| `time.Sleep` to "wait for" goroutines | A `WaitGroup` (or a channel, lesson 2) |
| `wg.Add(1)` inside the goroutine | `Add` before `go`, or just `wg.Go(f)` |
| Forgetting `Done` on an early `return` | `defer wg.Done()` as the first line, or `wg.Go` |
| Passing a `WaitGroup` by value | Pass `*sync.WaitGroup`, or keep it in one function |
| `results = append(results, v)` from goroutines | `results := make([]T, n)` up front, then `results[i] = v` |
| Writing one map from several goroutines | One local map per goroutine, merged after `Wait` (or a mutex, lesson 4) |
| Reading results before `Wait` returns | Read only after `wg.Wait()` |
| `recover` in `main` for a panic in a goroutine | `defer`/`recover` **inside** the goroutine |
| `i := i` copies in loops | Not needed since Go 1.22 |
| Expecting goroutines to run in start order | Order is unpredictable; put results at fixed indexes |

## Key takeaways

- `go f(x)` starts `f` concurrently and returns at once; arguments are evaluated immediately, return values are discarded.
- When `main` returns, all goroutines die. Wait explicitly, never with `time.Sleep`.
- `sync.WaitGroup`: `Add` before `go`, `defer Done()`, then `Wait`. Since Go 1.25, `wg.Go(f)` does all three.
- Loop variables are per iteration since Go 1.22, so closures in loops are safe.
- Return results through a **pre-sized slice**, one index per goroutine, read after `Wait`. Shared `append`, maps and counters are data races (lesson 4).
- Collect errors the same way and combine them with `errors.Join` in input order.
- A panic kills the program unless recovered **in its own goroutine**; recover at the task boundary and return it as an error.
- `runtime.NumGoroutine()` counts goroutines (watch for leaks); `GOMAXPROCS` limits parallel execution, not how many goroutines exist.
