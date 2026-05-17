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
