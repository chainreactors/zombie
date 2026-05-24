package core

import (
	"context"
	"testing"
)

func TestRunnerRunWithContextRejectsUnsupportedMod(t *testing.T) {
	runner := NewRunner(&RunnerOption{Mod: "not-a-mode"})

	if err := runner.RunWithContext(context.Background()); err == nil {
		t.Fatal("expected unsupported mode to return an error")
	}
}

func TestRunnerRunWithContextRejectsPitchforkWithoutAuth(t *testing.T) {
	runner := NewRunner(&RunnerOption{Mod: ModPitchFork})

	if err := runner.RunWithContext(context.Background()); err == nil {
		t.Fatal("expected pitchfork without auth to return an error")
	}
}
