package core

import "testing"

func TestParseKeyValueArgs(t *testing.T) {
	got, err := parseKeyValueArgs([]string{"cmd=id", "outfile=/tmp/a b.txt"})
	if err != nil {
		t.Fatalf("parseKeyValueArgs: %v", err)
	}
	if got["cmd"] != "id" {
		t.Fatalf("unexpected cmd: %#v", got["cmd"])
	}
	if got["outfile"] != "/tmp/a b.txt" {
		t.Fatalf("unexpected outfile: %#v", got["outfile"])
	}
}

func TestParseKeyValueArgsRejectsInvalid(t *testing.T) {
	if _, err := parseKeyValueArgs([]string{"cmd"}); err == nil {
		t.Fatal("expected invalid key=value to fail")
	}
}
