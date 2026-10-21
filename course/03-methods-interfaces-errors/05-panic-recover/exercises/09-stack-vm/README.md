# Stack VM

Build a tiny stack machine whose internals **panic** on failure, while the exported `Run` **recovers** and returns a typed error that names the failing instruction — the `encoding/json` style "panic inside, error outside".

```go
const MaxStack = 8

var (
	ErrStackUnderflow = errors.New("stack underflow")
	ErrStackOverflow  = errors.New("stack overflow")
	ErrBadInstruction = errors.New("bad instruction")
)

type VMError struct {
	Index int    // index of the failing instruction (len(program) for the final check)
	Instr string // the instruction text as given ("end" for the final check)
	Err   error  // the cause
}

func (e *VMError) Error() string // instruction <Index> (<Instr as %q>): <Err>
func (e *VMError) Unwrap() error // Err

func Run(program []string) (result int, err error)
```

| Instruction | Effect |
|---|---|
| `push N` | push the integer `N` |
| `pop` | discard the top value |
| `dup` | push a copy of the top value |
| `swap` | exchange the top two values |
| `neg` | replace the top value `a` by `-a` |
| `add` `sub` `mul` `div` `mod` | pop `b`, pop `a`, push `a+b`, `a-b`, `a*b`, `a/b`, `a%b` |

| Call | Result |
|---|---|
| `Run([]string{"push 2", "push 3", "add"})` | `5, nil` |
| `Run([]string{"push 1", "push 2", "swap", "sub"})` | `1, nil` |
| `Run([]string{"push 1", "add"})` | `instruction 1 ("add"): stack underflow` |
| `Run([]string{"push 1", "push 0", "div"})` | `instruction 2 ("div"): runtime error: integer divide by zero` |
| `Run([]string{"push x"})` | `instruction 0 ("push x"): bad instruction` |
| `Run(nil)` | `instruction 0 ("end"): stack underflow` |

## Rules
- Split each instruction with `strings.Fields`. `push` takes exactly one `strconv.Atoi`-parsable argument; every other op takes none. Anything else (unknown or upper-case op, wrong argument count, empty line) → `ErrBadInstruction`.
- Popping an empty stack → `ErrStackUnderflow`. Pushing when the stack already holds `MaxStack` values → `ErrStackOverflow`.
- `div` / `mod` compute `a / b` and `a % b` **directly**: a zero divisor causes Go's run-time panic, which `Run` reports as a `VMError` whose `Err` is that `runtime.Error` (so `errors.As(err, &re)` with `var re runtime.Error` works).
- After the last instruction, `result` is the top of the stack; an empty stack → `VMError{Index: len(program), Instr: "end", Err: ErrStackUnderflow}`.
- Every failure is returned as a `*VMError` with `result == 0`; `errors.Is(err, ErrStackUnderflow)` etc. must work.
- Suggested design: private helpers (`push`, `pop`, `exec`) **panic** with the sentinel error; one deferred function in `Run` recovers, wraps any `error` value into a `*VMError`, and **re-panics** values that aren't errors.
- `Run` must never let a panic escape for any input.
