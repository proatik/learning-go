package expr_test

import (
	"errors"
	"math"
	"testing"

	expr "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/10-expression-evaluator"
)

func parse(t *testing.T, s string) expr.Expr {
	t.Helper()
	e, err := expr.Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v, want nil", s, err)
	}
	if e == nil {
		t.Fatalf("Parse(%q) = nil, nil, want an expression", s)
	}
	return e
}

func TestNodes(t *testing.T) {
	// Build trees by hand: the node types are plain values that satisfy Expr.
	tests := []struct {
		e    expr.Expr
		str  string
		want float64
	}{
		{expr.Num(42), "42", 42},
		{expr.Num(2.5), "2.5", 2.5},
		{expr.Num(0.1), "0.1", 0.1},
		{expr.Neg{X: expr.Num(3)}, "(-3)", -3},
		{expr.Neg{X: expr.Neg{X: expr.Num(3)}}, "(-(-3))", 3},
		{expr.Binary{Op: '+', L: expr.Num(1), R: expr.Num(2)}, "(1 + 2)", 3},
		{expr.Binary{Op: '-', L: expr.Num(1), R: expr.Num(2)}, "(1 - 2)", -1},
		{expr.Binary{Op: '*', L: expr.Num(1.5), R: expr.Num(4)}, "(1.5 * 4)", 6},
		{expr.Binary{Op: '/', L: expr.Num(7), R: expr.Num(2)}, "(7 / 2)", 3.5},
		{
			expr.Binary{Op: '*', L: expr.Binary{Op: '+', L: expr.Num(1), R: expr.Num(2)}, R: expr.Neg{X: expr.Num(4)}},
			"((1 + 2) * (-4))", -12,
		},
	}
	for _, tt := range tests {
		if got := tt.e.String(); got != tt.str {
			t.Errorf("%#v.String() = %q, want %q", tt.e, got, tt.str)
		}
		got, err := tt.e.Eval()
		if err != nil || got != tt.want {
			t.Errorf("%s.Eval() = %v, %v, want %v, nil", tt.str, got, err, tt.want)
		}
	}
}

