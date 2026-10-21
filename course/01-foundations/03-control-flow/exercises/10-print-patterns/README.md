# Printing Patterns

Nested loops: the outer loop chooses the **row**, the inner loop prints that row's **characters**.

Implement:

```go
func PrintTriangle(n int) // left-aligned right triangle of '*'
func PrintPyramid(n int)  // centred pyramid of '*'
func PrintTable(n int)    // n×n multiplication table
```

All three print **nothing** when `n < 1`. No line may end with a space.

## `PrintTriangle(4)`
Row `i` (1-based) has `i` stars.
```
*
**
***
****
```

## `PrintPyramid(4)`
Row `i` has `n-i` spaces, then `2*i-1` stars — and nothing after the stars.
```
   *
  ***
 *****
*******
```

## `PrintTable(4)`
Row `i`, column `j` holds `i*j`, each printed right-aligned in a **4-character** field (`%4d`).
```
   1   2   3   4
   2   4   6   8
   3   6   9  12
   4   8  12  16
```

## Rules
- Print one character or cell at a time with `fmt.Print` / `fmt.Printf`, then end each row with `fmt.Println()`.
- Don't use `strings.Repeat` — build every row with loops (`for range k` repeats something `k` times).
