//go:build docker

package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const (
	commandTimeout = 5 * time.Minute
	probeTimeout   = 30 * time.Second
	cleanupTimeout = 90 * time.Second
)

type zombieResult struct {
	Extracteds []extracted `json:"extracteds"`
}

type extracted struct {
	Name          string   `json:"name"`
	ExtractResult []string `json:"extract_result"`
}

type dockerEnv struct {
	repoRoot     string
	fixture      string
	compose      string
	templates    string
	templatesExt string
	outDir       string
}

func TestServiceTemplateDocker(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker not found: %v", err)
	}
	if out, err := runCmdWithTimeout("", probeTimeout, "docker", "version"); err != nil {
		t.Skipf("docker is not available: %v\n%s", err, out)
	}

	env := newDockerEnv(t)
	t.Cleanup(func() {
		_ = os.RemoveAll(env.outDir)
		_, _ = runCmdWithTimeout(env.repoRoot, cleanupTimeout, "docker", "compose", "-f", env.compose, "down", "-v")
	})

	t.Run("Smoke", func(t *testing.T) {
		env.smoke(t)
	})
	t.Run("Payload", func(t *testing.T) {
		env.payload(t)
	})
	t.Run("PostExploit", func(t *testing.T) {
		env.postExploit(t)
	})
	t.Run("Existing", func(t *testing.T) {
		if _, err := os.Stat(env.templatesExt); err != nil {
			t.Skipf("external proton templates not found: %v", err)
		}
		env.existing(t)
	})
	t.Run("Exploit", func(t *testing.T) {
		env.exploit(t)
	})
}

