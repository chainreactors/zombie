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
	opt.Threads = 1
	opt.Timeout = 1

	if err := opt.Validate(); err == nil {
		t.Fatal("expected unsupported mode to be rejected")
	}
}

func TestOptionValidateRequiresPitchforkAuth(t *testing.T) {
	opt := &Option{}
	opt.IP = []string{"127.0.0.1"}
	opt.ServiceName = "redis"
	opt.Mod = ModPitchFork
	opt.Threads = 1
	opt.Timeout = 1

	if err := opt.Validate(); err == nil {
		t.Fatal("expected pitchfork without auth to be rejected")
	}

	opt.Auth = []string{"user::pass"}
	if err := opt.Validate(); err != nil {
		t.Fatalf("expected pitchfork with auth to pass validation: %v", err)
	}
}

func TestOptionValidateRejectsInvalidRuntimeLimits(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Option)
	}{
		{name: "zero threads", edit: func(opt *Option) { opt.Threads = 0 }},
		{name: "negative concurrency", edit: func(opt *Option) { opt.Concurrency = -1 }},
		{name: "zero timeout", edit: func(opt *Option) { opt.Timeout = 0 }},
		{name: "negative top", edit: func(opt *Option) { opt.Top = -1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Option{}
			opt.IP = []string{"127.0.0.1"}
			opt.ServiceName = "redis"
			opt.Mod = ModSniper
			opt.Threads = 1
			opt.Timeout = 1
			tt.edit(opt)

			if err := opt.Validate(); err == nil {
				t.Fatal("expected invalid runtime limit to be rejected")
			}
		})
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

func TestOptionPrepareDoesNotOwnOutputFiles(t *testing.T) {
	output := filepath.Join(t.TempDir(), "results.txt")
	opt := &Option{}
	opt.IP = []string{"127.0.0.1"}
	opt.ServiceName = "redis"
	opt.Mod = ModSniper
	opt.OutputFile = output

	_, err := opt.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("Prepare created output file: %v", err)
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
