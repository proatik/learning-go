package vm

import "errors"

// MaxStack is the maximum number of values the stack can hold.
const MaxStack = 8

// Sentinel causes reported inside a *VMError.
var (
	ErrStackUnderflow = errors.New("stack underflow")
	ErrStackOverflow  = errors.New("stack overflow")
	ErrBadInstruction = errors.New("bad instruction")
)

// VMError reports which instruction failed and why.
type VMError struct {
	Index int    // index of the failing instruction (len(program) for the final check)
	Instr string // the instruction text as given ("end" for the final check)
	Err   error  // the cause
}

// Error returns `instruction <Index> ("<Instr>"): <Err>`.
func (e *VMError) Error() string {
	// TODO
	return ""
}

// Unwrap returns the cause.
func (e *VMError) Unwrap() error {
	// TODO
	return nil
}

// Run executes program and returns the value on top of the stack. Any
// failure, including Go run-time panics, is returned as a *VMError.
func Run(program []string) (result int, err error) {
	// TODO
	return 0, nil
}
