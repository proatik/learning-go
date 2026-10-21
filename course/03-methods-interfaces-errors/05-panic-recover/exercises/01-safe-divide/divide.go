package divide

// SafeDivide returns a / b. A division by zero panics inside SafeDivide and
// is recovered into the error "safe divide: <panic value>".
func SafeDivide(a, b int) (q int, err error) {
	// TODO
	return 0, nil
}

// Average returns the truncated integer mean of nums. An empty slice causes
// a division by zero, recovered into the error "average: <panic value>".
func Average(nums []int) (avg int, err error) {
	// TODO
	return 0, nil
}
