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

func TestParsePayloadArgsPreservesRepeatedKeys(t *testing.T) {
	got, err := parsePayloadArgs([]string{"key=a", "key=b", "cmd=id"})
	if err != nil {
		t.Fatalf("parsePayloadArgs: %v", err)
	}
	keyVals, ok := got["key"].([]string)
	if !ok {
		t.Fatalf("expected key payload to be []string, got %#v", got["key"])
	}
	if len(keyVals) != 2 || keyVals[0] != "a" || keyVals[1] != "b" {
		t.Fatalf("unexpected key payload values: %#v", keyVals)
	}
	cmdVals, ok := got["cmd"].([]string)
	if !ok || len(cmdVals) != 1 || cmdVals[0] != "id" {
		t.Fatalf("unexpected cmd payload values: %#v", got["cmd"])
	}
}

func TestParsePayloadArgsRejectsInvalid(t *testing.T) {
	if _, err := parsePayloadArgs([]string{"cmd"}); err == nil {
		t.Fatal("expected invalid key=value to fail")
	}
}
