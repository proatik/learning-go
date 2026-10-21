package txn

import (
	"errors"
	"fmt"
)

// Tx is an open database transaction.
type Tx interface {
	Exec(stmt string) error
	Commit() error
	Rollback() error
}

// DB starts transactions.
type DB interface {
	Begin() (Tx, error)
}

// WithTx runs fn inside a transaction: it commits if fn returns nil and
// rolls back if fn returns an error or panics (re-panicking afterwards).
func WithTx(db DB, fn func(tx Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	committing := false
	defer func() {
		if committing {
			return // Commit was attempted: never roll back afterwards
		}
		if r := recover(); r != nil {
			_ = tx.Rollback() // best effort; the panic is the real problem
			panic(r)          // re-panic with the original value
		}
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Join(err, fmt.Errorf("rollback: %w", rerr))
		}
	}()

	if err := fn(tx); err != nil {
		return err // the deferred function rolls back
	}
	committing = true
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// ExecAll executes stmts in order inside a single transaction.
func ExecAll(db DB, stmts []string) error {
	return WithTx(db, func(tx Tx) error {
		for _, stmt := range stmts {
			if err := tx.Exec(stmt); err != nil {
				return fmt.Errorf("exec %q: %w", stmt, err)
			}
		}
		return nil
	})
}
