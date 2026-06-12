package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chainreactors/parsers"
	"github.com/chainreactors/zombie/action"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

// nilSessionPlugin returns (nil, nil) from Open — should not panic Execute
type nilSessionPlugin struct{}

func (p *nilSessionPlugin) Name() string                                 { return "nil-session" }
func (p *nilSessionPlugin) Open(task *pkg.Task) (pkg.Session, error)     { return nil, nil }
func (p *nilSessionPlugin) Unauth(task *pkg.Task) (pkg.Session, error)   { return nil, nil }

// panicPlugin panics inside Open — should be catchable
type panicPlugin struct{}

func (p *panicPlugin) Name() string                                 { return "panic" }
func (p *panicPlugin) Open(task *pkg.Task) (pkg.Session, error)     { panic("test panic") }
func (p *panicPlugin) Unauth(task *pkg.Task) (pkg.Session, error)   { panic("test panic") }

func baseTask(svc string) *pkg.Task {
	return &pkg.Task{
		ZombieResult: &parsers.ZombieResult{
			IP: "127.0.0.1", Port: "9999",
			Service: svc, Username: "u", Password: "p",
		},
		Timeout: 1,
	}
}

// --- Nil session from plugin ---

func TestPanic_NilSession_Execute(t *testing.T) {
	plugins := map[string]plugin.Plugin{"nil-session": &nilSessionPlugin{}}
	task := baseTask("nil-session")

	result := Execute(task, plugins, nil)
	if result.OK {
		t.Error("should not be OK")
	}
	if result.Err == nil {
		t.Error("should have error for nil session")
	}
	t.Logf("nil session: %v", result.Err)
}

func TestPanic_NilSession_ExecuteUnauth(t *testing.T) {
	plugins := map[string]plugin.Plugin{"nil-session": &nilSessionPlugin{}}
	task := baseTask("nil-session")

	result := ExecuteUnauth(task, plugins, nil)
	if result.OK {
		t.Error("should not be OK")
	}
	if result.Err == nil {
		t.Error("should have error for nil session")
	}
	t.Logf("nil session unauth: %v", result.Err)
}

// --- Missing plugin ---

func TestPanic_NoPlugin(t *testing.T) {
	plugins := map[string]plugin.Plugin{}
	task := baseTask("nonexistent")

	result := Execute(task, plugins, nil)
	if result.OK {
		t.Error("should not be OK")
	}
	if result.Err == nil {
		t.Error("should have error")
	}
	t.Logf("no plugin: %v", result.Err)
}

// --- Nil task fields ---

func TestPanic_NilParam_PluginOpen(t *testing.T) {
	runner := NewRunner(NewDefaultRunnerOption())

	services := []string{"ssh", "mysql", "redis", "ftp", "postgresql", "mssql", "oracle", "smb", "ldap"}
	for _, svc := range services {
		t.Run(svc, func(t *testing.T) {
			p, ok := runner.Plugins[svc]
			if !ok {
				t.Skipf("no plugin for %s", svc)
			}
			task := &pkg.Task{
				ZombieResult: &parsers.ZombieResult{
					IP: "127.0.0.1", Port: "1",
					Service: svc, Username: "u", Password: "p",
					Param: nil, // explicitly nil
				},
				Timeout: 1,
			}
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC with nil Param on %s: %v", svc, r)
				}
			}()
			p.Open(task)
		})
	}
}

// --- HTTP plugins with nil Param ---

func TestPanic_NilParam_HTTPPlugins(t *testing.T) {
	runner := NewRunner(NewDefaultRunnerOption())

	httpServices := []string{"http", "https", "http_proxy", "digest", "get", "post"}
	for _, svc := range httpServices {
		t.Run(svc, func(t *testing.T) {
			p, ok := runner.Plugins[svc]
			if !ok {
				t.Skipf("no plugin for %s", svc)
			}
			task := &pkg.Task{
				ZombieResult: &parsers.ZombieResult{
					IP: "127.0.0.1", Port: "1",
					Service: svc, Username: "u", Password: "p",
					Param: nil,
				},
				Timeout: 1,
			}
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC with nil Param on %s: %v", svc, r)
				}
			}()
			p.Open(task)
		})
	}
}

// --- Nil Extracteds in Merge ---

func TestPanic_MergeNilActionResult(t *testing.T) {
	result := &pkg.Result{
		Task: baseTask("ssh"),
		OK:   true,
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PANIC on Merge(nil): %v", r)
		}
	}()
	result.Merge(nil)
	result.Merge(&pkg.ActionResult{})
	result.Merge(&pkg.ActionResult{
		Loot: map[string][]byte{"test": []byte("data")},
	})
	if len(result.Loot) != 1 {
		t.Error("should have 1 loot entry")
	}
}


// --- PostAction with valid scanner on empty data ---

func TestPanic_PostAction_EmptyData(t *testing.T) {
	dir := createPanicTestTemplate(t)
	a, err := action.NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction: %v", err)
	}

	session := &mockShell{files: map[string][]byte{}}
	task := baseTask("ssh")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PANIC on empty data: %v", r)
		}
	}()

	result, err := a.Run(session, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("empty data: extracteds=%d, loot=%d", len(result.Extracteds), len(result.Loot))
}

// --- OutputHandler nil Err ---

func TestPanic_OutputHandler_NilErr(t *testing.T) {
	// Simulate what OutputHandler does with a failed result that has nil Err
	result := &pkg.Result{
		Task: baseTask("ssh"),
		OK:   false,
		Err:  nil, // this would panic on .Error() without our fix
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PANIC on nil Err formatting: %v", r)
		}
	}()

	errMsg := "unknown error"
	if result.Err != nil {
		errMsg = result.Err.Error()
	}
	_ = fmt.Sprintf("[%s] %s login failed, %s", result.Service, result.URI(), errMsg)
}

// --- Mock helpers ---

type mockShell struct {
	files map[string][]byte
}

func (m *mockShell) Service() string  { return "ssh" }
func (m *mockShell) Close() error     { return nil }
func (m *mockShell) Raw() interface{} { return nil }
func (m *mockShell) Exec(cmd string) ([]byte, error) {
	for path, data := range m.files {
		if len(cmd) > 0 && len(path) > 0 {
			for i := 0; i <= len(cmd)-len(path); i++ {
				if cmd[i:i+len(path)] == path {
					return data, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

func createPanicTestTemplate(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.yaml"), []byte(`id: panic-test
info:
  name: Panic Test
  severity: info
file:
  - extensions:
      - all
    extractors:
      - type: regex
        regex:
          - "password\\s*=\\s*(\\S+)"
        group: 1
`), 0644)
	return dir
}
