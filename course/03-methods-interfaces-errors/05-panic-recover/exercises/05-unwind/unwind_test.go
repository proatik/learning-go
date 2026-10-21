package unwind_test

import (
	"fmt"
	"os"

	unwind "learngo/course/03-methods-interfaces-errors/05-panic-recover/exercises/05-unwind"
)

// safe runs f and prints an escaped panic instead of crashing the test binary.
func safe(f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC ESCAPED:", r)
		}
	}()
	f()
}

func ExampleUnwind() {
	safe(func() { unwind.Unwind(os.Stdout, 3) })
	// Output:
	// enter 1
	// enter 2
	// enter 3
	// exit 3
	// exit 2
	// exit 1
	// recovered: bottom at 3
}

func ExampleUnwind_one() {
	safe(func() { unwind.Unwind(os.Stdout, 1) })
	// Output:
	// enter 1
	// exit 1
	// recovered: bottom at 1
}

func ExampleUnwind_zero() {
	safe(func() { unwind.Unwind(os.Stdout, 0) })
	safe(func() { unwind.Unwind(os.Stdout, -4) })
	// Output:
	// nothing to unwind
	// nothing to unwind
}

func ExampleAcquire() {
	safe(func() { unwind.Acquire(os.Stdout, []string{"db", "cache", "queue"}, 1) })
	// Output:
	// acquire db
	// acquire cache
	// release cache
	// release db
	// error: failed at cache
}

func ExampleAcquire_first() {
	safe(func() { unwind.Acquire(os.Stdout, []string{"db", "cache", "queue"}, 0) })
	// Output:
	// acquire db
	// release db
	// error: failed at db
}

func ExampleAcquire_last() {
	safe(func() { unwind.Acquire(os.Stdout, []string{"db", "cache", "queue"}, 2) })
	// Output:
	// acquire db
	// acquire cache
	// acquire queue
	// release queue
	// release cache
	// release db
	// error: failed at queue
}

func ExampleAcquire_noFailure() {
	safe(func() { unwind.Acquire(os.Stdout, []string{"db", "cache"}, -1) })
	// Output:
	// acquire db
	// acquire cache
	// all done
	// release cache
	// release db
}

func ExampleAcquire_failAtOutOfRange() {
	safe(func() { unwind.Acquire(os.Stdout, []string{"db"}, 5) })
	safe(func() { unwind.Acquire(os.Stdout, nil, 0) })
	// Output:
	// acquire db
	// all done
	// release db
	// all done
}
