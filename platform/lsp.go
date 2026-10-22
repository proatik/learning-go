package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// gopls powers "Full" code suggestions. It is installed once into
// workspace/.tools (built by the repo's Go toolchain, so it understands the
// newest syntax) and spoken to over stdio JSON-RPC (LSP).
const goplsVersion = "v0.24.0"

type lsp struct {
	ws *Workspace

	mu      sync.Mutex
	state   string // "off" until first use, then "starting", "ready" or "unavailable"
	errMsg  string
	stdin   io.WriteCloser
	nextID  int
	pending map[int]chan rpcResponse
	opened  map[string]int              // uri → version sent to gopls
	diags   map[string][]map[string]any // normalised uri → latest diagnostics
	diagVer map[string]int              // normalised uri → document version those diagnostics describe
	diagCh  chan struct{}               // closed and replaced whenever diagnostics arrive
	ready   chan struct{}
	writeMu sync.Mutex
	proc    *os.Process
}

// Stop ends gopls; it does not reliably exit on its own when the platform is killed.
func (l *lsp) Stop() {
	l.mu.Lock()
	p := l.proc
	l.mu.Unlock()
	if p != nil {
		l.write(map[string]any{"jsonrpc": "2.0", "id": -1, "method": "shutdown"})
		l.notify("exit", nil)
		time.AfterFunc(500*time.Millisecond, func() { p.Kill() })
		time.Sleep(600 * time.Millisecond)
	}
}

type rpcResponse struct {
	Result json.RawMessage
	Err    *struct {
		Message string `json:"message"`
	}
}

func newLSP(ws *Workspace) *lsp {
	return &lsp{ws: ws, state: "off", pending: map[int]chan rpcResponse{}, opened: map[string]int{},
		diags: map[string][]map[string]any{}, diagVer: map[string]int{}, diagCh: make(chan struct{}), ready: make(chan struct{})}
}

// Start installs and launches gopls in the background; it is safe to call more than once.
func (l *lsp) Start() {
	l.mu.Lock()
	if l.state != "off" {
		l.mu.Unlock()
		return
	}
	l.state = "starting"
	l.mu.Unlock()
	go func() {
		if err := l.launch(); err != nil {
			log.Println("code suggestions unavailable:", err)
			l.mu.Lock()
			l.state, l.errMsg = "unavailable", err.Error()
			l.mu.Unlock()
		}
	}()
}

