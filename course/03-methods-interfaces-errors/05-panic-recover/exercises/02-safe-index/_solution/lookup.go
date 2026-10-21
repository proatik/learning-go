package lookup

import "fmt"

// asError turns a recovered value into an error with the given prefix,
// wrapping it with %w when it already is an error.
func asError(prefix string, r any) error {
	if e, ok := r.(error); ok {
		return fmt.Errorf("%s: %w", prefix, e)
	}
	return fmt.Errorf("%s: %v", prefix, r)
}

// SafeIndex returns s[i], or "" and an error "safe index: ..." wrapping the
// run-time panic when i is out of range.
func SafeIndex(s []string, i int) (v string, err error) {
	defer func() {
		if r := recover(); r != nil {
			v, err = "", asError("safe index", r)
		}
	}()
	return s[i], nil
}

// SafeSet stores m[key] = val, or returns an error "safe set: ..." wrapping
// the run-time panic when m is nil.
func SafeSet(m map[string]int, key string, val int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = asError("safe set", r)
		}
	}()
	m[key] = val
	return nil
}

// StringAt returns row[i].(string), or "" and an error "string at <i>: ..."
// wrapping the run-time panic (bad index or wrong dynamic type).
func StringAt(row []any, i int) (s string, err error) {
	defer func() {
		if r := recover(); r != nil {
			s, err = "", asError(fmt.Sprintf("string at %d", i), r)
		}
	}()
	return row[i].(string), nil
}
