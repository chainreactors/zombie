package core

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chainreactors/parsers"
	"github.com/chainreactors/zombie/pkg"
)

// === CLI parsing & validation ===

func TestE2E_Version(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{"--version"}, RunOptions{
		Output:  &out,
		Version: "v2.0.0-test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "v2.0.0-test") {
		t.Fatalf("expected version, got: %q", out.String())
	}
}

func TestE2E_Help(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{"--help"}, RunOptions{Output: &out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "zombie") {
		t.Fatalf("expected help, got: %q", out.String())
	}
}

func TestE2E_ListServices(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{"-l"}, RunOptions{Output: &out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	for _, svc := range []string{"ssh", "mysql", "redis", "ftp", "smb", "ldap"} {
		if !strings.Contains(output, svc) {
			t.Errorf("service list missing %q", svc)
		}
	}
}

func TestE2E_NoTargetError(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{"-s", "ssh"}, RunOptions{Output: &out})
	if err == nil {
		t.Fatal("should error without target")
	}
}

func TestE2E_InvalidMod(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{
		"-i", "127.0.0.1", "-s", "ssh", "-m", "invalid",
	}, RunOptions{Output: &out})
	if err == nil {
		t.Fatal("should error on invalid mod")
	}
	if !strings.Contains(err.Error(), "unsupported mod") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestE2E_PitchforkWithoutAuth(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{
		"-i", "127.0.0.1", "-s", "ssh", "-m", "pitchfork",
	}, RunOptions{Output: &out})
	if err == nil {
		t.Fatal("pitchfork without -a should error")
	}
}

// === Proton flag validation ===

func TestE2E_ProtonWithoutTemplate(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{
		"-i", "127.0.0.1", "-s", "ssh", "-u", "root", "-p", "pass", "--proton",
	}, RunOptions{Output: &out})
	if err == nil {
		t.Fatal("--proton without --scan-template should error")
	}
	if !strings.Contains(err.Error(), "--scan-template") {
		t.Fatalf("error should mention --scan-template, got: %v", err)
	}
}

func TestE2E_ProtonWithInvalidTemplate(t *testing.T) {
	var out bytes.Buffer
	err := RunWithArgs(context.Background(), []string{
		"-i", "127.0.0.1", "-s", "ssh", "-u", "root", "-p", "pass",
		"--proton", "--scan-template", "/nonexistent/path",
	}, RunOptions{Output: &out})
	if err == nil {
		t.Fatal("--proton with bad template should error")
	}
}

// === Target URL parsing ===

func TestE2E_ParseURL_SSH(t *testing.T) {
	target, ok := ParseUrl("ssh://admin:pass@10.0.0.5:2222")
	if !ok {
		t.Fatal("should parse SSH URL")
	}
	assertTarget(t, target, "10.0.0.5", "2222", "ssh", "admin", "pass")
}

func TestE2E_ParseURL_MySQL(t *testing.T) {
	target, ok := ParseUrl("mysql://root:secret@db.host:3306")
	if !ok {
		t.Fatal("should parse MySQL URL")
	}
	assertTarget(t, target, "db.host", "3306", "mysql", "root", "secret")
}

func TestE2E_ParseURL_Redis(t *testing.T) {
	target, ok := ParseUrl("redis://:authpass@10.0.0.1:6379")
	if !ok {
		t.Fatal("should parse Redis URL")
	}
	if target.Service != "redis" {
		t.Errorf("Service = %q, want redis", target.Service)
	}
}

func TestE2E_ParseURL_PostgreSQL(t *testing.T) {
	target, ok := ParseUrl("postgresql://app:dbpass@pg.host:5432")
	if !ok {
		t.Fatal("should parse PostgreSQL URL")
	}
	if target.Service != "postgresql" {
		t.Errorf("Service = %q, want postgresql", target.Service)
	}
}

// === Brute via CLI (full RunWithArgs, closed port) ===

