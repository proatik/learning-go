# Writing Tests

So far you have **read** tests: every exercise came with a checker. In this lesson you **write** them. Go's testing tools are built in: the `testing` package plus `go test`. No framework and no assertion library.

## 1. Test files and test packages

- A test lives in a file whose name ends in **`_test.go`**, next to the code it tests (`grade.go` → `grade_test.go`).
- `go build` ignores `_test.go` files. Only `go test` compiles them, so tests never end up in your binary.
- `go test` finds functions by **name and signature**:

| Function | Signature | Purpose |
|---|---|---|
| Test | `func TestXxx(t *testing.T)` | checks behaviour, passes or fails |
| Benchmark | `func BenchmarkXxx(b *testing.B)` | measures speed and allocations (§11) |
| Fuzz test | `func FuzzXxx(f *testing.F)` | checks properties on generated inputs (§12) |
| Example | `func ExampleXxx()` | documentation that is also a test (§10) |
| TestMain | `func TestMain(m *testing.M)` | optional: setup/teardown for the whole package |

- After the prefix, `Xxx` must **not start with a lowercase letter**: `TestGrade` and `Test_grade` are tests, `Testgrade` is not.

Two kinds of test package can live in the same directory:

| `package grade` (internal test) | `package grade_test` (external test) |
|---|---|
| Same package as the code | A separate package that **imports** `grade` |
| Sees unexported names (`parseLine`, `cache`) | Sees only the exported API, like a real user |
| Use it to test internals and small helpers | Use it for API-level tests and Examples, and to avoid import cycles |

```go
// grade_test.go: an external test package
package grade_test

import (
	"testing"

	"example.com/school/grade"
)

func TestGrade(t *testing.T) {
	if got := grade.Grade(95); got != "A" {
		t.Errorf("Grade(95) = %q, want %q", got, "A")
	}
}
```

- Start with an internal test. Switch to `_test` when you want to prove the **public** API is enough.
- Names declared in an internal `_test.go` file are visible to the external test package in the same directory. This is how a std package exposes one internal hook to its external tests (a file called `export_test.go`).

## 2. Your first test: `t.Error` vs `t.Fatal`

```go
// abs.go
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// abs_test.go
func TestAbs(t *testing.T) {
	got := Abs(-3)
	if got != 3 {
		t.Errorf("Abs(-3) = %d, want 3", got)
	}
}
```

