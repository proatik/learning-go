package contracts

import "fmt"

// Namer has a name.
type Namer interface{ Name() string }

// Labeler has a display label.
type Labeler interface{ Label() string }

// Counter counts calls to Inc.
type Counter interface {
	Inc()
	Count() int
}

// Resetter can be reset to its zero state.
type Resetter interface{ Reset() }

// ResettableCounter is a Counter that can also be reset.
type ResettableCounter interface {
	Counter
	Resetter
}

// Tag is a short text tag.
type Tag struct{ Text string }

// Tally is a counter whose zero value is ready to use.
type Tally struct{ n int }

// Celsius is a temperature in degrees Celsius.
type Celsius float64

// Name returns the tag text.
func (t Tag) Name() string { return t.Text }

// Label returns the tag text prefixed with "#".
func (t Tag) Label() string { return "#" + t.Text }

// Inc adds one to the tally.
func (t *Tally) Inc() { t.n++ }

// Count returns the number of Inc calls since the last Reset.
func (t *Tally) Count() int { return t.n }

// Reset sets the tally back to zero.
func (t *Tally) Reset() { t.n = 0 }

// Label returns the temperature with one decimal, e.g. "21.5°C".
func (c Celsius) Label() string { return fmt.Sprintf("%.1f°C", float64(c)) }

var (
	_ Namer             = Tag{}
	_ Labeler           = Tag{}
	_ ResettableCounter = (*Tally)(nil)
	_ Labeler           = Celsius(0)
)
