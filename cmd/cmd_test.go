package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunReturnsZeroForVersion(t *testing.T) {
	var out bytes.Buffer

	if code := Run([]string{"--version"}, &out); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if strings.TrimSpace(out.String()) != "dev" {
		t.Fatalf("unexpected version output: %q", out.String())
	}
}

func TestRunReturnsNonZeroForInvalidFlag(t *testing.T) {
	var out bytes.Buffer

	if code := Run([]string{"--definitely-not-a-flag"}, &out); code == 0 {
		t.Fatal("expected non-zero exit code for invalid flag")
	}
	if !strings.Contains(out.String(), "unknown flag") {
		t.Fatalf("expected unknown flag output, got %q", out.String())
	}
}

func TestRunReturnsNonZeroForMissingInput(t *testing.T) {
	var out bytes.Buffer

	if code := Run(nil, &out); code == 0 {
		t.Fatal("expected non-zero exit code for missing input")
	}
	if !strings.Contains(out.String(), "please input ip") {
		t.Fatalf("expected missing input output, got %q", out.String())
	}
}
