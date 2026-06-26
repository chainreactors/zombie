package service

import (
	"testing"

	"github.com/chainreactors/neutron/operators"
	"github.com/chainreactors/neutron/protocols"
	"gopkg.in/yaml.v3"
)

// mockShellSession implements pkg.Session + pkg.ShellSession
type mockShellSession struct {
	svc     string
	outputs map[string]string
}

func (m *mockShellSession) Service() string  { return m.svc }
func (m *mockShellSession) Close() error     { return nil }
func (m *mockShellSession) Raw() interface{} { return nil }
func (m *mockShellSession) Exec(cmd string) ([]byte, error) {
	if out, ok := m.outputs[cmd]; ok {
		return []byte(out), nil
	}
	return []byte(""), nil
}

// mockKVSession implements pkg.Session + pkg.KVSession + RawCommander
type mockKVSession struct {
	svc   string
	data  map[string]string
	cmds  map[string]string
	calls []string
}

func (m *mockKVSession) Service() string  { return m.svc }
func (m *mockKVSession) Close() error     { return nil }
func (m *mockKVSession) Raw() interface{} { return nil }
func (m *mockKVSession) Get(key string) ([]byte, error) {
	return []byte(m.data[key]), nil
}
func (m *mockKVSession) Keys(pattern string) ([]string, error) {
	var keys []string
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}
func (m *mockKVSession) Command(name string, args ...string) (interface{}, error) {
	key := name
	for _, a := range args {
		key += " " + a
	}
	m.calls = append(m.calls, key)
	if out, ok := m.cmds[key]; ok {
		return out, nil
	}
	return "OK", nil
}

func TestRequestCompile(t *testing.T) {
	yamlData := `
ops:
  - shell: "id"
    name: whoami
matchers:
  - type: word
    part: whoami
    words: ["root"]
extractors:
  - type: regex
    name: user
    part: whoami
    regex: ['uid=\d+\((\w+)\)']
    group: 1
`
	var req Request
	if err := yaml.Unmarshal([]byte(yamlData), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(req.Ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(req.Ops))
	}
	if req.Ops[0].Shell != "id" {
		t.Fatalf("expected shell='id', got %q", req.Ops[0].Shell)
	}
	if req.CompiledOperators == nil {
		t.Fatal("expected compiled operators")
	}
}

