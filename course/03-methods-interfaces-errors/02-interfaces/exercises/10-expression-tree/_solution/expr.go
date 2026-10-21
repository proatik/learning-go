package expr

import (
	"maps"
	"slices"
	"strconv"
)

// Expr is a node of an arithmetic expression tree.
type Expr interface {
	Eval(env map[string]float64) float64
	String() string
}

// Num is a constant.
type Num float64

// Var is a variable whose value comes from the environment.
type Var string

// Add is L + R.
type Add struct{ L, R Expr }

// Mul is L * R.
type Mul struct{ L, R Expr }

// Neg is -X.
type Neg struct{ X Expr }

var (
	_ Expr = Num(0)
	_ Expr = Var("")
	_ Expr = Add{}
	_ Expr = Mul{}
	_ Expr = Neg{}
)

// Eval returns the constant.
func (n Num) Eval(env map[string]float64) float64 { return float64(n) }

// String formats the constant with strconv's shortest 'g' form.
func (n Num) String() string { return strconv.FormatFloat(float64(n), 'g', -1, 64) }

// Eval returns env[v], or 0 if it is missing.
func (v Var) Eval(env map[string]float64) float64 { return env[string(v)] }

// String returns the variable name.
func (v Var) String() string { return string(v) }

// Eval returns L + R.
func (a Add) Eval(env map[string]float64) float64 { return a.L.Eval(env) + a.R.Eval(env) }

// String returns "(L + R)".
func (a Add) String() string { return "(" + a.L.String() + " + " + a.R.String() + ")" }

// Eval returns L * R.
func (m Mul) Eval(env map[string]float64) float64 { return m.L.Eval(env) * m.R.Eval(env) }

// String returns "(L * R)".
func (m Mul) String() string { return "(" + m.L.String() + " * " + m.R.String() + ")" }

// Eval returns -X.
func (n Neg) Eval(env map[string]float64) float64 { return -n.X.Eval(env) }

// String returns "-X".
func (n Neg) String() string { return "-" + n.X.String() }

// isNum reports whether e is the constant v.
func isNum(e Expr, v float64) bool {
	n, ok := e.(Num)
	return ok && float64(n) == v
}

// Simplify returns an equivalent, simplified expression.
func Simplify(e Expr) Expr {
	switch x := e.(type) {
	case Add:
		l, r := Simplify(x.L), Simplify(x.R)
		a, aok := l.(Num)
		b, bok := r.(Num)
		switch {
		case aok && bok:
			return a + b
		case isNum(l, 0):
			return r
		case isNum(r, 0):
			return l
		}
		return Add{l, r}
	case Mul:
		l, r := Simplify(x.L), Simplify(x.R)
		a, aok := l.(Num)
		b, bok := r.(Num)
		switch {
		case aok && bok:
			return a * b
		case isNum(l, 0), isNum(r, 0):
			return Num(0)
		case isNum(l, 1):
			return r
		case isNum(r, 1):
			return l
		}
		return Mul{l, r}
	case Neg:
		switch in := Simplify(x.X).(type) {
		case Num:
			return -in
		case Neg:
			return in.X
		default:
			return Neg{in}
		}
	default: // Num, Var and unknown node types
		return e
	}
}

// collect adds every variable name in e to seen.
func collect(e Expr, seen map[string]bool) {
	switch x := e.(type) {
	case Var:
		seen[string(x)] = true
	case Add:
		collect(x.L, seen)
		collect(x.R, seen)
	case Mul:
		collect(x.L, seen)
		collect(x.R, seen)
	case Neg:
		collect(x.X, seen)
	}
}

// Vars returns the sorted, distinct variable names used in e.
func Vars(e Expr) []string {
	seen := map[string]bool{}
	collect(e, seen)
	return slices.Sorted(maps.Keys(seen))
}

// Size returns the number of nodes in e.
func Size(e Expr) int {
	switch x := e.(type) {
	case Add:
		return 1 + Size(x.L) + Size(x.R)
	case Mul:
		return 1 + Size(x.L) + Size(x.R)
	case Neg:
		return 1 + Size(x.X)
	default:
		return 1
	}
}
