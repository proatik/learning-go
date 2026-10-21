package txn_test

import (
	"errors"
	"runtime"
	"slices"
	"testing"

	txn "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/10-transactions"
)

var (
	errBegin    = errors.New("pool exhausted")
	errFn       = errors.New("fn failed")
	errCommit   = errors.New("disk full")
	errRollback = errors.New("conn lost")
	errExec     = errors.New("constraint violated")
)

// fakeDB records every call in log.
type fakeDB struct {
	log         []string
	beginErr    error
	commitErr   error
	rollbackErr error
	failExec    string // Exec of this statement returns errExec
	panicExec   string // Exec of this statement panics
}

type fakeTx struct{ db *fakeDB }

func (d *fakeDB) Begin() (txn.Tx, error) {
	d.log = append(d.log, "begin")
	if d.beginErr != nil {
		return nil, d.beginErr
	}
	return fakeTx{d}, nil
}

func (t fakeTx) Exec(stmt string) error {
	t.db.log = append(t.db.log, "exec "+stmt)
	if stmt == t.db.panicExec {
		panic("driver crashed")
	}
	if stmt == t.db.failExec {
		return errExec
	}
	return nil
}

func (t fakeTx) Commit() error {
	t.db.log = append(t.db.log, "commit")
	return t.db.commitErr
}

func (t fakeTx) Rollback() error {
	t.db.log = append(t.db.log, "rollback")
	return t.db.rollbackErr
}

// catch runs f and returns the value it panicked with (nil if it did not).
func catch(f func()) (r any) {
	defer func() { r = recover() }()
	f()
	return nil
}

func execA(tx txn.Tx) error { return tx.Exec("a") }

func TestWithTx(t *testing.T) {
	tests := []struct {
		name    string
		db      fakeDB
		fn      func(tx txn.Tx) error
		wantErr []error // each must match errors.Is; nil means success
		wantMsg string
		wantLog []string
	}{
		{
			name: "success commits",
			fn: func(tx txn.Tx) error {
				_ = tx.Exec("a")
				return tx.Exec("b")
			},
			wantLog: []string{"begin", "exec a", "exec b", "commit"},
		},
		{
			name:    "empty fn commits",
			fn:      func(txn.Tx) error { return nil },
			wantLog: []string{"begin", "commit"},
		},
		{
			name: "fn error rolls back",
			fn: func(tx txn.Tx) error {
				_ = tx.Exec("a")
				return errFn
			},
			wantErr: []error{errFn},
			wantMsg: "fn failed",
			wantLog: []string{"begin", "exec a", "rollback"},
		},
		{
			name: "fn error and rollback error are joined",
			db:   fakeDB{rollbackErr: errRollback},
			fn: func(tx txn.Tx) error {
				_ = tx.Exec("a")
				return errFn
			},
			wantErr: []error{errFn, errRollback},
			wantMsg: "fn failed\nrollback: conn lost",
			wantLog: []string{"begin", "exec a", "rollback"},
		},
		{
			name:    "commit error: no rollback afterwards",
			db:      fakeDB{commitErr: errCommit, rollbackErr: errRollback},
			fn:      execA,
			wantErr: []error{errCommit},
			wantMsg: "commit: disk full",
			wantLog: []string{"begin", "exec a", "commit"},
		},
		{
			name:    "begin error: fn not called",
			db:      fakeDB{beginErr: errBegin},
			fn:      execA,
			wantErr: []error{errBegin},
			wantMsg: "begin: pool exhausted",
			wantLog: []string{"begin"},
		},
	}
	for _, tt := range tests {
		db := tt.db
		var err error
		if r := catch(func() { err = txn.WithTx(&db, tt.fn) }); r != nil {
			t.Errorf("%s: WithTx panicked: %v", tt.name, r)
			continue
		}
		if tt.wantErr == nil && err != nil {
			t.Errorf("%s: WithTx error = %v, want <nil>", tt.name, err)
		}
		for _, target := range tt.wantErr {
			if !errors.Is(err, target) {
				t.Errorf("%s: WithTx error = %v, want %v", tt.name, err, target)
			}
		}
		if tt.wantErr != nil && err != nil && err.Error() != tt.wantMsg {
			t.Errorf("%s: WithTx error = %q, want %q", tt.name, err.Error(), tt.wantMsg)
		}
		if !slices.Equal(db.log, tt.wantLog) {
			t.Errorf("%s: WithTx calls = %q, want %q", tt.name, db.log, tt.wantLog)
		}
	}
}

func TestWithTxReturnsCallbackErrorUnchanged(t *testing.T) {
	db := &fakeDB{}
	var err error
	if r := catch(func() { err = txn.WithTx(db, func(txn.Tx) error { return errFn }) }); r != nil {
		t.Fatalf("WithTx panicked: %v", r)
	}
	if err != errFn {
		t.Errorf("WithTx error = %v (%T), want exactly errFn", err, err)
	}
}

