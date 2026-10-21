# panic, recover & defer Patterns

## 1. What a panic is

A **panic** is Go's way of stopping a goroutine that cannot go on: a bug, a broken invariant, an impossible state.

```go
func main() {
	fmt.Println("start")
	panic("something impossible happened")
	fmt.Println("never printed") // unreachable
}
```

```text
start
panic: something impossible happened

goroutine 1 [running]:
main.main()
	/path/main.go:5 +0x...
exit status 2
```

- `panic(v)` accepts **any value** (`panic(v any)`): a string, an `error`, a struct…
- An unrecovered panic prints `panic: <v>` plus a **stack trace** and exits with status **2**.
- Panics are **not** Go's exception system. Expected failures (bad input, missing file, network down) are `error` values (lesson 04).

## 2. Run-time panics

The runtime panics by itself when an operation is impossible. Each value implements **`runtime.Error`** (an `error` with an extra marker method `RuntimeError()`).

| Operation | Panic message (`%v` of the value) |
|---|---|
| `var m map[string]int; m["a"] = 1` | `assignment to entry in nil map` |
| `var p *Point; _ = p.X` | `runtime error: invalid memory address or nil pointer dereference` |
| `var f func(); f()` | same nil-dereference message |
| `s[5]` with `len(s) == 3` | `runtime error: index out of range [5] with length 3` |
| `s[-1]` (index in a variable) | `runtime error: index out of range [-1]` |
| `s[1:5]` with `cap(s) == 3` | `runtime error: slice bounds out of range [:5] with capacity 3` |
| `var x any = "hi"; x.(int)` | `interface conversion: interface {} is string, not int` |
| `a / b` with `b == 0` (integers) | `runtime error: integer divide by zero` |
| `close(ch)` twice (concurrency module) | `close of closed channel` |

```go
var x any = "hi"
n, ok := x.(int) // comma-ok form: never panics → 0 false
n = x.(int)      // single-value form: panics on mismatch
```

- Reading a **nil map** is fine (zero value); only *writing* panics.
- Floating-point division by zero does **not** panic: `1.0 / zero` is `+Inf`.
- `math.MinInt / -1` does not panic either: it wraps to `math.MinInt`.
- Constant expressions are checked at compile time: `10 / 0` is a compile error, not a panic.
- Some failures are **fatal errors**, not panics, and can't be recovered: stack overflow (`fatal error: stack overflow`), out of memory, concurrent map writes.

## 3. Unwinding: deferred calls still run

When `F` panics, normal execution stops, then:

1. `F`'s deferred calls run (LIFO), then `F` returns to its caller **still panicking**.
2. The same happens in the caller, and its caller… up the stack (**unwinding**).
3. If nothing recovers, the program crashes after the last deferred call of the goroutine has run.

```go
func level(n int) {
	fmt.Println("enter", n)
	defer fmt.Println("exit", n)
	if n == 3 {
		panic("bottom reached")
	}
	level(n + 1)
}

func main() {
	defer fmt.Println("main's defer still runs")
	level(1)
}
```

```text
enter 1
enter 2
enter 3
exit 3
exit 2
exit 1
main's defer still runs
panic: bottom reached
...
```

This is why `defer` is the right place for cleanup: it runs on **every** exit path, `return` or panic.

## 4. `recover`

`recover()` stops the unwinding and returns the value passed to `panic`. It returns `nil` when the goroutine is not panicking.

```go
func safely(f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	f()
	fmt.Println("f returned normally") // skipped if f panics
}

safely(func() { panic("boom") }) // recovered: boom
fmt.Println("program continues")  // program continues
```

After a successful recover, the function containing the `defer` **returns normally** to its caller (with whatever its results are at that moment).

⚠️ `recover` only works when called **directly by a deferred function**:

| Code | Recovers? |
|---|---|
| `defer func() { recover() }()` | ✅ the deferred function calls it |
| `defer handle()` where `handle` calls `recover()` | ✅ `handle` *is* the deferred function |
| `defer func() { handle() }()` where `handle` calls `recover()` | ❌ one call too deep → `nil` |
| `recover()` in normal (non-deferred) code | ❌ always `nil` |
| `defer recover()` | ❌ `recover` itself is the deferred call, not called *by* a deferred function |

- A panic in **another goroutine** cannot be recovered here: `recover` only sees the panic of its own goroutine (concurrency module).

## 5. Turning a panic into an error

