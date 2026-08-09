package pkg

import "testing"

func TestParseMethodUsesTaskRawMode(t *testing.T) {
	method, value := ParseMethod("pk:key-data", false)
	if method != "pk" || value != "key-data" {
		t.Fatalf("parsed method = %q, value = %q", method, value)
	}

	method, value = ParseMethod("pk:key-data", true)
	if method != "" || value != "pk:key-data" {
		t.Fatalf("raw method = %q, value = %q", method, value)
	}
}
