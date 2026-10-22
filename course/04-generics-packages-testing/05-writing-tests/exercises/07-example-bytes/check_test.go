package bytesize

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// FormatBytes formats a size in bytes for humans, with binary units.
// During the bug hunt it is swapped for a broken version, and your Example
// must notice each one.
func FormatBytes(n int64) string { return impl(n) }

var impl = func(n int64) string { return format(n, 1024, "%.1f %s", 4) }

// format divides n by base until it is below base, using at most maxUnit units.
func format(n int64, base float64, layout string, maxUnit int) string {
	if float64(n) < base {
		return fmt.Sprintf("%d B", n)
	}
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	v := float64(n) / base
	i := 0
	for v >= base && i < maxUnit-1 {
		v /= base
		i++
	}
	return fmt.Sprintf(layout, v, units[i])
}

var bugs = map[string]func(int64) string{
	"divides by 1000":   func(n int64) string { return format(n, 1000, "%.1f %s", 4) },
	"drops the decimal": func(n int64) string { return format(n, 1024, "%.0f %s", 4) },
	"stops at MiB":      func(n int64) string { return format(n, 1024, "%.1f %s", 2) },
	"1024 stays in bytes": func(n int64) string {
		if n == 1024 {
			return "1024 B"
		}
		return format(n, 1024, "%.1f %s", 4)
	},
}

const bugEnv = "LEARNGO_BUG"

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

// TestYourExampleCatchesEveryBug re-runs your Example once per bug: each run
// must fail.
func TestYourExampleCatchesEveryBug(t *testing.T) {
	if os.Getenv(bugEnv) != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourExampleCatchesEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	// -test.list only shows examples that have an output comment.
	list, _ := run("", "-test.list=^Example")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return strings.HasPrefix(n, "Example") }) {
		t.Fatal("no Example functions with an // Output: comment found: write one")
	}
	if out, err := run("", "-test.run=^Example"); err != nil {
		t.Fatalf("your Example fails on the correct FormatBytes:\n%s", out)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		if _, err := run(name, "-test.run=^Example"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
