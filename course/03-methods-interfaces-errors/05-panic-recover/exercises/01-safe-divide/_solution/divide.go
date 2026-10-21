package divide

import "fmt"

// SafeDivide returns a / b. A division by zero panics inside SafeDivide and
// is recovered into the error "safe divide: <panic value>".
func SafeDivide(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			q, err = 0, fmt.Errorf("safe divide: %v", r)
		}
	}()
	return a / b, nil
}

// Average returns the truncated integer mean of nums. An empty slice causes
// a division by zero, recovered into the error "average: <panic value>".
func Average(nums []int) (avg int, err error) {
	defer func() {
		if r := recover(); r != nil {
			avg, err = 0, fmt.Errorf("average: %v", r)
		}
	}()
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum / len(nums), nil
}
