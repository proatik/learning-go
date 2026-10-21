package size

import "fmt"

// Binary byte sizes. The expression is repeated on each line with iota
// increasing: 1<<10, 1<<20, 1<<30, 1<<40. These are untyped constants, so
// they can be compared with any integer type, including int64.
const (
	_  = iota // skip iota == 0 (1 << 0 would be 1)
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

// HumanSize formats a byte count using the largest fitting unit
// (B, KB, MB, GB or TB), with one decimal for units above bytes.
func HumanSize(n int64) string {
	f := float64(n)
	if n >= TB {
		return fmt.Sprintf("%.1f TB", f/TB)
	} else if n >= GB {
		return fmt.Sprintf("%.1f GB", f/GB)
	} else if n >= MB {
		return fmt.Sprintf("%.1f MB", f/MB)
	} else if n >= KB {
		return fmt.Sprintf("%.1f KB", f/KB)
	}
	return fmt.Sprintf("%d B", n)
}
