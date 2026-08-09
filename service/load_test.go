package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chainreactors/neutron/protocols"
	"gopkg.in/yaml.v3"
)

func TestLoadAllTemplates(t *testing.T) {
	templatesDir := "../../proton/templates/services"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Skipf("templates dir not found: %s", templatesDir)
	}

	var total, passed, failed int
	filepath.WalkDir(templatesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}
		total++

		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			failed++
			return nil
		}

		var tmpl Template
		if err := yaml.Unmarshal(data, &tmpl); err != nil {
			t.Errorf("unmarshal %s: %v", path, err)
			failed++
			return nil
		}

		if tmpl.Id == "" {
			t.Errorf("%s: missing id", path)
			failed++
			return nil
		}
		if len(tmpl.Service) == 0 {
			t.Errorf("%s: missing service filter", path)
			failed++
			return nil
		}
		if len(tmpl.RequestsService) == 0 {
			t.Errorf("%s: no service request blocks", path)
			failed++
			return nil
		}

		if err := tmpl.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
			t.Errorf("%s: compile failed: %v", path, err)
			failed++
			return nil
		}

		for i, req := range tmpl.RequestsService {
			if len(req.Ops) == 0 {
				t.Errorf("%s: services[%d] has no ops", path, i)
				failed++
				return nil
			}
			for j, op := range req.Ops {
				if !hasAction(op) {
					t.Errorf("%s: services[%d].ops[%d] has no action field set", path, i, j)
					failed++
					return nil
				}
			}
		}

		passed++
		t.Logf("OK  %s (id=%s, service=%v, blocks=%d)", filepath.Base(path), tmpl.Id, tmpl.Service, len(tmpl.RequestsService))
		return nil
	})

	t.Logf("\n--- Summary: %d total, %d passed, %d failed ---", total, passed, failed)
	if failed > 0 {
		t.Fatalf("%d templates failed validation", failed)
	}
}

func TestLoadLootTemplates(t *testing.T) {
	lootDir := "../templates/zombie/loot"
	if _, err := os.Stat(lootDir); os.IsNotExist(err) {
		t.Skipf("loot templates dir not found: %s", lootDir)
	}

	var total, passed int
	filepath.WalkDir(lootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}
		total++

		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			return nil
		}

		var raw map[string]interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			t.Errorf("unmarshal %s: %v", path, err)
			return nil
		}
		if raw["id"] == nil {
			t.Errorf("%s: missing id", path)
			return nil
		}
		if raw["file"] == nil {
			t.Errorf("%s: missing file section", path)
			return nil
		}

		passed++
		t.Logf("OK  %s (id=%v)", filepath.Base(path), raw["id"])
		return nil
	})

	t.Logf("\n--- Loot templates: %d total, %d passed ---", total, passed)
	if total == 0 {
		t.Error("no loot templates found")
	}
}

func hasAction(op *Op) bool {
	return op.Shell != "" || op.DB != "" || op.KV != "" ||
		(op.File != nil && (op.File.List != "" || op.File.Read != "")) ||
		op.LDAP != nil ||
		op.Exec != "" || op.Query != "" ||
		op.Get != "" || op.Keys != "" || op.Cmd != "" ||
		op.List != "" || op.Read != "" || op.Search != nil
}
