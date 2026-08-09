package mongo

import (
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMatchingDocumentFieldsFiltersByFieldName(t *testing.T) {
	raw, err := bson.Marshal(bson.M{
		"username": "alice",
		"password": "secret",
		"profile": bson.M{
			"email": "alice@example.com",
			"city":  "Shanghai",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	fields, err := matchingDocumentFields(raw, []string{"password", "email"})
	if err != nil {
		t.Fatalf("matchingDocumentFields: %v", err)
	}
	joined := strings.Join(fields, "\n")
	for _, want := range []string{"password: secret", "profile.email: alice@example.com"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("matched fields %q do not contain %q", joined, want)
		}
	}
	for _, unwanted := range []string{"username", "profile.city"} {
		if strings.Contains(joined, unwanted) {
			t.Fatalf("matched fields %q unexpectedly contain %q", joined, unwanted)
		}
	}
}

func TestMatchingDocumentFieldsReturnsEmptyWithoutPatterns(t *testing.T) {
	raw, err := bson.Marshal(bson.M{"password": "secret"})
	if err != nil {
		t.Fatal(err)
	}
	fields, err := matchingDocumentFields(raw, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 0 {
		t.Fatalf("matched fields = %#v, want empty", fields)
	}
}
