package numeric

// Number is the set of signed integer and float types,
// including named types whose underlying type is one of them.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Sum returns the total of xs.
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// Abs returns the absolute value of x.
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Scale returns a new slice of the same type as s with every element multiplied by k.
func Scale[S ~[]E, E Number](s S, k E) S {
	out := make(S, len(s))
	for i, v := range s {
		out[i] = v * k
	}
	return out
}
