package pkg

import "testing"

func TestSplitUserDomain(t *testing.T) {
	user, domain := SplitUserDomain("DOMAIN/user")
	if user != "user" || domain != "DOMAIN" {
		t.Fatalf("unexpected split result: user=%q domain=%q", user, domain)
	}
}

func TestSplitUserDomainWithoutDomain(t *testing.T) {
	user, domain := SplitUserDomain("user")
	if user != "user" || domain != "" {
		t.Fatalf("unexpected split result: user=%q domain=%q", user, domain)
	}
}
