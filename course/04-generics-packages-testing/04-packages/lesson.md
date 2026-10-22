# Packages & API Design

## 1. Package = directory, module = tree of packages

- A **package** is one directory of `.go` files that all start with the same `package name` clause.
- A **module** is a tree of packages with a `go.mod` at its root. The `module` line is the prefix of every import path inside it.
- **Import path** = module path + directory. **Package name** = the word in the `package` clause; by convention the last element of the path.

```text
learngo/                      go.mod: module learngo
├── shop/                     import "learngo/shop"        package shop
│   ├── cart.go
│   └── price.go              same directory → same package (shared unexported names)
└── shop/money/               import "learngo/shop/money"  package money
    └── money.go
```

```go
package shop

import "learngo/shop/money" // import by path…

func Total() money.Amount { // …use by package name
	return money.New(1250)
}
```

| Fact | Consequence |
|---|---|
| All files in a directory form **one** package | Split a big package into several files freely; they see each other's unexported names |
| A sub-directory is a **different** package | `shop` and `shop/money` share nothing unexported; `shop/money` is not "inside" `shop` |
| Only **exported** (Capitalised) names cross a package boundary | The package's API is exactly its exported names |
| Directories named `testdata`, or starting with `_` or `.`, are ignored by `./...` | That's why exercise solutions live in `_solution/` |

## 2. Package names

The package name is part of every call site: `bytes.Buffer`, `strings.Split`, `http.Get`. Pick it so that **name + identifier** reads well.

| Rule | ✅ | ❌ |
|---|---|---|
| Short, lowercase, one word | `money`, `slug`, `httptest` | `Money`, `money_utils`, `moneyUtils` |
| Don't repeat the package name in its identifiers (**stutter**) | `bytes.Buffer`, `bufio.Reader`, `ring.New` | `bytes.BytesBuffer`, `bufio.BufReader` |
| Name what it **provides**, not "where helpers live" | `slug`, `money`, `auth` | `util`, `common`, `misc`, `helpers` |
| The package's main type and constructor can lean on the name | `list.List` + `list.New()` | `list.NewList()` |
| Don't steal a good variable name | `bufio` | `buf` (now nobody can name a variable `buf`) |

- `util` / `common` packages grow forever and get imported by everything, which invites **import cycles** (§6). Move each helper next to the code that uses it, or give it a real name.
- Rename an import only on a clash: `import mrand "math/rand"`.

## 3. Keep the exported API small

Every exported name is a **promise**: once other code uses it, changing it breaks them. Start unexported; export only what callers need.

```go
package money

// Amount is an amount of money in cents. The zero value is 0.00.
type Amount struct {
	cents int64 // unexported: callers can't create an invalid Amount
}

// New returns an Amount of the given cents.
func New(cents int64) Amount { return Amount{cents: cents} }

// Cents returns the amount in cents.
func (a Amount) Cents() int64 { return a.cents }
```

From another package:

```text
money.Amount{cents: 5}  →  cannot refer to unexported field cents in struct literal of type money.Amount
a.cents                 →  a.cents undefined (type money.Amount has no field or method cents, but does have method Cents)
```

| Export | Keep unexported |
|---|---|
| Types and functions callers really use | Helpers, caches, intermediate types |
| Sentinel errors callers compare with `errors.Is` | Fields that must stay consistent (balances, counters) |
| Interfaces **consumers** of other packages need | "Just in case" knobs |

## 4. Doc comments and `go doc`

- A **package comment** sits right above the `package` clause and starts with `Package name`. Write it in **one** file (a long one often goes in `doc.go`).
- A **doc comment** on every exported name starts with that name.
- A paragraph starting with `Deprecated:` marks an old API (§10).

```go
// Package money represents amounts of money as integer cents.
package money

// Amount is an amount of money in cents.
type Amount struct{ cents int64 }
```

`go doc` (and pkg.go.dev) prints only the exported API:

```text
$ go doc ./money
package money // import "example.com/app/money"

Package money represents amounts of money as integer cents.

type Amount struct{ ... }
    func New(cents int64) Amount

$ go doc ./money Amount
package money // import "example.com/app/money"

type Amount struct {
	// Has unexported fields.
}
    Amount is an amount of money in cents.

func New(cents int64) Amount
func (a Amount) Cents() int64
```

- Constructors are listed **under** the type they return, so a reader finds them next to the type.
- Read your own `go doc` output: if it's long and confusing, the API is too big.

## 5. `internal/`: shareable inside, private outside

An import of a path containing the element `internal` is allowed only from code **inside the tree rooted at the parent of `internal`**.

```text
example.com/app/
├── shop/
│   ├── shop.go             ✅ may import example.com/app/shop/internal/tax
│   ├── checkout/           ✅ also inside the shop/ tree
│   └── internal/tax/
└── other/other.go          ❌ outside shop/
```

