# The Nil Interface Trap

The validators below look right, but `CheckName(validUser) == nil` is **false**. Fix the typed-nil bug, then build on the fixed checks.

```go
type Problem interface {
	Problem() string
}

type FieldProblem struct {
	Field, Reason string
}

func (p *FieldProblem) Problem() string // "<Field>: <Reason>"

type User struct {
	Name  string
	Age   int
	Email string
}

func CheckName(u User) Problem
func CheckAge(u User) Problem
func CheckEmail(u User) Problem
func Validate(u User) Problem
func Collect(u User) []Problem
func IsNil(p Problem) bool
```

## Examples

| Call | Result |
|---|---|
| `(&FieldProblem{"age", "out of range"}).Problem()` | `"age: out of range"` |
| `CheckName(User{Name: "Ann"})` | `nil` (a true nil: `== nil` is `true`) |
| `CheckName(User{Name: "  "})` | `&FieldProblem{"name", "required"}` |
| `CheckAge(User{Age: 151})` | `&FieldProblem{"age", "out of range"}` |
| `CheckEmail(User{Email: "ann.example.com"})` | `&FieldProblem{"email", "missing @"}` |
| `Validate(User{Name: "", Age: -1})` | the name problem (first failing check) |
| `Collect(User{Age: 200})` | `[name problem, age problem, email problem]` |
| `IsNil(nil)` / `IsNil((*FieldProblem)(nil))` / `IsNil(&FieldProblem{})` | `true` / `true` / `false` |

## Rules
- Every check returns a `*FieldProblem` on failure and a **nil interface** on success (`== nil` must be `true`):

  | Check | Fails when | Field / Reason |
  |---|---|---|
  | `CheckName` | `strings.TrimSpace(Name) == ""` | `name` / `required` |
  | `CheckAge` | `Age < 0` or `Age > 150` | `age` / `out of range` |
  | `CheckEmail` | `Email` has no `@` | `email` / `missing @` |

- `Validate` runs the checks in the order name, age, email and returns the **first** problem, or `nil`.
- `Collect` returns every problem in that order; no problems → an empty result (`nil` is fine).
- `IsNil` reports whether `p` is a nil interface **or** holds a nil `*FieldProblem`.