func newDockerEnv(t *testing.T) *dockerEnv {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve caller")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
	fixture := filepath.Join(repoRoot, "testdata", "service-template")
	templatesExt := os.Getenv("ZOMBIE_SERVICE_TEMPLATES_DIR")
	if templatesExt == "" {
		templatesExt = filepath.Join(repoRoot, "..", "proton", "templates", "services")
	}
	outDir := filepath.Join(fixture, "out")
	if err := os.RemoveAll(outDir); err != nil {
		t.Fatalf("clean output dir: %v", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	return &dockerEnv{
		repoRoot:     repoRoot,
		fixture:      fixture,
		compose:      filepath.Join(fixture, "compose.yaml"),
		templates:    filepath.Join(fixture, "templates"),
		templatesExt: filepath.Clean(templatesExt),
		outDir:       outDir,
	}
}

func (e *dockerEnv) smoke(t *testing.T) {
	e.composeUp(t, "redis", "postgres", "mysql", "ssh")
	waitTCP(t, "127.0.0.1:16379")
	waitTCP(t, "127.0.0.1:15432")
	waitTCP(t, "127.0.0.1:13306")
	waitTCP(t, "127.0.0.1:10022")
	e.initRedis(t)

	out := e.out(t, "smoke")
	tmpl := filepath.Join(e.templates, "smoke")
	assertExtraction(t, e.runZombie(t, out, "redis", "redis://:zombie_redis_pass@127.0.0.1:16379", tmpl), "local-redis-smoke:redis_value", "zombie-template-ok")
	assertExtraction(t, e.runZombie(t, out, "postgres", "postgresql://zombie:zombie_pg_pass@127.0.0.1:15432", tmpl), "local-postgres-smoke:postgres_marker", "zombie-postgres-ok")
	assertExtraction(t, e.runZombie(t, out, "mysql", "mysql://zombie:zombie_mysql_pass@127.0.0.1:13306", tmpl), "local-mysql-smoke:mysql_marker", "zombie-mysql-ok")
	assertExtraction(t, e.runZombie(t, out, "ssh", "ssh://zombie:zombie_ssh_pass@127.0.0.1:10022", tmpl), "local-ssh-smoke:ssh_marker", "zombie-ssh-ok")
}

func (e *dockerEnv) payload(t *testing.T) {
	e.composeUp(t, "redis")
	e.container(t, "redis", "REDISCLI_AUTH=zombie_redis_pass redis-cli DEL zombie:payload:cli-a zombie:payload:cli-b >/dev/null")

	result := e.runZombie(
		t,
		e.out(t, "payload"),
		"redis-payload-cli",
		"redis://:zombie_redis_pass@127.0.0.1:16379",
		filepath.Join(e.templates, "payload", "redis-payload-cli.yaml"),
		"--payload", "payload_key=zombie:payload:cli-a",
		"--payload", "payload_key=zombie:payload:cli-b",
	)
	assertExtractionCount(t, result, "local-redis-payload-cli:redis_payload_value", "payload-cli-ok", 2)
	e.container(t, "redis", `
REDISCLI_AUTH=zombie_redis_pass redis-cli GET zombie:payload:cli-a | grep -Fx payload-cli-ok &&
REDISCLI_AUTH=zombie_redis_pass redis-cli GET zombie:payload:cli-b | grep -Fx payload-cli-ok
`)
}

func (e *dockerEnv) postExploit(t *testing.T) {
	e.composeUp(t, "redis", "postgres", "mysql", "ssh")
	e.initRedis(t)
	e.container(t, "redis", "rm -f /tmp/zombie_redis_post_exploit.rdb")
	e.container(t, "mysql", "rm -f /tmp/zombie_mysql_post_exploit.txt")
	e.container(t, "postgres", `rm -f /tmp/zombie_pg_post_exploit.txt
psql -U zombie -d postgres -c "DROP TABLE IF EXISTS zombie_post_exploit;"`)
	e.container(t, "ssh", `
mkdir -p /home/zombie/app /home/zombie/.aws /home/zombie/.kube /home/zombie/.ssh
echo 'APP_SECRET=zombie-env-secret' > /home/zombie/app/.env
printf '[default]\naws_access_key_id = AKIAZOMBIETEST\naws_secret_access_key = zombie\n' > /home/zombie/.aws/credentials
printf 'apiVersion: v1\nclusters:\n- cluster:\n    server: https://kube.local\n' > /home/zombie/.kube/config
printf '%s\n' '-----BEGIN OPENSSH PRIVATE KEY-----' 'zombie-test-key' '-----END OPENSSH PRIVATE KEY-----' > /home/zombie/.ssh/id_rsa
rm -f /tmp/zombie_ssh_post_exploit.txt
chown -R zombie:zombie /home/zombie
`)

	out := e.out(t, "post-exploit")
	dir := filepath.Join(e.templates, "post-exploit")
	result := e.runZombie(t, out, "redis-post-exploit", "redis://:zombie_redis_pass@127.0.0.1:16379", filepath.Join(dir, "redis-post-exploit-local.yaml"))
	assertExtraction(t, result, "local-redis-post-exploit:redis_sensitive_value", "zombie-secret-pass")
	assertExtraction(t, result, "local-redis-post-exploit:redis_file_write_result", "OK")
	e.container(t, "redis", "test -f /tmp/zombie_redis_post_exploit.rdb && grep -aF zombie-redis-post-exploit /tmp/zombie_redis_post_exploit.rdb")

	result = e.runZombie(t, out, "mysql-post-exploit", "mysql://root:zombie_mysql_root@127.0.0.1:13306", filepath.Join(dir, "mysql-post-exploit-local.yaml"))
	assertExtraction(t, result, "local-mysql-post-exploit:mysql_app_credential", "demo:token")
	assertExtraction(t, result, "local-mysql-post-exploit:mysql_passwd_root", "root:")
	assertExtraction(t, result, "local-mysql-post-exploit:mysql_outfile_content", "zombie-mysql-post-exploit")
	e.container(t, "mysql", "test -f /tmp/zombie_mysql_post_exploit.txt && grep -F zombie-mysql-post-exploit /tmp/zombie_mysql_post_exploit.txt")

	result = e.runZombie(t, out, "postgres-post-exploit", "postgresql://zombie:zombie_pg_pass@127.0.0.1:15432", filepath.Join(dir, "postgres-post-exploit-local.yaml"))
	assertExtraction(t, result, "local-postgres-post-exploit:pg_passwd_root", "root:")
	assertExtraction(t, result, "local-postgres-post-exploit:pg_program_output", "zombie-pg-post-exploit")
	e.container(t, "postgres", "test -f /tmp/zombie_pg_post_exploit.txt && grep -F zombie-pg-post-exploit /tmp/zombie_pg_post_exploit.txt")

	result = e.runZombie(t, out, "ssh-post-exploit", "ssh://zombie:zombie_ssh_pass@127.0.0.1:10022", filepath.Join(dir, "ssh-post-exploit-local.yaml"))
	assertExtraction(t, result, "local-ssh-post-exploit:ssh_username", "zombie")
	assertExtraction(t, result, "local-ssh-post-exploit:ssh_file_write", "zombie-ssh-post-exploit")
	assertExtraction(t, result, "local-ssh-post-exploit:ssh_env_secret", "zombie-env-secret")
	assertExtraction(t, result, "local-ssh-post-exploit:ssh_aws_credentials_path", ".aws/credentials")
	e.container(t, "ssh", "test -f /tmp/zombie_ssh_post_exploit.txt && grep -F zombie-ssh-post-exploit /tmp/zombie_ssh_post_exploit.txt")
}

func (e *dockerEnv) existing(t *testing.T) {
	e.composeUp(t, "redis", "postgres", "mysql", "ssh")
	e.initExistingData(t)

	out := e.out(t, "existing")
	cases := []struct {
		name       string
		target     string
		template   string
		expectName string
		contains   string
	}{
		{"redis-info-gather", "redis://:zombie_redis_pass@127.0.0.1:16379", "redis/redis-info-gather.yaml", "redis-info-gather:redis_version", ""},
		{"redis-config-check", "redis://:zombie_redis_pass@127.0.0.1:16379", "redis/redis-config-check.yaml", "redis-config-check:redis_dir", "/var/www"},
		{"redis-sensitive-keys", "redis://:zombie_redis_pass@127.0.0.1:16379", "redis/redis-sensitive-keys.yaml", "redis-sensitive-keys:sensitive_keys", "app:password"},
		{"redis-mdut-rogue-prereq-check", "redis://:zombie_redis_pass@127.0.0.1:16379", "redis/redis-mdut-rogue-prereq-check.yaml", "redis-mdut-rogue-prereq-check:redis_replica_read_only", "yes"},
		{"postgresql-info-gather", "postgresql://zombie:zombie_pg_pass@127.0.0.1:15432", "postgresql/postgresql-info-gather.yaml", "postgresql-info-gather:pg_version", ""},
		{"postgresql-mdut-capability-check", "postgresql://zombie:zombie_pg_pass@127.0.0.1:15432", "postgresql/postgresql-mdut-capability-check.yaml", "postgresql-mdut-capability-check:pg_copy_program", "yes"},
		{"mysql-info-gather", "mysql://zombie:zombie_mysql_pass@127.0.0.1:13306", "mysql/mysql-info-gather.yaml", "mysql-info-gather:mysql_version", ""},
		{"mysql-credential-columns", "mysql://zombie:zombie_mysql_pass@127.0.0.1:13306", "mysql/mysql-credential-columns.yaml", "mysql-credential-columns:found_columns", "app_credentials"},
		{"mysql-user-enum", "mysql://root:zombie_mysql_root@127.0.0.1:13306", "mysql/mysql-user-enum.yaml", "mysql-user-enum:mysql_users", "root"},
		{"mysql-udf-check", "mysql://root:zombie_mysql_root@127.0.0.1:13306", "mysql/mysql-udf-check.yaml", "mysql-udf-check:mysql_plugin_dir", ""},
		{"mysql-mdut-capability-check", "mysql://root:zombie_mysql_root@127.0.0.1:13306", "mysql/mysql-mdut-capability-check.yaml", "mysql-mdut-capability-check:mysql_file_privilege", "yes"},
		{"ssh-info-gather", "ssh://zombie:zombie_ssh_pass@127.0.0.1:10022", "ssh/ssh-info-gather.yaml", "ssh-info-gather-linux:username", ""},
		{"ssh-env-files", "ssh://zombie:zombie_ssh_pass@127.0.0.1:10022", "ssh/ssh-env-files.yaml", "ssh-env-files:env_files", "/home/zombie/app/.env"},
		{"ssh-credential-files", "ssh://zombie:zombie_ssh_pass@127.0.0.1:10022", "ssh/ssh-credential-files.yaml", "ssh-credential-files-linux:aws_key", "AKIAZOMBIETEST"},
		{"ssh-docker-escape", "ssh://zombie:zombie_ssh_pass@127.0.0.1:10022", "ssh/ssh-docker-escape.yaml", "ssh-docker-escape:cap_flags", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := e.runZombie(t, out, tc.name, tc.target, filepath.Join(e.templatesExt, filepath.FromSlash(tc.template)))
			assertExtraction(t, result, tc.expectName, tc.contains)
		})
	}
}

func (e *dockerEnv) exploit(t *testing.T) {
	e.composeUp(t, "redis", "postgres", "mysql")
	e.container(t, "redis", `mkdir -p /var/www/html &&
chmod 777 /var/www/html &&
rm -f /var/www/html/shell.php`)
	e.container(t, "mysql", "rm -f /tmp/zombie_mysql_outfile.txt")
	e.container(t, "postgres", `rm -f /tmp/zombie_pg_cmd.txt
psql -U zombie -d postgres -c "DROP TABLE IF EXISTS zombie_cmd_exec;"`)

	out := e.out(t, "exploit")
	dir := filepath.Join(e.templates, "exploit")
	result := e.runZombie(t, out, "redis-write-webshell-local", "redis://:zombie_redis_pass@127.0.0.1:16379", filepath.Join(dir, "redis-write-webshell-local.yaml"))
	assertExtraction(t, result, "local-redis-write-webshell:redis_webshell_written", "OK")
	e.container(t, "redis", "test -f /var/www/html/shell.php && grep -aF '<?php system($_GET[0]);?>' /var/www/html/shell.php")

	result = e.runZombie(t, out, "mysql-outfile-local", "mysql://root:zombie_mysql_root@127.0.0.1:13306", filepath.Join(dir, "mysql-outfile-local.yaml"))
	assertExtraction(t, result, "local-mysql-outfile-write:mysql_outfile_content", "zombie-mysql-outfile-ok")
	e.container(t, "mysql", "test -f /tmp/zombie_mysql_outfile.txt && grep -F zombie-mysql-outfile-ok /tmp/zombie_mysql_outfile.txt")

	result = e.runZombie(t, out, "postgres-copy-program-local", "postgresql://zombie:zombie_pg_pass@127.0.0.1:15432", filepath.Join(dir, "postgres-copy-program-local.yaml"))
	assertExtraction(t, result, "local-postgres-copy-program:pg_program_output", "uid=")
	e.container(t, "postgres", "test -f /tmp/zombie_pg_cmd.txt && grep -F uid= /tmp/zombie_pg_cmd.txt")
}

func (e *dockerEnv) composeUp(t *testing.T, services ...string) {
	t.Helper()
	args := []string{"compose", "-f", e.compose, "up", "-d", "--wait", "--wait-timeout", "180", "--build"}
	args = append(args, services...)
	run(t, e.repoRoot, "docker", args...)
}

func (e *dockerEnv) container(t *testing.T, service, command string) {
	t.Helper()
	run(t, e.repoRoot, "docker", "compose", "-f", e.compose, "exec", "-T", service, "sh", "-lc", command)
}

func (e *dockerEnv) initRedis(t *testing.T) {
	t.Helper()
	run(t, e.repoRoot, "docker", "compose", "-f", e.compose, "run", "--rm", "redis-init")
}

func (e *dockerEnv) initExistingData(t *testing.T) {
	t.Helper()
	e.initRedis(t)
	run(t, e.repoRoot, "docker", "compose", "-f", e.compose, "exec", "-T", "mysql", "mysql", "-uroot", "-pzombie_mysql_root", "zombie", "-e", `
CREATE TABLE IF NOT EXISTS app_credentials (
  id INT PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  password_hash VARCHAR(128) NOT NULL,
  api_token VARCHAR(128) NOT NULL
);
INSERT IGNORE INTO app_credentials (id, username, password_hash, api_token)
VALUES (1, 'demo', 'hash', 'token');
`)
	e.container(t, "ssh", `
mkdir -p /opt /srv /var/www /home/zombie/app /home/zombie/.aws /home/zombie/.kube /home/zombie/.ssh
echo 'APP_SECRET=zombie-env-secret' > /home/zombie/app/.env
printf '[default]\naws_access_key_id = AKIAZOMBIETEST\naws_secret_access_key = zombie\n' > /home/zombie/.aws/credentials
printf 'apiVersion: v1\nclusters:\n- cluster:\n    server: https://kube.local\n' > /home/zombie/.kube/config
printf '%s\n' '-----BEGIN OPENSSH PRIVATE KEY-----' 'zombie-test-key' '-----END OPENSSH PRIVATE KEY-----' > /home/zombie/.ssh/id_rsa
chown -R zombie:zombie /home/zombie
`)
}

func (e *dockerEnv) runZombie(t *testing.T, outDir, name, target, template string, extraArgs ...string) zombieResult {
	t.Helper()
	outFile := filepath.Join(outDir, name+".json")
	args := []string{
		"run", ".",
		"-i", target,
		"--no-honeypot",
		"--no-unauth",
		"--service-template", template,
		"--timeout", "10",
	}
	args = append(args, extraArgs...)
	args = append(args, "-f", outFile, "-O", "json", "-o", "full")
	run(t, e.repoRoot, "go", args...)

	raw, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("%s output file: %v", name, err)
	}
	var result zombieResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("%s output json: %v\n%s", name, err, raw)
	}
	return result
}

