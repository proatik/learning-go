package lookup

import (
	"errors"
	"runtime"
	"testing"
)

// noPanic runs f and reports a panic that escaped instead of crashing.
func noPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked: %v", name, r)
		}
	}()
	f()
}

// checkErr verifies an error's exact text and that it wraps a runtime.Error.
func checkErr(t *testing.T, call string, err error, want string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s error = <nil>, want %q", call, want)
		return
	}
	if err.Error() != want {
		t.Errorf("%s error = %q, want %q", call, err.Error(), want)
	}
	var re runtime.Error
	if !errors.As(err, &re) {
		t.Errorf("%s error = %v, want it to wrap the runtime.Error with %%w", call, err)
	}
}

func TestSafeIndex(t *testing.T) {
	s := []string{"a", "b", "ü"}
	for i, want := range s {
		noPanic(t, "SafeIndex", func() {
			got, err := SafeIndex(s, i)
			if got != want || err != nil {
				t.Errorf("SafeIndex(%q, %d) = %q, %v; want %q, <nil>", s, i, got, err, want)
			}
		})
	}
}

func TestSafeIndexOutOfRange(t *testing.T) {
	tests := []struct {
		s    []string
		i    int
		want string
	}{
		{[]string{"a", "b"}, 5, "safe index: runtime error: index out of range [5] with length 2"},
		{[]string{"a", "b"}, 2, "safe index: runtime error: index out of range [2] with length 2"},
		{[]string{"a"}, -1, "safe index: runtime error: index out of range [-1]"},
		{nil, 0, "safe index: runtime error: index out of range [0] with length 0"},
		{[]string{}, 0, "safe index: runtime error: index out of range [0] with length 0"},
	}
	for _, tt := range tests {
		noPanic(t, "SafeIndex", func() {
			got, err := SafeIndex(tt.s, tt.i)
			checkErr(t, "SafeIndex", err, tt.want)
			if got != "" {
				t.Errorf("SafeIndex(%q, %d) = %q, want \"\"", tt.s, tt.i, got)
			}
		})
	}
}

func TestSafeSet(t *testing.T) {
	m := map[string]int{"a": 1}
	noPanic(t, "SafeSet", func() {
		if err := SafeSet(m, "b", 2); err != nil {
			t.Errorf("SafeSet(m, %q, 2) error = %v, want <nil>", "b", err)
		}
		if err := SafeSet(m, "a", 10); err != nil {
			t.Errorf("SafeSet(m, %q, 10) error = %v, want <nil>", "a", err)
		}
	})
	if m["a"] != 10 || m["b"] != 2 || len(m) != 2 {
		t.Errorf("after SafeSet calls m = %v, want map[a:10 b:2]", m)
	}
}

func TestSafeSetNilMap(t *testing.T) {
	noPanic(t, "SafeSet", func() {
		err := SafeSet(nil, "x", 1)
		checkErr(t, `SafeSet(nil, "x", 1)`, err, "safe set: assignment to entry in nil map")
	})
}

func TestStringAt(t *testing.T) {
	row := []any{"id", 7, "名前", nil, 3.5, ""}
	tests := []struct {
		i    int
		want string
	}{
		{0, "id"},
		{2, "名前"},
		{5, ""},
	}
	for _, tt := range tests {
		noPanic(t, "StringAt", func() {
			got, err := StringAt(row, tt.i)
			if got != tt.want || err != nil {
				t.Errorf("StringAt(%v, %d) = %q, %v; want %q, <nil>", row, tt.i, got, err, tt.want)
			}
		})
	}
}

func TestStringAtErrors(t *testing.T) {
	row := []any{"id", 7, nil, 3.5}
	tests := []struct {
		i      int
		want   string
		assert bool // the error must wrap a *runtime.TypeAssertionError
	}{
		{1, "string at 1: interface conversion: interface {} is int, not string", true},
		{2, "string at 2: interface conversion: interface {} is nil, not string", true},
		{3, "string at 3: interface conversion: interface {} is float64, not string", true},
		{4, "string at 4: runtime error: index out of range [4] with length 4", false},
		{-2, "string at -2: runtime error: index out of range [-2]", false},
	}
	for _, tt := range tests {
		noPanic(t, "StringAt", func() {
			got, err := StringAt(row, tt.i)
			checkErr(t, "StringAt", err, tt.want)
			if got != "" {
				t.Errorf("StringAt(%v, %d) = %q, want \"\"", row, tt.i, got)
			}
			var tae *runtime.TypeAssertionError
			if tt.assert && err != nil && !errors.As(err, &tae) {
				t.Errorf("StringAt(%v, %d) error = %v, want it to wrap a *runtime.TypeAssertionError", row, tt.i, err)
			}
		})
	}
}
