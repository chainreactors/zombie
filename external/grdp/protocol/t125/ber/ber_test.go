package ber

import (
	"bytes"
	"testing"
)

func TestReadIntegerThreeBytes(t *testing.T) {
	got, err := ReadInteger(bytes.NewReader([]byte{TAG_INTEGER, 0x03, 0x01, 0x02, 0x03}))
	if err != nil {
		t.Fatal(err)
	}
	if got != 0x010203 {
		t.Fatalf("expected 0x010203, got %#x", got)
	}
}
