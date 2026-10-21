# Swap and Rotate

A Go function can return **several values**: `func F() (int, int)` is called as `x, y := F()`.

Implement:

```go
func Swap(a, b int) (int, int)            // returns b, a
func Rotate(a, b, c int) (int, int, int)  // returns b, c, a
func SwapStrings(a, b string) (string, string)
```

| Call | Result |
|---|---|
| `Swap(1, 2)` | `2, 1` |
| `Rotate(1, 2, 3)` | `2, 3, 1` |
| `SwapStrings("go", "lang")` | `"lang", "go"` |

## Rules
- Inside each function, **reassign the parameters with one parallel assignment** (`a, b = b, a`), then `return a, b`. No temporary variable.
- `return b, a` would also be correct Go, but the point here is to practise `=` with several variables.
