package core

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils/fileutils"
	"github.com/chainreactors/utils/parsers"
	"github.com/chainreactors/zombie/pkg"
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

func TestCLIResultHandlerOwnsFormattingAndForwardsEveryResult(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "results.jsonl")
	outputFile, err := fileutils.NewFile(outputPath, fileutils.ModeAppend, false, false)
	if err != nil {
		t.Fatal(err)
	}

	oldLog := logs.Log
	logs.Log = logs.NewLogger(oldLog.Level)
	logs.Log.SetOutput(&bytes.Buffer{})
	defer func() { logs.Log = oldLog }()

	var forwarded []*pkg.Result
	handler := cliResultHandler(parsers.ZombieFormatString, parsers.ZombieFormatJSONLine, outputFile, func(result *pkg.Result) {
		forwarded = append(forwarded, result)
	})
	success := pkg.NewResult(&pkg.Task{ZombieResult: &parsers.ZombieResult{IP: "127.0.0.1", Port: "6379", Service: "redis"}}, nil)
	failure := pkg.NewResult(&pkg.Task{ZombieResult: &parsers.ZombieResult{IP: "127.0.0.1", Port: "1", Service: "redis"}}, errors.New("connection refused"))
	handler(success)
	handler(failure)
	if err := outputFile.Close(); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), `"ok":true`) || strings.Contains(string(content), "connection refused") {
		t.Fatalf("unexpected success-only file output: %s", content)
	}
	if len(forwarded) != 2 || forwarded[0] != success || forwarded[1] != failure {
		t.Fatalf("forwarded results = %#v", forwarded)
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
