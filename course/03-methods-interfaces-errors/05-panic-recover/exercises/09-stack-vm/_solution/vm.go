package vm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
	return fmt.Sprintf("instruction %d (%q): %v", e.Index, e.Instr, e.Err)
}

// Unwrap returns the cause.
func (e *VMError) Unwrap() error {
	return e.Err
}

// machine is the VM state. Its methods panic on failure; Run recovers.
type machine struct {
	stack []int
}

func (m *machine) push(v int) {
	if len(m.stack) == MaxStack {
		panic(ErrStackOverflow)
	}
	m.stack = append(m.stack, v)
}

func (m *machine) pop() int {
	if len(m.stack) == 0 {
		panic(ErrStackUnderflow)
	}
	v := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]
	return v
}

// exec runs one instruction.
func (m *machine) exec(instr string) {
	f := strings.Fields(instr)
	if len(f) == 2 && f[0] == "push" {
		n, err := strconv.Atoi(f[1])
		if err != nil {
			panic(ErrBadInstruction)
		}
		m.push(n)
		return
	}
	if len(f) != 1 {
		panic(ErrBadInstruction)
	}
	switch f[0] {
	case "pop":
		m.pop()
	case "dup":
		v := m.pop()
		m.push(v)
		m.push(v)
	case "swap":
		b, a := m.pop(), m.pop()
		m.push(b)
		m.push(a)
	case "neg":
		m.push(-m.pop())
	case "add", "sub", "mul", "div", "mod":
		b, a := m.pop(), m.pop()
		switch f[0] {
		case "add":
			m.push(a + b)
		case "sub":
			m.push(a - b)
		case "mul":
			m.push(a * b)
		case "div":
			m.push(a / b) // b == 0: run-time panic, recovered by Run
		case "mod":
			m.push(a % b)
		}
	default:
		panic(ErrBadInstruction)
	}
}

// Run executes program and returns the value on top of the stack. Any
// failure, including Go run-time panics, is returned as a *VMError.
func Run(program []string) (result int, err error) {
	var m machine
	index, instr := 0, ""
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		cause, ok := r.(error) // our sentinels and runtime.Error values
		if !ok {
			panic(r) // not something we know how to report
		}
		result, err = 0, &VMError{Index: index, Instr: instr, Err: cause}
	}()
	for index, instr = range program {
		m.exec(instr)
	}
	index, instr = len(program), "end"
	return m.pop(), nil
}
