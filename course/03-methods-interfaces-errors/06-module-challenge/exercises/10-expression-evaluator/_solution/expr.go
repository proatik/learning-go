package expr

import (
	"errors"
	"fmt"
	"strconv"
)

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

// Compile-time checks.
var (
	_ Expr = Num(0)
	_ Expr = Neg{}
	_ Expr = Binary{}
)

// Eval returns the number itself.
func (n Num) Eval() (float64, error) { return float64(n), nil }

// String formats the number in the shortest exact form.
func (n Num) String() string { return strconv.FormatFloat(float64(n), 'g', -1, 64) }

// Eval returns -X.
func (n Neg) Eval() (float64, error) {
	v, err := n.X.Eval()
	if err != nil {
		return 0, err
	}
	return -v, nil
}

// String returns "(-X)".
func (n Neg) String() string { return "(-" + n.X.String() + ")" }

// Eval evaluates both sides, then applies Op.
func (b Binary) Eval() (float64, error) {
	l, err := b.L.Eval()
	if err != nil {
		return 0, err
	}
	r, err := b.R.Eval()
	if err != nil {
		return 0, err
	}
	switch b.Op {
	case '+':
		return l + r, nil
	case '-':
		return l - r, nil
	case '*':
		return l * r, nil
	case '/':
		if r == 0 {
			return 0, &DivisionByZeroError{Expr: b.String()}
		}
		return l / r, nil
	}
	return 0, fmt.Errorf("expr: unknown operator %q", b.Op)
}

// String returns "(L op R)".
func (b Binary) String() string {
	return "(" + b.L.String() + " " + string(b.Op) + " " + b.R.String() + ")"
}

// DivisionByZeroError reports the division whose right side was zero.
type DivisionByZeroError struct {
	Expr string // the String() of the failing division
}

// Error returns "expr: division by zero in EXPR".
func (e *DivisionByZeroError) Error() string {
	return "expr: division by zero in " + e.Expr
}

// Is makes errors.Is(err, ErrDivisionByZero) true.
func (e *DivisionByZeroError) Is(target error) bool { return target == ErrDivisionByZero }

// SyntaxError reports where parsing failed.
type SyntaxError struct {
	Pos int // byte offset in the input
	Msg string
}

// Error returns "expr: syntax error at POS: MSG".
func (e *SyntaxError) Error() string {
	return fmt.Sprintf("expr: syntax error at %d: %s", e.Pos, e.Msg)
}

// token is a number, an operator or a parenthesis.
type token struct {
	text string
	pos  int
}

func isDigit(c byte) bool { return '0' <= c && c <= '9' }

// tokenize splits s into tokens, skipping spaces.
func tokenize(s string) ([]token, error) {
	var toks []token
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == ' ' || c == '\t':
			i++
		case isDigit(c) || c == '.':
			start := i
			for i < len(s) && (isDigit(s[i]) || s[i] == '.') {
				i++
			}
			toks = append(toks, token{s[start:i], start})
		case c == '+' || c == '-' || c == '*' || c == '/' || c == '(' || c == ')':
			toks = append(toks, token{string(c), i})
			i++
		default:
			return nil, &SyntaxError{Pos: i, Msg: fmt.Sprintf("unexpected %q", string(c))}
		}
	}
	return toks, nil
}

// parser is a recursive-descent parser over the tokens.
type parser struct {
	toks []token
	i    int
	end  int // len(input), the position reported at end of input
}

func (p *parser) peek() (token, bool) {
	if p.i < len(p.toks) {
		return p.toks[p.i], true
	}
	return token{}, false
}

func (p *parser) errHere() error {
	if t, ok := p.peek(); ok {
		return &SyntaxError{Pos: t.pos, Msg: fmt.Sprintf("unexpected %q", t.text)}
	}
	return &SyntaxError{Pos: p.end, Msg: "unexpected end of input"}
}

// expr = term { ("+" | "-") term }
func (p *parser) expr() (Expr, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || (t.text != "+" && t.text != "-") {
			return left, nil
		}
		p.i++
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		left = Binary{Op: t.text[0], L: left, R: right}
	}
}

// term = unary { ("*" | "/") unary }
func (p *parser) term() (Expr, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
	for {
		t, ok := p.peek()
		if !ok || (t.text != "*" && t.text != "/") {
			return left, nil
		}
		p.i++
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		left = Binary{Op: t.text[0], L: left, R: right}
	}
}

// unary = "-" unary | primary
func (p *parser) unary() (Expr, error) {
	if t, ok := p.peek(); ok && t.text == "-" {
		p.i++
		x, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Neg{X: x}, nil
	}
	return p.primary()
}

// primary = number | "(" expr ")"
func (p *parser) primary() (Expr, error) {
	t, ok := p.peek()
	if !ok {
		return nil, p.errHere()
	}
	switch {
	case t.text == "(":
		p.i++
		e, err := p.expr()
		if err != nil {
			return nil, err
		}
		if c, ok := p.peek(); !ok || c.text != ")" {
			return nil, p.errHere()
		}
		p.i++
		return e, nil
	case isDigit(t.text[0]) || t.text[0] == '.':
		f, err := strconv.ParseFloat(t.text, 64)
		if err != nil {
			return nil, &SyntaxError{Pos: t.pos, Msg: fmt.Sprintf("invalid number %q", t.text)}
		}
		p.i++
		return Num(f), nil
	}
	return nil, p.errHere()
}

// Parse builds an expression tree from s.
func Parse(s string) (Expr, error) {
	toks, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	p := &parser{toks: toks, end: len(s)}
	e, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok := p.peek(); ok {
		return nil, p.errHere()
	}
	return e, nil
}

// Eval parses and evaluates s.
func Eval(s string) (float64, error) {
	e, err := Parse(s)
	if err != nil {
		return 0, err
	}
	return e.Eval()
}
