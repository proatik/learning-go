# Set Operations

Go has no set type; `map[string]struct{}` is the idiomatic set (the empty struct takes no memory). This exercise names it:

```go
type Set map[string]struct{}

func NewSet(items ...string) Set
func Union(a, b Set) []string
func Intersection(a, b Set) []string
func Difference(a, b Set) []string
func IsSubset(a, b Set) bool
```

With `a := NewSet("go", "rust", "zig")` and `b := NewSet("go", "c", "zig", "c")`:

| Call | Result |
|---|---|
| `len(b)` | `3` (duplicates collapse) |
| `Union(a, b)` | `[c go rust zig]` |
| `Intersection(a, b)` | `[go zig]` |
| `Difference(a, b)` | `[rust]` (in `a`, not in `b`) |
| `Difference(b, a)` | `[c]` |
| `IsSubset(NewSet("go"), a)` | `true` (every element of the first is in the second) |
| `IsSubset(a, b)` | `false` |
| `IsSubset(NewSet(), a)` | `true` |

## Rules
- `NewSet()` with no items returns an empty, **non-nil** set.
- `Union`, `Intersection` and `Difference` return a **sorted** slice (they range over maps, whose order is random). An empty result may be `nil`.
- Every function must accept `nil` sets (a nil map reads as empty) and must **not modify** `a` or `b`.
- Test membership with comma-ok: `_, ok := b[x]`.
