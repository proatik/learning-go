# Two Sum

Find two **different positions** whose values add up to `target`:

```go
func TwoSum(nums []int, target int) (i, j int, ok bool)
```

Return indices `i < j` with `nums[i] + nums[j] == target` and `ok == true`. If no such pair exists, return `0, 0, false`.

| `nums` | `target` | Result |
|---|---|---|
| `[2 7 11 15]` | `9` | `0, 1, true` |
| `[3 2 4]` | `6` | `1, 2, true` (not `0, 0`: one element can't be used twice) |
| `[3 3]` | `6` | `0, 1, true` |
| `[1 5 1 5]` | `6` | `0, 1, true` |
| `[-3 4 3 90]` | `0` | `0, 2, true` |
| `[1 2 3]` | `7` | `0, 0, false` |
| `[]` | `0` | `0, 0, false` |

## Which pair?
When several pairs work, return the one with the **smallest `j`**; for that `j`, the **smallest `i`**. A single left-to-right pass produces exactly this.

## Rules
- **One pass**, O(n): no nested loops.
- Keep a `map[int]int` from *value* → *index where it was first seen*.
- For each `nums[j]`, look up `target - nums[j]` with comma-ok **before** storing `nums[j]` itself.
- `nums` must not be modified.
