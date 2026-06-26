package pkg

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/chainreactors/parsers"
)

func TestResultFormatJSONIncludesExtracteds(t *testing.T) {
	result := &Result{
		Task: &Task{ZombieResult: &parsers.ZombieResult{
			IP:       "127.0.0.1",
			Port:     "6379",
			Service:  "redis",
			Scheme:   "redis",
			Username: "default",
			Password: "pass",
			Mod:      parsers.ZombieModBrute,
		}},
		OK: true,
		Extracteds: parsers.Extracteds{
			{Name: "local-redis-smoke:redis_value", ExtractResult: []string{"zombie-template-ok"}},
		},
	}

	formatted := result.Format(parsers.ZombieFormatJSON)
	if !strings.Contains(formatted, "extracteds") {
		t.Fatalf("formatted result should contain extracteds, got %s", formatted)
	}

	var decoded Result
	if err := json.Unmarshal([]byte(formatted), &decoded); err != nil {
		t.Fatalf("unmarshal formatted result: %v", err)
	}
	if len(decoded.Extracteds) != 1 {
		t.Fatalf("decoded Extracteds = %d, want 1", len(decoded.Extracteds))
	}
	if decoded.Extracteds[0].Name != "local-redis-smoke:redis_value" {
		t.Fatalf("decoded extractor name = %q", decoded.Extracteds[0].Name)
	}
}
