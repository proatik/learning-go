# context

## 1. Why context exists

A request to a server starts work in many places: handlers, database calls, HTTP calls to other services, goroutines. When the client hangs up or a deadline passes, **all** of that work should stop. Lesson 2's done channel does this for one producer; `context.Context` is the standard version that works **across API boundaries**.

| A `Context` carries… | Through… |
|---|---|
| a **cancellation signal**: "stop, nobody needs the result" | `Done()` channel + `Err()` |
| a **deadline**: "stop at 12:00:05" | `Deadline()` |
| **request-scoped values**: request ID, auth token | `Value(key)` |

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

- Contexts form a **tree**. You never cancel a context you were given: you **derive** a child (`WithCancel`, `WithTimeout`, …) and cancel that. Canceling a parent cancels every descendant; canceling a child never touches the parent.
- A `Context` is immutable and **safe to share between goroutines**.
- Every library that does I/O takes one: `http.NewRequestWithContext`, `db.QueryContext(ctx, …)`, `exec.CommandContext`. Your microservice code will pass `ctx` through almost every function.

## 2. The roots: `Background` and `TODO`

Every tree starts at an empty root that is never canceled, has no deadline and holds no values.

```go
fmt.Println(context.Background()) // context.Background
fmt.Println(context.TODO())       // context.TODO

_, ok := context.Background().Deadline()
fmt.Println(ok, context.Background().Done() == nil, context.Background().Err()) // false true <nil>
```

| Root | Use it… |
|---|---|
| `context.Background()` | at the **top**: `main`, `init`, tests, a server's startup code |
| `context.TODO()` | as a placeholder where a ctx **should** come from a caller but the code isn't plumbed yet |

- Both behave identically; `TODO` is a marker for humans and tools.
- **Never pass `nil`** as a context, even if a function seems to allow it. Pass `context.TODO()`.
- `Done()` of a root is a `nil` channel: receiving from it blocks forever, so in a `select` that case is simply never chosen (lesson 3).

## 3. `WithCancel`: stop on demand

```go
ctx, cancel := context.WithCancel(context.Background())
fmt.Println(ctx.Err()) // <nil>
cancel()
<-ctx.Done()           // closed now: receive returns immediately
fmt.Println(ctx.Err()) // context canceled
cancel()               // calling it again is a no-op
```

```go
parent, cancelParent := context.WithCancel(context.Background())
child, cancelChild := context.WithCancel(parent)
defer cancelChild()
cancelParent()
fmt.Println(parent.Err(), child.Err()) // context canceled context canceled
```

- `Done()` returns a channel that is **closed** on cancellation, so any number of goroutines can wait on it (a close wakes every receiver, lesson 2).
- **Always call `cancel`**, normally with `defer cancel()` right after creating the context, even if the work finished. Each derived context is registered with its parent (and may own a timer); `cancel` releases that. Forgetting it leaks the child until the parent ends.
- `go vet` catches the classic slip (the **lostcancel** check):

```go
ctx, _ := context.WithTimeout(context.Background(), time.Second)
// vet: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
```

## 4. `Done()` and `Err()`: why did it stop?

| `ctx.Err()` | Meaning |
|---|---|
| `nil` | still running (`Done()` not closed yet) |
| `context.Canceled` | someone called `cancel` (or a parent was canceled) |
| `context.DeadlineExceeded` | the deadline / timeout passed |

```go
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()
select {
case <-time.After(time.Second):
	fmt.Println("work done")
case <-ctx.Done():
	fmt.Println("gave up:", ctx.Err()) // gave up: context deadline exceeded
}
err := ctx.Err()
fmt.Println(errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled)) // true false
```

- Compare with `errors.Is`, because callers usually **wrap** it: `fmt.Errorf("load user: %w", ctx.Err())`.
- `DeadlineExceeded` also has a `Timeout() bool` method returning `true`, which network code checks.
- Once `Err()` is non-nil it never changes. The first reason wins.
- Use `Done()` in a `select` (to wait); use `Err()` for a quick "should I stop?" check, or to report why.

## 5. `WithTimeout` and `WithDeadline`

```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second) // relative: now + 2s
defer cancel()

ctx2, cancel2 := context.WithDeadline(parent, time.Now().Add(2*time.Second)) // absolute: the same deadline
defer cancel2()
```

- `WithTimeout(p, d)` is exactly `WithDeadline(p, time.Now().Add(d))`.
- `ctx.Deadline()` returns the deadline and `true` (or the zero time and `false` if there is none). Use it to decide whether there is time left for something (`time.Until(deadline)`).
- **A child can only shorten** the deadline, never extend it:

```go
short, c1 := context.WithTimeout(context.Background(), time.Second)
defer c1()
long, c2 := context.WithTimeout(short, time.Hour) // asks for 1h…
defer c2()
d1, _ := short.Deadline()
d2, _ := long.Deadline()
fmt.Println(d1.Equal(d2)) // true: …but still ends after 1s, with its parent
```

- `defer cancel()` still matters with a timeout: if the work finishes in 10ms, `cancel` stops the timer now instead of in 2s.
- A context never stops your code by force. It only closes `Done()`; **your code must notice** and return.

