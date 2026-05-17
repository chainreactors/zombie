package core

import (
	"os"
	"path/filepath"
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
	if runner.OutFunc == nil {
		t.Fatal("expected output function")
	}

	runner.OutFunc("ok\n")
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
