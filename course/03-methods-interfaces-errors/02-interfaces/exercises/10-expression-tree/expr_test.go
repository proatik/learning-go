package expr_test

import (
	"math"
	"slices"
	"testing"

	expr "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/10-expression-tree"
)

type (
	num = expr.Num
	v   = expr.Var
	add = expr.Add
	mul = expr.Mul
	neg = expr.Neg
)

// pow is an Expr node the package knows nothing about.
type pow struct{ B, E expr.Expr }

func (p pow) Eval(env map[string]float64) float64 { return math.Pow(p.B.Eval(env), p.E.Eval(env)) }
func (p pow) String() string                      { return "(" + p.B.String() + " ^ " + p.E.String() + ")" }

// sample is ((2 * x) + -3).
var sample = add{mul{num(2), v("x")}, neg{num(3)}}

func TestString(t *testing.T) {
	tests := []struct {
		e    expr.Expr
		want string
	}{
		{num(3), "3"},
		{num(2.5), "2.5"},
		{num(-1), "-1"},
		{num(1e21), "1e+21"},
		{v("y"), "y"},
		{v("größe"), "größe"},
		{add{num(1), num(2)}, "(1 + 2)"},
		{mul{v("a"), v("b")}, "(a * b)"},
		{neg{v("x")}, "-x"},
		{neg{add{v("x"), num(1)}}, "-(x + 1)"},
		{sample, "((2 * x) + -3)"},
		{add{pow{v("x"), num(2)}, num(1)}, "((x ^ 2) + 1)"},
	}
	for _, tt := range tests {
		if got := tt.e.String(); got != tt.want {
			t.Errorf("%#v.String() = %q, want %q", tt.e, got, tt.want)
		}
	}
}

func TestEval(t *testing.T) {
	env := map[string]float64{"x": 5, "y": -2}
	tests := []struct {
		e    expr.Expr
		env  map[string]float64
		want float64
	}{
		{num(4), nil, 4},
		{v("x"), env, 5},
		{v("missing"), env, 0},
		{v("x"), nil, 0},
		{sample, env, 7},
		{sample, nil, -3},
		{add{v("x"), v("y")}, env, 3},
		{mul{v("x"), v("y")}, env, -10},
		{neg{neg{v("y")}}, env, -2},
		{mul{add{num(1), num(2)}, neg{v("x")}}, env, -15},
		{add{pow{v("x"), num(2)}, num(1)}, env, 26},
		{mul{num(0.5), num(0.5)}, nil, 0.25},
	}
	for _, tt := range tests {
		if got := tt.e.Eval(tt.env); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("%s.Eval(%v) = %v, want %v", tt.e, tt.env, got, tt.want)
		}
	}
}

func TestSimplify(t *testing.T) {
	tests := []struct {
		in, want expr.Expr
	}{
		{num(7), num(7)},
		{v("x"), v("x")},
		{add{num(1), num(2)}, num(3)},
		{add{v("x"), num(0)}, v("x")},
		{add{num(0), v("x")}, v("x")},
		{add{v("x"), v("y")}, add{v("x"), v("y")}},
		{mul{num(2), num(3)}, num(6)},
		{mul{v("x"), num(0)}, num(0)},
		{mul{num(0), add{v("x"), v("y")}}, num(0)},
		{mul{num(1), v("x")}, v("x")},
		{mul{v("x"), num(1)}, v("x")},
		{mul{add{num(1), num(2)}, v("x")}, mul{num(3), v("x")}},
		{neg{num(3)}, num(-3)},
		{neg{neg{v("x")}}, v("x")},
		{neg{neg{neg{v("x")}}}, neg{v("x")}},
		{neg{v("x")}, neg{v("x")}},
		{sample, add{mul{num(2), v("x")}, num(-3)}},
		{add{mul{v("x"), add{num(0), num(1)}}, mul{v("y"), num(0)}}, v("x")}, // (x*(0+1)) + (y*0) → x
		{mul{add{num(1), neg{num(1)}}, v("z")}, num(0)},                      // (1 + -1) * z → 0
		{neg{add{neg{v("a")}, num(0)}}, v("a")},                              // -(-a + 0) → a
		{add{pow{add{num(1), num(1)}, num(2)}, num(0)}, pow{add{num(1), num(1)}, num(2)}}, // unknown node untouched
		{pow{v("x"), num(1)}, pow{v("x"), num(1)}},
	}
	for _, tt := range tests {
		got := expr.Simplify(tt.in)
		if got == nil {
			t.Errorf("Simplify(%s) = nil, want %s", tt.in, tt.want)
			continue
		}
		if got != tt.want {
			t.Errorf("Simplify(%s) = %s (%T), want %s (%T)", tt.in, got, got, tt.want, tt.want)
		}
	}
}

func TestVars(t *testing.T) {
	tests := []struct {
		e    expr.Expr
		want []string
	}{
		{num(1), nil},
		{v("x"), []string{"x"}},
		{sample, []string{"x"}},
		{add{v("y"), mul{v("x"), v("y")}}, []string{"x", "y"}},
		{neg{add{v("b"), neg{v("a")}}}, []string{"a", "b"}},
		{mul{v("Z"), v("a")}, []string{"Z", "a"}},
		{pow{v("hidden"), num(2)}, nil},
	}
	for _, tt := range tests {
		if got := expr.Vars(tt.e); !slices.Equal(got, tt.want) {
			t.Errorf("Vars(%s) = %q, want %q", tt.e, got, tt.want)
		}
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		e    expr.Expr
		want int
	}{
		{num(1), 1},
		{v("x"), 1},
		{neg{v("x")}, 2},
		{add{num(1), num(2)}, 3},
		{sample, 6},
		{neg{neg{neg{v("x")}}}, 4},
		{pow{v("x"), num(2)}, 1},
		{mul{pow{v("x"), num(2)}, num(3)}, 3},
	}
	for _, tt := range tests {
		if got := expr.Size(tt.e); got != tt.want {
			t.Errorf("Size(%s) = %d, want %d", tt.e, got, tt.want)
		}
	}
}