func TestE2E_Brute_Sniper_ClosedPort(t *testing.T) {
	port := findFreePort(t)
	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunWithArgs(ctx, []string{
		"-i", fmt.Sprintf("127.0.0.1:%d", port),
		"-s", "ssh", "-u", "root", "-p", "test",
		"-m", "sniper", "--timeout", "2",
		"-q", "-f", os.DevNull,
	}, RunOptions{Output: &out})
	t.Logf("sniper: err=%v", err)
}

func TestE2E_Brute_ClusterBomb_ClosedPort(t *testing.T) {
	port := findFreePort(t)
	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunWithArgs(ctx, []string{
		"-i", fmt.Sprintf("127.0.0.1:%d", port),
		"-s", "redis", "-u", "default", "-p", "test",
		"-m", "clusterbomb", "--timeout", "2",
		"--no-honeypot", "--no-unauth",
		"-q", "-f", os.DevNull,
	}, RunOptions{Output: &out})
	t.Logf("clusterbomb: err=%v", err)
}

func TestE2E_Brute_Pitchfork_ClosedPort(t *testing.T) {
	port := findFreePort(t)
	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunWithArgs(ctx, []string{
		"-i", fmt.Sprintf("127.0.0.1:%d", port),
		"-s", "mysql", "-a", "root::password",
		"-m", "pitchfork", "--timeout", "2",
		"-q", "-f", os.DevNull,
	}, RunOptions{Output: &out})
	t.Logf("pitchfork: err=%v", err)
}

// === Multiple services via CLI ===

func TestE2E_AllServices_ClosedPort(t *testing.T) {
	port := findFreePort(t)
	services := []string{"ssh", "mysql", "redis", "ftp", "postgresql", "mssql"}

	for _, svc := range services {
		t.Run(svc, func(t *testing.T) {
			var out bytes.Buffer
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := RunWithArgs(ctx, []string{
				"-i", fmt.Sprintf("127.0.0.1:%d", port),
				"-s", svc, "-u", "test", "-p", "test",
				"-m", "sniper", "--timeout", "1",
				"--no-honeypot",
				"-q", "-f", os.DevNull,
			}, RunOptions{Output: &out})
			t.Logf("%s: err=%v", svc, err)
		})
	}
}

// === Proton pipeline via CLI ===

func TestE2E_Proton_ClosedPort(t *testing.T) {
	port := findFreePort(t)
	tmplDir := createE2ETemplate(t)

	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := RunWithArgs(ctx, []string{
		"-i", fmt.Sprintf("127.0.0.1:%d", port),
		"-s", "ssh", "-u", "root", "-p", "test",
		"-m", "sniper", "--timeout", "1",
		"--proton", "--scan-template", tmplDir,
		"-q", "-f", os.DevNull,
	}, RunOptions{Output: &out})
	t.Logf("proton pipeline: err=%v", err)
}

// === Runner API ===

func TestE2E_RunnerAPI_DefaultPipeline(t *testing.T) {
	runner := NewRunner(NewDefaultRunnerOption())
	if err := runner.BuildPipeline(); err != nil {
		t.Fatalf("default pipeline should not error: %v", err)
	}
	if len(runner.Pipeline) != 0 {
		t.Fatal("default pipeline should be empty")
	}
}

func TestE2E_RunnerAPI_ProtonPipeline(t *testing.T) {
	tmplDir := createE2ETemplate(t)
	opt := NewDefaultRunnerOption()
	opt.Proton = true
	opt.ScanTemplates = []string{tmplDir}

	runner := NewRunner(opt)
	if err := runner.BuildPipeline(); err != nil {
		t.Fatalf("build failed: %v", err)
	}
	if len(runner.Pipeline) != 1 {
		t.Fatalf("expected 1 action, got %d", len(runner.Pipeline))
	}
	if runner.Pipeline[0].Name() != "post" {
		t.Errorf("name = %q, want post", runner.Pipeline[0].Name())
	}
}

