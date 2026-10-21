# Integer Stack

Methods on a **slice type**: `Push` must update the caller's slice header, so every method uses a pointer receiver. Then use the stack to evaluate Reverse Polish Notation.

```go
type Stack []int

func (s *Stack) Push(v int)
func (s *Stack) Pop() (int, bool)
func (s *Stack) Peek() (int, bool)
func (s *Stack) Len() int
func (s *Stack) Values() []int // bottom → top, a copy

func EvalRPN(tokens []string) (int, bool)
```

## Examples

| Code | Result |
|---|---|
| `var s Stack; s.Push(1); s.Push(2)` | `s.Values()` = `[1 2]`, `s.Len()` = `2` |
| `s.Peek()` | `2, true` (stack unchanged) |
| `s.Pop()` | `2, true` (then `s.Values()` = `[1]`) |
| `var e Stack; e.Pop()` | `0, false` |
| `EvalRPN([]string{"2", "3", "+", "4", "*"})` | `20, true` — (2+3)*4 |
| `EvalRPN([]string{"10", "3", "-"})` | `7, true` |
| `EvalRPN([]string{"1", "+"})` | `0, false` |
| `EvalRPN([]string{"4", "0", "/"})` | `0, false` |

## Rules
- The zero value (`var s Stack`, a nil slice) is an empty, usable stack.
- `Pop` / `Peek` on an empty stack return `0, false`.
- `Values` returns a **new** slice (changing it must not affect the stack); empty stack → empty or nil slice.
- `EvalRPN`: tokens are integers (may be negative, like `"-3"`) or one of `+ - * /` (integer division, truncated).
  - For an operator, pop `b` then `a` and push `a op b`.
  - Return `0, false` if an operator lacks two operands, a token is invalid, you divide by zero, or the stack doesn't hold **exactly one** value at the end (including no tokens).
