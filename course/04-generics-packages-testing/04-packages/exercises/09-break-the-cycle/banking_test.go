package banking_test

import (
	"errors"
	"go/parser"
	"go/token"
	"slices"
	"strconv"
	"strings"
	"testing"

	banking "learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/audit"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/bank"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/event"
)

// The compiler checks that *audit.Log satisfies event.Recorder.
var _ event.Recorder = (*audit.Log)(nil)

func TestOpen(t *testing.T) {
	acc, log := banking.Open("acc-1")
	if acc == nil || log == nil {
		t.Fatalf("Open(%q) = %v, %v, want an account and a log", "acc-1", acc, log)
	}
	steps := []struct {
		op   string
		n    int
		want error
	}{
		{"deposit", 100, nil},
		{"withdraw", 30, nil},
		{"withdraw", 500, bank.ErrFunds},
		{"deposit", 0, bank.ErrAmount},
		{"withdraw", -5, bank.ErrAmount},
		{"withdraw", 70, nil},
		{"deposit", 15, nil},
	}
	for _, s := range steps {
		var err error
		if s.op == "deposit" {
			err = acc.Deposit(s.n)
		} else {
			err = acc.Withdraw(s.n)
		}
		if !errors.Is(err, s.want) {
			t.Errorf("%s(%d) error = %v, want %v", s.op, s.n, err, s.want)
		}
	}
	if got := acc.Balance(); got != 15 {
		t.Errorf("Balance() = %d, want 15", got)
	}
	if got := acc.ID(); got != "acc-1" {
		t.Errorf("ID() = %q, want %q", got, "acc-1")
	}
	want := []string{"acc-1 deposit 100", "acc-1 withdraw 30", "acc-1 rejected 500", "acc-1 withdraw 70", "acc-1 deposit 15"}
	if got := log.Lines(); !slices.Equal(got, want) {
		t.Errorf("Lines() = %q, want %q", got, want)
	}
	for kind, want := range map[string]int{"deposit": 115, "withdraw": 100, "rejected": 500, "other": 0} {
		if got := log.Total(kind); got != want {
			t.Errorf("Total(%q) = %d, want %d", kind, got, want)
		}
	}
}

// fake is a test double: bank works with any event.Recorder.
type fake struct{ got []event.Event }

func (f *fake) Record(e event.Event) { f.got = append(f.got, e) }

func TestBankWithFakeRecorder(t *testing.T) {
	f := &fake{}
	acc := bank.New("x", f)
	acc.Deposit(10)
	acc.Withdraw(20)
	want := []event.Event{{Account: "x", Kind: "deposit", Amount: 10}, {Account: "x", Kind: "rejected", Amount: 20}}
	if !slices.Equal(f.got, want) {
		t.Errorf("recorded events = %+v, want %+v", f.got, want)
	}
}

func TestBankWithoutRecorder(t *testing.T) {
	acc := bank.New("quiet", nil)
	if err := acc.Deposit(5); err != nil {
		t.Errorf("Deposit(5) with nil recorder error = %v, want <nil>", err)
	}
	if got := acc.Balance(); got != 5 {
		t.Errorf("Balance() = %d, want 5", got)
	}
}

func TestLogDirectly(t *testing.T) {
	var log audit.Log
	if got := log.Lines(); len(got) != 0 {
		t.Errorf("empty Lines() = %q, want []", got)
	}
	log.Record(event.Event{Account: "ü-1", Kind: "deposit", Amount: 7})
	lines := log.Lines()
	if !slices.Equal(lines, []string{"ü-1 deposit 7"}) {
		t.Fatalf("Lines() = %q, want %q", lines, []string{"ü-1 deposit 7"})
	}
	lines[0] = "changed"
	if got := log.Lines()[0]; got != "ü-1 deposit 7" {
		t.Errorf("Lines()[0] after editing the copy = %q, want %q", got, "ü-1 deposit 7")
	}
}

func imports(t *testing.T, file string) []string {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse %s: %v", file, err)
	}
	var paths []string
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		paths = append(paths, p)
	}
	return paths
}

func TestDependencies(t *testing.T) {
	tests := []struct {
		file      string
		forbidden string
	}{
		{"bank/bank.go", "/audit"},
		{"audit/audit.go", "/bank"},
		{"event/event.go", "/bank"},
		{"event/event.go", "/audit"},
	}
	for _, tt := range tests {
		for _, p := range imports(t, tt.file) {
			if strings.HasSuffix(p, "09-break-the-cycle"+tt.forbidden) {
				t.Errorf("%s imports = %q, want no %s", tt.file, p, tt.forbidden[1:])
			}
		}
	}
}