func TestExecuteShell(t *testing.T) {
	session := &mockShellSession{
		svc: "ssh",
		outputs: map[string]string{
			"id":       "uid=0(root) gid=0(root)",
			"uname -a": "Linux box 5.15.0 x86_64",
		},
	}

	yamlData := `
ops:
  - shell: "id"
    name: whoami
  - shell: "uname -a"
    name: uname
matchers:
  - type: word
    part: whoami
    words: ["root"]
extractors:
  - type: regex
    name: user
    part: whoami
    regex: ['uid=\d+\((\w+)\)']
    group: 1
`
	var req Request
	if err := yaml.Unmarshal([]byte(yamlData), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}

	payloads := map[string]interface{}{"_session": session}
	scanCtx := protocols.NewScanContext("10.0.0.1:22", payloads)

	var result *operators.Result
	err := req.ExecuteWithResults(scanCtx, nil, nil, func(event *protocols.InternalWrappedEvent) {
		if event.OperatorsResult != nil {
			result = event.OperatorsResult
		}
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil {
		t.Fatal("expected result")
	}
	if !result.Matched {
		t.Error("expected match on 'root'")
	}
	if len(result.Extracts["user"]) == 0 {
		t.Error("expected extraction of user")
	} else if result.Extracts["user"][0] != "root" {
		t.Errorf("expected extracted user='root', got %q", result.Extracts["user"][0])
	}
}

func TestExecuteKV(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{"password_key": "s3cret"},
		cmds: map[string]string{
			"CONFIG GET dir": "dir\n/var/lib/redis",
		},
	}

	yamlData := `
ops:
  - kv: "CONFIG GET dir"
    name: dir
matchers:
  - type: word
    part: dir
    words: ["/var"]
extractors:
  - type: regex
    name: redis_dir
    part: dir
    regex: ['dir\s+(.+)']
    group: 1
`
	var req Request
	if err := yaml.Unmarshal([]byte(yamlData), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}

	payloads := map[string]interface{}{"_session": session}
	scanCtx := protocols.NewScanContext("10.0.0.1:6379", payloads)

	var result *operators.Result
	err := req.ExecuteWithResults(scanCtx, nil, nil, func(event *protocols.InternalWrappedEvent) {
		if event.OperatorsResult != nil {
			result = event.OperatorsResult
		}
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil {
		t.Fatal("expected result")
	}
	if !result.Matched {
		t.Error("expected match on '/var'")
	}
	if len(result.Extracts["redis_dir"]) == 0 {
		t.Error("expected extraction of redis_dir")
	} else if result.Extracts["redis_dir"][0] != "/var/lib/redis" {
		t.Errorf("expected '/var/lib/redis', got %q", result.Extracts["redis_dir"][0])
	}
}

func TestExecuteKVQuotedValue(t *testing.T) {
	session := &mockKVSession{
		svc: "redis",
		cmds: map[string]string{
			"SET x <?php system($_GET[0]);?>": "OK",
		},
	}

	out, err := execKV(session, `SET x "<?php system($_GET[0]);?>"`)
	if err != nil {
		t.Fatalf("execKV: %v", err)
	}
	if out != "OK" {
		t.Fatalf("expected OK, got %q", out)
	}
}

func TestParseCommandFieldsEscapes(t *testing.T) {
	fields, err := parseCommandFields(`SET x "\nline two\n"`)
	if err != nil {
		t.Fatalf("parseCommandFields: %v", err)
	}
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d: %#v", len(fields), fields)
	}
	if fields[2] != "\nline two\n" {
		t.Fatalf("unexpected escaped payload: %#v", fields[2])
	}
}

func TestFormatCommandResultArray(t *testing.T) {
	out := formatCommandResult([]interface{}{"dir", "/var/www"})
	if out != "dir\n/var/www" {
		t.Fatalf("unexpected formatted result: %q", out)
	}
}

func TestExecuteKVGetKeys(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{"secret": "val1", "password": "val2"},
	}

	yamlData := `
ops:
  - kv: "GET secret"
    name: secret_val
  - kv: "KEYS *"
    name: all_keys
matchers:
  - type: word
    part: secret_val
    words: ["val1"]
`
	var req Request
	if err := yaml.Unmarshal([]byte(yamlData), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}

	payloads := map[string]interface{}{"_session": session}
	scanCtx := protocols.NewScanContext("10.0.0.1:6379", payloads)

	var result *operators.Result
	err := req.ExecuteWithResults(scanCtx, nil, nil, func(event *protocols.InternalWrappedEvent) {
		if event.OperatorsResult != nil {
			result = event.OperatorsResult
		}
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || !result.Matched {
		t.Error("expected match on 'val1'")
	}
}

func TestTemplateMultiBlock(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{},
		cmds: map[string]string{
			"CONFIG GET dir":                 "dir\n/var/lib/redis",
			"CONFIG GET dbfilename":          "dbfilename\ndump.rdb",
			"CONFIG SET dir /tmp":            "OK",
			"CONFIG SET dbfilename test.rdb": "OK",
		},
	}

	yamlData := `
id: redis-multi-block
service: [redis]
info:
  name: Redis Multi Block Test
  severity: info

services:
  - ops:
      - kv: "CONFIG GET dir"
        name: dir
    extractors:
      - type: regex
        name: orig_dir
        internal: true
        part: dir
        regex: ['dir\s+(.+)']
        group: 1

  - ops:
      - kv: "CONFIG SET dir /tmp"
    matchers:
      - type: word
        words: ["OK"]
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}
	if !tmpl.Match("redis") {
		t.Fatal("expected match on redis")
	}
	if tmpl.Match("mysql") {
		t.Fatal("expected no match on mysql")
	}

	result, err := tmpl.Execute(session, "10.0.0.1:6379")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || !result.Matched {
		t.Error("expected matched result from second block")
	}
}

func TestTemplateVariablesAndCLIOverride(t *testing.T) {
	session := &mockShellSession{
		svc: "ssh",
		outputs: map[string]string{
			"id":     "uid=1000(zombie)",
			"whoami": "root",
		},
	}

	yamlData := `
id: variable-template
service: [ssh]
variables:
  cmd: whoami
services:
  - ops:
      - shell: "{{cmd}}"
        name: command_output
    matchers:
      - type: word
        part: command_output
        words: ["zombie"]
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := tmpl.ExecuteWithVariables(session, "127.0.0.1:22", map[string]interface{}{"cmd": "id"})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || !result.Matched {
		t.Fatal("expected CLI variable override to match id output")
	}
}

func TestTemplatePayloadsClusterbomb(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{},
	}

	yamlData := `
id: payload-template
service: [redis]
services:
  - attack: clusterbomb
    payloads:
      key:
        - a
        - b
      value:
        - "1"
        - "2"
    ops:
      - kv: "SET §key§ {{value}}"
        name: set_result
    extractors:
      - type: regex
        name: set_ok
        part: set_result
        regex: ['(OK)']
        group: 1
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := tmpl.Execute(session, "127.0.0.1:6379")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || len(result.Extracts["set_ok"]) != 4 {
		t.Fatalf("expected 4 payload executions, got result %#v", result)
	}
	if len(session.calls) != 4 {
		t.Fatalf("expected 4 redis commands, got %d: %#v", len(session.calls), session.calls)
	}
	expected := map[string]struct{}{
		"SET a 1": {},
		"SET a 2": {},
		"SET b 1": {},
		"SET b 2": {},
	}
	for _, call := range session.calls {
		if _, ok := expected[call]; !ok {
			t.Fatalf("unexpected call %q in %#v", call, session.calls)
		}
	}
}

func TestTemplatePayloadCLIOverride(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{},
	}

	yamlData := `
id: payload-override-template
service: [redis]
services:
  - attack: pitchfork
    payloads:
      key:
        - a
        - b
    ops:
      - kv: "SET §key§ 1"
        name: set_result
    extractors:
      - type: regex
        name: set_ok
        part: set_result
        regex: ['(OK)']
        group: 1
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := tmpl.ExecuteWithVariables(session, "127.0.0.1:6379", map[string]interface{}{"key": "cli"})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || len(result.Extracts["set_ok"]) != 1 {
		t.Fatalf("expected one overridden payload execution, got %#v", result)
	}
	if len(session.calls) != 1 || session.calls[0] != "SET cli 1" {
		t.Fatalf("unexpected calls: %#v", session.calls)
	}
}

func TestTemplateExplicitPayloadCLIOverride(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{},
	}

	yamlData := `
id: payload-explicit-override-template
service: [redis]
services:
  - attack: pitchfork
    payloads:
      key:
        - a
        - b
    ops:
      - kv: "SET §key§ 1"
        name: set_result
    extractors:
      - type: regex
        name: set_ok
        part: set_result
        regex: ['(OK)']
        group: 1
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := tmpl.ExecuteWithOptions(session, "127.0.0.1:6379", nil, map[string]interface{}{
		"key": []string{"cli-a", "cli-b"},
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || len(result.Extracts["set_ok"]) != 2 {
		t.Fatalf("expected two overridden payload executions, got %#v", result)
	}
	expected := []string{"SET cli-a 1", "SET cli-b 1"}
	if len(session.calls) != len(expected) {
		t.Fatalf("unexpected call count: %#v", session.calls)
	}
	for i, call := range expected {
		if session.calls[i] != call {
			t.Fatalf("unexpected calls: %#v", session.calls)
		}
	}
}

func TestTemplateExplicitPayloadBeatsVarOverride(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{},
	}

	yamlData := `
id: payload-precedence-template
service: [redis]
services:
  - attack: pitchfork
    payloads:
      key:
        - a
    ops:
      - kv: "SET §key§ 1"
        name: set_result
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}

	_, err := tmpl.ExecuteWithOptions(
		session,
		"127.0.0.1:6379",
		map[string]interface{}{"key": "var-value"},
		map[string]interface{}{"key": []string{"payload-value"}},
	)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(session.calls) != 1 || session.calls[0] != "SET payload-value 1" {
		t.Fatalf("unexpected calls: %#v", session.calls)
	}
}

