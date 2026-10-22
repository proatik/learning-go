package pair

import "fmt"

// Pair holds two values of possibly different types.
type Pair[A, B any] struct {
	First  A
	Second B
}

// New returns the pair (a, b). Unlike a Pair literal, a call can infer A and
// B from its arguments.
func New[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{First: a, Second: b}
}

// Swap returns the pair with its fields exchanged: note the result type
// Pair[B, A].
func (p Pair[A, B]) Swap() Pair[B, A] {
	return Pair[B, A]{First: p.Second, Second: p.First}
}

// String formats the pair as "(First, Second)".
func (p Pair[A, B]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

// Zip pairs as[i] with bs[i], stopping at the shorter slice.
func Zip[A, B any](as []A, bs []B) []Pair[A, B] {
	n := min(len(as), len(bs))
	out := make([]Pair[A, B], n)
	for i := range n {
		out[i] = New(as[i], bs[i])
	}
	return out
}

// Unzip splits pairs back into a slice of firsts and a slice of seconds.
func Unzip[A, B any](ps []Pair[A, B]) ([]A, []B) {
	as := make([]A, len(ps))
	bs := make([]B, len(ps))
	for i, p := range ps {
		as[i], bs[i] = p.First, p.Second
	}
	return as, bs
}
