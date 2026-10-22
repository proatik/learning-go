// Package banking wires accounts to their audit log. It is the only package
// that imports both bank and audit.
package banking

import (
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/audit"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/bank"
)

// Open returns a new account whose activity is recorded in a new audit log.
func Open(id string) (*bank.Account, *audit.Log) {
	log := &audit.Log{}
	return bank.New(id, log), log
}