```text
package example.com/app/other
	other/other.go:3:8: use of internal package example.com/app/shop/internal/tax not allowed
```

- `internal` is enforced by the **go command**, not by convention: other modules can never import it.
- So you may **change an internal package freely**: its only users are in your own tree.
- Use it for helpers that several of your packages share but that are not part of your public API.

## 6. Import cycles

Package `a` imports `b` and `b` imports `a`: Go refuses to build.

```text
package example.com/app/a
	imports example.com/app/b from a.go
	imports example.com/app/a from b.go: import cycle not allowed
```

Why it's forbidden: packages are compiled and **initialised** in dependency order (§7). A cycle has no first package. It's also a design smell: two packages that need each other are really one.

| Fix | When |
|---|---|
| **Move an interface to the consumer** | `a` only needs a few methods of `b`'s type: declare `a.Notifier` and let `b`'s type satisfy it implicitly |
| **Extract a smaller package** both import | Both need shared types (`Event`, `ID`): move them into a third, lower package |
| **Merge** | The two packages are one concept split too early |

Dependency inversion with an interface:

```go
// ❌ account imports audit (to log)  AND  audit imports account (to read the ID)

// ✅ package event: the small, low-level package both depend on
package event

type Event struct {
	Account string
	Kind    string
	Amount  int
}

type Recorder interface{ Record(e Event) }

// package account: imports only event
type Account struct {
	id  string
	rec event.Recorder // any recorder: audit.Log, a test fake…
}

// package audit: imports only event; *Log satisfies event.Recorder
func (l *Log) Record(e event.Event) { l.lines = append(l.lines, e.Account+" "+e.Kind) }
```

The arrows now point one way: `account → event ← audit`. A `main` (or root) package imports both and **wires** them: `account.New("acc-1", &audit.Log{})`.

## 7. Package initialization

Initialization happens before `main` runs, in this order:

1. **Imported packages** first (each only once, however many importers).
2. **Package-level variables**, in dependency order: a variable is initialised once everything it refers to is ready; otherwise in source order.
3. Every **`init()`** function, in the order they appear (files are presented to the compiler sorted by name).

```go
var (
	a = c + b // == 9
	b = f()   // == 4
	c = f()   // == 5
	d = 3     // == 5 after initialization has finished
)

func f() int {
	d++
	return d
}

func init() { fmt.Println("init:", a, b, c, d) }

func main() { fmt.Println("main") }
// init: 9 4 5 5
// main
```

- The order is `d, b, c, a`: `a` waits for `b` and `c`, and `b`, `c` wait for `d`.
- `init` takes no arguments and returns nothing. You may write several, even in one file.
- You can't call it: `init()` → `undefined: init`.

**Avoid `init` side effects** in packages you share:

| `init` that… | Problem |
|---|---|
| opens files, reads env vars, dials a database | Runs in **every** importer, including tests, and can't return an error |
| registers into a global map | Hidden coupling: behaviour depends on which packages happen to be imported |
| sets mutable package-level state | Tests can't get a fresh copy; two users share one state |

Prefer **explicit** calls: `cfg, err := config.Load(...)`, `reg := render.NewRegistry(); builtin.RegisterAll(reg)`. `init` is fine for cheap, local work that can't fail (precomputing a lookup table). The standard library uses the `import _ "image/png"` idiom (a **blank import**, kept only for its `init`) to register formats; recognise it, but don't copy it for your own APIs.

## 8. Constructors, unexported fields, useful zero values

| Situation | Design |
|---|---|
| Any combination of field values is valid | Exported fields, no constructor needed (`config.Config{Port: 80}`) |
| Some values are invalid | Unexported fields + `New(...) (T, error)` that **validates** |
| A ready state needs setup (maps, defaults) | `New(...) *T`, **or** make the zero value work by lazy setup |

```go
// ✅ the zero value is ready to use, like bytes.Buffer or sync.Mutex
type Counter struct {
	counts map[string]int
}

func (c *Counter) Add(key string) {
	if c.counts == nil { // lazy init: var c Counter just works
		c.counts = map[string]int{}
	}
	c.counts[key]++
}
```

- Name the constructor `New` when the package has one main type (`ring.New`), `NewX` otherwise (`strings.NewReader`).
- **Return concrete types** (`*Server`), accept interfaces (Module 3).
- Validate once in the constructor; methods can then trust the invariants.

## 9. Functional options

A constructor with many optional settings becomes `NewServer("addr", 0, 0, nil, false)`: unreadable, and every new setting breaks callers. **Functional options** keep the required arguments positional and make the rest named and optional.

