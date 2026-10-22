package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// course/ holds read-only content. Names starting with "_" (_meta.json,
// _solution/) are ignored by the go tool; learner code lives in workspace/.

type Course struct {
	Modules   []*Module
	exercises map[string]*Exercise
	lessons   map[string]*Lesson
}

type Module struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Short   string    `json:"short"` // sidebar label
	Summary string    `json:"summary"`
	Lessons []*Lesson `json:"lessons"`
}

type Lesson struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Short     string      `json:"short"` // sidebar label
	Dir       string      `json:"-"`
	Exercises []*Exercise `json:"exercises"`
}

type Exercise struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Dir        string   `json:"-"`
	RelDir     string   `json:"-"` // slash-separated, relative to the repo root
	Hints      []string `json:"-"`
	Files      []string `json:"-"` // learner-editable files, slash paths relative to Dir
	Checkers   []string `json:"-"` // read-only *_test.go files that judge the solution
	LessonID   string   `json:"-"`
}

type exerciseMeta struct {
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Hints      []string `json:"hints"`
	Editable   []string `json:"editable"` // test files the learner writes (testing lessons)
}

type moduleMeta struct {
	Title   string            `json:"title"`
	Short   string            `json:"short"`
	Summary string            `json:"summary"`
	Lessons map[string]string `json:"lessons"` // lesson dir → short sidebar label
}

// shortTitle is the fallback sidebar label: the part of a title before ':'.
func shortTitle(title string) string {
	before, _, _ := strings.Cut(title, ":")
	return strings.TrimSpace(before)
}

func loadCourse(root string) (*Course, error) {
	c := &Course{exercises: map[string]*Exercise{}, lessons: map[string]*Lesson{}}
	courseDir := filepath.Join(root, "course")

	moduleDirs, err := subdirs(courseDir)
	if err != nil {
		return nil, err
	}
	for _, md := range moduleDirs {
		m := &Module{ID: md, Title: md}
		var mm moduleMeta
		if err := readJSON(filepath.Join(courseDir, md, "module.json"), &mm); err == nil {
			m.Title, m.Short, m.Summary = mm.Title, mm.Short, mm.Summary
		}
		if m.Short == "" {
			m.Short = m.Title
		}

		lessonDirs, err := subdirs(filepath.Join(courseDir, md))
		if err != nil {
			return nil, err
		}
		for _, ld := range lessonDirs {
			dir := filepath.Join(courseDir, md, ld)
			l := &Lesson{ID: md + "/" + ld, Dir: dir, Title: markdownTitle(filepath.Join(dir, "lesson.md"), ld)}
			l.Short = mm.Lessons[ld]
			if l.Short == "" {
				l.Short = shortTitle(l.Title)
			}

			exDirs, _ := subdirs(filepath.Join(dir, "exercises"))
			for _, ed := range exDirs {
				ex, err := loadExercise(root, filepath.Join(dir, "exercises", ed), l.ID+"/"+ed)
				if err != nil {
					log.Printf("skipping exercise %s/%s: %v", l.ID, ed, err)
					continue
				}
				ex.LessonID = l.ID
				l.Exercises = append(l.Exercises, ex)
				c.exercises[ex.ID] = ex
			}
			m.Lessons = append(m.Lessons, l)
			c.lessons[l.ID] = l
		}
		c.Modules = append(c.Modules, m)
	}
	return c, nil
}

func loadExercise(root, dir, id string) (*Exercise, error) {
	var meta exerciseMeta
	if err := readJSON(filepath.Join(dir, "_meta.json"), &meta); err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return nil, err
	}
	ex := &Exercise{
		ID:         id,
		Title:      meta.Title,
		Difficulty: meta.Difficulty,
		Dir:        dir,
		RelDir:     filepath.ToSlash(rel),
		Hints:      meta.Hints,
	}
	files, err := walkGo(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") && !slices.Contains(meta.Editable, f) {
			ex.Checkers = append(ex.Checkers, f)
		} else {
			ex.Files = append(ex.Files, f)
		}
	}
	for _, f := range meta.Editable {
		if !slices.Contains(files, f) {
			return nil, errors.New("_meta.json: editable file not found: " + f)
		}
	}
	return ex, nil
}

func subdirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), "_") && !strings.HasPrefix(e.Name(), ".") && e.Name() != "exercises" {
			out = append(out, e.Name())
		}
	}
	slices.Sort(out)
	return out, nil
}

func readJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func markdownTitle(path, fallback string) string {
	f, err := os.Open(path)
	if err != nil {
		return fallback
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if t, ok := strings.CutPrefix(sc.Text(), "# "); ok {
			return strings.TrimSpace(t)
		}
	}
	return fallback
}

func copyFile(src, dst string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, b, 0o644)
}

// walkGo lists .go files under dir as sorted slash paths, including nested
// packages but skipping "_" and "." directories (_solution/ etc.).
func walkGo(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != dir && (strings.HasPrefix(d.Name(), "_") || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".go") {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			out = append(out, filepath.ToSlash(rel))
		}
		return nil
	})
	slices.Sort(out)
	return out, err
}

// fromSlash joins a slash path from walkGo onto base.
func fromSlash(base, rel string) string { return filepath.Join(base, filepath.FromSlash(rel)) }
