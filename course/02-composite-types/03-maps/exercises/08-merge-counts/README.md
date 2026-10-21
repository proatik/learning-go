# Merge, Prune and Print Counts

Three everyday operations on `map[string]int` tallies, each with a different rule about **mutation**.

```go
func MergeCounts(a, b map[string]int) map[string]int
func Prune(m map[string]int, minCount int)
func PrintCounts(m map[string]int)
```

## 1. `MergeCounts`: returns a new map, inputs untouched

Sum the counts of both maps.

| `a` | `b` | Result |
|---|---|---|
| `map[go:2 rust:1]` | `map[go:3 zig:4]` | `map[go:5 rust:1 zig:4]` |
| `map[x:1]` | `nil` | `map[x:1]` |
| `nil` | `nil` | empty, non-nil map |

Callers keep using `a` and `b` afterwards: changing them is a bug. Remember that maps are reference-like, so the result must be a **different** map (writing to it must not affect `a`).

## 2. `Prune`: modifies `m` in place

Delete every entry whose count is **less than** `minCount`.

| `m` before | `minCount` | `m` after |
|---|---|---|
| `map[a:1 b:5 c:3]` | `3` | `map[b:5 c:3]` |
| `map[a:1]` | `0` | `map[a:1]` |

## 3. `PrintCounts`: deterministic output

Print one `key: value` line per entry, **sorted by key**. An empty (or nil) map prints `(none)`.

```
go: 5
rust: 1
zig: 4
```

## Rules
- Use the `maps` package: `maps.Clone` / `maps.Copy` in `MergeCounts`, `maps.DeleteFunc` in `Prune`, and `slices.Sorted(maps.Keys(m))` in `PrintCounts`.
- Don't print the map with `fmt.Println(m)`: produce the exact format above.
