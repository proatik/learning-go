# Print a Profile with Printf

Implement:

```go
func PrintProfile(name string, year int, active bool)
```

It **prints** (to standard output) exactly three lines. For `PrintProfile("Gopher", 2009, true)`:

```
Name:   Gopher
Born:   2009
Active: true
```

## Rules
- Use `fmt.Printf` with the verbs `%s` (string), `%d` (integer) and `%t` (bool).
- `Printf` does **not** add a newline — end each format with `\n`.
- Note the spaces after the colons so the values line up.

## Verb cheat-sheet
| Verb | Meaning |
|---|---|
| `%v` | default format of any value |
| `%s` | string |
| `%d` | integer (base 10) |
| `%t` | bool |
| `%q` | quoted string |
| `%T` | the **type** of the value |
