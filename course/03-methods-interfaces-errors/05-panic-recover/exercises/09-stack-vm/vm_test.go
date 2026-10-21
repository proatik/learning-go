package vm_test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	vm "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/09-stack-vm"
)

// run calls vm.Run and turns an escaped panic into a test error.
func run(t *testing.T, program []string) (result int, err error, ok bool) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run(%q) panicked: %v", program, r)
			ok = false
		}
	}()
	result, err = vm.Run(program)
	return result, err, true
}

func TestVMErrorMethods(t *testing.T) {
	e := &vm.VMError{Index: 3, Instr: "div", Err: vm.ErrStackUnderflow}
	if got, want := e.Error(), `instruction 3 ("div"): stack underflow`; got != want {
		t.Errorf("(%+v).Error() = %q, want %q", *e, got, want)
	}
	if got := e.Unwrap(); got != vm.ErrStackUnderflow {
		t.Errorf("(%+v).Unwrap() = %v, want %v", *e, got, vm.ErrStackUnderflow)
	}
	e = &vm.VMError{Index: 0, Instr: "say \"hi\"", Err: vm.ErrBadInstruction}
	if got, want := e.Error(), `instruction 0 ("say \"hi\""): bad instruction`; got != want {
		t.Errorf("(%+v).Error() = %q, want %q", *e, got, want)
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		program []string
		want    int
	}{
		{[]string{"push 2", "push 3", "add"}, 5},
		{[]string{"push 7", "push 2", "sub"}, 5},
		{[]string{"push 6", "push -7", "mul"}, -42},
		{[]string{"push 7", "push 2", "div"}, 3},
		{[]string{"push -7", "push 2", "div"}, -3},
		{[]string{"push -7", "push 2", "mod"}, -1},
		{[]string{"push 3", "dup", "mul"}, 9},
		{[]string{"push 1", "push 2", "swap", "sub"}, 1},
		{[]string{"push 5", "neg"}, -5},
		{[]string{"push 1", "push 2", "pop"}, 1},
		{[]string{"  push   4  "}, 4},
		{[]string{"push +12"}, 12},
		{[]string{"push 0"}, 0},
		{[]string{"push 1", "push 2", "push 3", "push 4", "push 5", "push 6", "push 7", "push 8"}, 8},
		{[]string{"push 1", "push 2", "push 3", "push 4", "push 5", "push 6", "push 7", "push 8", "add", "push 9"}, 9},
		{[]string{"push 2", "push 3", "push 4", "mul", "add", "dup", "add"}, 28},
	}
	for _, tt := range tests {
		got, err, ok := run(t, tt.program)
		if !ok {
			continue
		}
		if got != tt.want || err != nil {
			t.Errorf("Run(%q) = %d, %v; want %d, <nil>", tt.program, got, err, tt.want)
		}
	}
}

