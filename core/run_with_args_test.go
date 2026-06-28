package core

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
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

func TestRunWithArgsWithoutOutputFileDoesNotDeadlock(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		var out bytes.Buffer
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		done <- RunWithArgs(ctx, []string{
			"-i", "127.0.0.1:1",
			"-s", "redis",
			"-m", ModSniper,
			"-u", "default",
			"-p", "test",
			"--timeout", "1",
			"-q",
		}, RunOptions{Output: &out})
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("RunWithArgs deadlocked without -f")
	}
}
