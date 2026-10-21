package store_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	store "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/07-close-capture"
)

var (
	errOpen  = errors.New("no such store")
	errWrite = errors.New("disk full")
	errClose = errors.New("flush failed")
)

// fake records every call in a shared log.
type fake struct {
	name      string
	log       *[]string
	failWrite string // Write of this line returns errWrite
	panicOn   string // Write of this line panics
	closeErr  error
}

func (f *fake) Write(line string) error {
	*f.log = append(*f.log, "write "+f.name+" "+line)
	if line == f.panicOn {
		panic("disk on fire")
	}
	if line == f.failWrite {
		return errWrite
	}
	return nil
}

func (f *fake) Close() error {
	*f.log = append(*f.log, "close "+f.name)
	return f.closeErr
}

// catch runs f and returns the value it panicked with (nil if it did not).
func catch(f func()) (r any) {
	defer func() { r = recover() }()
	f()
	return nil
}

func TestSave(t *testing.T) {
	tests := []struct {
		name    string
		fake    fake
		openErr error
		lines   []string
		wantErr error  // nil means success
		wantMsg string // exact error text when wantErr != nil
		wantLog []string
	}{
		{
			name:    "success",
			lines:   []string{"a", "b"},
			wantLog: []string{"write s a", "write s b", "close s"},
		},
		{
			name:    "no lines",
			wantLog: []string{"close s"},
		},
		{
			name:    "close fails",
			fake:    fake{closeErr: errClose},
			lines:   []string{"a"},
			wantErr: errClose,
			wantMsg: "close: flush failed",
			wantLog: []string{"write s a", "close s"},
		},
		{
			name:    "write fails",
			fake:    fake{failWrite: "b"},
			lines:   []string{"a", "b", "c"},
			wantErr: errWrite,
			wantMsg: `write "b": disk full`,
			wantLog: []string{"write s a", "write s b", "close s"},
		},
		{
			name:    "write and close fail: first error wins",
			fake:    fake{failWrite: "a", closeErr: errClose},
			lines:   []string{"a", "b"},
			wantErr: errWrite,
			wantMsg: `write "a": disk full`,
			wantLog: []string{"write s a", "close s"},
		},
		{
			name:    "open fails",
			openErr: errOpen,
			lines:   []string{"a"},
			wantErr: errOpen,
			wantMsg: "open: no such store",
			wantLog: nil,
		},
	}
	for _, tt := range tests {
		var log []string
		f := tt.fake
		f.name, f.log = "s", &log
		open := func() (store.Store, error) {
			if tt.openErr != nil {
				return nil, tt.openErr
			}
			return &f, nil
		}
		var err error
		if r := catch(func() { err = store.Save(open, tt.lines) }); r != nil {
			t.Errorf("%s: Save panicked: %v", tt.name, r)
			continue
		}
		switch {
		case tt.wantErr == nil && err != nil:
			t.Errorf("%s: Save(%q) error = %v, want <nil>", tt.name, tt.lines, err)
		case tt.wantErr != nil && !errors.Is(err, tt.wantErr):
			t.Errorf("%s: Save(%q) error = %v, want %v", tt.name, tt.lines, err, tt.wantErr)
		case tt.wantErr != nil && err.Error() != tt.wantMsg:
			t.Errorf("%s: Save(%q) error = %q, want %q", tt.name, tt.lines, err.Error(), tt.wantMsg)
		}
		if !slices.Equal(log, tt.wantLog) {
			t.Errorf("%s: Save(%q) calls = %q, want %q", tt.name, tt.lines, log, tt.wantLog)
		}
	}
}

func TestSaveClosesOnPanic(t *testing.T) {
	var log []string
	f := &fake{name: "s", log: &log, panicOn: "b"}
	open := func() (store.Store, error) { return f, nil }
	r := catch(func() { _ = store.Save(open, []string{"a", "b", "c"}) })
	if r != "disk on fire" {
		t.Errorf("Save with a panicking Write: panic = %v, want %q (Save must not recover)", r, "disk on fire")
	}
	want := []string{"write s a", "write s b", "close s"}
	if !slices.Equal(log, want) {
		t.Errorf("Save with a panicking Write: calls = %q, want %q", log, want)
	}
}

func TestSaveBatches(t *testing.T) {
	tests := []struct {
		name      string
		batches   [][]string
		failWrite map[int]string
		closeErr  map[int]error
		openErr   map[int]error
		wantErr   error
		wantMsg   string
		wantLog   []string
	}{
		{
			name:    "all succeed: each store closed before the next opens",
			batches: [][]string{{"a"}, {"b", "c"}, {}},
			wantLog: []string{
				"open 0", "write 0 a", "close 0",
				"open 1", "write 1 b", "write 1 c", "close 1",
				"open 2", "close 2",
			},
		},
		{
			name:    "no batches",
			wantLog: nil,
		},
		{
			name:     "close of batch 1 fails",
			batches:  [][]string{{"a"}, {"b"}, {"c"}},
			closeErr: map[int]error{1: errClose},
			wantErr:  errClose,
			wantMsg:  "batch 1: close: flush failed",
			wantLog:  []string{"open 0", "write 0 a", "close 0", "open 1", "write 1 b", "close 1"},
		},
		{
			name:      "write in batch 0 fails",
			batches:   [][]string{{"a", "x"}, {"b"}},
			failWrite: map[int]string{0: "x"},
			wantErr:   errWrite,
			wantMsg:   `batch 0: write "x": disk full`,
			wantLog:   []string{"open 0", "write 0 a", "write 0 x", "close 0"},
		},
		{
			name:    "open of batch 2 fails",
			batches: [][]string{{"a"}, {"b"}, {"c"}},
			openErr: map[int]error{2: errOpen},
			wantErr: errOpen,
			wantMsg: "batch 2: open: no such store",
			wantLog: []string{
				"open 0", "write 0 a", "close 0",
				"open 1", "write 1 b", "close 1",
				"open 2",
			},
		},
	}
	for _, tt := range tests {
		var log []string
		open := func(i int) (store.Store, error) {
			log = append(log, fmt.Sprintf("open %d", i))
			if err := tt.openErr[i]; err != nil {
				return nil, err
			}
			return &fake{name: fmt.Sprint(i), log: &log, failWrite: tt.failWrite[i], closeErr: tt.closeErr[i]}, nil
		}
		var err error
		if r := catch(func() { err = store.SaveBatches(open, tt.batches) }); r != nil {
			t.Errorf("%s: SaveBatches panicked: %v", tt.name, r)
			continue
		}
		switch {
		case tt.wantErr == nil && err != nil:
			t.Errorf("%s: SaveBatches(%q) error = %v, want <nil>", tt.name, tt.batches, err)
		case tt.wantErr != nil && !errors.Is(err, tt.wantErr):
			t.Errorf("%s: SaveBatches(%q) error = %v, want %v", tt.name, tt.batches, err, tt.wantErr)
		case tt.wantErr != nil && err.Error() != tt.wantMsg:
			t.Errorf("%s: SaveBatches(%q) error = %q, want %q", tt.name, tt.batches, err.Error(), tt.wantMsg)
		}
		if !slices.Equal(log, tt.wantLog) {
			t.Errorf("%s: SaveBatches(%q) calls = %q, want %q", tt.name, tt.batches, log, tt.wantLog)
		}
	}
}
