package mixsum

// Valuer is any type that knows its numeric value.
type Valuer interface {
	Value() float64
}

// Sum adds every numeric element of vals (int, int64, float32, float64,
// non-nil *int and Valuer) and reports how many elements were skipped.
func Sum(vals []any) (total float64, skipped int) {
	// TODO
	return 0, 0
}

// Strings returns the string elements of vals in order.
func Strings(vals []any) []string {
	// TODO
	return nil
}
