package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chainreactors/utils/parsers"
	"github.com/chainreactors/zombie/pkg"
)

// --- Mock Sessions ---

type mockShellSession struct {
	files map[string][]byte
}

func (m *mockShellSession) Service() string  { return "ssh" }
func (m *mockShellSession) Close() error     { return nil }
func (m *mockShellSession) Raw() interface{} { return nil }
func (m *mockShellSession) Exec(cmd string) ([]byte, error) {
	for path, data := range m.files {
		if containsSubstr(cmd, path) {
			return data, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

type mockSQLSession struct {
	service string
	rows    map[string][][]string
}

func (m *mockSQLSession) Service() string  { return m.service }
func (m *mockSQLSession) Close() error     { return nil }
func (m *mockSQLSession) Raw() interface{} { return nil }
func (m *mockSQLSession) Query(query string, args ...any) ([][]string, error) {
	for key, rows := range m.rows {
		if containsSubstr(query, key) {
			return rows, nil
		}
	}
	return nil, fmt.Errorf("no results")
}
func (m *mockSQLSession) Databases() ([]string, error) {
	return []string{"testdb", "production"}, nil
}

type mockKVSession struct{}

func (m *mockKVSession) Service() string  { return "redis" }
func (m *mockKVSession) Close() error     { return nil }
func (m *mockKVSession) Raw() interface{} { return nil }
func (m *mockKVSession) Get(key string) ([]byte, error) {
	if key == "user:token" {
		return []byte("ghp_abcdefghij1234567890abcdefghij1234"), nil
	}
	return nil, nil
}
func (m *mockKVSession) Keys(pattern string) ([]string, error) {
	if pattern == "*" || pattern == "*token*" {
		return []string{"user:token"}, nil
	}
	return nil, nil
}

type mockFileSession struct{}

func (m *mockFileSession) Service() string  { return "ftp" }
func (m *mockFileSession) Close() error     { return nil }
func (m *mockFileSession) Raw() interface{} { return nil }
func (m *mockFileSession) List(path string) ([]string, error) {
	return []string{".env", "config.yaml", "data.csv"}, nil
}
func (m *mockFileSession) Read(path string) ([]byte, error) {
	if path == "/.env" {
		return []byte("DB_PASSWORD=SuperSecret123\nAPI_KEY=sk_live_abc123\n"), nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *mockFileSession) Write(path string, data []byte) error { return nil }

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func loadAndCreateServiceAction(t *testing.T, dir string) *ServiceAction {
	t.Helper()
	tmpls, err := LoadServiceTemplatesFromPaths([]string{dir})
	if err != nil {
		t.Fatalf("load templates: %v", err)
	}
	a, err := NewServiceAction(tmpls, nil)
	if err != nil {
		t.Fatalf("NewServiceAction: %v", err)
	}
	return a
}

func mockTask() *pkg.Task {
	return &pkg.Task{
		ZombieResult: &parsers.ZombieResult{
			IP:      "10.0.0.1",
			Port:    "22",
			Service: "ssh",
		},
		Timeout: 5,
	}
}

func createTestTemplate(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	tmpl := `id: test-secret-scan
info:
  name: Test Secret Scanner
  severity: high
file:
  - extensions:
      - all
    extractors:
      - type: regex
        regex:
          - "(?i)password\\s*[=:]\\s*(\\S+)"
        group: 1
      - type: regex
        regex:
          - "ghp_[A-Za-z0-9]{36}"
    matchers:
      - type: word
        words:
          - "password"
          - "ghp_"
`
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(tmpl), 0644)
	return dir
}

// --- PostAction Tests ---

func TestPostAction_ScanData(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir})
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	results := a.ScanData([]byte("password = hunter2\nclean line\n"), "test:label")
	if len(results) == 0 {
		t.Fatal("should find password in test data")
	}
	found := false
	for _, e := range results {
		for _, v := range e.ExtractResult {
			if v == "hunter2" {
				found = true
			}
		}
	}
	if !found {
		t.Error("should extract 'hunter2'")
	}
}

func TestPostAction_GitHubToken(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir})
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	token := "ghp_abcdefghijklmnopqrstuvwxyz1234567890"
	results := a.ScanData([]byte("GITHUB_TOKEN="+token+"\n"), "test:github")
	if len(results) == 0 {
		t.Fatal("should find GitHub token")
	}
	found := false
	for _, e := range results {
		for _, v := range e.ExtractResult {
			if v == token {
				found = true
			}
		}
	}
	if !found {
		t.Error("should extract GitHub token")
	}
}

