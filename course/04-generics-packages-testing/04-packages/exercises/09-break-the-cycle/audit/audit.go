// Package audit keeps a log of account activity.
package audit

import (
	"fmt"
	"slices"

	// TODO: audit only needs an account ID, a kind and an amount. Importing
	// bank for that is what closes the cycle.
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/bank"
)

// Log is an in-memory audit log. The zero value is an empty log.
type Log struct {
	lines  []string
	totals map[string]int
}

// Record appends "<account> <kind> <amount>" to the log.
// TODO: take an event.Event instead, so *Log satisfies event.Recorder.
func (l *Log) Record(a *bank.Account, kind string, amount int) {
	l.lines = append(l.lines, fmt.Sprintf("%s %s %d", a.ID(), kind, amount))
	if l.totals == nil {
		l.totals = map[string]int{}
	}
	l.totals[kind] += amount
}

// Lines returns a copy of the log lines, oldest first.
func (l *Log) Lines() []string { return slices.Clone(l.lines) }

// Total returns the summed amount of all entries of the given kind.
func (l *Log) Total(kind string) int { return l.totals[kind] }
