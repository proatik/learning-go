# Module 1 Challenge: Putting It Together

No new concepts here. This page is a one-page recap of Module 1; the ten exercises each combine several of these topics. If you can solve them all without peeking at hints, you have mastered the foundations.

## 1. Programs, packages, tooling

| Thing | Recap |
|---|---|
| Program | `package main` + `func main()`; run with `go run .` |
| Library | any other package name; one package per directory |
| Exported | `Capitalized` names are visible outside the package |
| Tests | `go test ./path/` runs `TestXxx(t *testing.T)` and `ExampleXxx` (`// Output:`) |
| Hygiene | `gofmt -l .` must print nothing; `go vet ./...` must be clean |

## 2. Variables, types, constants

```go
var n int            // zero value: 0, "", false, 0.0
x := 3               // short declaration (inside functions only)
f := float64(x) / 2  // explicit conversion; no implicit numeric conversions
b := byte('A')       // byte = uint8, rune = int32

const Pi = 3.14159   // untyped constant: takes the type its use needs

type Weekday int     // named type used as an enum
const (
	Sunday Weekday = iota // 0
	Monday                // 1
)
```

| Rule | Example |
|---|---|
| Integer division truncates toward zero | `7 / 2 == 3`, `-7 / 2 == -3` |
| `%` takes the sign of the dividend | `-7 % 3 == -1`; wrap safely with `((a % n) + n) % n` |
| Integers overflow silently (wrap around) | `var i int8 = 127; i++` → `-128`; for `x := math.MinInt`, `-x == x` |
| `strconv` returns `(value, error)` | `n, err := strconv.Atoi("42")` |

## 3. Control flow

```go
if v, err := strconv.Atoi(s); err != nil { ... } else { ... } // init statement
for i := 0; i < n; i++ { }   // classic
for i := range n { }         // 0..n-1 (Go 1.22+)
for cond { }                 // while
for { break }                // forever
switch {                     // no fallthrough by default
case x < 0:
case x == 0, x == 1:         // several values
default:
}
outer:
for ... { for ... { continue outer } } // labels
lo, hi := min(a, b, c), max(a, b)      // built-ins (Go 1.21+)
```

## 4. Functions

```go
func divmod(a, b int) (q, r int) { q, r = a/b, a%b; return } // named results
func parse(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty input")
	}
	return 0, fmt.Errorf("bad input %q", s) // formatted error
}
func sum(nums ...int) int { ... }     // variadic: nums is a []int
apply := func(x int) int { return x * 2 } // function value
func counter() func() int { c := 0; return func() int { c++; return c } } // closure
defer fmt.Println("runs when the function returns") // LIFO order
```

Convention: return the error **last**, check it **immediately**, and return zero values alongside a non-nil error.

## 5. Strings

| Fact | Code |
|---|---|
| A string is read-only **bytes** (UTF-8) | `len("é") == 2` |
| `s[i]` is a **byte**; `range s` yields **runes** | `for i, r := range s { ... }` |
| Count characters | `utf8.RuneCountInString(s)` |
| Digit byte → number | `int(c - '0')` when `'0' <= c && c <= '9'` |
| Build efficiently | `var sb strings.Builder; sb.WriteString("x"); sb.WriteByte('-'); sb.WriteRune('é'); sb.String()` |
| Useful `strings` funcs | `Fields`, `Split`, `Join`, `Repeat`, `ReplaceAll`, `TrimSpace`, `HasPrefix`, `Contains`, `ToUpper` |
| Useful `unicode` funcs | `IsDigit`, `IsLetter`, `IsSpace`, `ToUpper` |

| Verb | Output |
|---|---|
| `%d` / `%5d` / `%-5d` / `%05d` | `42` / `   42` / `42   ` / `00042` |
| `%.2f` / `%8.2f` | `3.14` / `    3.14` |
| `%s` / `%q` / `%v` / `%T` | text / `"quoted"` / default / type |
| `%x` / `%b` / `%c` / `%U` | hex / binary / character / `U+00E9` |

## 6. How to approach a problem

1. **Read the README and the test file first.** The table of cases is the spec — note every edge case and error condition.
2. **Write the signature and the edge cases** (empty input, zero, negatives, out of range) before the main logic.
3. **Write the simplest correct version.** Clarity beats cleverness.
4. **Run the tests**, read the first failure message carefully, fix, repeat.
5. **Refactor** once green: remove duplication, name things well, add doc comments.
6. **Finish with `gofmt` and `go vet`** — your editor usually does the first for you.

## Key takeaways

- Every Module 1 tool fits in this page — the challenge is combining them.
- Spell out edge cases before writing the happy path.
- `(value, error)` and table-driven tests are the everyday rhythm of Go.
