package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// workspace/ keeps learner code out of course/ so the repo only holds content.
// It is a nested module that mirrors course/ paths under the same module path,
// so external test packages that import an exercise still resolve.
const workspaceDirName = "workspace"

type Workspace struct {
	root string
	dir  string
}

func newWorkspace(root string) *Workspace {
	return &Workspace{root: root, dir: filepath.Join(root, workspaceDirName)}
}

func (w *Workspace) exDir(ex *Exercise) string {
	return filepath.Join(w.dir, filepath.FromSlash(ex.RelDir))
}

func (w *Workspace) progressPath() string { return filepath.Join(w.dir, "progress.json") }

const workspaceReadme = `# Your workspace

- This folder holds **your code** and **progress.json**.
- It is created by the platform and ignored by git. Never commit it.
- Delete it to start over from scratch (stop the platform first).
`

func (w *Workspace) Ensure(c *Course) error {
	if err := os.MkdirAll(w.dir, 0o755); err != nil {
		return err
	}
	goMod, err := os.ReadFile(filepath.Join(w.root, "go.mod"))
	if err != nil {
		return err
	}
	if err := writeIfChanged(filepath.Join(w.dir, "go.mod"), goMod); err != nil {
		return err
	}
	if err := writeIfMissing(filepath.Join(w.dir, "README.md"), []byte(workspaceReadme)); err != nil {
		return err
	}
	for _, ex := range c.exercises {
		if err := w.EnsureExercise(ex); err != nil {
			return err
		}
	}
	return nil
}

// EnsureExercise refreshes tests but never overwrites existing learner code.
func (w *Workspace) EnsureExercise(ex *Exercise) error {
	dst := w.exDir(ex)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, f := range ex.Files {
		src, err := os.ReadFile(filepath.Join(ex.Dir, f))
		if err != nil {
			return err
		}
		if err := writeIfMissing(filepath.Join(dst, f), src); err != nil {
			return err
		}
	}
	files, err := goFiles(ex.Dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(ex.Dir, f))
		if err != nil {
			return err
		}
		if err := writeIfChanged(filepath.Join(dst, f), b); err != nil {
			return err
		}
	}
	return nil
}

func (w *Workspace) ResetExercise(ex *Exercise) error {
	if err := w.EnsureExercise(ex); err != nil {
		return err
	}
	for _, f := range ex.Files {
		if err := copyFile(filepath.Join(ex.Dir, f), filepath.Join(w.exDir(ex), f)); err != nil {
			return err
		}
	}
	return nil
}

func writeIfMissing(path string, b []byte) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func writeIfChanged(path string, b []byte) error {
	if cur, err := os.ReadFile(path); err == nil && bytes.Equal(cur, b) {
		return nil
	}
	return os.WriteFile(path, b, 0o644)
}
