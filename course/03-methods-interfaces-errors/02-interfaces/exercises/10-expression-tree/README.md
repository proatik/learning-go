# Expression Tree

Model arithmetic expressions as a tree of nodes that all satisfy one `Expr` interface, then walk the tree with type switches to simplify it.

```go
type Expr interface {
	Eval(env map[string]float64) float64
	String() string
}

type Num float64           // a constant
type Var string            // a variable, looked up in env
type Add struct{ L, R Expr }
type Mul struct{ L, R Expr }
type Neg struct{ X Expr }

func Simplify(e Expr) Expr
func Vars(e Expr) []string
func Size(e Expr) int
```

## Examples

With `e := Add{Mul{Num(2), Var("x")}, Neg{Num(3)}}`:

| Call | Result |
|---|---|
| `e.String()` | `"((2 * x) + -3)"` |
| `e.Eval(map[string]float64{"x": 5})` | `7` |
| `e.Eval(nil)` | `-3` (missing variables are `0`) |
| `Num(2.5).String()` / `Var("y").String()` | `"2.5"` / `"y"` |
| `Simplify(Add{Var("x"), Num(0)})` | `Var("x")` |
| `Simplify(Mul{Add{Num(1), Num(2)}, Var("x")})` | `Mul{Num(3), Var("x")}` |
| `Simplify(Neg{Neg{Var("x")}})` | `Var("x")` |
| `Vars(Add{Var("y"), Mul{Var("x"), Var("y")}})` | `["x" "y"]` |
| `Size(e)` | `6` |

## Rules
- All methods use value receivers. You may assume no child is `nil`.
- `String`:

  | Node | Format |
  |---|---|
  | `Num` | `strconv.FormatFloat(f, 'g', -1, 64)` (`"3"`, `"2.5"`, `"-1"`) |
  | `Var` | the name |
  | `Add` / `Mul` | `"(" + L + " + " + R + ")"` / `"(" + L + " * " + R + ")"` |
  | `Neg` | `"-" + X` |

- `Eval`: a `Var` missing from `env` (or a nil `env`) evaluates to `0`.
- `Simplify` simplifies the children **first** (bottom-up), then applies the first matching rule:

  | Node (after simplifying children) | Result |
  |---|---|
  | `Add{Num a, Num b}` | `Num(a + b)` |
  | `Add{Num 0, r}` / `Add{l, Num 0}` | `r` / `l` |
  | `Mul{Num a, Num b}` | `Num(a * b)` |
  | `Mul` with either side `Num 0` | `Num(0)` |
  | `Mul{Num 1, r}` / `Mul{l, Num 1}` | `r` / `l` |
  | `Neg{Num a}` | `Num(-a)` |
  | `Neg{Neg{x}}` | `x` |
  | anything else | the node rebuilt with its simplified children |

  `Num` and `Var` are returned unchanged, and so is any `Expr` type the package doesn't know.
- `Vars` returns the distinct variable names, sorted (none → empty, `nil` is fine). Unknown node types contribute nothing.
- `Size` counts the nodes: `Num` and `Var` are 1; `Add`/`Mul` are 1 + both sides; `Neg` is 1 + its child; unknown types count as 1.
