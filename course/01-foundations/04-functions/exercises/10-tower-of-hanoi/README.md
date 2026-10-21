# Tower of Hanoi

Move a tower of `n` disks from peg `from` to peg `to`, using peg `via` as a spare. Only one disk moves at a time, and a larger disk may never sit on a smaller one.

```go
func Hanoi(n int, from, to, via string, move func(from, to string)) int
```

- Call `move(src, dst)` once for **every** disk move, in order.
- Return the total number of moves.
- If `move` is `nil`, still solve the puzzle and return the count, just don't call it.
- `n <= 0` means no disks: zero moves.

The recursive idea for `n` disks:
1. Move the top `n-1` disks from `from` to `via` (using `to` as the spare).
2. Move the largest disk from `from` to `to`.
3. Move the `n-1` disks from `via` to `to` (using `from` as the spare).

| Call | Moves passed to `move` | Returns |
|---|---|---|
| `Hanoi(1, "A", "C", "B", f)` | `A→C` | `1` |
| `Hanoi(2, "A", "C", "B", f)` | `A→B`, `A→C`, `B→C` | `3` |
| `Hanoi(3, "A", "C", "B", f)` | 7 moves, starting `A→C`, `A→B`, `C→B`, `A→C` | `7` |
| `Hanoi(20, "A", "C", "B", nil)` | (none) | `1048575` |

## Example: recording moves with a closure

```go
log := ""
count := Hanoi(2, "A", "C", "B", func(from, to string) {
	log += from + "->" + to + " "
})
// log == "A->B A->C B->C ", count == 3
```

The tests check that every move is legal and that the tower ends up on `to`.
