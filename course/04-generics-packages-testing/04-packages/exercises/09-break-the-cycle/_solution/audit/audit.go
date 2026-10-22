// Package audit keeps a log of account activity.
package audit

import (
	"fmt"
	"slices"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/event"
)

// Log is an in-memory audit log. The zero value is an empty log.
type Log struct {
	lines  []string
	totals map[string]int
}

var _ event.Recorder = (*Log)(nil)

// Record appends "<account> <kind> <amount>" to the log.
func (l *Log) Record(e event.Event) {
	l.lines = append(l.lines, fmt.Sprintf("%s %s %d", e.Account, e.Kind, e.Amount))
	if l.totals == nil {
		l.totals = map[string]int{}
	}
	l.totals[e.Kind] += e.Amount
}

// Lines returns a copy of the log lines, oldest first.
func (l *Log) Lines() []string { return slices.Clone(l.lines) }

// Total returns the summed amount of all entries of the given kind.
func (l *Log) Total(kind string) int { return l.totals[kind] }
