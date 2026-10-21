# Errors

## 1. Errors are values

`error` is a predeclared **interface** with one method:

```go
type error interface {
	Error() string
}
```

- Any type with an `Error() string` method is an error. Nothing else is special: an error is an ordinary value you can store, pass, compare, wrap and return.
- `nil` means "no error". A non-nil `error` means the other results are not reliable, unless the docs say otherwise (`io.Reader` can return `n > 0` *and* `io.EOF`, see L3).
- Go has no exceptions for ordinary failures. `panic` is reserved for bugs and truly unrecoverable situations (next lesson).

```go
n, err := strconv.Atoi("abc")
if err != nil {
	fmt.Println(err) // strconv.Atoi: parsing "abc": invalid syntax
	return
}
fmt.Println(n)
```

"Errors are values" (Rob Pike, Go blog) means you **program** with errors like any other value. You don't just check them. For example, a type can remember the first error and skip the remaining work:

```go
type writer struct {
	w   io.Writer
	err error // the first error; once set, every later write is skipped
}

func (w *writer) write(s string) {
	if w.err != nil {
		return
	}
	_, w.err = io.WriteString(w.w, s)
}

// w.write("a"); w.write("b"); w.write("c")
// if w.err != nil { ... } // one check instead of three
```

## 2. The idioms

| Idiom | Looks like |
|---|---|
| `error` is the **last** result | `func Open(name string) (*File, error)` |
| Check **immediately**, before using the other results | `f, err := Open(n); if err != nil { return err }` |
| On error, return **zero values** for the other results | `return Config{}, err`, `return 0, err`, `return nil, err` |
| Keep the happy path at the left margin (**guard clauses**) | `if err != nil { return ... }`, never `else` after it |
| **Handle once**: log it *or* return it, never both | Logging then returning prints the same failure twice up the stack |
| Messages are **lowercase**, with **no trailing punctuation** | `errors.New("empty name")`, not `"Empty name."` |
| Add **context** that the caller doesn't already have | `fmt.Errorf("open config %s: %w", path, err)` |
| Don't start a message with "error" or "failed to" | Every error is a failure. Say what you were doing: `"read header: ..."` |

Lowercase messages matter because errors get chained. Each layer adds a prefix:

```text
load user 42: query db: dial tcp 10.0.0.5:5432: connection refused
```

A capital letter or a `.` in the middle of that chain looks wrong.

Guard clauses keep the code flat:

```go
// ✅ each failure exits early; the success path reads straight down
func load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("load %s: %w", path, err)
	}
	cfg, err := parse(data)
	if err != nil {
		return Config{}, fmt.Errorf("load %s: %w", path, err)
	}
	return cfg, nil
}
```

The short form `if err := f(); err != nil { ... }` scopes `err` to the `if` when `f` returns only an error.

## 3. Creating errors: `errors.New` and `fmt.Errorf`

```go
err1 := errors.New("empty name")               // fixed message
err2 := fmt.Errorf("age %d out of range", age) // formatted message
```

- Every call to `errors.New` returns a **distinct** value, even with the same text:

```go
errors.New("x") == errors.New("x") // false: two different *errorString pointers
```

- That is why errors are compared by **identity**, never by text. Sentinels (next section) give you a single value to compare against.

## 4. Sentinel errors

A **sentinel** is a package-level error variable that callers compare against.

```go
package store

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate key")
)

func (s *Store) Get(key string) (string, error) {
	v, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}
```

```go
v, err := s.Get("k")
if errors.Is(err, store.ErrNotFound) { // §7: works even if err was wrapped
	v = "default"
} else if err != nil {
	return err
}
```

- Name them `ErrXxx`, or `errXxx` if unexported. Known std examples: `io.EOF`, `fs.ErrNotExist`, `sql.ErrNoRows`, `strconv.ErrSyntax`.
- **Export a sentinel only if callers need to react to it differently.** An exported sentinel is part of your API forever: you can't change or remove it without breaking callers.
- Sentinels carry no data. If the caller needs details (which field? how much money?), use a custom type.
- `io.EOF` is the classic exception: it is returned **unwrapped**, and callers compare it with `==`. This is because it signals a normal end, not a failure.

## 5. Custom error types

Use a struct when the caller needs **data** about the failure.

```go
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Reason
}

func Validate(u User) error {
	if u.Name == "" {
		return &ValidationError{Field: "name", Reason: "must not be empty"}
	}
	return nil
}
```

