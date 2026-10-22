// Package event holds what bank and audit share, so that neither has to
// import the other.
package event

// Event describes one piece of account activity.
type Event struct {
	Account string
	Kind    string // "deposit", "withdraw" or "rejected"
	Amount  int
}

// Recorder receives account events.
type Recorder interface {
	Record(e Event)
}
