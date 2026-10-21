# Byte Sizes with `iota`

## Tasks
1. In `size.go`, rewrite the constant block so that

| Constant | Value |
|---|---|
| `KB` | 1 024 (2¹⁰) |
| `MB` | 1 048 576 (2²⁰) |
| `GB` | 2³⁰ |
| `TB` | 2⁴⁰ |

using **one** `iota` expression of the form `1 << (10 * iota)`, skipping the `iota == 0` line with `_`.

2. Implement:

```go
func HumanSize(n int64) string
```

Use the largest unit that is `<= n` and print one decimal; below 1 KB print whole bytes:

| Call | Result |
|---|---|
| `HumanSize(0)` | `"0 B"` |
| `HumanSize(1023)` | `"1023 B"` |
| `HumanSize(1024)` | `"1.0 KB"` |
| `HumanSize(1536)` | `"1.5 KB"` |
| `HumanSize(5 * MB)` | `"5.0 MB"` |
| `HumanSize(3*GB + 512*MB)` | `"3.5 GB"` |
| `HumanSize(2 * TB)` | `"2.0 TB"` |

## Rules
- Use `fmt.Sprintf` with `%d` and `%.1f`. Convert to `float64` before dividing.
- Check the biggest unit first with an `if … else if … else` chain:

```go
if n >= TB {
	return ...
} else if n >= GB {
	return ...
}
```

- `n` is never negative.
