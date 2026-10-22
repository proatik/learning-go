// Package ids generates sequential string IDs.
package ids

import "strconv"

// defaultGen backs the deprecated package-level Next.
var defaultGen Generator

// Next returns the next ID from the package's default generator: "1", "2", …
//
// Deprecated: Use a Generator from New instead; the default generator is
// shared by every caller in the program.
func Next() string {
	return defaultGen.Next()
}

// Generator produces IDs prefix+"1", prefix+"2", … Each Generator counts on
// its own. The zero value is ready to use with an empty prefix.
type Generator struct {
	prefix string
	n      int
}

// New returns a Generator whose IDs start with prefix.
func New(prefix string) *Generator {
	return &Generator{prefix: prefix}
}

// Next returns the generator's next ID.
func (g *Generator) Next() string {
	g.n++
	return g.prefix + strconv.Itoa(g.n)
}

// Issued reports how many IDs the generator has produced.
func (g *Generator) Issued() int {
	return g.n
}
