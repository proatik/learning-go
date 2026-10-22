package appconf

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var errMissing = errors.New("config file missing")

// counting returns a loader that counts its calls and returns cfg, err.
func counting(cfg Config, err error) (load func() (Config, error), calls *atomic.Int64) {
	calls = new(atomic.Int64)
	load = func() (Config, error) {
		calls.Add(1)
		return cfg, err
	}
	return load, calls
}

func TestNewIsLazy(t *testing.T) {
	load, calls := counting(Config{Host: "localhost", Port: 8080}, nil)
	a := New(load)
	if a == nil {
		t.Fatalf("New(load) = nil, want an *App")
	}
	if got := calls.Load(); got != 0 {
		t.Errorf("New(load): load calls = %d, want 0 (load on first Config call)", got)
	}
}

func TestLoadsOnce(t *testing.T) {
	want := Config{Host: "localhost", Port: 8080}
	load, calls := counting(want, nil)
	a := New(load)
	for i := range 3 {
		got, err := a.Config()
		if got != want || err != nil {
			t.Errorf("Config() call %d = (%+v, %v), want (%+v, <nil>)", i+1, got, err, want)
		}
	}
	addr, err := a.Addr()
	if addr != "localhost:8080" || err != nil {
		t.Errorf("Addr() = (%q, %v), want (%q, <nil>)", addr, err, "localhost:8080")
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("3 Config() + 1 Addr(): load calls = %d, want 1", got)
	}
}

func TestErrorIsCached(t *testing.T) {
	load, calls := counting(Config{}, errMissing)
	a := New(load)
	for i := range 3 {
		if _, err := a.Config(); !errors.Is(err, errMissing) {
			t.Errorf("Config() call %d error = %v, want %v", i+1, err, errMissing)
		}
	}
	addr, err := a.Addr()
	if addr != "" || !errors.Is(err, errMissing) {
		t.Errorf("Addr() = (%q, %v), want (\"\", %v)", addr, err, errMissing)
	}
	if err != nil && err.Error() != "appconf: config file missing" {
		t.Errorf("Addr() error = %q, want %q", err.Error(), "appconf: config file missing")
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("failing loader, 3 Config() + 1 Addr(): load calls = %d, want 1 (errors are not retried)", got)
	}
}

func TestAddr(t *testing.T) {
	tests := []struct {
		cfg  Config
		want string
	}{
		{Config{Host: "localhost", Port: 8080}, "localhost:8080"},
		{Config{Host: "10.0.0.1", Port: 443}, "10.0.0.1:443"},
		{Config{Host: "", Port: 9000}, ":9000"},
		{Config{Host: "api.example.com", Port: 0}, "api.example.com:0"},
	}
	for _, tt := range tests {
		load, _ := counting(tt.cfg, nil)
		got, err := New(load).Addr()
		if got != tt.want || err != nil {
			t.Errorf("Addr() with %+v = (%q, %v), want (%q, <nil>)", tt.cfg, got, err, tt.want)
		}
	}
}

func TestConcurrentCallersShareOneLoad(t *testing.T) {
	want := Config{Host: "db", Port: 5432}
	var calls atomic.Int64
	started := make(chan struct{}, 100)
	release := make(chan struct{})
	a := New(func() (Config, error) {
		calls.Add(1)
		started <- struct{}{}
		<-release // hold every other caller until the test lets go
		return want, nil
	})

	const callers = 100
	results := make([]Config, callers)
	errs := make([]error, callers)
	var wg sync.WaitGroup
	for i := range callers {
		wg.Go(func() { results[i], errs[i] = a.Config() })
	}
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case <-started:
	case <-allDone:
		t.Fatalf("100 concurrent Config() calls returned without calling load")
	case <-time.After(2 * time.Second):
		close(release)
		t.Fatalf("100 concurrent Config(): load was not called within 2s")
	}
	time.Sleep(10 * time.Millisecond) // let the other callers pile up
	close(release)
	select {
	case <-allDone:
	case <-time.After(2 * time.Second):
		t.Fatalf("100 concurrent Config() calls did not all return within 2s")
	}

	if got := calls.Load(); got != 1 {
		t.Errorf("100 concurrent Config(): load calls = %d, want 1", got)
	}
	for i := range callers {
		if results[i] != want || errs[i] != nil {
			t.Fatalf("concurrent Config() caller %d = (%+v, %v), want (%+v, <nil>)", i, results[i], errs[i], want)
		}
	}
}

func TestAppsAreIndependent(t *testing.T) {
	loadA, callsA := counting(Config{Host: "a", Port: 1}, nil)
	loadB, callsB := counting(Config{Host: "b", Port: 2}, nil)
	a, b := New(loadA), New(loadB)
	ca, _ := a.Config()
	cb, _ := b.Config()
	if ca.Host != "a" || cb.Host != "b" {
		t.Errorf("two Apps: Config().Host = %q and %q, want %q and %q", ca.Host, cb.Host, "a", "b")
	}
	if callsA.Load() != 1 || callsB.Load() != 1 {
		t.Errorf("two Apps: load calls = %d and %d, want 1 and 1 (each App has its own once)", callsA.Load(), callsB.Load())
	}
}

func TestPanicIsRepeated(t *testing.T) {
	var calls atomic.Int64
	a := New(func() (Config, error) {
		calls.Add(1)
		panic("bad config")
	})
	for i := range 2 {
		func() {
			defer func() {
				if r := recover(); r != "bad config" {
					t.Errorf("Config() call %d with a panicking loader: panic = %v, want %q", i+1, r, "bad config")
				}
			}()
			a.Config()
		}()
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("panicking loader, 2 Config(): load calls = %d, want 1", got)
	}
}