## 6. Checking ctx in loops

Long-running work must check the context **between steps** and at **every blocking operation**.

```go
// Process handles jobs from in until in is closed or ctx is done.
func Process(ctx context.Context, in <-chan Job) (int, error) {
	n := 0
	for {
		select {
		case <-ctx.Done():
			return n, ctx.Err()
		case job, ok := <-in:
			if !ok {
				return n, nil
			}
			handle(job)
			n++
		}
	}
}
```

- Every send or receive that might block goes in a `select` with a `case <-ctx.Done():`. A bare `out <- v` can block forever after the receiver has gone away: a **goroutine leak**.
- For CPU-only loops with no channel to wait on, poll cheaply: `if err := ctx.Err(); err != nil { return err }`.
- ⚠️ `select` picks **randomly** among ready cases (lesson 3). With a canceled ctx **and** jobs waiting, the loop above keeps processing a random number of extra jobs (run 5 times on a canceled ctx with 100 buffered jobs, it returned 4, 1, 1, 2, 1). If cancellation must win, check `ctx.Err()` first at the top of each iteration.

## 7. Cancellation causes

`Canceled` says *that* it stopped, not *why*. Since Go 1.20/1.21 you can attach a **cause**.

```go
var ErrShutdown = errors.New("server shutting down")

ctx, cancel := context.WithCancelCause(context.Background())
cancel(ErrShutdown)
cancel(errors.New("second")) // ignored: the first cancellation wins
fmt.Println(ctx.Err(), "|", context.Cause(ctx)) // context canceled | server shutting down

child, cancelChild := context.WithCancel(ctx)
defer cancelChild()
fmt.Println(context.Cause(child)) // server shutting down (inherited)
```

```go
var ErrSlowDB = errors.New("database too slow")

ctx, cancel := context.WithTimeoutCause(context.Background(), 10*time.Millisecond, ErrSlowDB)
defer cancel()
<-ctx.Done()
fmt.Println(ctx.Err(), "|", context.Cause(ctx)) // context deadline exceeded | database too slow
```

| Function | Cause when it fires |
|---|---|
| `WithCancelCause(p)` → `cancel(err)` | `err` (`cancel(nil)` gives `context.Canceled`) |
| `WithTimeoutCause(p, d, err)` / `WithDeadlineCause(p, t, err)` | `err` when the **deadline** fires. Calling the returned `cancel` gives plain `context.Canceled` |
| `WithCancel` / `WithTimeout` (no cause) | same as `ctx.Err()` |

- `context.Cause(ctx)` returns `nil` while the context is live.
- `ctx.Err()` is **unchanged** (`Canceled` / `DeadlineExceeded`), so existing `errors.Is` checks keep working. Only `Cause` reveals the reason.
- Return `context.Cause(ctx)` instead of `ctx.Err()` when the caller benefits from knowing *why*.

## 8. `WithValue`: request-scoped data

```go
// ctxKey is unexported, so no other package can create the same key.
type ctxKey struct{}

// WithRequestID returns a copy of ctx carrying id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// RequestID returns the request ID stored in ctx, if any.
func RequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKey{}).(string)
	return id, ok
}

ctx := WithRequestID(context.Background(), "req-42")
fmt.Println(RequestID(ctx)) // req-42 true
id, ok := RequestID(context.Background())
fmt.Printf("%q %v\n", id, ok) // "" false
```

- `Value` compares keys with `==` **including their type**. A `string` key `"user"` from your package and `"user"` from a library collide; `ctxKey{}` can't collide with anything, because only your package can name the type.
- Hide the key behind **two accessor functions** (`WithX` / `X`). Callers never touch the key, and the comma-ok assertion means a missing value is just `"", false`, never a panic.
- `Value` walks **up** the chain, so an inner `WithValue` with the same key shadows the outer one. Lookups are O(depth): fine for a handful of values.

| Put in a context value | Keep out (pass as a parameter) |
|---|---|
| request ID, trace ID | database handles, loggers you could inject |
| the authenticated user for this request | optional function parameters, config, flags |
| data that crosses API boundaries with the request | anything a function needs to do its job |

Rule of thumb: if a function **can't work** without it, it's a parameter.

## 9. `AfterFunc`: run something on cancellation

`context.AfterFunc(ctx, f)` (Go 1.21) runs `f` in its own goroutine once `ctx` is done. It's how you connect a context to an API that **doesn't** accept one, e.g. close a connection so a blocked read returns.

```go
ctx, cancel := context.WithCancel(context.Background())
done := make(chan struct{})
stop := context.AfterFunc(ctx, func() {
	fmt.Println("cleanup ran")
	close(done)
})
cancel()
<-done              // cleanup ran
fmt.Println(stop()) // false: f had already started
```

```go
stop := context.AfterFunc(ctx, func() { fmt.Println("never") })
fmt.Println(stop(), stop()) // true false: the first call un-registered f
cancel()                    // f does not run
```

| `stop()` returns | Meaning |
|---|---|
| `true` | you stopped `f` **before** it ran; it never will |
| `false` | `f` has already been started (ctx was canceled), or it was already stopped |