```go
type Server struct {
	addr     string
	timeout  time.Duration
	maxConns int
}

// Option configures a Server.
type Option func(*Server)

func WithTimeout(d time.Duration) Option { return func(s *Server) { s.timeout = d } }
func WithMaxConns(n int) Option          { return func(s *Server) { s.maxConns = n } }

func New(addr string, opts ...Option) (*Server, error) {
	s := &Server{addr: addr, timeout: 30 * time.Second, maxConns: 100} // defaults
	for _, opt := range opts {
		opt(s) // later options override earlier ones
	}
	if s.maxConns < 1 { // validate after all options are applied
		return nil, errors.New("server: maxConns must be at least 1")
	}
	return s, nil
}

s, err := New(":8080", WithTimeout(5*time.Second))
```

- Adding `WithTLS(...)` later is **not** a breaking change: old calls still compile.
- Defaults live in one place; zero options give a working server.
- Alternative: a `Config` struct argument (`New(addr, Config{Timeout: 5 * time.Second})`). Simpler, but the zero value of a field must then mean "default".

## 10. Evolving an API without breaking users

| Change | Breaks callers? |
|---|---|
| Add a new function, type, method or option | ✅ safe |
| Add a field to a struct | ✅ safe (unless callers use unkeyed literals `T{1, 2}`) |
| Rename / remove an exported name, change a signature | ❌ breaks |
| Add a method to an **interface** | ❌ breaks every type that implemented it |

To replace an API: add the new one, keep the old one working (often by calling the new one), and mark it:

```go
// Next returns the next ID from the package's default generator.
//
// Deprecated: Use a Generator instead; package-level state is shared by every caller.
func Next() string { return defaultGen.Next() }
```

- `Deprecated:` must start a **paragraph** of the doc comment. Editors strike the name through; linters warn on use.
- Say **what to use instead**.

## 11. Module layout (go.dev/doc/modules/layout)

| Project | Layout |
|---|---|
| A library | Package at the module root: `go.mod`, `money.go`, `money_test.go` |
| A library with helpers | Root package + `internal/…` for the helpers |
| A command | `package main` at the root; split into files before splitting into packages |
| Several commands, or commands + library | `cmd/<prog>/main.go` for each program, shared code in `internal/` |
| A server | `cmd/api-server/main.go` + the logic in `internal/…`; servers rarely export packages |

```text
todo-api/
├── go.mod
├── cmd/api/main.go           wiring only: config, store, service, HTTP
└── internal/
    ├── store/                persistence (interface + implementations)
    ├── service/              business rules
    └── httpapi/              handlers
```

- The guide recommends keeping packages **in `internal` as much as possible**, because you can then refactor them freely.
- **Don't build deep trees up front.** Start with one package; split when a piece has a clear name and a small API. `pkg/`, `src/`, `models/`, `utils/` directories are not Go conventions.
- `main` stays thin: read config, build the dependencies, wire them together (**dependency injection**), start.

## Common mistakes

| Mistake | Fix |
|---|---|
| `package util` / `common` / `helpers` | Name the package after what it provides, or keep the helper where it's used |
| `cache.CacheStore`, `user.NewUser` (stutter) | `cache.Store`, `user.New` |
| Exporting struct fields that must stay consistent | Unexported fields + a validating constructor + methods |
| Constructor that silently accepts invalid input | Return `(T, error)` and validate |
| Two packages importing each other | Interface on the consumer side, or extract the shared types into a third package |
| Package-level mutable state (`var counter int`) | A type (`Generator`) the caller creates; keep a `Default` only for compatibility |
| `init()` that reads files, env vars or registers plugins | An explicit `Load` / `Register` call that returns an error |
| Long positional constructors `New(a, 0, "", nil, false)` | Functional options or a `Config` struct |
| Renaming an exported function in place | Add the new one; keep the old one with `// Deprecated: Use X instead.` |
| Importing another tree's `internal/` package | Only code under the parent of `internal/` may import it; ask the owner to export an API |
| Deep `pkg/models/utils` trees before the code needs them | Start flat; `cmd/` + `internal/` when there are several programs |

## Key takeaways

- A **package** is a directory; a **module** is a tree of them; import by **path**, use by **name**.
- Short, lowercase package names; no stutter; no `util`.
- The exported names **are** the API: keep them few, document each, and check with `go doc`.
- `internal/` is importable only from the tree rooted at its parent, so you can refactor it freely.
- Import cycles don't compile: depend on **interfaces declared by the consumer**, or extract a lower package; let `main` wire things together.
- Initialization: imports → package variables (dependency order) → `init()`. Keep `init` free of side effects; prefer explicit calls.
- Unexported fields + validating constructors protect invariants; a useful **zero value** needs no constructor.
- **Functional options** give defaults and extensibility without breaking callers.
- Evolve by **adding**; mark old APIs `// Deprecated: Use X instead.`
- Layout: start flat; `cmd/` for programs, `internal/` for everything you don't want to support publicly.
