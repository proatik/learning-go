# Duplicates with a Set

Use a map as a **set** of things already seen.

```go
func HasDuplicate(nums []int) bool
func Dedupe(words []string) []string
```

## 1. `HasDuplicate`

Report whether any value appears more than once.

| Input | Output |
|---|---|
| `[]int{1, 2, 3, 1}` | `true` |
| `[]int{1, 2, 3}` | `false` |
| `[]int{}` or `nil` | `false` |

## 2. `Dedupe`

Return the words with duplicates removed, keeping the **first** occurrence of each, in the original order.

| Input | Output |
|---|---|
| `{"go", "is", "go", "fun", "is"}` | `{"go", "is", "fun"}` |
| `{"a", "A", "a"}` | `{"a", "A"}` (case-sensitive) |
| `{}` or `nil` | an empty slice (`nil` is fine) |

## Rules
- One pass each: O(n), no nested loops.
- `HasDuplicate` uses a `map[int]bool` and returns as soon as it finds a repeat.
- `Dedupe` uses a `map[string]struct{}` and comma-ok to test membership.
- `Dedupe` must **not modify** `words`; build a new slice.
