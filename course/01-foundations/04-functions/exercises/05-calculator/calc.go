package calc

// BinOp is a binary operation on two ints.
type BinOp func(a, b int) int

// Lookup returns the operation for the operator op ("+", "-", "*", "/" or "%").
// It returns an error for any other operator.
func Lookup(op string) (BinOp, error) {
	// TODO
	return nil, nil
}

// Calc applies the operator op to a and b.
// It returns an error for an unknown operator or for division by zero.
func Calc(op string, a, b int) (int, error) {
	// TODO
	return 0, nil
}
