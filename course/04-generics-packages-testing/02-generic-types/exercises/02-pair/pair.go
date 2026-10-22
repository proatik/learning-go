package pair

// Pair holds two values of possibly different types.
type Pair[A, B any] struct {
	First  A
	Second B
}

// New returns the pair (a, b). Unlike a Pair literal, a call can infer A and
// B from its arguments.
func New[A, B any](a A, b B) Pair[A, B] {
	// TODO
	return Pair[A, B]{}
}

// Swap returns the pair with its fields exchanged: note the result type
// Pair[B, A].
func (p Pair[A, B]) Swap() Pair[B, A] {
	// TODO
	return Pair[B, A]{}
}

// String formats the pair as "(First, Second)".
func (p Pair[A, B]) String() string {
	// TODO
	return ""
}

// Zip pairs as[i] with bs[i], stopping at the shorter slice.
func Zip[A, B any](as []A, bs []B) []Pair[A, B] {
	// TODO
	return nil
}

// Unzip splits pairs back into a slice of firsts and a slice of seconds.
func Unzip[A, B any](ps []Pair[A, B]) ([]A, []B) {
	// TODO
	return nil, nil
}