- If `ctx` is **already** done, `f` starts immediately, so `stop()` returns `false`.
- `stop()` does **not** wait for `f` to finish. If you need that, coordinate with a channel.
- Always call `stop` when the guarded operation finishes, or `f` may still run later. `stop()`'s result tells you whether cancellation got there first.

## 10. `WithoutCancel`: outlive the request

Sometimes work must **finish** even though the request is over: writing an audit log, flushing metrics. `context.WithoutCancel(ctx)` (Go 1.21) keeps the **values** but drops cancellation and the deadline.

```go
req, cancel := context.WithTimeout(WithRequestID(context.Background(), "req-42"), time.Hour)
bg := context.WithoutCancel(req)
cancel()
_, hasDeadline := bg.Deadline()
id, _ := RequestID(bg)
fmt.Println(req.Err(), bg.Err(), bg.Done() == nil, hasDeadline, id)
// context canceled <nil> true false req-42
```

- Its `Done()` is `nil`, so give that background work **its own** timeout: `context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)`.
- `context.Cause` on it returns `nil`.

## 11. The rules for APIs

```go
func FetchUser(ctx context.Context, id int) (*User, error) // ✅ ctx first, named ctx

type Client struct {
	ctx context.Context // ❌ don't store contexts in structs
}
```

- `ctx context.Context` is the **first** parameter, named **`ctx`**.
- **Don't store a context in a struct.** A struct outlives one call, a context belongs to one call. Pass it to each method: `func (c *Client) Get(ctx context.Context, url string)`.
- Pass the ctx you received **down** to everything you call. Creating a fresh `context.Background()` deep in the stack cuts the cancellation chain.
- A function that **starts** a goroutine should make it stop when ctx is done, and the function that **creates** a cancelable ctx is the one that calls `cancel`.

## 12. Contexts in tests

```go
func TestSleepTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		start := time.Now()
		err := Sleep(ctx, time.Hour)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Sleep(ctx, 1h) error = %v, want %v", err, context.DeadlineExceeded)
		}
		if got := time.Since(start); got != 2*time.Second {
			t.Errorf("Sleep(ctx, 1h) took %v, want 2s", got)
		}
	})
}
```

- `t.Context()` (Go 1.24) is canceled **just before** the test's `Cleanup` functions run. Use it instead of `context.Background()`, so goroutines that watch it shut down when the test ends.
- Inside a `synctest` bubble, context deadlines use the **fake clock** (lesson 3): the 2s timeout above fires instantly and `time.Since` is exactly `2s`.

## Common mistakes

| Mistake | Fix |
|---|---|
| `ctx, _ := context.WithTimeout(…)` | Keep `cancel` and `defer cancel()`; `go vet` reports the leak |
| Calling `cancel` only on the error path | `defer cancel()` right after creating the context |
| `out <- v` or `<-in` without a `ctx.Done()` case | Put every blocking operation in a `select` with `case <-ctx.Done():` |
| Expecting a context to stop code by itself | Your code must check `Done()` / `Err()` and return |
| `err == context.DeadlineExceeded` on a wrapped error | `errors.Is(err, context.DeadlineExceeded)` |
| `context.WithValue(ctx, "user", u)` | An unexported key type: `type ctxKey struct{}` + accessor functions |
| `ctx.Value(k).(string)` without comma-ok | `v, ok := ctx.Value(k).(string)`: a missing value is `nil`, and the one-value form panics |
| Passing dependencies (DB, logger, config) in ctx values | Explicit parameters or struct fields |
| Storing a `ctx` in a struct field | Pass `ctx` as the first parameter of each method |
| `context.Background()` deep inside a request | Pass the caller's ctx down; use `WithoutCancel` if work must outlive it |
| Passing `nil` as a context | `context.TODO()` |
| Reporting only `Canceled` when you know why | `WithCancelCause` / `WithTimeoutCause` + `context.Cause` |
| Forgetting `stop()` after `AfterFunc` | Call it when the operation ends; its result says whether `f` already ran |

## Key takeaways

- A `Context` carries cancellation, a deadline and request-scoped values across API boundaries. Contexts form a tree: canceling a parent cancels all its children.
- Start from `context.Background()` (or `TODO()` as a placeholder); derive with `WithCancel`, `WithTimeout`, `WithDeadline`, and **always `defer cancel()`**.
- `Done()` is closed on cancellation; `Err()` says why: `Canceled` or `DeadlineExceeded` (check with `errors.Is`). A child can shorten its parent's deadline, never extend it.
- Check ctx in every loop and at every blocking send/receive with `select`; `select` is random, so check `ctx.Err()` first when cancellation must win.
- `WithCancelCause`, `WithTimeoutCause` and `context.Cause` record **why** a context stopped, without changing `Err()`.
- `WithValue` only for request-scoped data, with an unexported key type behind accessor functions.
- `AfterFunc` runs cleanup on cancellation (`stop()` reports whether it got there first); `WithoutCancel` keeps values but drops cancellation.
- `ctx` is the first parameter, never stored in a struct, never `nil`. In tests use `t.Context()`, and `synctest` for anything with deadlines.