- Convention: name it `XxxError`, use a **pointer receiver**, and return `&XxxError{...}`.
  - With a pointer, every check matches the same way (`errors.As` with a `*XxxError` target). Mixing `XxxError` and `*XxxError` values makes `errors.As` miss.
  - It also avoids copying big structs.
- Use the **error interface** as the result type (`error`), not `*ValidationError`. §10 explains why.
- `Error()` should only *format*. It must not fail, do I/O or change state.

Std examples: `*fs.PathError{Op, Path, Err}`, `*strconv.NumError{Func, Num, Err}`, `*json.SyntaxError{Offset}`.

```go
_, err := os.Open("missing.txt")
fmt.Println(err) // open missing.txt: no such file or directory  (text is OS-specific)
```

## 6. Wrapping: adding context with `%w`

`fmt.Errorf` with the `%w` verb returns a new error that **wraps** its operand. The text includes the original message, and the original value stays reachable through a method `Unwrap() error`.

```go
func loadConfig(name string) (Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Config{}, fmt.Errorf("load %s: %w", name, err)
	}
	// ...
}
```

Errors form a **chain**: each wrapper points to the error inside it.

```go
base := errors.New("not found")
a := fmt.Errorf("find user 7: %w", base)
b := fmt.Errorf("get profile: %w", a)

fmt.Println(b)                           // get profile: find user 7: not found
fmt.Println(errors.Unwrap(b) == a)       // true
fmt.Println(errors.Unwrap(a) == base)    // true
fmt.Println(errors.Unwrap(base))         // <nil>: end of the chain
```

- `errors.Unwrap(err)` calls `err.Unwrap()` once and returns `nil` if there is no such method. In practice you rarely call it: `errors.Is` and `errors.As` walk the whole chain for you.
- Your own type joins the chain by adding an `Unwrap` method:

```go
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string { return "query " + e.Query + ": " + e.Err.Error() }
func (e *QueryError) Unwrap() error { return e.Err }
```

### `%w` vs `%v`

| | `fmt.Errorf("load: %w", err)` | `fmt.Errorf("load: %v", err)` |
|---|---|---|
| Message | `load: not found` | `load: not found` (identical) |
| `errors.Is(result, ErrNotFound)` | `true` | `false`: the chain is cut |
| What callers can depend on | the wrapped error **becomes part of your API** | only your own text |

**When NOT to wrap:** wrapping exposes the inner error to callers. If your `UserService` wraps a `*pq.Error` from its Postgres driver, callers start writing `errors.As(err, &pqErr)`. Then you can never switch databases without breaking them.

- Wrap (`%w`) errors you **document**, like your own sentinels or `fs.ErrNotExist` from a file API.
- Use `%v`, or return your own sentinel, at an abstraction boundary where the cause is an implementation detail.

```go
// boundary: hide the driver, expose our own sentinel
if err != nil {
	return User{}, fmt.Errorf("get user %d: %w (%v)", id, ErrUnavailable, err)
}
```

## 7. `errors.Is`: "is this error (or anything it wraps) X?"

```go
err := fmt.Errorf("load config: %w", fs.ErrNotExist)

err == fs.ErrNotExist          // false: err is the wrapper
errors.Is(err, fs.ErrNotExist) // true: found inside the chain
```

`errors.Is(err, target)` walks the chain. For each error `e` it reports `true` if:
1. `e == target`, or
2. `e` has a method `Is(error) bool` and `e.Is(target)` returns `true`.

A **custom `Is` method** lets an error match a whole *category*:

```go
var ErrClient = errors.New("client error")

type HTTPError struct{ Code int }

func (e *HTTPError) Error() string { return fmt.Sprintf("http %d", e.Code) }

func (e *HTTPError) Is(target error) bool {
	return target == ErrClient && e.Code >= 400 && e.Code < 500
}

err := fmt.Errorf("fetch: %w", &HTTPError{Code: 404})
errors.Is(err, ErrClient) // true
```

- `Is` should only compare. Don't unwrap inside it: `errors.Is` does the walking.
- Rule of thumb: use `errors.Is` for **sentinels**, `errors.As`/`errors.AsType` for **types**.

## 8. `errors.As` and `errors.AsType`: "give me the error of type T"

`errors.As(err, &target)` finds the first error in the chain that can be assigned to `target`, **stores it there**, and returns `true`.

