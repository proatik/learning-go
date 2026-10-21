# Expression Evaluator

Parse arithmetic like `2 * (3 + -4)` into a **tree of values that share one interface**, then evaluate and print it. Each node type knows how to `Eval` and `String` itself; errors are typed and positioned.

```go
var ErrDivisionByZero = errors.New("expr: division by zero")

type Expr interface {
	Eval() (float64, error)
	String() string
}

type Num float64                        // 3.5
type Neg struct{ X Expr }               // -X
type Binary struct {                    // L Op R
	Op   byte // '+', '-', '*' or '/'
	L, R Expr
}

type DivisionByZeroError struct{ Expr string }
func (e *DivisionByZeroError) Error() string
func (e *DivisionByZeroError) Is(target error) bool

type SyntaxError struct {
	Pos int // byte offset
	Msg string
}
func (e *SyntaxError) Error() string

func Parse(s string) (Expr, error)
func Eval(s string) (float64, error)
```

## Examples

| Code | Result |
|---|---|
| `Parse("1 + 2 * 3").String()` | `"(1 + (2 * 3))"` |
| `Parse("-2 * 3").String()` | `"((-2) * 3)"` |
| `Eval("(1 + 2) * 3")` | `9, nil` |
| `Eval("1 / (2 - 2)")` | `0`, `"expr: division by zero in (1 / (2 - 2))"` |
| `Parse("1 + * 2")` | `*SyntaxError{4, "unexpected \"*\""}` → `"expr: syntax error at 4: unexpected \"*\""` |
| `Parse("(1 + 2")` | `*SyntaxError{6, "unexpected end of input"}` |

## Grammar

```text
expr    = term    { ("+" | "-") term }      left-associative
term    = unary   { ("*" | "/") unary }     left-associative
unary   = "-" unary | primary
primary = number | "(" expr ")"
```

## Rules
- **Tokens**: spaces and tabs are skipped; a **number** is a maximal run of digits and `.` converted with `strconv.ParseFloat`; the other tokens are `+ - * / ( )`.
- **Parse** returns the tree built from **values** — `Num`, `Neg`, `Binary`, not pointers — or `nil` and a `*SyntaxError`:
  - an unknown character, or a token where it doesn't fit (including leftovers after a complete expression) → `Pos` = that token's offset, `Msg` = `fmt.Sprintf("unexpected %q", text)`
  - input ends too early → `Pos = len(s)`, `Msg = "unexpected end of input"`
  - a number `ParseFloat` rejects (`"1.2.3"`, `"."`) → its offset, `Msg = fmt.Sprintf("invalid number %q", text)`
  - there is no unary `+`
- **SyntaxError.Error**: `"expr: syntax error at %d: %s"`.
- **String**: `Num` uses `strconv.FormatFloat(f, 'g', -1, 64)`; `Neg` is `"(-X)"`; `Binary` is `"(L op R)"` with single spaces.
- **Eval**: evaluate `L` then `R`; return the first error **unchanged**. Dividing by zero returns `0, &DivisionByZeroError{Expr: b.String()}` for that `Binary`.
- **DivisionByZeroError**: `Error()` is `"expr: division by zero in " + Expr`; `Is` returns `target == ErrDivisionByZero`.
- **Eval(s)** = `Parse` then `Eval`, returning either step's error.