func TestWithTxPanicRollsBackAndRepanics(t *testing.T) {
	tests := []struct {
		name    string
		db      fakeDB
		fn      func(tx txn.Tx) error
		check   func(r any) bool
		want    string
		wantLog []string
	}{
		{
			name: "string panic",
			fn: func(tx txn.Tx) error {
				_ = tx.Exec("a")
				panic("boom")
			},
			check:   func(r any) bool { return r == "boom" },
			want:    "boom",
			wantLog: []string{"begin", "exec a", "rollback"},
		},
		{
			name: "run-time panic",
			fn: func(txn.Tx) error {
				var m map[string]int
				m["x"]++
				return nil
			},
			check: func(r any) bool {
				_, ok := r.(runtime.Error)
				return ok
			},
			want:    "a runtime.Error",
			wantLog: []string{"begin", "rollback"},
		},
		{
			name: "panic and rollback error: still re-panics",
			db:   fakeDB{rollbackErr: errRollback},
			fn: func(tx txn.Tx) error {
				panic(errFn)
			},
			check:   func(r any) bool { return r == errFn },
			want:    "errFn",
			wantLog: []string{"begin", "rollback"},
		},
	}
	for _, tt := range tests {
		db := tt.db
		r := catch(func() { _ = txn.WithTx(&db, tt.fn) })
		if r == nil {
			t.Errorf("%s: WithTx did not panic, want a re-panic with %s", tt.name, tt.want)
		} else if !tt.check(r) {
			t.Errorf("%s: WithTx panicked with %v, want %s", tt.name, r, tt.want)
		}
		if !slices.Equal(db.log, tt.wantLog) {
			t.Errorf("%s: WithTx calls = %q, want %q", tt.name, db.log, tt.wantLog)
		}
	}
}

func TestExecAll(t *testing.T) {
	tests := []struct {
		name    string
		db      fakeDB
		stmts   []string
		wantErr error
		wantMsg string
		wantLog []string
	}{
		{
			name:    "all statements",
			stmts:   []string{"a", "b", "c"},
			wantLog: []string{"begin", "exec a", "exec b", "exec c", "commit"},
		},
		{
			name:    "no statements",
			wantLog: []string{"begin", "commit"},
		},
		{
			name:    "failing statement rolls back",
			db:      fakeDB{failExec: "b"},
			stmts:   []string{"a", "b", "c"},
			wantErr: errExec,
			wantMsg: `exec "b": constraint violated`,
			wantLog: []string{"begin", "exec a", "exec b", "rollback"},
		},
		{
			name:    "commit fails",
			db:      fakeDB{commitErr: errCommit},
			stmts:   []string{"a"},
			wantErr: errCommit,
			wantMsg: "commit: disk full",
			wantLog: []string{"begin", "exec a", "commit"},
		},
		{
			name:    "begin fails",
			db:      fakeDB{beginErr: errBegin},
			stmts:   []string{"a"},
			wantErr: errBegin,
			wantMsg: "begin: pool exhausted",
			wantLog: []string{"begin"},
		},
	}
	for _, tt := range tests {
		db := tt.db
		var err error
		if r := catch(func() { err = txn.ExecAll(&db, tt.stmts) }); r != nil {
			t.Errorf("%s: ExecAll(%q) panicked: %v", tt.name, tt.stmts, r)
			continue
		}
		switch {
		case tt.wantErr == nil && err != nil:
			t.Errorf("%s: ExecAll(%q) error = %v, want <nil>", tt.name, tt.stmts, err)
		case tt.wantErr != nil && !errors.Is(err, tt.wantErr):
			t.Errorf("%s: ExecAll(%q) error = %v, want %v", tt.name, tt.stmts, err, tt.wantErr)
		case tt.wantErr != nil && err.Error() != tt.wantMsg:
			t.Errorf("%s: ExecAll(%q) error = %q, want %q", tt.name, tt.stmts, err.Error(), tt.wantMsg)
		}
		if !slices.Equal(db.log, tt.wantLog) {
			t.Errorf("%s: ExecAll(%q) calls = %q, want %q", tt.name, tt.stmts, db.log, tt.wantLog)
		}
	}
}

func TestExecAllPanicPropagates(t *testing.T) {
	db := &fakeDB{panicExec: "b"}
	r := catch(func() { _ = txn.ExecAll(db, []string{"a", "b", "c"}) })
	if r != "driver crashed" {
		t.Errorf("ExecAll with a panicking Exec: panic = %v, want %q", r, "driver crashed")
	}
	want := []string{"begin", "exec a", "exec b", "rollback"}
	if !slices.Equal(db.log, want) {
		t.Errorf("ExecAll with a panicking Exec: calls = %q, want %q", db.log, want)
	}
}
