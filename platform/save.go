package main

import (
	"encoding/json"
	"go/format"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

const maxSaveBytes = 1 << 20 // 1 MiB per request is plenty for exercise code

type saveRequest struct {
	Files  map[string]string `json:"files"`
	Format bool              `json:"format"` // run gofmt before writing
}

type saveResponse struct {
	Exercise     exerciseDetail    `json:"exercise"`
	FormatErrors map[string]string `json:"formatErrors,omitempty"` // file name → syntax error
}

// handleSave only writes the exercise's starter files, never its tests.
func (s *server) handleSave(w http.ResponseWriter, r *http.Request) {
	c, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	var req saveRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxSaveBytes)).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp := saveResponse{FormatErrors: map[string]string{}}
	for name, content := range req.Files {
		if !slices.Contains(ex.Files, name) {
			http.Error(w, "not an editable file: "+name, http.StatusBadRequest)
			return
		}
		src := []byte(content)
		if req.Format {
			// Like gofmt-on-save. A syntax error leaves the code as typed.
			if out, err := format.Source(src); err == nil {
				src = out
			} else {
				resp.FormatErrors[name] = err.Error()
			}
		}
		if err := os.WriteFile(filepath.Join(s.ws.exDir(ex), name), src, 0o644); err != nil {
			httpError(w, err)
			return
		}
	}

	d, err := s.detail(c, ex)
	if err != nil {
		httpError(w, err)
		return
	}
	resp.Exercise = d
	writeJSON(w, resp)
}

type resetAllRequest struct {
	Code bool `json:"code"` // also restore every exercise's starter files
}

func (s *server) handleResetAll(w http.ResponseWriter, r *http.Request) {
	var req resetAllRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	c, ok := s.course(w)
	if !ok {
		return
	}
	if req.Code {
		for _, ex := range c.exercises {
			if err := s.ws.ResetExercise(ex); err != nil {
				httpError(w, err)
				return
			}
		}
	}
	s.mu.Lock()
	clear(s.lastResults)
	s.mu.Unlock()
	if err := s.progress.ResetAll(); err != nil {
		httpError(w, err)
		return
	}
	writeJSON(w, map[string]any{"ok": true, "exercises": len(c.exercises), "codeReset": req.Code})
}
