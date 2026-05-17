package core

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsListsServices(t *testing.T) {
	var out bytes.Buffer

	if err := RunWithArgs(context.Background(), []string{"-l"}, RunOptions{Output: &out}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "support service list") || !strings.Contains(out.String(), "ssh") {
		t.Fatalf("unexpected service list output: %q", out.String())
	}
}

func TestRunWithArgsRejectsUnsupportedMod(t *testing.T) {
	var out bytes.Buffer

	err := RunWithArgs(context.Background(), []string{"-i", "127.0.0.1", "-s", "redis", "-m", "not-a-mode"}, RunOptions{Output: &out})
	if err == nil {
		t.Fatal("expected unsupported mode to return an error")
	}
	if !strings.Contains(err.Error(), "unsupported mod") {
		t.Fatalf("unexpected error: %v", err)
	}
}