- **No assertions.** A test is plain Go: compute `got`, compare it with `want` using `if`, and report with a `t` method.
- A test **passes** if it returns without reporting a failure (and doesn't panic).

| Method | Marks the test failed? | Keeps running? | Use it when |
|---|---|---|---|
| `t.Error(args...)`, `t.Errorf(format, args...)` | yes | **yes** | the other checks are still meaningful: report them all in one run |
| `t.Fatal(args...)`, `t.Fatalf(format, args...)` | yes | **no**, stops this test now | continuing is pointless or would panic: setup failed, `err != nil`, a nil or too-short result |
| `t.Log`, `t.Logf` | no | yes | extra context, printed only on failure or with `-v` |
| `t.Skip`, `t.Skipf` | no | no, the test is skipped | the test can't run here (missing tool, `-short` mode) |

```go
u, err := Parse("https://go.dev/doc")
if err != nil {
	t.Fatalf("Parse() error = %v", err) // u is unusable: stop here
}
if u.Host != "go.dev" {
	t.Errorf("Host = %q, want %q", u.Host, "go.dev") // reported...
}
if u.Path != "/doc" {
	t.Errorf("Path = %q, want %q", u.Path, "/doc") // ...and so is this
}
```

- `Fatal` stops the test by ending its goroutine. Call it only from the test goroutine, never from a goroutine you started.

## 3. Failure messages that explain themselves

When a test fails, the message is all the reader gets. The Go convention:

```text
Func(input) = got, want want
```

```go
t.Errorf("Abs(%d) = %d, want %d", in, got, want)               // numbers
t.Errorf("Slugify(%q) = %q, want %q", in, got, want)           // strings: %q shows spaces and ""
t.Errorf("ParsePort(%q) error = %v, want %v", in, err, ErrEmpty) // errors
```

Output:

```text
--- FAIL: TestAbs (0.00s)
    abs_test.go:8: Abs(-3) = -3, want 3
```

| ❌ Unhelpful | ✅ Helpful |
|---|---|
| `t.Error("wrong")` | `t.Errorf("Abs(%d) = %d, want %d", in, got, want)` |
| `t.Errorf("expected %d got %d", want, got)` | **got first, then want**, as in the convention |
| `t.Errorf("Slugify = %s", got)` | include the **input**: `Slugify(%q) = %q` |
| `%v` for strings | `%q`, so `""`, `" "` and `"\n"` are visible |

- `go test` prefixes every message with `file_test.go:LINE:`, so you don't need to say where.

## 4. Table-driven tests

One test function, a **table** of cases, one loop. This is the standard Go style.

```go
func TestIsLeap(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2023, false}, // not divisible by 4
		{2024, true},  // divisible by 4
		{1900, false}, // century
		{2000, true},  // 400th year
	}
	for _, tt := range tests {
		if got := IsLeap(tt.year); got != tt.want {
			t.Errorf("IsLeap(%d) = %t, want %t", tt.year, got, tt.want)
		}
	}
}
```

- Adding a case is one line. The cases are easy to review: you can see what is covered and what is missing.
- The table is a slice of an **anonymous struct**: one field per input, plus `want` (and `wantErr` for errors, §8).
- Use `t.Errorf`, not `t.Fatalf`, in the loop, so one bad case doesn't hide the others.

## 5. Subtests with `t.Run`

`t.Run(name, func(t *testing.T) {...})` runs a **named subtest** with its own `t`.

```go
func TestGrade(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{"lowest A", 90, "A"},
		{"highest B", 89, "B"},
		{"negative", -1, "invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Grade(tt.score); got != tt.want {
				t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}
```

```text
--- FAIL: TestGrade (0.00s)
    --- FAIL: TestGrade/lowest_A (0.00s)
        grade_test.go:18: Grade(90) = "B", want "A"
```

- The full name is `Parent/child`. Spaces become `_`: `"lowest A"` → `TestGrade/lowest_A`.
- **Run one case** with `-run`. The pattern is split at `/`, one regexp per level:

| Flag | Runs |
|---|---|
| `-run TestGrade` | `TestGrade` and all its subtests |
| `-run 'TestGrade/lowest_A'` | only that case |
| `-run 'TestGrade/lowest'` | every case whose name contains `lowest` |
| `-run '/negative'` | the `negative` case of every test |

- `t.Fatal` inside a subtest stops **only that subtest**; the loop moves on to the next case.
- Since Go 1.22 each loop iteration has its own `tt`, so the closure captures the right case.
- Name cases by **what they check** (`"lowest A"`, `"empty input"`), not `"case 1"`.

## 6. Helpers and `t.Helper()`

Repeated checks belong in a helper. Mark it with `t.Helper()`, so a failure points at the **caller's** line, not at the helper's.

```go
func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestInitials(t *testing.T) {
	assertEqual(t, Initials("ada lovelace"), "AL") // line 13
	assertEqual(t, Initials("grace hopper"), "GH") // line 14
}
```

With a buggy `Initials` that forgets to upper-case (the `t.Errorf` is on line 8):

| Without `t.Helper()` | With `t.Helper()` |
|---|---|
| `initials_test.go:8: got al, want AL` | `initials_test.go:13: got al, want AL` |
| `initials_test.go:8: got gh, want GH` | `initials_test.go:14: got gh, want GH` |
| both point at the `t.Errorf` inside the helper | each points at the failing call |

- Call `t.Helper()` as the **first line** of every helper that reports failures.
- Helpers take `t *testing.T` as the **first** parameter. Use `testing.TB` instead if benchmarks share the helper.
- A helper may call `t.Fatalf` for setup that must work (`mustOpen(t, path)`). The test stops, and the message points at the caller.

## 7. Setup and cleanup: `t.Cleanup`, `t.TempDir`, `t.Context`

| Method | What it does |
|---|---|
| `t.Cleanup(f)` | runs `f` after the test and all its subtests finish, **even after `Fatal`**. Several cleanups run **last-in, first-out** |
| `t.TempDir()` | returns a new empty directory, deleted automatically after the test |
| `t.Context()` (Go 1.24) | a `context.Context` that is canceled **just before** the cleanups run |
| `t.Setenv(k, v)` | sets an environment variable, restored after the test |
| `t.Chdir(dir)` (Go 1.24) | changes the working directory, restored after the test |

- `t.Cleanup` vs `defer`: `defer` runs when the **current function** returns. `t.Cleanup` runs when the **test** ends, so a helper can register it for the test that called it:

```go
// openLog opens a log in a fresh temp dir and closes it when the test ends.
func openLog(t *testing.T) (*Log, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "app.log")
	l, err := Open(path)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", path, err)
	}
	t.Cleanup(func() { l.Close() })
	return l, path
}
```

- Cleanup order, verified:

```go
func TestOrder(t *testing.T) {
	t.Cleanup(func() { fmt.Println("cleanup 1") })
	t.Cleanup(func() { fmt.Println("cleanup 2") })
	ctx := t.Context()
	t.Cleanup(func() { fmt.Println("ctx canceled:", ctx.Err() != nil) })
	fmt.Println("test body")
}
// test body
// ctx canceled: true
// cleanup 2
// cleanup 1
```

- **Never write test files into the package directory.** Use `t.TempDir()`: every test gets its own directory, and nothing is left behind.
- On Windows an **open file can't be deleted**, so a forgotten `Close` makes the `TempDir` cleanup fail: `TempDir RemoveAll cleanup: ... being used by another process`.
- `t.Setenv` and `t.Chdir` change the whole process, so they can't be used in parallel tests.

## 8. Testing errors

Put the expected error in the table, then check it with `errors.Is`:

```go
tests := []struct {
	in      string
	want    int
	wantErr error // nil for success
}{
	{"8080", 8080, nil},
	{"", 0, ErrEmpty},
	{"0", 0, ErrOutOfRange},
}
for _, tt := range tests {
	got, err := ParsePort(tt.in)
	if !errors.Is(err, tt.wantErr) {
		t.Errorf("ParsePort(%q) error = %v, want %v", tt.in, err, tt.wantErr)
		continue
	}
	if got != tt.want {
		t.Errorf("ParsePort(%q) = %d, want %d", tt.in, got, tt.want)
	}
}
```

- `errors.Is(err, nil)` is true only when `err == nil`, so one check covers success **and** failure cases.
- **Check the value on errors too.** The convention is the zero value (`0`, `""`, `nil`) alongside a non-nil error.
- For an error **type**, use `errors.As` and then check the fields:

```go
var se *SyntaxError
if !errors.As(err, &se) {
	t.Fatalf("Parse(%q) error = %v, want a *SyntaxError", in, err)
}
if se.Pos != 1 {
	t.Errorf("Parse(%q) error Pos = %d, want 1", in, se.Pos)
}
```

| ❌ Brittle | ✅ Robust |
|---|---|
| `err == ErrEmpty` (fails once the error is wrapped) | `errors.Is(err, ErrEmpty)` |
| `err.Error() == "port \"\": empty port"` (breaks when the wording changes) | `errors.Is` / `errors.As`. Compare text only if the message itself is the contract |
| `if err != nil` only (any error passes) | check **which** error |

## 9. Time and other dependencies: inject them

Code that calls `time.Now()` directly can only be tested by **sleeping**: slow and flaky. Instead, let the caller pass the clock in:

```go
type Session struct {
	start time.Time
	ttl   time.Duration
	now   func() time.Time // time.Now in production, a fake in tests
}

func NewSession(ttl time.Duration, now func() time.Time) *Session {
	return &Session{start: now(), ttl: ttl, now: now}
}

func (s *Session) Expired() bool { return !s.now().Before(s.start.Add(s.ttl)) }
```

```go
func TestSessionExpired(t *testing.T) {
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	s := NewSession(time.Hour, func() time.Time { return now })

	now = now.Add(59 * time.Minute) // "move the clock"
	if s.Expired() {
		t.Error("Expired() after 59m = true, want false")
	}
	now = now.Add(time.Minute) // exactly at the deadline
	if !s.Expired() {
		t.Error("Expired() after 1h = false, want true")
	}
}
```

- The closure reads the variable `now` on every call. Assigning to `now` moves time for the session too.
- Test **exactly at** the boundary, one tick before and one after. `>` vs `>=` bugs live there.
- The same idea works for every outside dependency: write output to an `io.Writer` (pass a `*bytes.Buffer` in tests), read input from an `io.Reader` (`strings.NewReader`), call a small interface instead of the network.
- For code that uses timers and goroutines, `testing/synctest` (Go 1.25) runs it with a fake clock. It is covered in the concurrency module.

## 10. Example functions

An Example is **documentation that is also a test**. `go test` runs it and compares what it prints with the `// Output:` comment.

```go
func ExampleGrade() {
	fmt.Println(Grade(95))
	fmt.Println(Grade(85))
	// Output:
	// A
	// B
}
```

- The comparison ignores leading and trailing space.
- **No output comment, no run:** the example is compiled but never executed, so it checks nothing.
- `// Unordered output:` accepts the lines in any order, e.g. when ranging over a map:

```go
func ExampleGrade_unordered() {
	for _, s := range map[string]int{"ann": 95, "bob": 85} {
		fmt.Println(Grade(s))
	}
	// Unordered output:
	// B
	// A
}
```

| Name | Documents |
|---|---|
| `Example()` | the package |
| `ExampleGrade()` | the function `Grade` |
| `ExampleSession()` | the type `Session` |
| `ExampleSession_Expired()` | the method `Session.Expired` |
| `ExampleGrade_boundaries()` | a second example for `Grade` (suffix starts lowercase) |

- `go doc` and pkg.go.dev show examples next to the API they document. Write them in an external `_test` package, so they read exactly like user code (`grade.Grade(95)`).

## 11. Benchmarks with `b.Loop`

```go
func BenchmarkJoinInts(b *testing.B) {
	nums := make([]int, 100) // setup: not measured
	for i := range nums {
		nums[i] = i * 1000
	}
	b.ReportAllocs()
	for b.Loop() {
		JoinInts(nums, ",") // measured
	}
}
```

```text
BenchmarkJoinInts-8   	   99050	     12717 ns/op	    2432 B/op	     107 allocs/op
```

| Column | Meaning |
|---|---|
| `-8` | `GOMAXPROCS` during the run |
| `99050` | how many iterations ran |
| `12717 ns/op` | time per iteration |
| `2432 B/op`, `107 allocs/op` | memory per iteration (from `b.ReportAllocs()` or the `-benchmem` flag) |

- `107 allocs/op` for 100 numbers is a hint: something allocates once per number. That is the kind of thing a benchmark is for.

- **`for b.Loop()`** (Go 1.24) is the modern form. It times only the loop, runs the setup once, and keeps the call's results alive, so the compiler can't optimize the work away.
- The older form is `for range b.N { ... }`. You will see it in existing code. It needs `b.ResetTimer()` after setup and a sink variable for results. Don't mix it with `b.Loop`.
- `go test` **does not run benchmarks** unless you pass `-bench`. Results vary between runs; compare runs with the `benchstat` tool (`golang.org/x/perf/cmd/benchstat`) rather than trusting one number.

## 12. Fuzzing

A fuzz test has **seed inputs** and a function that checks **properties**: facts that must hold for *every* input, since there is no fixed `want`.

```go
func FuzzReverse(f *testing.F) {
	f.Add("hello")      // seed corpus
	f.Add("héllo, 世界")
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip("Reverse is defined on valid UTF-8 only")
		}
		r := Reverse(s)
		if back := Reverse(r); back != s {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, back, s) // round trip
		}
		if !utf8.ValidString(r) {
			t.Errorf("Reverse(%q) = %q, not valid UTF-8", s, r)
		}
	})
}
```

- **Plain `go test` runs only the seeds**, once each, like a table test (`FuzzReverse/seed#0`, `seed#1`, ...). Fast and deterministic.
- **`go test -fuzz=FuzzReverse -fuzztime=30s`** keeps mutating the seeds to explore new inputs until it finds a failure or time runs out. The regexp must match exactly one fuzz test.
- A failing input is saved under `testdata/fuzz/FuzzReverse/<hash>`. From then on, plain `go test` runs it as a seed: a permanent regression test. Commit that file.
- For example, with only the `"hello"` seed and a `Reverse` that swaps **bytes**, the seeds pass, but `-fuzz` finds a failing input in under a second:

```text
--- FAIL: FuzzReverse (0.65s)
    --- FAIL: FuzzReverse (0.00s)
        rev_test.go:19: Reverse("矟") = "\x9f\x9f\xe7", not valid UTF-8

    Failing input written to testdata\fuzz\FuzzReverse\8d7271b7de71140d
    To re-run:
    go test -run=FuzzReverse/8d7271b7de71140d
```

- Fuzz arguments can be `string`, `[]byte`, `bool`, `rune`, `byte`, and the int, uint and float types.
- Useful properties:

| Property | Example |
|---|---|
| **Round trip** | `Decode(Encode(x)) == x`, `Reverse(Reverse(s)) == s` |
| **Invariants** | the output is valid UTF-8, has the same length, is sorted |
| **Never panics** | just calling the function on any input |
| **Agrees with a simple version** | a fast `Sort` gives the same result as a slow obvious one |

- One property is rarely enough. A byte-wise `Reverse` passes the round trip and only fails the UTF-8 check.

## 13. Coverage

With the `Grade` from §5 (bands A to F) and a table that only checks `95`, `85` and `-1`:

```text
$ go test -cover ./grade
ok  	example.com/school/grade	2.600s	coverage: 57.1% of statements

$ go test -coverprofile=c.out ./grade
$ go tool cover -func=c.out      # per function
example.com/school/grade/grade.go:4:	Grade		57.1%
total:					(statements)	57.1%

$ go tool cover -html=c.out      # source in the browser, uncovered lines in red
```

- The missing 42.9% is the `C`, `D` and `F` branches: no case reaches them.
- Coverage shows which lines **ran**, not which were **checked**. A test that calls every function and asserts nothing has high coverage and catches nothing.
- Use it to find **untested code** (an error branch nobody reaches). Don't chase 100%.

## 14. What makes a test valuable

Ask of every test: **"Would this test fail if the code were wrong?"**

- **Mutation testing** answers that question mechanically. Plant a small bug (`>` → `>=`, drop a line, return the wrong error), run the tests, and see if they fail. A bug that **survives** shows a missing case.
- This lesson's exercises work exactly like that. The function under test is already written. You write the tests, and the checker re-runs them against several broken versions. Your tests must **pass on the correct code and fail on every bug**.

Where bugs hide, so where your cases should go:

| Look at | Cases |
|---|---|
| Boundaries | each edge, one below and one above (`89`, `90`, `100`, `101`) |
| Empty and zero | `""`, `nil`, `0`, a single element |
| Every rule in the spec | one case per rule that **only** that rule decides (`1900` for the century rule) |
| Error paths | every sentinel, plus the value returned with it |
| Text | spaces, punctuation, upper case, multi-byte runes (`é`, `世`) |
| State | call order: after `Close`, after a reopen, twice in a row |

Good tests are also:

- **Deterministic:** no real clock, no `time.Sleep`, no network, no dependence on map order.
- **Fast:** milliseconds, so you run them all the time.
- **About behaviour, not implementation:** test what the function promises, so you can refactor without rewriting the tests.
- **Independent:** each test sets up its own data (`t.TempDir`, fresh values) and doesn't rely on another test running first.

## 15. `go test` flags you will use

| Flag | Effect |
|---|---|
| `-v` | print every test name (`=== RUN`) and every `t.Log` |
| `-run regexp` | run only matching tests, examples and fuzz seeds (`/` separates subtest levels) |
| `-skip regexp` | skip matching tests |
| `-count=1` | ignore the test cache and really run again |
| `-failfast` | stop after the first failing test |
| `-short` | ask long tests to skip themselves (`if testing.Short() { t.Skip(...) }`) |
| `-list regexp` | list matching tests, benchmarks, fuzz tests and examples without running them |
| `-cover`, `-coverprofile=c.out` | measure coverage (§13) |
| `-bench regexp`, `-benchmem`, `-benchtime=2s` or `=100x` | run benchmarks (§11) |
| `-fuzz regexp`, `-fuzztime=30s` | fuzz one target (§12) |
| `-race` | detect data races (concurrency module) |

- `go test` also runs a subset of `go vet` first, so common mistakes (a bad `Errorf` format, a misnamed Example) fail the build.

## Common mistakes

| Mistake | Fix |
|---|---|
| `t.Fatalf` inside a table loop | `t.Errorf` (or a subtest), so the other cases still run |
| `t.Errorf` before using a nil or short result | `t.Fatalf` first, then index or call methods |
| `"expected %v got %v"` with no input | `"F(%v) = %v, want %v"`: input, got, then want |
| Only happy-path cases | add boundaries, empty input and every error path |
| `err == ErrX` on a wrapped error | `errors.Is(err, ErrX)` |
| Comparing `err.Error()` text | `errors.Is` / `errors.As` |
| Helper failures all point at the helper | `t.Helper()` as its first line |
| `time.Sleep` to test expiry | inject `now func() time.Time` and move a fake clock |
| Example without `// Output:` | add the comment: otherwise it never runs |
| Writing files into the package folder | `t.TempDir()` |
| Forgetting to close files in tests | `t.Cleanup(func() { f.Close() })` right after opening |
| `for i := 0; i < b.N; i++` with setup inside | `for b.Loop()` with setup before it |
| Fuzz test with only a round-trip check | add invariants: some bugs survive a round trip |
| Chasing 100% coverage | ask "would this fail if the code were wrong?" |
| Tests that depend on map order or run order | sort, use `// Unordered output:`, keep tests independent |

## Key takeaways

- Tests are plain Go in `_test.go` files: `TestXxx(t *testing.T)`, an `if`, and `t.Errorf`.
- `t.Error` reports and continues; `t.Fatal` reports and stops the current (sub)test.
- Messages follow `F(input) = got, want want`.
- Table-driven tests + `t.Run` subtests are the default shape. Run one case with `-run 'TestX/case_name'`.
- `t.Helper()` makes failures point at the caller; `t.Cleanup` and `t.TempDir` clean up for you, even after `Fatal`.
- Check errors with `errors.Is` / `errors.As`, and check the value returned with them.
- Inject time and other dependencies instead of sleeping or touching the real world.
- Examples are tested docs, benchmarks use `for b.Loop()`, fuzz tests check properties on seeds and, with `-fuzz`, on generated inputs.
- Coverage tells you what ran. Mutation thinking tells you what is **checked**: a good test fails when the code is wrong.
