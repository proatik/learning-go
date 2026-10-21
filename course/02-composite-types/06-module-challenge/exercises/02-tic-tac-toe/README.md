# Tic-Tac-Toe Referee

Parse a tic-tac-toe board into a fixed-size array and decide the state of the game.

```go
func Parse(rows []string) ([3][3]rune, error)
func Winner(b [3][3]rune) (rune, error)
func Status(b [3][3]rune) (string, error)
```

Cells are `'X'`, `'O'` or `'.'` (empty). `b[row][col]`, row 0 at the top. **X always moves first.**

## Examples

```go
b, _ := Parse([]string{"XXX", "OO.", "..."})
Winner(b) // 'X', nil
Status(b) // "X wins", nil

b, _ = Parse([]string{"XO.", "...", "..."})
Winner(b) // '.', nil
Status(b) // "X to move", nil
```

## Rules

**Parse**
- Exactly 3 rows, each exactly 3 runes, every rune one of `X`, `O`, `.`. Otherwise return the zero array and an error.
- Parse checks only shape and characters, not whether the position is legal.

**Winner** returns `'X'` or `'O'` if that player has three in a row (row, column or either diagonal), otherwise `'.'`. It returns `0` and an error if the board is **impossible**:
- a cell holds anything other than `X`, `O`, `.`;
- the number of X's minus the number of O's is not `0` or `1`;
- both players have a line;
- X has a line but X's count is not O's count + 1 (O kept playing after X won);
- O has a line but the counts are not equal (X kept playing after O won).

A single player completing two lines with one move (e.g. an X in the centre finishing a row and a diagonal) is legal.

**Status** returns, for a legal board, one of `"X wins"`, `"O wins"`, `"draw"` (full board, no winner), `"X to move"`, `"O to move"`; for an impossible board, `""` and the error from `Winner`.