func TestE2E_RunnerAPI_PluginRegistry(t *testing.T) {
	runner := NewRunner(NewDefaultRunnerOption())
	required := []string{"ssh", "mysql", "redis", "ftp", "smb", "ldap", "postgresql", "mssql", "oracle", "neutron"}
	for _, svc := range required {
		if _, ok := runner.Plugins[svc]; !ok {
			t.Errorf("registry missing %q", svc)
		}
	}
}

// === Worker Execute (direct, no runner) ===

func TestE2E_WorkerExecute_ClosedPort(t *testing.T) {
	runner := NewRunner(NewDefaultRunnerOption())
	port := findFreePort(t)

	task := &pkg.Task{
		ZombieResult: &parsers.ZombieResult{
			IP: "127.0.0.1", Port: fmt.Sprintf("%d", port),
			Service: "ssh", Username: "root", Password: "test",
		},
		Timeout: 1,
	}

	result := Execute(task, runner.Plugins, runner.Pipeline)
	if result.OK {
		t.Error("should not succeed on closed port")
	}
	if result.Err == nil {
		t.Error("should have error")
	}
	t.Logf("Execute: OK=%v, Err=%v", result.OK, result.Err)
}

func TestE2E_WorkerExecute_WithProton_ClosedPort(t *testing.T) {
	tmplDir := createE2ETemplate(t)
	opt := NewDefaultRunnerOption()
	opt.Proton = true
	opt.ScanTemplates = []string{tmplDir}

	runner := NewRunner(opt)
	runner.BuildPipeline()
	port := findFreePort(t)

	task := &pkg.Task{
		ZombieResult: &parsers.ZombieResult{
			IP: "127.0.0.1", Port: fmt.Sprintf("%d", port),
			Service: "ssh", Username: "root", Password: "test",
		},
		Timeout: 1,
	}

	result := Execute(task, runner.Plugins, runner.Pipeline)
	if result.OK {
		t.Error("should not succeed on closed port")
	}
	if len(result.ActionResults) > 0 {
		t.Error("no actions should run when Open fails")
	}
}

func TestE2E_WorkerExecute_MultipleServices_ClosedPort(t *testing.T) {
	runner := NewRunner(NewDefaultRunnerOption())
	port := findFreePort(t)

	services := []string{"ssh", "mysql", "redis", "ftp", "postgresql", "mssql", "smb", "ldap"}
	for _, svc := range services {
		t.Run(svc, func(t *testing.T) {
			task := &pkg.Task{
				ZombieResult: &parsers.ZombieResult{
					IP: "127.0.0.1", Port: fmt.Sprintf("%d", port),
					Service: svc, Username: "test", Password: "test",
				},
				Timeout: 1,
			}
			result := Execute(task, runner.Plugins, runner.Pipeline)
			if result.OK {
				t.Errorf("%s should not succeed on closed port", svc)
			}
			t.Logf("%s: OK=%v, Err=%v", svc, result.OK, result.Err)
		})
	}
}

// === Helpers ===

func findFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func createE2ETemplate(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	tmpl := `id: e2e-test
info:
  name: E2E Test
  severity: info
file:
  - extensions:
      - all
    extractors:
      - type: regex
        regex:
          - "password\\s*=\\s*(\\S+)"
        group: 1
`
	os.WriteFile(filepath.Join(dir, "test.yaml"), []byte(tmpl), 0644)
	return dir
}

func assertTarget(t *testing.T, target *Target, ip, port, service, user, pass string) {
	t.Helper()
	if target.IP != ip {
		t.Errorf("IP = %q, want %q", target.IP, ip)
	}
	if target.Port != port {
		t.Errorf("Port = %q, want %q", target.Port, port)
	}
	if target.Service != service {
		t.Errorf("Service = %q, want %q", target.Service, service)
	}
	if target.Username != user {
		t.Errorf("Username = %q, want %q", target.Username, user)
	}
	if target.Password != pass {
		t.Errorf("Password = %q, want %q", target.Password, pass)
	}
}
