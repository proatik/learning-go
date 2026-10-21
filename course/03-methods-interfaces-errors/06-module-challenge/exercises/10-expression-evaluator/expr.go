package expr

import "errors"

// ErrDivisionByZero matches every *DivisionByZeroError via errors.Is.
var ErrDivisionByZero = errors.New("expr: division by zero")

// Expr is a node of an expression tree.
type Expr interface {
	Eval() (float64, error)
	String() string
}

// Num is a number literal.
type Num float64

// Neg is unary minus: -X.
type Neg struct {
	X Expr
}

// Binary is L Op R, where Op is one of + - * /.
type Binary struct {
	Op   byte
	L, R Expr
}

// Eval returns the number itself.
func (n Num) Eval() (float64, error) {
	// TODO
	return 0, nil
}

// String formats the number in the shortest exact form.
func (n Num) String() string {
	// TODO
	return ""
}

// Eval returns -X.
func (n Neg) Eval() (float64, error) {
	// TODO
	return 0, nil
}

// String returns "(-X)".
func (n Neg) String() string {
	// TODO
	return ""
}

// Eval evaluates both sides, then applies Op.
func (b Binary) Eval() (float64, error) {
	// TODO
	return 0, nil
}

// String returns "(L op R)".
func (b Binary) String() string {
	// TODO
	return ""
}

// DivisionByZeroError reports the division whose right side was zero.
type DivisionByZeroError struct {
	Expr string // the String() of the failing division
}

// Error returns "expr: division by zero in EXPR".
func (e *DivisionByZeroError) Error() string {
	// TODO
	return ""
}

// Is makes errors.Is(err, ErrDivisionByZero) true.
func (e *DivisionByZeroError) Is(target error) bool {
	// TODO
	return false
}

// SyntaxError reports where parsing failed.
type SyntaxError struct {
	Pos int // byte offset in the input
	Msg string
}

// Error returns "expr: syntax error at POS: MSG".
func (e *SyntaxError) Error() string {
	// TODO
	return ""
}

// Parse builds an expression tree from s.
func Parse(s string) (Expr, error) {
	// TODO: tokenize, then parse with one function per precedence level
	return nil, nil
}

// Eval parses and evaluates s.
func Eval(s string) (float64, error) {
	// TODO
	return 0, nil
}