```go
var ve *ValidationError // the target: a nil pointer of the type you want
if errors.As(err, &ve) { // pass a POINTER to the target
	fmt.Println("bad field:", ve.Field)
}
```

- The second argument must be a **non-nil pointer** to an error type or interface type. Passing `ve` instead of `&ve` panics at run time, and `go vet` (`errorsas` check) catches it.
- The target type must match the way the error was created. A `*ValidationError` target finds `&ValidationError{...}`, but not a `ValidationError` value.

**Go 1.26+: `errors.AsType`** is the generic, type-safe form. `[T]` is a *type argument*: you pass a type in square brackets, and the function returns a value of that type.

```go
if ve, ok := errors.AsType[*ValidationError](err); ok {
	fmt.Println("bad field:", ve.Field) // ve is *ValidationError, scoped to the if
}
```

| | `errors.As(err, &ve)` | `errors.AsType[*ValidationError](err)` |
|---|---|---|
| Target variable | declare it first | none: returned as a result |
| Scope | the variable leaks into the outer scope | scoped to the `if` |
| Mistakes | wrong target → **run-time panic** (vet may catch it) | wrong type → **compile error** |
| Speed | uses reflection | faster, no reflection |

**Prefer `errors.AsType` in new code.** You will still read `errors.As` everywhere in existing code.

## 9. Multiple errors: `errors.Join` and several `%w`

```go
err := errors.Join(ErrNameRequired, nil, ErrEmailInvalid) // nil values are dropped
fmt.Println(err)
// name is required
// email is invalid

errors.Is(err, ErrEmailInvalid) // true: Is/As search every branch
errors.Join(nil, nil) == nil    // true: Join of only nils is nil
```

- `errors.Join` separates the messages with a **newline**.
- `fmt.Errorf` accepts **several `%w`** verbs: `fmt.Errorf("sync: %w; cleanup: %w", errA, errB)`.
- Both produce an error with `Unwrap() []error` (a **tree**, not a chain). `errors.Is`/`As`/`AsType` search it depth-first.
- `errors.Unwrap` only knows `Unwrap() error`, so it returns **`nil`** for joined errors.

The **error-accumulating validator** is the main use case. Report *every* problem at once, not one per submit:

```go
func ValidateSignup(f Form) error {
	var errs []error
	if f.Name == "" {
		errs = append(errs, ErrNameRequired)
	}
	if !strings.Contains(f.Email, "@") {
		errs = append(errs, fmt.Errorf("email %q: %w", f.Email, ErrEmailInvalid))
	}
	return errors.Join(errs...) // nil when errs is empty
}
```

## 10. The typed-nil trap, revisited

From L2: an interface is `nil` only if **both** its type and value are nil. A nil `*MyErr` stored in an `error` is a **non-nil error**.

```go
func find(ok bool) error {
	var e *NotFoundError // nil pointer
	if !ok {
		e = &NotFoundError{Name: "x"}
	}
	return e // ❌ when ok, this returns error(type=*NotFoundError, value=nil)
}

err := find(true)
fmt.Println(err == nil) // false! the caller thinks it failed
```

Fix: return the literal `nil` on success, and **declare results as `error`**, never as `*MyErr`:

```go
func find(ok bool) error {
	if !ok {
		return &NotFoundError{Name: "x"}
	}
	return nil // ✅ a truly nil interface
}
```

The same bug hides in helpers like `func validate() *ValidationError` whose result is returned through an `error`. Keep concrete error types out of signatures.

## 11. Behaviour, not identity: checking for a method

Sometimes the caller only cares whether an error *can* be retried, not where it came from. Ask for a **behaviour** through an interface:

```go
type temporary interface {
	Temporary() bool
}

func isTemporary(err error) bool {
	var t temporary
	return errors.As(err, &t) && t.Temporary() // As accepts a pointer to an interface type
}
```

With `AsType` the type argument must implement `error`, so embed it in the interface:

```go
t, ok := errors.AsType[interface {
	error
	Temporary() bool
}](err)
retry := ok && t.Temporary()
```

- This decouples packages. A network package's error can say "retry me" without the caller importing that package's types.
- The same idea appears in the standard library: `net.Error` has a `Timeout() bool` method.

## 12. Error handling patterns

**Constructors and validation return `(T, error)`:**

```go
func NewAccount(owner string, balance int) (*Account, error) {
	if owner == "" {
		return nil, errors.New("account: empty owner")
	}
	if balance < 0 {
		return nil, fmt.Errorf("account: negative balance %d", balance)
	}
	return &Account{owner: owner, balance: balance}, nil
}
```

