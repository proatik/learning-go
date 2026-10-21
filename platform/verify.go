package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// verifyCourse checks every exercise (content-authoring tool):
//   - the reference solution must pass the tests
//   - the pristine starter must NOT pass the tests
//
// Each check runs in a throwaway module so the learner's files are untouched.
func verifyCourse(root string, c *Course, filter string) bool {
	goMod, err := tempGoMod(root)
	if err != nil {
		fmt.Println("verify:", err)
		return false
	}

	var ids []string
	for id := range c.exercises {
		if strings.Contains(id, filter) {
			ids = append(ids, id)
		}
	}
	slices.Sort(ids)

	type result struct {
		id       string
		problems []string
	}
	results := make([]result, len(ids))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = result{id: id, problems: verifyExercise(goMod, c.exercises[id])}
		})
	}
	wg.Wait()

	ok := true
	for _, r := range results {
		if len(r.problems) == 0 {
			fmt.Println("ok    ", r.id)
			continue
		}
		ok = false
		fmt.Println("FAIL  ", r.id)
		for _, p := range r.problems {
			fmt.Println("      ", strings.ReplaceAll(p, "\n", "\n       "))
		}
	}
	fmt.Printf("\n%d exercises checked\n", len(ids))
	return ok
}

func verifyExercise(goMod string, ex *Exercise) []string {
	var problems []string
	if ex.Title == "" {
		problems = append(problems, "_meta.json: missing title")
	}
	if len(ex.Hints) < 2 {
		problems = append(problems, "_meta.json: expected at least 2 hints")
	}
	if _, err := os.Stat(filepath.Join(ex.Dir, "README.md")); err != nil {
		problems = append(problems, "missing README.md")
	}

	if out, err := runInTemp(goMod, ex, "_solution"); err != nil {
		problems = append(problems, "solution does not pass:\n"+tail(out, 30))
	}
	if _, err := runInTemp(goMod, ex, "."); err == nil {
		problems = append(problems, "starter code already passes the tests")
	}
	return problems
}

// runInTemp copies the exercise tests plus the .go files from srcSub into a
// fresh copy of the module (same module path and directory, so external test
// packages can import the exercise) and runs go test there.
func runInTemp(goMod string, ex *Exercise, srcSub string) (string, error) {
	tmpRoot, err := os.MkdirTemp("", "learngo-verify-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpRoot)

	if err := os.WriteFile(filepath.Join(tmpRoot, "go.mod"), []byte(goMod), 0o644); err != nil {
		return "", err
	}
	tmp := filepath.Join(tmpRoot, filepath.FromSlash(ex.RelDir))
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		return "", err
	}
	files, err := goFiles(ex.Dir)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			if err := copyFile(filepath.Join(ex.Dir, f), filepath.Join(tmp, f)); err != nil {
				return "", err
			}
		}
	}
	srcDir := filepath.Join(ex.Dir, srcSub)
	src, err := goFiles(srcDir)
	if err != nil {
		return "", err
	}
	for _, f := range src {
		if srcSub == "." && strings.HasSuffix(f, "_test.go") {
			continue
		}
		if err := copyFile(filepath.Join(srcDir, f), filepath.Join(tmp, f)); err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-count=1", "./"+ex.RelDir+"/")
	cmd.Dir = tmpRoot
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// tempGoMod returns the repo's go.mod, so verification modules use the same
// module path and go/toolchain versions.
func tempGoMod(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	return string(b), err
}

func tail(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}