func (l *lsp) launch() error {
	bin, err := l.ensureBinary()
	if err != nil {
		return err
	}
	cmd := exec.Command(bin)
	cmd.Dir = l.ws.dir
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	l.stdin = stdin
	l.mu.Lock()
	l.proc = cmd.Process
	l.mu.Unlock()
	go l.readLoop(bufio.NewReader(stdout))
	go func() { // gopls exits when the platform closes its stdin
		err := cmd.Wait()
		l.mu.Lock()
		l.state, l.errMsg = "unavailable", fmt.Sprint("gopls exited: ", err)
		l.mu.Unlock()
	}()

	init := map[string]any{
		"processId": os.Getpid(),
		"rootUri":   fileURI(l.ws.dir),
		"capabilities": map[string]any{"textDocument": map[string]any{
			"completion":    map[string]any{"completionItem": map[string]any{"snippetSupport": true, "documentationFormat": []string{"markdown"}}},
			"hover":         map[string]any{"contentFormat": []string{"markdown"}},
			"signatureHelp": map[string]any{"signatureInformation": map[string]any{"documentationFormat": []string{"markdown"}}},
		}},
		"initializationOptions": map[string]any{"usePlaceholders": true},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if _, err := l.call(ctx, "initialize", init); err != nil {
		return err
	}
	if err := l.notify("initialized", map[string]any{}); err != nil {
		return err
	}
	l.mu.Lock()
	l.state = "ready"
	l.mu.Unlock()
	close(l.ready)
	log.Println("code suggestions ready (gopls " + goplsVersion + ")")
	return nil
}

// ensureBinary installs gopls into workspace/.tools on first use.
func (l *lsp) ensureBinary() (string, error) {
	dir := filepath.Join(l.ws.dir, ".tools", goplsVersion)
	name := "gopls"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(dir, name)
	if _, err := os.Stat(bin); err == nil {
		return bin, nil
	}
	log.Println("installing gopls " + goplsVersion + " for code suggestions (first run only)")
	cmd := exec.Command("go", "install", "golang.org/x/tools/gopls@"+goplsVersion)
	cmd.Dir = l.ws.root // the repo's go.mod pins the toolchain gopls is built with
	cmd.Env = append(os.Environ(), "GOBIN="+dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("installing gopls: %v\n%s", err, out)
	}
	return bin, nil
}

func (l *lsp) readLoop(r *bufio.Reader) {
	for {
		length := 0
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			if v, ok := strings.CutPrefix(strings.ToLower(line), "content-length:"); ok {
				length, _ = strconv.Atoi(strings.TrimSpace(v))
			}
		}
		body := make([]byte, length)
		if _, err := io.ReadFull(r, body); err != nil {
			return
		}
		var msg struct {
			ID     *json.RawMessage `json:"id"`
			Method string           `json:"method"`
			Params json.RawMessage  `json:"params"`
			Result json.RawMessage  `json:"result"`
			Error  *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &msg) != nil {
			continue
		}
		switch {
		case msg.Method == "textDocument/publishDiagnostics":
			var p struct {
				URI         string           `json:"uri"`
				Version     int              `json:"version"`
				Diagnostics []map[string]any `json:"diagnostics"`
			}
			if json.Unmarshal(msg.Params, &p) == nil {
				l.mu.Lock()
				l.diags[normURI(p.URI)] = p.Diagnostics
				l.diagVer[normURI(p.URI)] = p.Version
				close(l.diagCh)
				l.diagCh = make(chan struct{})
				l.mu.Unlock()
			}
		case msg.Method != "" && msg.ID != nil:
			// server → client request (configuration, progress, ...): a null answer is enough
			l.write(map[string]any{"jsonrpc": "2.0", "id": msg.ID, "result": nil})
		case msg.Method == "" && msg.ID != nil:
			id, _ := strconv.Atoi(string(*msg.ID))
			l.mu.Lock()
			ch := l.pending[id]
			delete(l.pending, id)
			l.mu.Unlock()
			if ch != nil {
				ch <- rpcResponse{Result: msg.Result, Err: msg.Error}
			}
		}
	}
}

func (l *lsp) write(msg any) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	l.writeMu.Lock()
	defer l.writeMu.Unlock()
	if l.stdin == nil {
		return errors.New("gopls is not running")
	}
	_, err = fmt.Fprintf(l.stdin, "Content-Length: %d\r\n\r\n%s", len(b), b)
	return err
}

func (l *lsp) notify(method string, params any) error {
	return l.write(map[string]any{"jsonrpc": "2.0", "method": method, "params": params})
}

func (l *lsp) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	ch := make(chan rpcResponse, 1)
	l.mu.Lock()
	l.nextID++
	id := l.nextID
	l.pending[id] = ch
	l.mu.Unlock()
	if err := l.write(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params}); err != nil {
		return nil, err
	}
	select {
	case r := <-ch:
		if r.Err != nil {
			return nil, errors.New(r.Err.Message)
		}
		return r.Result, nil
	case <-ctx.Done():
		l.mu.Lock()
		delete(l.pending, id)
		l.mu.Unlock()
		l.notify("$/cancelRequest", map[string]any{"id": id})
		return nil, ctx.Err()
	}
}

// diagnostics waits until gopls has checked the latest synced text, so an error that was
// just fixed is not reported again. gopls skips publishing when nothing changed, so after
// the timeout the last diagnostics are still the right answer.
func (l *lsp) diagnostics(ctx context.Context, uri string) []map[string]any {
	timeout := time.After(3 * time.Second)
	for {
		l.mu.Lock()
		d, done, changed := l.diags[normURI(uri)], l.diagVer[normURI(uri)] >= l.opened[uri], l.diagCh
		l.mu.Unlock()
		if done {
			return d
		}
		select {
		case <-changed:
		case <-timeout:
			return d
		case <-ctx.Done():
			return d
		}
	}
}