The classic pattern: a **named result** plus a deferred closure that sets it.

```go
func ParseConfig(src string) (cfg Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse config: %v", r)
		}
	}()
	cfg = mustParse(src) // deep internal code may panic
	return cfg, nil
}
```

- The deferred closure runs **after** the panic stopped the body, so it is the only place that can still set `err`.
- Without a **named** result there is nothing to assign: the caller would get the zero values.
- Keep the panic value inspectable: if `r` is an `error`, wrap it with `%w`.

```go
defer func() {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			err = fmt.Errorf("parse config: %w", e) // errors.Is / errors.As still work
		} else {
			err = fmt.Errorf("parse config: %v", r)
		}
	}
}()
```

Inside the standard library, `encoding/json` and `text/template` use exactly this: internal code panics with a private error type, and the exported function recovers it and returns an `error`. **Panics never cross the package API.**

## 6. Re-panicking unexpected values

Only recover what you **expect**. Anything else is a real bug and must keep crashing: re-panic it.

```go
type parseError struct{ msg string }

func Parse(s string) (n int, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return // no panic
		}
		pe, ok := r.(parseError)
		if !ok {
			panic(r) // not ours: keep unwinding with the same value
		}
		err = errors.New(pe.msg)
	}()
	return parse(s), nil // parse panics with parseError on bad input
}
```

The crash report then shows both stages: `panic: <value> [recovered, repanicked]`.

Distinguish run-time panics with the `runtime.Error` interface:

```go
if re, ok := r.(runtime.Error); ok {
	err = fmt.Errorf("internal error: %w", re)
}
```

## 7. `panic(nil)` since Go 1.21

Before Go 1.21, `panic(nil)` made `recover()` return `nil`, indistinguishable from "no panic". Now it panics with a **`*runtime.PanicNilError`**:

```go
defer func() {
	r := recover()
	fmt.Printf("%T: %v\n", r, r)
	// *runtime.PanicNilError: panic called with nil argument
}()
panic(nil)
```

- So `if r := recover(); r != nil` is reliable again.
- `*runtime.PanicNilError` is a `runtime.Error`: `errors.As(err, &pnErr)` works once wrapped with `%w`.
- Old behaviour returns only with `GODEBUG=panicnil=1`. Never panic with `nil` on purpose.

## 8. When to panic, when to return an error

| Return an `error` | `panic` |
|---|---|
| Anything the **caller** can cause or expect: bad input, missing file, timeout | A **programmer bug**: impossible `switch` case, broken invariant |
| Library APIs | Unrecoverable setup in `init` or `main` (bad embedded config) |
| Failures you want callers to handle | `Must` helpers on **constant** inputs: `regexp.MustCompile` |

```go
func weekdayName(d int) string {
	switch d {
	case 0:
		return "Sunday"
	// ...
	case 6:
		return "Saturday"
	}
	panic(fmt.Sprintf("weekdayName: invalid day %d", d)) // only reachable through a bug
}
```

- A good panic message names the function and the bad value.
- Never use panic/recover as `try/catch` for normal control flow (e.g. to exit a loop early): it's slow, hides bugs and surprises readers.

## 9. `Must` helpers

A `MustX` function wraps an `(T, error)` function and **panics instead of returning the error**. Use it only where an error means "the program itself is wrong": package-level variables, tests, constant inputs.

```go
var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // a typo in the pattern crashes at startup
var page = template.Must(template.New("p").Parse("Hi {{.}}")) // template.Must(t, err)

// MustAtoi is Atoi for inputs that are known to be valid.
func MustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("MustAtoi(%q): %w", s, err)) // an error value: errors.Is still works after recovery
	}
	return n
}

port := MustAtoi("8080") // fine
MustAtoi("80a")          // panic: MustAtoi("80a"): strconv.Atoi: parsing "80a": invalid syntax
```

- Convention: the `Must` prefix **warns the caller** that the function panics.
- The standard library pairs them: `regexp.Compile` / `regexp.MustCompile`, `t.Parse` / `template.Must`.
- With generics (later module) one `Must[T any](v T, err error) T` covers every type; the idea is identical.

## 10. Defer patterns

### Cleanup next to acquisition

```go
f, err := os.Open(name)
if err != nil {
	return err
}
defer f.Close() // runs on every return path below, and on panic
```

- Defer **after** the error check: if `Open` failed, `f` is `nil`.
- Unlocking works the same way: `mu.Lock(); defer mu.Unlock()` (mutexes come in the concurrency module).