func (e *dockerEnv) out(t *testing.T, name string) string {
	t.Helper()
	dir := filepath.Join(e.outDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("create output dir %s: %v", name, err)
	}
	return dir
}

func assertExtraction(t *testing.T, result zombieResult, name, contains string) {
	t.Helper()
	assertExtractionCount(t, result, name, contains, 1)
}

func assertExtractionCount(t *testing.T, result zombieResult, name, contains string, minCount int) {
	t.Helper()
	count := 0
	for _, item := range result.Extracteds {
		if item.Name != name {
			continue
		}
		if contains == "" && len(item.ExtractResult) > 0 {
			count += len(item.ExtractResult)
			continue
		}
		for _, value := range item.ExtractResult {
			if strings.Contains(value, contains) {
				count++
			}
		}
	}
	if count < minCount {
		b, _ := json.MarshalIndent(result, "", "  ")
		t.Fatalf("expected extraction %s containing %q at least %d time(s), got %d\n%s", name, contains, minCount, count, b)
	}
}

func waitTCP(t *testing.T, address string) {
	t.Helper()
	deadline := time.Now().Add(120 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, time.Second)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("%s did not open in time", address)
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	out, err := runCmd(dir, name, args...)
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return out
}

func runCmd(dir, name string, args ...string) (string, error) {
	return runCmdWithTimeout(dir, commandTimeout, name, args...)
}

func runCmdWithTimeout(dir string, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return string(out), fmt.Errorf("%s %s timed out after %s", name, strings.Join(args, " "), timeout)
	}
	if err != nil {
		return string(out), fmt.Errorf("%w", err)
	}
	if len(out) > 0 && strings.Contains(string(out), "error during connect") {
		return string(out), errors.New("docker connection failed")
	}
	return string(out), nil
}
