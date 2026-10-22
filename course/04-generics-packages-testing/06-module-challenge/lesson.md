# Module 4 Challenge: Reusable, Tested Code

No new concepts here. This page is a one-page recap of Module 4; the ten exercises each combine generics, iterators, package design and testing into small reusable libraries — and three of them ask **you** to write the tests.

## 1. Generic functions and constraints

```go
func Map[T, U any](s []T, f func(T) U) []U { ... }    // type parameters in [ ]
func Max[T cmp.Ordered](a, b T) T { ... }             // constrained: supports < > <= >=
func Index[T comparable](s []T, v T) int { ... }      // supports == and !=

out := Map(nums, strconv.Itoa)                        // T and U inferred from the arguments
zero := Zero[int]()                                   // no arguments to infer from: write [int]
```

| Constraint | Allows | Example types |
|---|---|---|
| `any` | assign, pass, store — no operators | everything |
| `comparable` | `==`, `!=`, map keys | `int`, `string`, structs of comparables |
| `cmp.Ordered` | `<` `<=` `>` `>=`, plus `==` | integers, floats, `string` |
| `interface{ ~int \| ~float64 }` | the operators **all** listed types share (`+ - * /`, `<`) | `int`, `float64`, `type Celsius float64` |
| `interface{ String() string }` | calling the method | any type with that method |
| `interface{ ~int; String() string }` | both: type set **and** method | a named int with a `String` method |

- `~T` means "any type whose **underlying** type is `T`"; without `~`, `type Celsius float64` is rejected.
- `var zero T` is the zero value of a type parameter — use it for "not found" results.
- Constraints with a type list (`~int | ~float64`) can **only** be used as constraints, never as ordinary variable types.

## 2. Generic types and methods

```go
type Stack[T any] struct{ items []T }                 // type parameter on the type

func NewStack[T any]() *Stack[T] { return &Stack[T]{} }
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) } // receiver repeats [T], no constraint

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

type Box[T any] struct{ V T }

func (b Box[T]) Map[U any](f func(T) U) Box[U] { return Box[U]{f(b.V)} } // Go 1.27 generic method

s := NewStack[string]()                               // must instantiate: nothing to infer from
b := Box[int]{V: 3}.Map(strconv.Itoa)                 // U inferred as string → Box[string]
```

- Methods use the type's parameters; a **generic method** (Go 1.27) adds its **own** in `[U any]` after the name.
- A generic method never satisfies an interface method — interfaces cannot declare type parameters.
- A type with a `less func(a, b T) bool` field works for **any** `T`; a `T cmp.Ordered` constraint is simpler but narrower. Offer both: `New(less)` and `NewOrdered[T cmp.Ordered]()`.

## 3. Iterators

```go
type Seq[V any] func(yield func(V) bool)              // iter.Seq
type Seq2[K, V any] func(yield func(K, V) bool)       // iter.Seq2

func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s.items) - 1; i >= 0; i-- {
			if !yield(s.items[i]) {                     // the yield contract:
				return                                  // false → stop at once, never call yield again
			}
		}
	}
}

for v := range s.All() { ... }                        // break makes yield return false
```

| Tool | Does |
|---|---|
| `slices.Values(s)` / `slices.All(s)` | slice → `Seq[V]` / `Seq2[int, V]` |
| `maps.Keys(m)` / `maps.All(m)` | map → `Seq[K]` / `Seq2[K, V]` (random order) |
| `slices.Collect(seq)` / `slices.Sorted(seq)` | `Seq[V]` → slice (sorted) |
| `next, stop := iter.Pull(seq); defer stop()` | push → pull: `v, ok := next()` one value at a time |
| Adapter `Filter(seq, keep) iter.Seq[V]` | wraps a sequence; lazy — does no work until ranged |

- Calling `yield` after it returned `false` **panics** at run time.
- Iterators are **lazy**: `Take(Map(seq, f), 3)` calls `f` only three times.
- Return an empty sequence, never a `nil` func: `return func(yield func(V) bool) {}`.