func TestParseStructure(t *testing.T) {
	tests := []struct {
		in   string
		want string // fully parenthesised String()
	}{
		{"1", "1"},
		{"  7  ", "7"},
		{"3.25", "3.25"},
		{".5", "0.5"},
		{"1+2", "(1 + 2)"},
		{"1 + 2 * 3", "(1 + (2 * 3))"},   // * binds tighter than +
		{"(1 + 2) * 3", "((1 + 2) * 3)"}, // parentheses override
		{"1 - 2 - 3", "((1 - 2) - 3)"},   // left-associative
		{"8 / 4 / 2", "((8 / 4) / 2)"},   // left-associative
		{"2 * 3 + 4 * 5", "((2 * 3) + (4 * 5))"},
		{"-3", "(-3)"},
		{"--3", "(-(-3))"},
		{"2 * -3", "(2 * (-3))"},
		{"-(1 + 2)", "(-(1 + 2))"},
		{"-2 * 3", "((-2) * 3)"}, // unary minus binds tighter than *
		{"((((5))))", "5"},
		{"1-(-1)", "(1 - (-1))"},
		{"\t10 /\t( 3 - 1 )", "(10 / (3 - 1))"},
	}
	for _, tt := range tests {
		if got := parse(t, tt.in).String(); got != tt.want {
			t.Errorf("Parse(%q).String() = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseNodeTypes(t *testing.T) {
	e := parse(t, "1 + -2")
	b, ok := e.(expr.Binary)
	if !ok {
		t.Fatalf("Parse(%q) = %T, want expr.Binary (a value, not a pointer)", "1 + -2", e)
	}
	if b.Op != '+' {
		t.Errorf("Parse(%q).Op = %q, want '+'", "1 + -2", b.Op)
	}
	if n, ok := b.L.(expr.Num); !ok || n != 1 {
		t.Errorf("Parse(%q).L = %#v, want expr.Num(1)", "1 + -2", b.L)
	}
	if _, ok := b.R.(expr.Neg); !ok {
		t.Errorf("Parse(%q).R = %T, want expr.Neg", "1 + -2", b.R)
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{"0", 0},
		{"1 + 2 * 3", 7},
		{"(1 + 2) * 3", 9},
		{"10 - 4 - 3", 3},
		{"2 * (3 + 4) * 5", 70},
		{"1 / 4", 0.25},
		{"100 / 10 / 5", 2},
		{"-(2 + 3) * -2", 10},
		{"0.1 + 0.2", 0.30000000000000004}, // float64 rounding, not 0.3
		{"3.5 * 2 - .5", 6.5},
		{"((2))", 2},
		{"0 / 5", 0},
	}
	for _, tt := range tests {
		got, err := expr.Eval(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("Eval(%q) = %v, %v, want %v, nil", tt.in, got, err, tt.want)
		}
	}
	if got, err := expr.Eval("1 / 3"); err != nil || math.Abs(got-1.0/3) > 1e-15 {
		t.Errorf("Eval(%q) = %v, %v, want %v, nil", "1 / 3", got, err, 1.0/3)
	}
}

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		in   string
		expr string // the failing division
	}{
		{"1 / 0", "(1 / 0)"},
		{"1 / (2 - 2)", "(1 / (2 - 2))"},
		{"5 + 3 / (1 * 0)", "(3 / (1 * 0))"},
		{"(1 / 0) / 0", "(1 / 0)"}, // the inner error is returned unchanged
		{"0 / 0", "(0 / 0)"},
		{"2 / -0", "(2 / (-0))"},
	}
	for _, tt := range tests {
		got, err := expr.Eval(tt.in)
		if err == nil {
			t.Errorf("Eval(%q) = %v, nil, want a division-by-zero error", tt.in, got)
			continue
		}
		if !errors.Is(err, expr.ErrDivisionByZero) {
			t.Errorf("Eval(%q) error = %v, want %v", tt.in, err, expr.ErrDivisionByZero)
		}
		de, ok := errors.AsType[*expr.DivisionByZeroError](err)
		if !ok {
			t.Errorf("Eval(%q) error = %v (%T), want *DivisionByZeroError", tt.in, err, err)
			continue
		}
		if de.Expr != tt.expr {
			t.Errorf("Eval(%q) error Expr = %q, want %q", tt.in, de.Expr, tt.expr)
		}
		if want := "expr: division by zero in " + tt.expr; err.Error() != want {
			t.Errorf("Eval(%q) error = %q, want %q", tt.in, err.Error(), want)
		}
		if got != 0 {
			t.Errorf("Eval(%q) value = %v, want 0 with an error", tt.in, got)
		}
	}
	if errors.Is(&expr.SyntaxError{}, expr.ErrDivisionByZero) {
		t.Errorf("a SyntaxError must not match ErrDivisionByZero")
	}
}

func TestSyntaxErrors(t *testing.T) {
	tests := []struct {
		in  string
		pos int
		msg string
	}{
		{"", 0, "unexpected end of input"},
		{"   ", 3, "unexpected end of input"},
		{"1 +", 3, "unexpected end of input"},
		{"(1 + 2", 6, "unexpected end of input"},
		{"-", 1, "unexpected end of input"},
		{"1 2", 2, `unexpected "2"`},
		{"1 + * 2", 4, `unexpected "*"`},
		{")", 0, `unexpected ")"`},
		{"(1 + 2))", 7, `unexpected ")"`},
		{"()", 1, `unexpected ")"`},
		{"2 (3)", 2, `unexpected "("`},
		{"1 + $", 4, `unexpected "$"`},
		{"x", 0, `unexpected "x"`},
		{"1.2.3 + 1", 0, `invalid number "1.2.3"`},
		{"4 * .", 4, `invalid number "."`},
		{"+1", 0, `unexpected "+"`}, // no unary plus
	}
	for _, tt := range tests {
		e, err := expr.Parse(tt.in)
		se, ok := errors.AsType[*expr.SyntaxError](err)
		if !ok {
			t.Errorf("Parse(%q) = %v, %v, want a *SyntaxError", tt.in, e, err)
			continue
		}
		if se.Pos != tt.pos || se.Msg != tt.msg {
			t.Errorf("Parse(%q) error = {Pos:%d Msg:%s}, want {Pos:%d Msg:%s}", tt.in, se.Pos, se.Msg, tt.pos, tt.msg)
		}
		if e != nil {
			t.Errorf("Parse(%q) expression = %v, want nil with an error", tt.in, e)
		}
		_, evalErr := expr.Eval(tt.in)
		if es, ok := errors.AsType[*expr.SyntaxError](evalErr); !ok || *es != *se {
			t.Errorf("Eval(%q) error = %v, want %v", tt.in, evalErr, err)
		}
	}
	err := &expr.SyntaxError{Pos: 4, Msg: `unexpected "*"`}
	if want := `expr: syntax error at 4: unexpected "*"`; err.Error() != want {
		t.Errorf("SyntaxError.Error() = %q, want %q", err.Error(), want)
	}
}
