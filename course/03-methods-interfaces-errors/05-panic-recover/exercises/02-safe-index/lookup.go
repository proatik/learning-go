package lookup

// SafeIndex returns s[i], or "" and an error "safe index: ..." wrapping the
// run-time panic when i is out of range.
func SafeIndex(s []string, i int) (v string, err error) {
	// TODO
	return "", nil
}

// SafeSet stores m[key] = val, or returns an error "safe set: ..." wrapping
// the run-time panic when m is nil.
func SafeSet(m map[string]int, key string, val int) (err error) {
	// TODO
	return nil
}

// StringAt returns row[i].(string), or "" and an error "string at <i>: ..."
// wrapping the run-time panic (bad index or wrong dynamic type).
func StringAt(row []any, i int) (s string, err error) {
	// TODO
	return "", nil
}
