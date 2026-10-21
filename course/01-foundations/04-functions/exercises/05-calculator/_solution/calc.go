package calc

import (
	"errors"
	"fmt"
)

// BinOp is a binary operation on two ints.
type BinOp func(a, b int) int

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }

// Lookup returns the operation for the operator op ("+", "-", "*", "/" or "%").
// It returns an error for any other operator.
func Lookup(op string) (BinOp, error) {
	switch op {
	case "+":
		return add, nil // a named function is a value too
	case "-":
		return sub, nil
	case "*":
		return func(a, b int) int { return a * b }, nil
	case "/":
		return func(a, b int) int { return a / b }, nil
	case "%":
		return func(a, b int) int { return a % b }, nil
	default:
		return nil, fmt.Errorf("unknown operator %q", op)
	}
}

// Calc applies the operator op to a and b.
// It returns an error for an unknown operator or for division by zero.
func Calc(op string, a, b int) (int, error) {
	f, err := Lookup(op)
	if err != nil {
		return 0, err
	}
	if (op == "/" || op == "%") && b == 0 {
		return 0, errors.New("division by zero")
	}
	return f(a, b), nil
}
