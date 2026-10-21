# Rotate in Place

Rotate a slice's elements in place, without allocating:

```go
func RotateLeft(s []int, k int)
func RotateRight(s []int, k int)
```

Rotating left by `k` moves `s[k]` to the front; the first `k` elements wrap around to the end.

| Before | Call | After |
|---|---|---|
| `[1 2 3 4 5]` | `RotateLeft(s, 2)` | `[3 4 5 1 2]` |
| `[1 2 3 4 5]` | `RotateRight(s, 2)` | `[4 5 1 2 3]` |
| `[1 2 3 4 5]` | `RotateLeft(s, 7)` | `[3 4 5 1 2]` (7 ≡ 2 mod 5) |
| `[1 2 3 4 5]` | `RotateLeft(s, -1)` | `[5 1 2 3 4]` (negative = the other way) |
| `[1 2 3]` | `RotateLeft(s, 3)` | `[1 2 3]` |
| `[]` / `nil` | `RotateLeft(s, 4)` | unchanged, no panic |

## Rules
- Modify `s` in place; return nothing.
- **Zero allocations** (checked with `testing.AllocsPerRun`): no `make`, no `append` to a new slice.
- Any `k` is valid, including `k > len(s)` and negative `k`. Watch out: in Go, `-1 % 5 == -1`.
- Don't use the `slices` package. Writing a small unexported `reverse` helper is encouraged.
