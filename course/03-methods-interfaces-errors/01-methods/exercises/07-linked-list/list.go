package list

// List is a node of a singly linked list of ints. A nil *List is the empty
// list, and every method accepts a nil receiver.
type List struct {
	Val  int
	Next *List
}

// FromSlice builds a list holding vals in order (nil for no values).
func FromSlice(vals []int) *List {
	// TODO
	return nil
}

// Len returns the number of nodes.
func (l *List) Len() int {
	// TODO
	return 0
}

// Sum returns the sum of all values.
func (l *List) Sum() int {
	// TODO
	return 0
}

// Contains reports whether v is in the list.
func (l *List) Contains(v int) bool {
	// TODO
	return false
}

// Values returns the values in order, or nil for the empty list.
func (l *List) Values() []int {
	// TODO
	return nil
}

// String formats the list like "[1 -> 2 -> 3]" ("[]" when empty).
func (l *List) String() string {
	// TODO
	return ""
}

// Push adds v at the front and returns the new head.
func (l *List) Push(v int) *List {
	// TODO
	return nil
}

// Append adds v at the end and returns the head of the list.
func (l *List) Append(v int) *List {
	// TODO
	return nil
}

// Reverse reverses the list in place and returns the new head.
func (l *List) Reverse() *List {
	// TODO
	return nil
}

// RemoveAll unlinks every node holding v and returns the new head.
func (l *List) RemoveAll(v int) *List {
	// TODO
	return nil
}
