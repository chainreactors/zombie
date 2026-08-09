package pkg

import "testing"

func TestEmbeddedAuditPatternsLoad(t *testing.T) {
	data := LoadEmbeddedConfig("zombie_audit")
	if len(data) == 0 {
		t.Fatal("embedded zombie_audit config is empty")
	}

	DefaultAuditPatterns = nil
	if err := LoadAuditConfig(); err != nil {
		t.Fatalf("LoadAuditConfig: %v", err)
	}
	if len(DefaultAuditPatterns) == 0 {
		t.Fatal("audit patterns were not loaded")
	}
}
