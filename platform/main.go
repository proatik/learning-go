// Command platform serves the local Go learning platform.
//
//	go run ./platform              # start at http://localhost:8080
//	go run ./platform -verify      # authoring check: solutions pass, starters fail
package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:embed web
var webFS embed.FS

type server struct {
	root     string
	ws       *Workspace
	progress *ProgressStore

	mu          sync.Mutex
	lastResults map[string]TestResult // exercise id → latest run (in memory only)
}

func main() {
	addr := flag.String("addr", "localhost:8080", "listen address")
	root := flag.String("root", ".", "repository root (contains go.mod and course/)")
	verify := flag.Bool("verify", false, "verify all exercises and exit")
	only := flag.String("only", "", "with -verify: only exercises whose id contains this text")
	flag.Parse()

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(absRoot, "course")); err != nil {
		log.Fatalf("no course/ directory in %s; run from the repository root", absRoot)
	}

	course, err := loadCourse(absRoot)
	if err != nil {
		log.Fatal(err)
	}
	if *verify {
		if !verifyCourse(absRoot, course, *only) {
			os.Exit(1)
		}
		return
	}
	ws := newWorkspace(absRoot)
	if err := ws.Ensure(course); err != nil {
		log.Fatal("preparing workspace/: ", err)
	}

	progress, err := openProgress(ws.progressPath())
	if err != nil {
		log.Fatal(err)
	}
	s := &server{root: absRoot, ws: ws, progress: progress, lastResults: map[string]TestResult{}}

	index, err := fs.ReadFile(webFS, "web/index.html")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	// single-page app: every non-API URL gets the page, which routes client-side
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(index)
	})
	mux.HandleFunc("GET /api/", http.NotFound) // unknown API paths are 404, not the page
	mux.HandleFunc("GET /api/course", s.handleCourse)
	mux.HandleFunc("GET /api/lesson/{id...}", s.handleLesson)
	mux.HandleFunc("GET /api/exercise/{id...}", s.handleExercise)
	mux.HandleFunc("POST /api/hint/{id...}", s.handleHint)
	mux.HandleFunc("POST /api/run/{id...}", s.handleRun)
	mux.HandleFunc("POST /api/solution/{id...}", s.handleSolution)
	mux.HandleFunc("POST /api/reset/{id...}", s.handleReset)
	mux.HandleFunc("POST /api/save/{id...}", s.handleSave)
	mux.HandleFunc("POST /api/reset-all", s.handleResetAll)

	srv := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Go learning platform running at http://%s", *addr)
	log.Fatal(srv.ListenAndServe())
}

type exerciseSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Status     string `json:"status"`
}

type lessonSummary struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Short     string            `json:"short"`
	Exercises []exerciseSummary `json:"exercises"`
}

type moduleSummary struct {
	ID      string          `json:"id"`
	Title   string          `json:"title"`
	Short   string          `json:"short"`
	Summary string          `json:"summary"`
	Lessons []lessonSummary `json:"lessons"`
}

func (s *server) handleCourse(w http.ResponseWriter, r *http.Request) {
	c, ok := s.course(w)
	if !ok {
		return
	}
	var out []moduleSummary
	for _, m := range c.Modules {
		ms := moduleSummary{ID: m.ID, Title: m.Title, Short: m.Short, Summary: m.Summary}
		for _, l := range m.Lessons {
			ms.Lessons = append(ms.Lessons, s.lessonSummary(l))
		}
		out = append(out, ms)
	}
	writeJSON(w, out)
}

func (s *server) lessonSummary(l *Lesson) lessonSummary {
	ls := lessonSummary{ID: l.ID, Title: l.Title, Short: l.Short}
	for _, ex := range l.Exercises {
		ls.Exercises = append(ls.Exercises, exerciseSummary{
			ID: ex.ID, Title: ex.Title, Difficulty: ex.Difficulty,
			Status: s.progress.Get(ex.ID).Status,
		})
	}
	return ls
}

