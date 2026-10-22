package joiner

import (
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// JoinInts formats nums in decimal, separated by sep. During the bug hunt it
// is swapped for a broken version, and your tests must notice each one.
func JoinInts(nums []int, sep string) string {
	if probe != "" {
		probeOnce.Do(func() { os.WriteFile(probe, nil, 0o644) })
	}
	return impl(nums, sep)
}

var impl = joinInts

func joinInts(nums []int, sep string) string {
	var b strings.Builder
	for i, n := range nums {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(strconv.Itoa(n))
	}
	return b.String()
}

var bugs = map[string]func([]int, string) string{
	"adds a trailing separator": func(nums []int, sep string) string {
		if len(nums) == 0 {
			return ""
		}
		return joinInts(nums, sep) + sep
	},
	"drops the first number": func(nums []int, sep string) string {
		if len(nums) == 0 {
			return ""
		}
		return joinInts(nums[1:], sep)
	},
	"loses the minus sign": func(nums []int, sep string) string {
		return strings.ReplaceAll(joinInts(nums, sep), "-", "")
	},
	"ignores the separator": func(nums []int, sep string) string {
		return joinInts(nums, ",")
	},
}

const (
	bugEnv    = "LEARNGO_BUG"
	probeEnv  = "LEARNGO_PROBE" // a file that JoinInts creates when it is called
	learnFile = "joiner_test.go"
)

var (
	probe     = os.Getenv(probeEnv)
	probeOnce sync.Once
)

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

// usesLoop reports whether a Benchmark function in your file calls b.Loop().
func usesLoop(t *testing.T) bool {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), learnFile, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", learnFile, err)
	}
	found := false
	for _, d := range f.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(fn.Name.Name, "Benchmark") || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if c, ok := n.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok && s.Sel.Name == "Loop" {
					found = true
				}
			}
			return true
		})
	}
	return found
}

// TestYourTestsCatchEveryBug re-runs your tests once per bug (each run must
// fail), then runs your benchmark once.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if os.Getenv(bugEnv) != "" || probe != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, env []string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourTestsCatchEveryBug$"}, args...)...)
		cmd.Env = append(append(os.Environ(), bugEnv+"="+bug), env...)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", nil, "-test.list=.")
	names := strings.Fields(list)
	if !slices.ContainsFunc(names, func(n string) bool { return n != "TestYourTestsCatchEveryBug" && strings.HasPrefix(n, "Test") }) {
		t.Fatal("no Test functions found: write at least one test")
	}
	if out, err := run("", nil, "-test.run=^Test"); err != nil {
		t.Fatalf("your tests fail on the correct JoinInts:\n%s", out)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		if _, err := run(name, nil, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}

	if !slices.ContainsFunc(names, func(n string) bool { return strings.HasPrefix(n, "Benchmark") }) {
		t.Fatal("no Benchmark functions found: write one")
	}
	if !usesLoop(t) {
		t.Error("benchmark loop = none, want for b.Loop() { ... }")
	}
	marker := filepath.Join(t.TempDir(), "called")
	out, err := run("", []string{probeEnv + "=" + marker}, "-test.run=^$", "-test.bench=.", "-test.benchtime=1x")
	if err != nil {
		t.Fatalf("your benchmark fails:\n%s", out)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Error("JoinInts calls in your benchmark = 0, want at least 1")
	}
	if !strings.Contains(out, "allocs/op") {
		t.Errorf("benchmark output has no allocs/op column, want one (call b.ReportAllocs()):\n%s", out)
	}
}
