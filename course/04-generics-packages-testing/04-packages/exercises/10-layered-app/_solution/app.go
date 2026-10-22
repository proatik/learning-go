// Package todo wires the to-do application together.
package todo

import (
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/service"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/store"
)

// New returns a to-do service backed by an in-memory store, using now as
// its clock (nil means time.Now).
func New(now func() time.Time) *service.Service {
	return service.New(store.NewMemory(), now)
}