**Wrap at layer boundaries, decide at the top.** Lower layers add context, and one place decides what to do:

```go
// repository: returns the fact
return User{}, fmt.Errorf("find user %d: %w", id, ErrNotFound)

// service: adds what it was doing
return User{}, fmt.Errorf("get profile: %w", err)

// handler: maps errors to a response, ONCE
func status(err error) int {
	if err == nil {
		return 200
	}
	if errors.Is(err, ErrNotFound) {
		return 404
	}
	if _, ok := errors.AsType[*ValidationError](err); ok {
		return 400
	}
	return 500 // unknown: log the details, show the user a generic message
}
```

**Never inspect `err.Error()` text:**

```go
if strings.Contains(err.Error(), "not found") { ... } // ❌ breaks when anyone rewords a message
if errors.Is(err, ErrNotFound) { ... }               // ✅
```

The message is for **humans**. Programs use sentinels, types and behaviours.

**Tooling helps:**
- `go vet` flags `errors.As(err, ve)` without `&`.
- **Go 1.27:** `go vet` flags `fmt.Errorf("...: %w", p)` when `p` is a `*E` but the *value type* `E` is what implements `error` (`%w wants operand of error type E, not pointer type *E (defeats errors.Is)`). The conventional pointer-receiver `*E` is fine.

## 13. Choosing the right tool

| You want callers to… | Use |
|---|---|
| just know it failed | `errors.New` / `fmt.Errorf` with `%v` |
| test for one specific condition | exported sentinel `ErrX` + `errors.Is` |
| read details (field, amount, position) | `*XxxError` struct + `errors.AsType` |
| match a category of errors | custom `Is(target error) bool` method |
| ask "can I retry?" | behaviour interface (`Temporary() bool`) + `errors.As` |
| see every problem at once | `errors.Join` |
| keep the original cause | wrap with `%w` (and document it) |

## Common mistakes

| Mistake | Fix |
|---|---|
| `if err.Error() == "not found"` / `strings.Contains(err.Error(), …)` | `errors.Is(err, ErrNotFound)` |
| `err == ErrNotFound` after the error was wrapped | `errors.Is(err, ErrNotFound)` |
| Using the result before checking `err` | Check `err` first; the other results may be garbage |
| `log.Println(err); return err` | Handle once: log at the top **or** return with context |
| `errors.New("Failed to open file.")` | `errors.New("open file")`: lowercase, no punctuation, no "failed to" |
| `fmt.Errorf("load: %v", err)` when callers need `errors.Is` | `%w` keeps the chain |
| Wrapping a third-party driver error at your API boundary | `%v` or your own sentinel: don't leak implementation details |
| `errors.As(err, ve)` | `errors.As(err, &ve)`, or better `errors.AsType[*VE](err)` |
| Creating `ValidationError{}` values but checking `*ValidationError` | Pointer receiver, always return `&ValidationError{...}` |
| `var e *MyErr; …; return e` from a func returning `error` | `return nil` explicitly on success |
| Declaring results as `*MyErr` instead of `error` | Always return the `error` interface |
| Returning `nil, nil` or a half-built value with an error | On error: zero values + the error |
| `errors.Unwrap(joined)` expecting the parts | Use `errors.Is`/`AsType`; `Unwrap` returns `nil` for `Unwrap() []error` |
| `panic` for invalid input | Return an error; `panic` is for bugs (next lesson) |

## Key takeaways

- `error` is an interface (`Error() string`). Errors are values: return them, store them, program with them.
- Return `error` last, check it immediately, return zero values with it, and handle it once.
- Messages: lowercase, no punctuation, add context, so the chain reads `a: b: c`.
- **Sentinel** (`var ErrX = errors.New(...)`) for conditions, **`*XxxError` type** for details, **behaviour interface** for "retryable?". Export only what callers must act on.
- `fmt.Errorf("context: %w", err)` wraps. `%v` deliberately cuts the chain at abstraction boundaries.
- `errors.Is` = sentinel through the chain (and custom `Is`). `errors.As(err, &t)` / **`errors.AsType[*T](err)`** (Go 1.26+, preferred) = typed error through the chain.
- `errors.Join` and multiple `%w` build trees. `Is`/`As` search them all, `errors.Unwrap` doesn't.
- Never return a typed nil pointer as an `error`, and never compare `err.Error()` strings.
