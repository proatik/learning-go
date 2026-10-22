// Package event holds what bank and audit share, so that neither has to
// import the other.
package event

// TODO: declare the two shared pieces here:
//
//   Event:    a struct with Account string, Kind string and Amount int
//   Recorder: an interface with one method, Record(e Event)
//
// Then make bank depend on Recorder instead of *audit.Log, and make
// audit.Log.Record take an Event instead of a *bank.Account.
