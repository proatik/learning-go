# Module 3 Challenge: Designing with Types

No new concepts here. This page is a one-page recap of Module 3; the ten exercises each combine methods, interfaces, the standard interfaces and errors to design small but real APIs.

## 1. Methods and receivers

```go
type Counter struct{ n int }

func (c *Counter) Inc()         { c.n++ }      // pointer receiver: mutates
func (c Counter) Value() int    { return c.n } // value receiver: reads a copy
func NewCounter(start int) *Counter { return &Counter{n: start} } // constructor
```

| Situation | Receiver |
|---|---|
| Method modifies the receiver | `*T` |
| Struct is large, or holds a map/slice you mean to reassign | `*T` |
| Type has identity (account, cache, logger) or contains a lock | `*T` |
| Small immutable value (`Point`, `Money`, `Level`) | `T` |
| Any method on the type uses `*T` | `*T` for **all** of them (consistency) |

- `v.M()` on an **addressable** `v` auto-takes `&v` for pointer methods; map elements and return values are not addressable.
- Embedding promotes fields **and** methods: `type Admin struct{ User }` gives `a.Name()`. The outer type can shadow a promoted method.
- Method value `f := c.Inc` binds the receiver now; method expression `(*Counter).Inc` takes it as the first argument.
- Encapsulate with lowercase fields + exported methods; the constructor sets the invariants (`make` the map, check arguments).

## 2. Interfaces

| Rule | Detail |
|---|---|
| Satisfaction is implicit | a type implements `I` by having its methods — no `implements` keyword |
| Method set of `T` | value-receiver methods only |
| Method set of `*T` | value **and** pointer-receiver methods |
| Consequence | if `Write` has a `*T` receiver, store `&t` in the interface, not `t` |
| Compile-time check | `var _ Shape = (*Circle)(nil)` |
| Interface value | a pair *(dynamic type, value)*; it is `nil` only if **both** are nil |
| `any` | `interface{}` — holds anything; get the value back with an assertion |
| Embedding | `type ReadWriter interface{ Reader; Writer }` |

```go
s, ok := v.(fmt.Stringer)  // comma-ok assertion: never panics
switch x := v.(type) {     // type switch: x has the case's type
case nil:                  // the interface itself is nil
case int, int64:           // several types: x stays the interface type
case *Circle:
	_ = x.R
case error:                // interface cases match any implementer
default:
}
```

## 3. Standard interfaces

| Interface | Method | Contract |
|---|---|---|
| `fmt.Stringer` | `String() string` | used by `%v`/`%s`/`Println`; don't call `fmt.Sprint(t)` inside it (infinite recursion) |
| `error` | `Error() string` | lowercase, no trailing punctuation, `pkg: context: cause` |
| `io.Reader` | `Read(p []byte) (n int, err error)` | handle `p[:n]` **before** looking at `err`; `io.EOF` = clean end; may return `n > 0` together with `io.EOF` |
| `io.Writer` | `Write(p []byte) (n int, err error)` | must return a non-nil error if `n < len(p)`; must not modify or keep `p` |
| `sort.Interface` | `Len`, `Less(i, j)`, `Swap(i, j)` | `sort.Sort(x)`, `sort.Stable(x)` |

```go
io.Copy(dst, src)                 // reads until EOF; returns bytes copied, EOF is not an error
io.MultiWriter(w1, w2)            // fan-out writer
io.LimitReader(r, n)              // stops after n bytes
sc := bufio.NewScanner(r)         // for sc.Scan() { sc.Text() }; then check sc.Err()
bw := bufio.NewWriter(w); defer bw.Flush()
fmt.Fprintf(w, "%s=%d\n", k, v)   // print to any io.Writer
```

## 4. Error toolkit

