package initials

import (
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// Initials returns the upper-cased first letter of every word in name.
// During the bug hunt it is swapped for a broken version, and your tests must
// notice each one.
func Initials(name string) string { return impl(name) }

var impl = initials

func initials(name string) string {
	var b strings.Builder
	for _, w := range strings.Fields(name) {
		r := []rune(w)[0]
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

var bugs = map[string]func(string) string{
	"keeps lowercase": func(name string) string {
		var b strings.Builder
		for _, w := range strings.Fields(name) {
			b.WriteRune([]rune(w)[0])
		}
		return b.String()
	},
	"only the first word": func(name string) string {
		if f := strings.Fields(name); len(f) > 0 {
			return initials(f[0])
		}
		return ""
	},
	"drops the last word": func(name string) string {
		if f := strings.Fields(name); len(f) > 1 {
			return initials(strings.Join(f[:len(f)-1], " "))
		}
		return initials(name)
	},
	"takes the first byte, not the first rune": func(name string) string {
		var b strings.Builder
		for _, w := range strings.Fields(name) {
			b.WriteRune(unicode.ToUpper(rune(w[0])))
		}
		return b.String()
	},
}

const (
	bugEnv     = "LEARNGO_BUG"
	learnFile  = "initials_test.go"
	helperName = "assertEqual"
)

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

// helperLines returns the first and last line of the helper in your file.
func helperLines(t *testing.T) (first, last int) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, learnFile, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", learnFile, err)
	}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == helperName {
			return fset.Position(fn.Pos()).Line, fset.Position(fn.End()).Line
		}
	}
	t.Fatalf("no %s helper found in %s: write one", helperName, learnFile)
	return 0, 0
}

var failureLine = regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(learnFile) + `:(\d+):`)

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must
// fail, and each failure must point at a line outside the helper.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if os.Getenv(bugEnv) != "" {
		t.Skip("inside the bug hunt")
	}
	first, last := helperLines(t)
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourTestsCatchEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", "-test.list=^Test")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return n != "TestYourTestsCatchEveryBug" && strings.HasPrefix(n, "Test") }) {
		t.Fatal("no Test functions found: write at least one test")
	}
	if out, err := run("", "-test.run=^Test"); err != nil {
		t.Fatalf("your tests fail on the correct Initials:\n%s", out)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		out, err := run(name, "-test.run=^Test")
		if err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
			continue
		}
		for _, m := range failureLine.FindAllStringSubmatch(out, -1) {
			if line, _ := strconv.Atoi(m[1]); first <= line && line <= last {
				t.Errorf("failure line for bug %q = %s:%d (inside %s), want the caller's line (call t.Helper())", name, learnFile, line, helperName)
				break
			}
		}
	}
}
