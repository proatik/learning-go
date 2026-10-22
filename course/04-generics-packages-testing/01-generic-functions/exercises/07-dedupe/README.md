# Dedupe, Preserving Order

Remove duplicates while keeping the **first** occurrence in place, by value or by a key you choose.

```go
func Dedupe[T comparable](xs []T) []T
func DedupeBy[T any, K comparable](xs []T, key func(T) K) []T
func Intersect[T comparable](a, b []T) []T
```

## Examples

| Call | Result |
|---|---|
| `Dedupe([]int{3, 1, 3, 2, 1})` | `[3 1 2]` |
| `Dedupe([]string{"b", "a", "b"})` | `[b a]` |
| `DedupeBy(users, func(u User) string { return u.Email })` | the first user for each email, in order |
| `DedupeBy([]string{"Go", "go", "GO", "rust"}, strings.ToLower)` | `[Go rust]` |
| `Intersect([]int{1, 2, 3, 2}, []int{2, 3, 4})` | `[2 3]` |
| `Intersect([]int{1}, nil)` | `[]` |

## Rules
- `Dedupe`: keep each value the first time you see it, in the original order.
- `DedupeBy`: two elements are duplicates when `key` returns the same value for them. Keep the **first** element for each key. Call `key` exactly once per element.
  - `T` is `any`: the elements themselves may not be comparable (e.g. structs with a slice field). Only the key must be.
- `Intersect`: the values of `a` that also appear in `b`, in `a`'s order, without duplicates.
- All three return a **new** slice (length `0` when nothing is kept; nil is fine) and never modify the input.
- Use a `map[K]struct{}` or `map[K]bool` as a "seen" set. Don't sort: that would lose the order.