| Tool | Use it to |
|---|---|
| `errors.New("store: not found")` | declare a **sentinel**: `var ErrNotFound = ...` |
| `fmt.Errorf("get %q: %w", k, err)` | add context and keep the cause in the chain (`%v` would cut it) |
| `errors.Is(err, ErrNotFound)` | test for a sentinel anywhere in the chain |
| `var fe *FieldError; errors.As(err, &fe)` | extract a custom error type from the chain |
| `fe, ok := errors.AsType[*FieldError](err)` | the same, typed, no variable needed (Go 1.26+) |
| `errors.Join(e1, e2, ...)` | combine errors; `nil`s are dropped, all `nil` → `nil`; `Is`/`As` search every branch |
| `Unwrap() error` / `Unwrap() []error` | make your type part of the chain / tree |
| `Is(target error) bool` method | custom matching, e.g. any `*TransitionError` matches `ErrInvalidTransition` |

```go
type FieldError struct{ Field, Problem string }

func (e *FieldError) Error() string { return e.Field + ": " + e.Problem }

func (e *FieldError) Is(target error) bool { return target == ErrInvalid } // custom Is

func check(name string) error {
	var fe *FieldError           // typed-nil trap:
	if name == "" {
		fe = &FieldError{"name", "required"}
	}
	return fe                    // BUG: never == nil — return a literal nil on success
}
```

## 5. panic vs error, and defer

| Situation | Use |
|---|---|
| Expected failure: bad input, missing key, I/O, business rule | return an `error` |
| Programmer bug / broken invariant that can't continue | `panic` |
| `Must`-style helper for constants at start-up (`regexp.MustCompile`) | panic on error |
| Calling code you don't trust (plugins, callbacks) | `recover` at the boundary and turn it into an error |
| Across an API boundary | never let a panic escape — convert it |

```go
defer bw.Flush()                     // cleanup runs on every return path, LIFO order
defer func() {                       // recover only works directly inside a deferred func
	if r := recover(); r != nil {
		err = fmt.Errorf("callback panicked: %v", r) // named result → caller sees it
	}
}()
defer func() {                       // rollback on failure
	if err != nil {
		undo()
	}
}()
```

- Deferred calls' **arguments** are evaluated at the `defer` statement; the call runs later.
- A deferred closure can read and change **named** results.

## Designing with interfaces — checklist

- **Small interfaces.** One to three methods (`io.Reader`, `Store`, `Logger`). Compose bigger ones by embedding.
- **Define the interface where it is used**, not next to the implementation. The consumer states what it needs.
- **Accept interfaces, return concrete types.** `NewCached(s Store) *CachedStore` — callers get every method; the parameter stays flexible.
- **Decorators** wrap an interface and implement it again: `CachedStore`, `LevelFilter`, `upperReader`. They stack.
- **Function adapters** (`type HandlerFunc func(...)` with a method) turn a plain func into an interface value.
- **Errors are part of the API.** Export sentinels and types, document which one each method returns, wrap with `%w`, and let callers use `errors.Is` / `errors.As` — never compare `err.Error()` strings.
- **Constructors enforce invariants** — create maps, reject bad arguments — so methods never check for a half-built value.

## Common mistakes

| Mistake | Fix |
|---|---|
| `var s Shape = Circle{}` but `Area` has a `*Circle` receiver | store `&Circle{}`, or use value receivers |
| Returning a typed nil pointer as `error` | `return nil` explicitly |
| `fmt.Errorf("...: %v", err)` | `%w` keeps the cause for `errors.Is` / `As` |
| `err == ErrNotFound` on a wrapped error | `errors.Is(err, ErrNotFound)` |
| Ignoring `n` when `Read` returns an error | process `p[:n]` first, then handle `err` |
| `String()` calling `fmt.Sprint(t)` on itself | format the fields, not the receiver |
| `recover()` in a helper called by the deferred func | call `recover()` directly in the deferred func |

## Key takeaways

- Pick the receiver by "does it mutate or have identity?"; keep all receivers of a type consistent.
- Interfaces are satisfied implicitly; the method set of `T` lacks pointer methods.
- An interface holding a typed nil pointer is not `nil`.
- Wrap with `%w`, inspect with `errors.Is` / `errors.As` / `errors.AsType`, combine with `errors.Join`.
- Return errors for expected failures; recover panics only at a boundary, and turn them into errors.