// sync sends the editor's live text (full-document sync).
func (l *lsp) sync(uri, text string) error {
	l.mu.Lock()
	v, open := l.opened[uri]
	v++
	l.opened[uri] = v
	l.mu.Unlock()
	if !open {
		return l.notify("textDocument/didOpen", map[string]any{"textDocument": map[string]any{"uri": uri, "languageId": "go", "version": v, "text": text}})
	}
	return l.notify("textDocument/didChange", map[string]any{
		"textDocument":   map[string]any{"uri": uri, "version": v},
		"contentChanges": []map[string]any{{"text": text}},
	})
}

func (l *lsp) close(uri string) {
	l.mu.Lock()
	_, open := l.opened[uri]
	delete(l.opened, uri)
	delete(l.diags, normURI(uri))
	delete(l.diagVer, normURI(uri))
	l.mu.Unlock()
	if open {
		l.notify("textDocument/didClose", map[string]any{"textDocument": map[string]any{"uri": uri}})
	}
}

func fileURI(path string) string {
	abs, _ := filepath.Abs(path)
	// gopls works with the real path; Windows short names (ATIKUL~1) or
	// symlinks would make its "which package is this file in" lookup fail.
	if real, err := filepath.EvalSymlinks(abs); err == nil {
		abs = real
	} else if real, err := filepath.EvalSymlinks(filepath.Dir(abs)); err == nil {
		abs = filepath.Join(real, filepath.Base(abs)) // file not created yet
	}
	return (&url.URL{Scheme: "file", Path: "/" + strings.TrimPrefix(filepath.ToSlash(abs), "/")}).String()
}

// normURI makes gopls's URIs comparable with ours (drive-letter case, escaping).
func normURI(u string) string {
	if p, err := url.Parse(u); err == nil {
		u = p.Path
	}
	if runtime.GOOS == "windows" {
		u = strings.ToLower(u)
	}
	return u
}

// --- HTTP ---------------------------------------------------------------------

type lspRequest struct {
	File      string `json:"file"`
	Text      string `json:"text"`
	Line      int    `json:"line"` // 0-based, LSP positions
	Character int    `json:"character"`
	Close     bool   `json:"close"`
}

// allowed LSP requests: read-only help only; no code actions or edits that
// could solve an exercise for the learner.
var lspMethods = map[string]string{
	"completion":    "textDocument/completion",
	"hover":         "textDocument/hover",
	"signatureHelp": "textDocument/signatureHelp",
}

func (s *server) handleLSPStatus(w http.ResponseWriter, r *http.Request) {
	s.lsp.Start()
	s.lsp.mu.Lock()
	defer s.lsp.mu.Unlock()
	writeJSON(w, map[string]string{"state": s.lsp.state, "error": s.lsp.errMsg})
}

// handleLSP serves POST /api/lsp/{op}/{id...} for one learner file of one exercise.
func (s *server) handleLSP(w http.ResponseWriter, r *http.Request) {
	_, ex, ok := s.exercise(w, r)
	if !ok {
		return
	}
	var req lspRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxSaveBytes)).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if !slices.Contains(ex.Files, req.File) {
		http.Error(w, "not a learner file: "+req.File, http.StatusBadRequest)
		return
	}
	s.lsp.Start()
	select {
	case <-s.lsp.ready:
	default:
		writeJSON(w, map[string]any{"state": "starting"})
		return
	}
	uri := fileURI(fromSlash(s.ws.exDir(ex), req.File))
	op := r.PathValue("op")
	switch op {
	case "sync":
		if req.Close {
			s.lsp.close(uri)
		} else if err := s.lsp.sync(uri, req.Text); err != nil {
			httpError(w, err)
			return
		}
		writeJSON(w, map[string]any{"ok": true})
	case "diagnostics":
		writeJSON(w, map[string]any{"diagnostics": s.lsp.diagnostics(r.Context(), uri)})
	default:
		method, ok := lspMethods[op]
		if !ok {
			http.NotFound(w, r)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		res, err := s.lsp.call(ctx, method, map[string]any{
			"textDocument": map[string]any{"uri": uri},
			"position":     map[string]any{"line": req.Line, "character": req.Character},
		})
		if err != nil {
			writeJSON(w, map[string]any{"result": nil, "error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"result":%s}`, orNull(res))
	}
}

func orNull(b json.RawMessage) []byte {
	if len(b) == 0 {
		return []byte("null")
	}
	return b
}