### Timing a function

```go
func timed(name string) func() {
	start := time.Now()
	return func() { fmt.Println(name, "took", time.Since(start)) }
}

func work() {
	defer timed("work")() // timed("work") runs NOW, the returned func at exit
	// ...
}
```

The **extra `()`** matters: `defer timed("work")` would defer the whole call and measure nothing.

### Tracing enter / exit

```go
func trace(name string) string {
	fmt.Println("enter", name)
	return name
}

func un(name string) { fmt.Println("exit", name) }

func a() {
	defer un(trace("a")) // the argument trace("a") is evaluated immediately
	fmt.Println("in a")
}
// enter a
// in a
// exit a
```

### Capturing the error from `Close`

`defer f.Close()` **drops** `Close`'s error. That's fine for files you only read, but for writes `Close` may be the call that reports the failure (data flushed at close). Capture it into the named result, without overwriting an earlier error:

```go
func writeReport(name string, lines []string) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr // report the Close error only if everything else succeeded
		}
	}()
	for _, l := range lines {
		if _, err := fmt.Fprintln(f, l); err != nil {
			return err // this error wins; Close still runs
		}
	}
	return nil
}
```

Alternative when both matter: `err = errors.Join(err, cerr)`.

### `defer` in loops

Deferred calls run at **function** exit, so this keeps every file open until the loop is done:

```go
for _, name := range names {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close() // ❌ 10 000 names → 10 000 open files
	process(f)
}
```

Extract the body into a function so each `defer` runs per iteration:

```go
for _, name := range names {
	if err := processFile(name); err != nil { // processFile opens, defers Close, processes
		return err
	}
}
```

### `os.Exit` skips deferred calls

`os.Exit(code)` (and `log.Fatal`, which calls it) ends the program **immediately**: no deferred calls run, no unwinding. Keep `main` small: `if err := run(); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }` so that `run`'s defers execute first.

## 11. Testing code that panics

A test can assert that something panics with a recovering helper:

```go
func mustPanic(t *testing.T, f func()) (r any) {
	t.Helper()
	defer func() {
		r = recover()
		if r == nil {
			t.Errorf("did not panic")
		}
	}()
	f()
	return nil
}
```

The tests of this lesson's exercises use helpers like this, so a missing `recover` shows up as a failed test instead of a crashed test binary.

## Common mistakes

| Mistake | Fix |
|---|---|
| `recover()` in normal code or in a helper called by the deferred func | Call it **directly** in the deferred function |
| Recovering without a named result: the caller gets `0, nil` | `func F() (v int, err error)` and set `err` in the deferred closure |
| `recover()` that swallows every panic silently | Convert to an error, or re-panic values you don't expect |
| `panic` for expected failures (bad user input) | Return an `error` |
| `panic`/`recover` as control flow (breaking out of loops) | Use `return`, `break`, error values |
| Letting a panic escape a library's exported API | Recover at the boundary and return an `error` |
| `defer f.Close()` before checking the `Open` error | Check `err` first, then defer |
| `defer timed("x")` (missing `()`) | `defer timed("x")()` |
| `defer` inside a long loop | Move the body into its own function |
| Ignoring `Close` errors on written files | Capture into the named `err` if it is still `nil` |
| Expecting defers to run after `os.Exit` / `log.Fatal` | Return from `run()` first, exit in `main` |
| `panic(nil)` | Panic with a descriptive value (Go 1.21+ turns it into `*runtime.PanicNilError`) |

## Key takeaways

- A panic stops the function, runs deferred calls while **unwinding** the stack, and crashes the program (status 2) if nothing recovers.
- Run-time panics (nil map write, nil dereference, out-of-range index, failed assertion, integer division by zero) implement `runtime.Error`.
- `recover()` works only **directly inside a deferred function**; it returns the panic value, or `nil`.
- Convert panics into errors at API boundaries with a **named `err` result**; wrap `error` values with `%w`.
- Recover only what you expect; **re-panic** the rest.
- Since Go 1.21, `panic(nil)` panics with `*runtime.PanicNilError`, so `recover() != nil` is reliable.
- Panic for bugs and impossible states; return errors for everything a caller can expect. `MustX` helpers are for inputs known to be valid.
- Defer patterns: cleanup/unlock, `defer timed("x")()`, `defer un(trace("x"))`, capture `Close` errors, extract loop bodies. `os.Exit` skips defers.
