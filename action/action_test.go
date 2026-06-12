package action

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chainreactors/parsers"
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

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
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
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	result := &pkg.ActionResult{}
	a.scanData([]byte("password = hunter2\nclean line\n"), "test:label", result)

	if len(result.Extracteds) == 0 {
		t.Fatal("should find password in test data")
	}
	found := false
	for _, e := range result.Extracteds {
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
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	token := "ghp_abcdefghijklmnopqrstuvwxyz1234567890"
	result := &pkg.ActionResult{}
	a.scanData([]byte("GITHUB_TOKEN="+token+"\n"), "test:github", result)

	if len(result.Extracteds) == 0 {
		t.Fatal("should find GitHub token")
	}
	found := false
	for _, e := range result.Extracteds {
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

func TestPostAction_Shell(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	session := &mockShellSession{
		files: map[string][]byte{
			"hostname":  []byte("prodserver\n"),
			"id":        []byte("uid=0(root)\n"),
			"~/.my.cnf": []byte("[client]\npassword = dbpass123\n"),
		},
	}

	result, err := a.Run(session, mockTask())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Loot) == 0 {
		t.Fatal("should produce loot")
	}

	hasProtonFinding := false
	for _, e := range result.Extracteds {
		if containsSubstr(e.Name, "test-secret-scan") {
			hasProtonFinding = true
		}
	}
	if !hasProtonFinding {
		t.Error("should have proton scan findings")
	}
}

func TestPostAction_SQL(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	session := &mockSQLSession{
		service: "mysql",
		rows: map[string][][]string{
			"mysql.user": {
				{"user", "host"},
				{"root", "localhost"},
			},
		},
	}
	task := mockTask()
	task.Service = "mysql"
	task.Port = "3306"

	result, err := a.Run(session, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundDB := false
	for _, e := range result.Extracteds {
		if e.Name == "databases" {
			foundDB = true
		}
	}
	if !foundDB {
		t.Error("should have extracted databases")
	}
}

func TestPostAction_KV(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	session := &mockKVSession{}
	task := mockTask()
	task.Service = "redis"

	result, err := a.Run(session, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Extracteds) == 0 {
		t.Fatal("should find GitHub token in Redis key")
	}
}

func TestPostAction_File(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	session := &mockFileSession{}
	task := mockTask()
	task.Service = "ftp"

	result, err := a.Run(session, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Loot) == 0 {
		t.Error("should collect .env as loot")
	}
}

// --- Worker Integration Test ---

func TestWorkerExecute_WithPostAction(t *testing.T) {
	dir := createTestTemplate(t)
	a, err := NewPostAction([]string{dir}, 100)
	if err != nil {
		t.Fatalf("NewPostAction failed: %v", err)
	}

	session := &mockShellSession{
		files: map[string][]byte{
			"hostname":       []byte("testhost\n"),
			"/etc/shadow":    []byte("root:$6$hash:18000:0:99999:7:::\n"),
			"~/.vault-token": []byte("s.abcdefghij1234567890\n"),
		},
	}

	task := mockTask()
	result := &pkg.Result{Task: task, OK: true}

	ar, err := a.Run(session, task)
	if err != nil {
		t.Fatalf("action failed: %v", err)
	}
	result.Merge(ar)

	if !result.OK {
		t.Fatal("result should be OK")
	}
	if len(result.Loot) == 0 {
		t.Fatal("should have loot")
	}

	t.Logf("Worker: %d extracteds, %d loot, %d action results",
		len(result.Extracteds), len(result.Loot), len(result.ActionResults))
}