func TestTemplateExplicitPayloadCanDefinePayloadSet(t *testing.T) {
	session := &mockKVSession{
		svc:  "redis",
		data: map[string]string{},
	}

	yamlData := `
id: payload-cli-defined-template
service: [redis]
services:
  - attack: pitchfork
    ops:
      - kv: "SET §key§ 1"
        name: set_result
    extractors:
      - type: regex
        name: set_ok
        part: set_result
        regex: ['(OK)']
        group: 1
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(yamlData), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(nil); err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := tmpl.ExecuteWithOptions(session, "127.0.0.1:6379", nil, map[string]interface{}{
		"key": []string{"cli-a", "cli-b"},
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || len(result.Extracts["set_ok"]) != 2 {
		t.Fatalf("expected CLI-defined payload executions, got %#v", result)
	}
	expected := []string{"SET cli-a 1", "SET cli-b 1"}
	if len(session.calls) != len(expected) {
		t.Fatalf("unexpected call count: %#v", session.calls)
	}
	for i, call := range expected {
		if session.calls[i] != call {
			t.Fatalf("unexpected calls: %#v", session.calls)
		}
	}
}

func TestResponsePreserved(t *testing.T) {
	session := &mockShellSession{
		svc:     "ssh",
		outputs: map[string]string{"id": "uid=0(root)", "hostname": "box1"},
	}

	yamlData := `
ops:
  - shell: "id"
    name: whoami
  - shell: "hostname"
    name: host
matchers:
  - type: word
    part: whoami
    words: ["root"]
`
	var req Request
	if err := yaml.Unmarshal([]byte(yamlData), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}

	payloads := map[string]interface{}{"_session": session}
	scanCtx := protocols.NewScanContext("10.0.0.1:22", payloads)

	var result *operators.Result
	err := req.ExecuteWithResults(scanCtx, nil, nil, func(event *protocols.InternalWrappedEvent) {
		if event.OperatorsResult != nil {
			result = event.OperatorsResult
		}
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil {
		t.Fatal("expected result")
	}
	if result.Response == "" {
		t.Fatal("Response should contain raw op output")
	}
	if !containsStr(result.Response, "uid=0(root)") {
		t.Errorf("Response missing id output, got %q", result.Response)
	}
	if !containsStr(result.Response, "box1") {
		t.Errorf("Response missing hostname output, got %q", result.Response)
	}
}

func TestResponsePreservedWithoutMatch(t *testing.T) {
	session := &mockShellSession{
		svc:     "ssh",
		outputs: map[string]string{"echo hello": "hello"},
	}

	tmplYaml := `
id: no-match-tmpl
service: [ssh]
info:
  name: No Match
  severity: info
services:
  - ops:
      - shell: "echo hello"
        name: greeting
    matchers:
      - type: word
        part: greeting
        words: ["NOMATCH"]
`
	var tmpl Template
	if err := yaml.Unmarshal([]byte(tmplYaml), &tmpl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := tmpl.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}

	result, err := tmpl.Execute(session, "10.0.0.1:22")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Response == "" {
		t.Fatal("Response should be populated even when matchers don't match")
	}
	if !containsStr(result.Response, "hello") {
		t.Errorf("Response missing op output, got %q", result.Response)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && findSubstr(s, sub))
}

func findSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestLegacyOpsCompat(t *testing.T) {
	session := &mockShellSession{
		svc:     "ssh",
		outputs: map[string]string{"id": "uid=0(root)"},
	}

	yamlData := `
ops:
  - exec: "id"
    name: whoami
matchers:
  - type: word
    part: whoami
    words: ["root"]
`
	var req Request
	if err := yaml.Unmarshal([]byte(yamlData), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Compile(&protocols.ExecuterOptions{Options: &protocols.Options{}}); err != nil {
		t.Fatalf("compile: %v", err)
	}

	payloads := map[string]interface{}{"_session": session}
	scanCtx := protocols.NewScanContext("10.0.0.1:22", payloads)

	var result *operators.Result
	err := req.ExecuteWithResults(scanCtx, nil, nil, func(event *protocols.InternalWrappedEvent) {
		if event.OperatorsResult != nil {
			result = event.OperatorsResult
		}
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result == nil || !result.Matched {
		t.Error("legacy exec: field should still work via normalizeOp")
	}
}