func (s *server) handleLesson(w http.ResponseWriter, r *http.Request) {
	c, ok := s.course(w)
	if !ok {
		return
	}
	l, ok := c.lessons[r.PathValue("id")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	md, err := os.ReadFile(filepath.Join(l.Dir, "lesson.md"))
	if err != nil {
		httpError(w, err)
		return
	}
	writeJSON(w, map[string]any{
		"lesson":   s.lessonSummary(l),
		"markdown": string(md),
	})
}

type fileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type exerciseDetail struct {
	ID            string            `json:"id"`
	LessonID      string            `json:"lessonId"`
	Title         string            `json:"title"`
	Difficulty    string            `json:"difficulty"`
	Readme        string            `json:"readme"`
	Files         []fileInfo        `json:"files"`
	Tests         []fileInfo        `json:"tests"` // read-only
	Index         int               `json:"index"` // 1-based position within the lesson
	LessonTitle   string            `json:"lessonTitle"`
	LessonShort   string            `json:"lessonShort"`
	LessonCount   int               `json:"lessonCount"`
	Dir           string            `json:"dir"`
	TestCommand   string            `json:"testCommand"`
	HintsTotal    int               `json:"hintsTotal"`
	Hints         []string          `json:"hints"` // only the revealed ones
	Status        string            `json:"status"`
	Runs          int               `json:"runs"`
	GaveUp        bool              `json:"gaveUp"`
	SolutionShown bool              `json:"solutionShown"`
	Solution      map[string]string `json:"solution,omitempty"`
	LastResult    *TestResult       `json:"lastResult,omitempty"`
	Prev          string            `json:"prev,omitempty"`
	Next          string            `json:"next,omitempty"`
}

func (s *server) handleExercise(w http.ResponseWriter, r *http.Request) {
	c, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	d, err := s.detail(c, ex)
	if err != nil {
		httpError(w, err)
		return
	}
	writeJSON(w, d)
}

func (s *server) detail(c *Course, ex *Exercise) (exerciseDetail, error) {
	p := s.progress.Get(ex.ID)
	readme, err := os.ReadFile(filepath.Join(ex.Dir, "README.md"))
	if err != nil {
		return exerciseDetail{}, err
	}
	d := exerciseDetail{
		ID: ex.ID, LessonID: ex.LessonID, Title: ex.Title, Difficulty: ex.Difficulty,
		Readme:        string(readme),
		Dir:           ex.Dir,
		TestCommand:   joinCmd(testCommand(ex)),
		HintsTotal:    len(ex.Hints),
		Hints:         ex.Hints[:min(p.HintsShown, len(ex.Hints))],
		Status:        p.Status,
		Runs:          p.Runs,
		GaveUp:        p.GaveUp,
		SolutionShown: p.SolutionShown,
	}
	for _, f := range ex.Files {
		path := fromSlash(s.ws.exDir(ex), f) // the learner's copy
		content, err := os.ReadFile(path)
		if err != nil {
			return exerciseDetail{}, err
		}
		d.Files = append(d.Files, fileInfo{Name: f, Path: filepath.ToSlash(path), Content: string(content)})
	}
	for _, f := range ex.Checkers {
		content, err := os.ReadFile(fromSlash(ex.Dir, f))
		if err != nil {
			return exerciseDetail{}, err
		}
		d.Tests = append(d.Tests, fileInfo{Name: f, Content: string(content)})
	}
	if l, ok := c.lessons[ex.LessonID]; ok {
		d.LessonTitle, d.LessonShort, d.LessonCount = l.Title, l.Short, len(l.Exercises)
		for i, e := range l.Exercises {
			if e.ID == ex.ID {
				d.Index = i + 1
			}
		}
	}
	if p.SolutionShown {
		d.Solution = map[string]string{}
		dir := filepath.Join(ex.Dir, "_solution")
		files, err := walkGo(dir)
		if err != nil {
			return exerciseDetail{}, err
		}
		for _, f := range files {
			b, err := os.ReadFile(fromSlash(dir, f))
			if err != nil {
				return exerciseDetail{}, err
			}
			d.Solution[f] = string(b)
		}
	}
	s.mu.Lock()
	if r, ok := s.lastResults[ex.ID]; ok {
		d.LastResult = &r
	}
	s.mu.Unlock()
	d.Prev, d.Next = neighbours(c, ex)
	return d, nil
}

func (s *server) handleHint(w http.ResponseWriter, r *http.Request) {
	c, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	_, err := s.progress.Update(ex.ID, func(p *ExerciseProgress) {
		// Hints unlock strictly in order: one more per request.
		if p.HintsShown < len(ex.Hints) {
			p.HintsShown++
		}
	})
	if err != nil {
		httpError(w, err)
		return
	}
	s.writeDetail(w, c, ex)
}

func (s *server) handleRun(w http.ResponseWriter, r *http.Request) {
	c, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	res := runTests(r.Context(), s.ws.dir, ex)
	s.mu.Lock()
	s.lastResults[ex.ID] = res
	s.mu.Unlock()
	_, err := s.progress.Update(ex.ID, func(p *ExerciseProgress) {
		p.Runs++
		switch {
		case res.Passed && p.Status != StatusPassed:
			now := time.Now()
			p.Status, p.PassedAt = StatusPassed, &now
		case !res.Passed && p.Status == StatusNotStarted:
			p.Status = StatusAttempted
		}
	})
	if err != nil {
		httpError(w, err)
		return
	}
	d, err := s.detail(c, ex)
	if err != nil {
		httpError(w, err)
		return
	}
	writeJSON(w, map[string]any{"result": res, "exercise": d})
}

func (s *server) handleSolution(w http.ResponseWriter, r *http.Request) {
	c, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	_, err := s.progress.Update(ex.ID, func(p *ExerciseProgress) {
		if p.Status != StatusPassed {
			p.GaveUp = true
		}
		p.SolutionShown = true
	})
	if err != nil {
		httpError(w, err)
		return
	}
	s.writeDetail(w, c, ex)
}

func (s *server) handleReset(w http.ResponseWriter, r *http.Request) {
	c, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	if err := s.ws.ResetExercise(ex); err != nil {
		httpError(w, err)
		return
	}
	s.mu.Lock()
	delete(s.lastResults, ex.ID)
	s.mu.Unlock()
	s.writeDetail(w, c, ex)
}

// course reloads the course from disk on every request so new content shows
// up without restarting the server.
func (s *server) course(w http.ResponseWriter) (*Course, bool) {
	c, err := loadCourse(s.root)
	if err != nil {
		httpError(w, err)
		return nil, false
	}
	return c, true
}

func (s *server) exercise(w http.ResponseWriter, r *http.Request) (*Course, *Exercise, bool) {
	c, ok := s.course(w)
	if !ok {
		return nil, nil, false
	}
	// Only ids discovered on disk are accepted, so the path can't escape course/.
	ex, ok := c.exercises[r.PathValue("id")]
	if !ok {
		http.NotFound(w, r)
		return nil, nil, false
	}
	// New content may have appeared since startup: seed the learner's copy.
	if err := s.ws.EnsureExercise(ex); err != nil {
		httpError(w, err)
		return nil, nil, false
	}
	return c, ex, true
}

func (s *server) writeDetail(w http.ResponseWriter, c *Course, ex *Exercise) {
	d, err := s.detail(c, ex)
	if err != nil {
		httpError(w, err)
		return
	}
	writeJSON(w, d)
}

func neighbours(c *Course, ex *Exercise) (prev, next string) {
	var all []string
	for _, m := range c.Modules {
		for _, l := range m.Lessons {
			for _, e := range l.Exercises {
				all = append(all, e.ID)
			}
		}
	}
	for i, id := range all {
		if id != ex.ID {
			continue
		}
		if i > 0 {
			prev = all[i-1]
		}
		if i < len(all)-1 {
			next = all[i+1]
		}
	}
	return prev, next
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("write json:", err)
	}
}

func httpError(w http.ResponseWriter, err error) {
	log.Println("error:", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
