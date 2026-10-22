// Package ids generates sequential string IDs.
package ids

import "strconv"

// counter is shared by every caller in the program.
// TODO: remove it; package-level mutable state is what this exercise replaces.
var counter int

// Next returns the next ID: "1", "2", …
//
// TODO: keep Next for old callers, but make it use a package-level default
// Generator and add a "Deprecated:" paragraph to this comment.
func Next() string {
	counter++
	return strconv.Itoa(counter)
}

// Generator produces IDs prefix+"1", prefix+"2", … Each Generator counts on
// its own. The zero value is ready to use with an empty prefix.
type Generator struct {
	// TODO: add the fields you need
}

// New returns a Generator whose IDs start with prefix.
func New(prefix string) *Generator {
	// TODO
	return &Generator{}
}

// Next returns the generator's next ID.
func (g *Generator) Next() string {
	// TODO
	return ""
}

// Issued reports how many IDs the generator has produced.
func (g *Generator) Issued() int {
	// TODO
	return 0
}
