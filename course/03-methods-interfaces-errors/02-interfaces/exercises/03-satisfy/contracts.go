package contracts

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

// TODO: Tag must satisfy Namer. Is this method's name right?
func (t Tag) name() string { return t.Text }

// TODO: Tag (the value, not only *Tag) must satisfy Labeler.
func (t *Tag) Label() string { return "#" + t.Text }

// TODO: this changes a copy of the tally, not the tally itself.
func (t Tally) Inc() { t.n++ }

// TODO: Counter wants Count() int.
func (t *Tally) Count() int64 { return int64(t.n) }

// TODO: *Tally also needs Reset, and Celsius needs Label.
