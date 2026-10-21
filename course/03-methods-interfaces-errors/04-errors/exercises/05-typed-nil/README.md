# Fix the Typed-Nil Trap

The starter code **compiles and looks right**, but valid input still returns a non-nil `error`. Find and fix every typed-nil bug.

```go
type CheckError struct {
	Name    string
	Problem string
}

func (e *CheckError) Error() string // "<name>: <problem>"

func CheckHost(host string) error
func CheckPort(port int) error
func CheckAll(host string, port int) error           // errors.Join of both checks
func FirstFailure(results []*CheckError) error       // first non-nil result, or nil
```

## Examples

| Call | Should return | Starter returns |
|---|---|---|
| `CheckHost("go.dev")` | `nil` | `(*CheckError)(nil)`, which is **`!= nil`** |
| `CheckHost("")` | `&CheckError{"host", "must not be empty"}` | ✅ |
| `CheckPort(8080)` | `nil` | `(*CheckError)(nil)` |
| `CheckPort(0)` | `&CheckError{"port", "must be 1-65535"}` | ✅ |
| `CheckAll("go.dev", 443)` | `nil` | a non-nil joined error |
| `CheckAll("", 0)` | both errors joined | joined, but contains the typed nils |
| `FirstFailure([]*CheckError{nil, e, nil})` | `e` | ✅ |
| `FirstFailure([]*CheckError{nil, nil})` | `nil` | `(*CheckError)(nil)` |

## Rules
- **Don't change** the exported signatures (they must keep returning `error`). You may change the unexported helper `checkPort`.
- Every function returns a literal `nil` on success.
- Failure values and messages stay exactly as in the starter: `host: must not be empty`, `port: must be 1-65535`.
- `CheckHost` fails for a host that is empty after `strings.TrimSpace`. `CheckPort` fails outside `1..65535`.