func TestServiceActionChain(t *testing.T) {
	dir := t.TempDir()
	root := `id: root-chain
service: [ssh]
chain: [child-chain]
info:
  name: Root Chain
  severity: info
services:
  - ops:
      - shell: "detect-os"
        name: os_detect
    extractors:
      - type: regex
        name: os_type
        internal: true
        part: os_detect
        regex: ['(Linux)']
        group: 1
`
	child := `id: child-chain
service: [ssh]
info:
  name: Child Chain
  severity: info
services:
  - ops:
      - shell: "child-command"
        name: child_output
    extractors:
      - type: regex
        name: child_value
        part: child_output
        regex: ['(child-ok)']
        group: 1
`
	if err := os.WriteFile(filepath.Join(dir, "root-chain.yaml"), []byte(root), 0644); err != nil {
		t.Fatalf("write root template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "child-chain.yaml"), []byte(child), 0644); err != nil {
		t.Fatalf("write child template: %v", err)
	}

	a := loadAndCreateServiceAction(t, dir)
	session := &mockShellSession{
		files: map[string][]byte{
			"detect-os":     []byte("Linux\n"),
			"child-command": []byte("child-ok\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, extracted := range result.Extracteds {
		if extracted.Name == "child-chain:child_value" && len(extracted.ExtractResult) == 1 && extracted.ExtractResult[0] == "child-ok" {
			return
		}
	}
	t.Fatalf("expected chained extraction, got %#v", result.Extracteds)
}

func TestServiceAction_ChainTargetNotEntrypoint(t *testing.T) {
	dir := t.TempDir()
	// root chains to child; child should NOT run as a top-level entry point
	root := `id: entry
service: [ssh]
chain: [helper]
info:
  name: Entry
  severity: info
services:
  - ops:
      - shell: "echo entry"
        name: entry_out
    extractors:
      - type: regex
        name: entry_val
        part: entry_out
        regex: ['(entry)']
        group: 1
`
	helper := `id: helper
service: [ssh]
info:
  name: Helper
  severity: info
services:
  - ops:
      - shell: "echo helper"
        name: helper_out
    extractors:
      - type: regex
        name: helper_val
        part: helper_out
        regex: ['(helper)']
        group: 1
`
	os.WriteFile(filepath.Join(dir, "entry.yaml"), []byte(root), 0644)
	os.WriteFile(filepath.Join(dir, "helper.yaml"), []byte(helper), 0644)

	a := loadAndCreateServiceAction(t, dir)
	session := &mockShellSession{
		files: map[string][]byte{
			"echo entry":  []byte("entry\n"),
			"echo helper": []byte("helper\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	// Both should appear in results (entry as entry point, helper via chain)
	found := map[string]bool{}
	for _, e := range result.Extracteds {
		found[e.Name] = true
	}
	if !found["entry:entry_val"] {
		t.Error("missing entry extraction")
	}
	if !found["helper:helper_val"] {
		t.Error("missing helper extraction (should run via chain)")
	}

	// helper should appear exactly once (not duplicated as both entry point and chain)
	count := 0
	for _, e := range result.Extracteds {
		if e.Name == "helper:helper_val" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("helper executed %d times, want 1", count)
	}
}

func TestServiceAction_ServiceMismatchSkipsChain(t *testing.T) {
	dir := t.TempDir()
	root := `id: ssh-root
service: [ssh]
chain: [mysql-only]
info:
  name: SSH Root
  severity: info
services:
  - ops:
      - shell: "echo root"
        name: root_out
    extractors:
      - type: regex
        name: root_val
        part: root_out
        regex: ['(root)']
        group: 1
`
	mysqlOnly := `id: mysql-only
service: [mysql]
info:
  name: MySQL Only
  severity: info
services:
  - ops:
      - shell: "echo mysql"
        name: mysql_out
    extractors:
      - type: regex
        name: mysql_val
        part: mysql_out
        regex: ['(mysql)']
        group: 1
`
	os.WriteFile(filepath.Join(dir, "root.yaml"), []byte(root), 0644)
	os.WriteFile(filepath.Join(dir, "mysql.yaml"), []byte(mysqlOnly), 0644)

	a := loadAndCreateServiceAction(t, dir)
	// session is SSH, so mysql-only should be skipped
	session := &mockShellSession{
		files: map[string][]byte{
			"echo root":  []byte("root\n"),
			"echo mysql": []byte("mysql\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	for _, e := range result.Extracteds {
		if e.Name == "mysql-only:mysql_val" {
			t.Fatal("mysql-only template should NOT execute on SSH session")
		}
	}
	found := false
	for _, e := range result.Extracteds {
		if e.Name == "ssh-root:root_val" {
			found = true
		}
	}
	if !found {
		t.Error("ssh-root should have executed")
	}
}

func TestServiceAction_ServiceMismatchStopsChain(t *testing.T) {
	// When a chain target doesn't match the session's service, it returns nil
	// and its own chains (if any) should NOT execute.
	dir := t.TempDir()
	root := `id: ssh-entry
service: [ssh]
chain: [mysql-gate]
info:
  name: SSH Entry
  severity: info
services:
  - ops:
      - shell: "echo entry"
        name: entry_out
    extractors:
      - type: regex
        name: entry_val
        part: entry_out
        regex: ['(entry)']
        group: 1
`
	gate := `id: mysql-gate
service: [mysql]
chain: [after-gate]
info:
  name: MySQL Gate
  severity: info
services:
  - ops:
      - shell: "echo gate"
        name: gate_out
    extractors:
      - type: regex
        name: gate_val
        part: gate_out
        regex: ['(gate)']
        group: 1
`
	afterGate := `id: after-gate
service: [ssh]
info:
  name: After Gate
  severity: info
services:
  - ops:
      - shell: "echo after"
        name: after_out
    extractors:
      - type: regex
        name: after_val
        part: after_out
        regex: ['(after)']
        group: 1
`
	os.WriteFile(filepath.Join(dir, "entry.yaml"), []byte(root), 0644)
	os.WriteFile(filepath.Join(dir, "gate.yaml"), []byte(gate), 0644)
	os.WriteFile(filepath.Join(dir, "after.yaml"), []byte(afterGate), 0644)

	a := loadAndCreateServiceAction(t, dir)
	session := &mockShellSession{
		files: map[string][]byte{
			"echo entry": []byte("entry\n"),
			"echo gate":  []byte("gate\n"),
			"echo after": []byte("after\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	found := map[string]bool{}
	for _, e := range result.Extracteds {
		found[e.Name] = true
	}
	if !found["ssh-entry:entry_val"] {
		t.Error("ssh-entry should have executed")
	}
	if found["mysql-gate:gate_val"] {
		t.Error("mysql-gate should NOT execute on SSH session")
	}
	if found["after-gate:after_val"] {
		t.Error("after-gate should NOT execute because mysql-gate was skipped")
	}
}

func TestServiceAction_MultipleChains(t *testing.T) {
	dir := t.TempDir()
	root := `id: multi-root
service: [ssh]
chain: [branch-a, branch-b]
info:
  name: Multi Root
  severity: info
services:
  - ops:
      - shell: "echo root"
        name: root_out
    extractors:
      - type: regex
        name: root_val
        part: root_out
        regex: ['(root)']
        group: 1
`
	branchA := `id: branch-a
service: [ssh]
info:
  name: Branch A
  severity: info
services:
  - ops:
      - shell: "echo a"
        name: a_out
    extractors:
      - type: regex
        name: a_val
        part: a_out
        regex: ['(a)']
        group: 1
`
	branchB := `id: branch-b
service: [ssh]
info:
  name: Branch B
  severity: info
services:
  - ops:
      - shell: "echo b"
        name: b_out
    extractors:
      - type: regex
        name: b_val
        part: b_out
        regex: ['(b)']
        group: 1
`
	os.WriteFile(filepath.Join(dir, "root.yaml"), []byte(root), 0644)
	os.WriteFile(filepath.Join(dir, "a.yaml"), []byte(branchA), 0644)
	os.WriteFile(filepath.Join(dir, "b.yaml"), []byte(branchB), 0644)

	a := loadAndCreateServiceAction(t, dir)
	session := &mockShellSession{
		files: map[string][]byte{
			"echo root": []byte("root\n"),
			"echo a":    []byte("a\n"),
			"echo b":    []byte("b\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	found := map[string]bool{}
	for _, e := range result.Extracteds {
		found[e.Name] = true
	}
	for _, want := range []string{"multi-root:root_val", "branch-a:a_val", "branch-b:b_val"} {
		if !found[want] {
			t.Errorf("missing extraction %s, got %v", want, found)
		}
	}
}

// --- ServiceAction Loot Population ---

func TestServiceAction_LootPopulated(t *testing.T) {
	dir := t.TempDir()
	tmpl := `id: loot-test
service: [ssh]
info:
  name: Loot Test
  severity: info
services:
  - ops:
      - shell: "cat /etc/passwd"
        name: passwd
    extractors:
      - type: regex
        name: users
        part: passwd
        regex: ['(\w+):']
`
	os.WriteFile(filepath.Join(dir, "loot.yaml"), []byte(tmpl), 0644)
	a := loadAndCreateServiceAction(t, dir)
	session := &mockShellSession{
		files: map[string][]byte{
			"cat /etc/passwd": []byte("root:x:0:0:root:/root:/bin/bash\nwww:x:33:33:www-data:/var/www\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Loot == nil || len(result.Loot) == 0 {
		t.Fatal("Loot should be populated with raw response data")
	}
	data, ok := result.Loot["loot-test"]
	if !ok {
		t.Fatalf("Loot missing key 'loot-test', got keys: %v", lootKeys(result.Loot))
	}
	if !strings.Contains(string(data), "root:x:0:0") {
		t.Errorf("Loot should contain raw passwd output, got %q", string(data))
	}
}

func TestServiceAction_LootPopulatedWithoutMatch(t *testing.T) {
	dir := t.TempDir()
	tmpl := `id: nomatch-loot
service: [ssh]
info:
  name: No Match Loot
  severity: info
services:
  - ops:
      - shell: "whoami"
        name: user
    matchers:
      - type: word
        part: user
        words: ["WILL_NOT_MATCH"]
`
	os.WriteFile(filepath.Join(dir, "nomatch.yaml"), []byte(tmpl), 0644)
	a := loadAndCreateServiceAction(t, dir)
	session := &mockShellSession{
		files: map[string][]byte{
			"whoami": []byte("testuser\n"),
		},
	}
	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Loot == nil || len(result.Loot) == 0 {
		t.Fatal("Loot should be populated even when matchers don't match")
	}
	if !strings.Contains(string(result.Loot["nomatch-loot"]), "testuser") {
		t.Error("Loot should contain raw response regardless of match result")
	}
}

func lootKeys(loot map[string][]byte) []string {
	var keys []string
	for k := range loot {
		keys = append(keys, k)
	}
	return keys
}

// --- PostAction from Data ---

func TestNewPostActionFromData(t *testing.T) {
	yamlData := []byte(`- id: test-loot-rule
  info:
    name: Test Password
    severity: high
  file:
    - extensions: [all]
      matchers:
        - type: word
          words: ["password"]
      extractors:
        - type: regex
          regex: ['(?i)password\s*[=:]\s*(\S+)']
          group: 1
`)
	pa, err := NewPostActionFromData(yamlData)
	if err != nil {
		t.Fatalf("NewPostActionFromData: %v", err)
	}
	results := pa.ScanData([]byte("password = secret123\n"), "test")
	if len(results) == 0 {
		t.Fatal("should detect password in scanned data")
	}
	found := false
	for _, e := range results {
		for _, v := range e.ExtractResult {
			if v == "secret123" {
				found = true
			}
		}
	}
	if !found {
		t.Error("should extract 'secret123' from loot data")
	}
}

func TestNewPostActionFromData_Empty(t *testing.T) {
	_, err := NewPostActionFromData(nil)
	if err == nil {
		t.Error("expected error for nil data")
	}
	_, err = NewPostActionFromData([]byte{})
	if err == nil {
		t.Error("expected error for empty data")
	}
}

// --- Worker Integration Test ---

func TestPostAction_ScanLoot(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir})
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	results := a.ScanData([]byte("[client]\npassword = dbpass123\n"), "ssh:10.0.0.1:22:~/.my.cnf")
	if len(results) == 0 {
		t.Fatal("should find password in loot data")
	}
}

// --- Audit E2E ---

type mockAuditableSession struct {
	svc  string
	data map[string]string
}

func (m *mockAuditableSession) Service() string  { return m.svc }
func (m *mockAuditableSession) Close() error     { return nil }
func (m *mockAuditableSession) Raw() interface{} { return nil }
func (m *mockAuditableSession) Audit(patterns []string, limit int) (map[string]string, error) {
	return m.data, nil
}

func TestAuditAction_E2E(t *testing.T) {
	session := &mockAuditableSession{
		svc: "mysql",
		data: map[string]string{
			"app.users.phone":    "13800138000\n13912345678",
			"app.users.password": "admin123\nP@ssw0rd",
			"app.config.api_key": "AKIA1234567890ABCDEF",
		},
	}

	// Step 1: AuditAction produces loot
	audit := NewAuditAction()
	ar, err := audit.Run(session, &pkg.Task{ZombieResult: &parsers.ZombieResult{Service: "mysql"}})
	if err != nil {
		t.Fatalf("audit action: %v", err)
	}
	if ar == nil || len(ar.Loot) == 0 {
		t.Fatal("audit should produce loot")
	}
	if len(ar.Loot) != 3 {
		t.Fatalf("expected 3 loot entries, got %d", len(ar.Loot))
	}

	// Step 2: Merge into Result (simulates worker.go)
	result := &pkg.Result{Task: &pkg.Task{ZombieResult: &parsers.ZombieResult{Service: "mysql"}}, OK: true}
	result.Merge(ar)

	// Step 3: PostAction scans loot for PII (simulates worker.go postAction loop)
	lootDir := createLootTemplateDir(t)
	postAction, err := NewPostAction([]string{lootDir})
	if err != nil {
		t.Fatalf("post action: %v", err)
	}
	for label, data := range result.Loot {
		result.Extracteds = append(result.Extracteds, postAction.ScanData(data, label)...)
	}

	// Step 4: Verify findings include location labels
	if len(result.Extracteds) == 0 {
		t.Fatal("PostAction should find PII in audit loot")
	}
	var foundPhone, foundCloud bool
	for _, e := range result.Extracteds {
		if containsSubstr(e.Name, "phone") {
			foundPhone = true
		}
		if containsSubstr(e.Name, "cloud") || containsSubstr(e.Name, "credential") {
			foundCloud = true
		}
		t.Logf("finding: %s → %v", e.Name, e.ExtractResult)
	}
	if !foundPhone {
		t.Error("should detect phone numbers in app.users.phone")
	}
	if !foundCloud {
		t.Error("should detect cloud credential (AKIA) in app.config.api_key")
	}
}

func TestAuditAction_NonAuditableSession(t *testing.T) {
	session := &mockShellSession{files: map[string][]byte{}}
	audit := NewAuditAction()
	ar, err := audit.Run(session, &pkg.Task{ZombieResult: &parsers.ZombieResult{Service: "ssh"}})
	if err != nil {
		t.Fatalf("should not error: %v", err)
	}
	if ar != nil {
		t.Fatal("non-auditable session should return nil")
	}
}

func createLootTemplateDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	phone := `
id: loot-phone-cn
info:
  name: Chinese Phone Number
  severity: medium
file:
  - extensions:
      - all
    extractors:
      - type: regex
        regex:
          - '(?:\+?86[-\s]?)?1[3-9]\d{9}'
`
	cloud := `
id: loot-cloud-credential
info:
  name: Cloud Access Credential
  severity: high
file:
  - extensions:
      - all
    extractors:
      - type: regex
        regex:
          - '(?:AKIA|ASIA)[A-Z0-9]{16}'
`
	os.WriteFile(filepath.Join(dir, "phone.yaml"), []byte(phone), 0644)
	os.WriteFile(filepath.Join(dir, "cloud.yaml"), []byte(cloud), 0644)
	return dir
}
