package numeric

// Number is the set of signed integer and float types,
// including named types whose underlying type is one of them.
type Number interface {
	// TODO: replace `any` with a union of ~int, ~int8, ... ~float64.
	// `any` is only here so the starter compiles: it doesn't allow + or <.
	any
}

// Sum returns the total of xs.
func Sum[T Number](xs []T) T {
	// TODO
	var zero T
	return zero
}

// Abs returns the absolute value of x.
func Abs[T Number](x T) T {
	// TODO
	var zero T
	return zero
}

// Scale returns a new slice of the same type as s with every element multiplied by k.
func Scale[S ~[]E, E Number](s S, k E) S {
	// TODO: make(S, len(s)) keeps the caller's slice type
	return nil
}
