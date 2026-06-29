package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chainreactors/words"
)

func TestReadFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "words.txt")
	if err := os.WriteFile(filename, []byte("admin\nroot\n"), 0600); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := words.NewWorderWithFile(f)
	w.Run()
	got := w.All()
	if len(got) != 2 || got[0] != "admin" || got[1] != "root" {
		t.Fatalf("unexpected words: %#v", got)
	}
}

func TestOptionValidateRejectsUnsupportedMod(t *testing.T) {
	opt := &Option{}
	opt.IP = []string{"127.0.0.1"}
	opt.ServiceName = "redis"
	opt.Mod = "not-a-mode"

	if err := opt.Validate(); err == nil {
		t.Fatal("expected unsupported mode to be rejected")
	}
}

func TestOptionValidateRequiresPitchforkAuth(t *testing.T) {
	opt := &Option{}
	opt.IP = []string{"127.0.0.1"}
	opt.ServiceName = "redis"
	opt.Mod = ModPitchFork

	if err := opt.Validate(); err == nil {
		t.Fatal("expected pitchfork without auth to be rejected")
	}

	opt.Auth = []string{"user::pass"}
	if err := opt.Validate(); err != nil {
		t.Fatalf("expected pitchfork with auth to pass validation: %v", err)
	}
}

func TestOptionPrepareRejectsUnknownService(t *testing.T) {
	for _, service := range []string{"memcache", "postgres", "8080"} {
		opt := &Option{}
		opt.IP = []string{"127.0.0.1"}
		opt.ServiceName = service
		opt.Mod = ModSniper

		_, err := opt.Prepare()
		if err == nil {
			t.Fatalf("expected %q to be rejected", service)
		}
		if !strings.Contains(err.Error(), `unknown service`) {
			t.Fatalf("unexpected error for %q: %v", service, err)
		}
	}
}

func TestOptionPrepareCanonicalizesServiceAliases(t *testing.T) {
	tests := []struct {
		name string
		want string
		port string
	}{
		{name: "postgre", want: "postgresql", port: "5432"},
		{name: "mongodb", want: "mongo", port: "27017"},
		{name: "pop", want: "pop3", port: "110"},
	}

	for _, tt := range tests {
		opt := &Option{}
		opt.IP = []string{"127.0.0.1"}
		opt.ServiceName = tt.name
		opt.Mod = ModSniper

		runner, err := opt.Prepare()
		if err != nil {
			t.Fatalf("Prepare(%q): %v", tt.name, err)
		}
		if len(runner.Targets) != 1 {
			t.Fatalf("Prepare(%q) targets = %d, want 1", tt.name, len(runner.Targets))
		}
		target := runner.Targets[0]
		if target.Service != tt.want || target.Port != tt.port {
			t.Fatalf("Prepare(%q) target = %s:%s, want %s:%s",
				tt.name, target.Service, target.Port, tt.want, tt.port)
		}
	}
}

func TestOptionPrepareKeepsFileNilWithoutOutputFlag(t *testing.T) {
	opt := &Option{}
	opt.IP = []string{"127.0.0.1"}
	opt.ServiceName = "redis"
	opt.Mod = ModSniper

	runner, err := opt.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	if runner.File != nil {
		t.Fatal("Prepare without -f should leave File nil")
	}
}

func TestOptionPrepareCanonicalizesFilterServiceAliases(t *testing.T) {
	targets := `[{"ip":"127.0.0.1","port":"5432","service":"postgresql"}]`
	jsonFile := filepath.Join(t.TempDir(), "targets.json")
	if err := os.WriteFile(jsonFile, []byte(targets), 0600); err != nil {
		t.Fatal(err)
	}

	opt := &Option{}
	opt.JsonFile = jsonFile
	opt.FilterService = "postgre"
	opt.Mod = ModSniper

	runner, err := opt.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.Targets) != 1 {
		t.Fatalf("targets = %d, want 1", len(runner.Targets))
	}
	if runner.Targets[0].Service != "postgresql" {
		t.Fatalf("service = %q, want postgresql", runner.Targets[0].Service)
	}
}

func TestOptionPrepareOutputFileWriter(t *testing.T) {
	output := filepath.Join(t.TempDir(), "results.txt")
	opt := &Option{}
	opt.IP = []string{"127.0.0.1"}
	opt.ServiceName = "redis"
	opt.OutputFile = output
	opt.Mod = ModSniper

	runner, err := opt.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	if runner.File == nil {
		t.Fatal("expected output file writer")
	}

	if err := runner.File.SyncWrite("ok\n"); err != nil {
		t.Fatal(err)
	}
	if err := runner.File.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "ok\n" {
		t.Fatalf("unexpected output file content: %q", string(got))
	}
}
