package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chainreactors/zombie/pkg"
)

func writeWordsFile(t *testing.T, words string) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "words.txt")
	if err := os.WriteFile(filename, []byte(words), 0600); err != nil {
		t.Fatal(err)
	}
	return filename
}

func TestNewGeneratorWithFile(t *testing.T) {
	g, err := NewGeneratorWithFile(writeWordsFile(t, "admin\nroot\n"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = g.File.Close() })

	got := g.RunAsSlice()
	if len(got) != 2 || got[0] != "admin" || got[1] != "root" {
		t.Fatalf("unexpected generated words: %#v", got)
	}
}

func TestNewGeneratorWithRule(t *testing.T) {
	g, err := NewGeneratorWithFile(writeWordsFile(t, "admin\n"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = g.File.Close() })
	if err := g.SetRuleFile(filepath.Join("..", "templates", "zombie", "rule", "weakpass.rule")); err != nil {
		t.Fatal(err)
	}

	got := g.RunAsSlice()
	if len(got) == 0 {
		t.Fatal("expected weakpass rule to generate words")
	}
}

func TestNewGeneratorWithWord(t *testing.T) {
	g, err := NewGeneratorWithWord("{?#3}", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_ = g
}

func TestKeyword(t *testing.T) {
	if err := pkg.LoadKeyword(); err != nil {
		t.Fatal(err)
	}
	g, err := NewGeneratorWithWord("{@mysql_pwd}", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	got := g.RunAsSlice()
	if len(got) == 0 {
		t.Fatal("expected mysql_pwd keyword to produce words")
	}
}

func TestWeakPass(t *testing.T) {
	g, err := NewGeneratorWithFile(writeWordsFile(t, "admin\n"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = g.File.Close() })
	if err := g.SetRuleFile(filepath.Join("..", "templates", "zombie", "rule", "weakpass.rule")); err != nil {
		t.Fatal(err)
	}

	got := g.RunAsSlice()
	if len(got) == 0 {
		t.Fatal("expected weakpass rule to generate words")
	}
}