func TestRunErrors(t *testing.T) {
	nine := []string{"push 1", "push 2", "push 3", "push 4", "push 5", "push 6", "push 7", "push 8", "push 9"}
	tests := []struct {
		program []string
		index   int
		instr   string
		cause   error // nil: the cause must be a runtime.Error
		msg     string
	}{
		{nil, 0, "end", vm.ErrStackUnderflow, `instruction 0 ("end"): stack underflow`},
		{[]string{}, 0, "end", vm.ErrStackUnderflow, `instruction 0 ("end"): stack underflow`},
		{[]string{"push 1", "pop"}, 2, "end", vm.ErrStackUnderflow, `instruction 2 ("end"): stack underflow`},
		{[]string{"add"}, 0, "add", vm.ErrStackUnderflow, `instruction 0 ("add"): stack underflow`},
		{[]string{"push 1", "add"}, 1, "add", vm.ErrStackUnderflow, `instruction 1 ("add"): stack underflow`},
		{[]string{"dup"}, 0, "dup", vm.ErrStackUnderflow, `instruction 0 ("dup"): stack underflow`},
		{[]string{"push 1", "swap"}, 1, "swap", vm.ErrStackUnderflow, `instruction 1 ("swap"): stack underflow`},
		{[]string{"neg"}, 0, "neg", vm.ErrStackUnderflow, `instruction 0 ("neg"): stack underflow`},
		{nine, 8, "push 9", vm.ErrStackOverflow, `instruction 8 ("push 9"): stack overflow`},
		{[]string{"push 1", "push 0", "div"}, 2, "div", nil, `instruction 2 ("div"): runtime error: integer divide by zero`},
		{[]string{"push 1", "push 0", "mod"}, 2, "mod", nil, `instruction 2 ("mod"): runtime error: integer divide by zero`},
		{[]string{"push x"}, 0, "push x", vm.ErrBadInstruction, `instruction 0 ("push x"): bad instruction`},
		{[]string{"push 1", "push"}, 1, "push", vm.ErrBadInstruction, `instruction 1 ("push"): bad instruction`},
		{[]string{"push 1 2"}, 0, "push 1 2", vm.ErrBadInstruction, `instruction 0 ("push 1 2"): bad instruction`},
		{[]string{"push 1", "pop 1"}, 1, "pop 1", vm.ErrBadInstruction, `instruction 1 ("pop 1"): bad instruction`},
		{[]string{"push 1", "ADD"}, 1, "ADD", vm.ErrBadInstruction, `instruction 1 ("ADD"): bad instruction`},
		{[]string{"push 1", "jump 0"}, 1, "jump 0", vm.ErrBadInstruction, `instruction 1 ("jump 0"): bad instruction`},
		{[]string{"push 1", ""}, 1, "", vm.ErrBadInstruction, `instruction 1 (""): bad instruction`},
		{[]string{"push 2", "push 3", "sub", "neg", "push 0", "div", "push 1"}, 5, "div", nil, `instruction 5 ("div"): runtime error: integer divide by zero`},
	}
	for _, tt := range tests {
		got, err, ok := run(t, tt.program)
		if !ok {
			continue
		}
		var ve *vm.VMError
		if !errors.As(err, &ve) {
			t.Errorf("Run(%q) = %d, %v; want a *VMError %q", tt.program, got, err, tt.msg)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("Run(%q) error = %q, want %q", tt.program, err.Error(), tt.msg)
		}
		if ve.Index != tt.index || ve.Instr != tt.instr {
			t.Errorf("Run(%q) VMError Index, Instr = %d, %q; want %d, %q", tt.program, ve.Index, ve.Instr, tt.index, tt.instr)
		}
		if tt.cause != nil && !errors.Is(err, tt.cause) {
			t.Errorf("Run(%q) error = %v, want errors.Is(err, %v)", tt.program, err, tt.cause)
		}
		var re runtime.Error
		if tt.cause == nil && !errors.As(err, &re) {
			t.Errorf("Run(%q) error = %v, want it to wrap a runtime.Error", tt.program, err)
		}
		if got != 0 {
			t.Errorf("Run(%q) result = %d alongside an error, want 0", tt.program, got)
		}
	}
}

func TestRunIsReusable(t *testing.T) {
	// Each call starts with an empty stack.
	for i := range 3 {
		got, err, ok := run(t, []string{"push 1", "push 1", "push 1", "push 1", "push 1", "push 1", "push 1", "push 1"})
		if ok && (got != 1 || err != nil) {
			t.Errorf("call %d: Run(8 pushes) = %d, %v; want 1, <nil>", i, got, err)
		}
	}
	if got, err, ok := run(t, []string{"pop"}); ok && err == nil {
		t.Errorf("Run(%q) = %d, <nil>; want a stack underflow error", []string{"pop"}, got)
	}
}

func ExampleRun() {
	for _, p := range [][]string{
		{"push 6", "push 7", "mul"},
		{"push 1", "push 0", "div"},
		{"push 1", "sub"},
	} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("PANIC:", r)
				}
			}()
			fmt.Println(vm.Run(p))
		}()
	}
	// Output:
	// 42 <nil>
	// 0 instruction 2 ("div"): runtime error: integer divide by zero
	// 0 instruction 1 ("sub"): stack underflow
}
