package txn

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
	// TODO
	return nil
}

// ExecAll executes stmts in order inside a single transaction.
func ExecAll(db DB, stmts []string) error {
	// TODO
	return nil
}