## 4. Packages and API design

| Rule | Detail |
|---|---|
| Exported = capitalised | `Store`, `New` are API; `items`, `config` are hidden |
| One directory = one package | import path = module path + directory: `learngo/.../inventory/pricing` |
| `internal/` | `a/internal/x` can be imported only by code inside `a/` |
| No import cycles | if `a` imports `b`, `b` can't import `a` — move the shared code down into a third package |
| Package names | short, lowercase, no stutter: `graph.New`, not `graph.NewGraph` |
| Constructors | `New(...)` sets invariants (makes maps, validates) and returns `*T` |
| Functional options | `type Option func(*config)`; `New(opts ...Option)`; defaults first, then `for _, o := range opts { o(&c) }` |
| Errors are API | export sentinels (`ErrCycle`) and types (`*CycleError`); wrap with `%w` across layers |

## 5. Testing toolkit

| Tool | Use it to |
|---|---|
| table + `for _, tt := range tests` | list cases as data: name, input, want |
| `t.Run(tt.name, func(t *testing.T) {...})` | subtests: each case gets its own name and its own pass/fail |
| `t.Errorf` / `t.Fatalf` | record and continue / record and stop this test |
| `t.Helper()` | report failures at the caller's line, not inside the helper |
| `t.Cleanup(f)` / `t.TempDir()` | teardown after the test / a directory deleted afterwards |
| `func ExampleX() { ... // Output: ... }` | documentation that is checked against its printed output |
| `for b.Loop() { ... }` | benchmark body (Go 1.24+): setup above the loop isn't timed |
| `f.Add(seed); f.Fuzz(func(t *testing.T, s string) {...})` | fuzz test: seeds run on every `go test`; check **properties**, e.g. `Decode(Encode(s)) == s` |
| `errors.Is(err, ErrX)` in tests | check the error kind, never `err.Error()` text |

- A good test **fails for a broken implementation**. Hunt for edges: empty, one element, boundaries (`>` vs `>=`), duplicates, order, Unicode.
- Test the contract, not the implementation: compare results, not internal fields.

## Designing reusable code — checklist

- **Start concrete, then generalise.** Write it for `int`, then replace `int` with `T` and pick the **weakest** constraint that compiles.
- **Return iterators for "walk everything"** (`All`, `Keys`, `BFS`); return slices when callers need the whole result anyway.
- **Respect the yield contract** in every iterator and adapter: check every `yield`, `return` on `false`.
- **Small packages with one job.** The leaf package (`graph`, `pipe`) knows nothing about the caller (`build`, `report`).
- **Hide what may change** behind `internal/` and unexported fields; expose constructors, options and methods.
- **Make zero results safe.** Empty input → empty, non-nil output; missing key → `zero, false`.
- **Prove it with tests** that would catch the bugs you are most likely to write.

## Common mistakes

| Mistake | Fix |
|---|---|
| `func Max[T any](a, b T) T { if a > b ...}` | `any` has no `<`: use `cmp.Ordered` |
| `type Celsius float64` rejected by `interface{ float64 }` | add the tilde: `~float64` |
| Ignoring `yield`'s result | `if !yield(v) { return }` |
| Returning a `nil` `iter.Seq` from a stub | ranging over a nil func panics; return an empty func |
| `iter.Pull` without `stop` | `defer stop()` right after the call |
| `import "x/internal/y"` from outside `x/` | only code under `x/` may import it |
| Two packages importing each other | extract the shared types into a lower package |
| Tests that only check the "happy path" | add boundaries and error cases until every planted bug fails |

## Key takeaways

- Type parameters + the weakest constraint make code reusable without `any` casts.
- Generic types carry `[T]` on every method receiver; Go 1.27 methods may add their own type parameters.
- `iter.Seq` / `iter.Seq2` iterators are lazy functions; the yield contract is non-negotiable.
- Packages are API boundaries: export little, hide the rest in unexported names and `internal/`.
- Tests are part of the deliverable: table-driven, subtests, helpers, examples, benchmarks and fuzzing.
