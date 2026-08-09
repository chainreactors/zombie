package pkg

type Session interface {
	Service() string
	Close() error
}

type ShellSession interface {
	Session
	Exec(cmd string) ([]byte, error)
}

type SQLSession interface {
	Session
	Query(query string, args ...any) ([][]string, error)
}

type KVSession interface {
	Session
	Get(key string) ([]byte, error)
	Keys(pattern string) ([]string, error)
	Command(name string, args ...string) (interface{}, error)
}

type FileSession interface {
	Session
	List(path string) ([]string, error)
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
}

type DirectorySession interface {
	Session
	Search(baseDN, filter string, attrs []string) ([]map[string][]string, error)
}

// AuditableSession can discover and sample sensitive data automatically.
// Each database plugin implements its own discovery and sampling logic.
type AuditableSession interface {
	Session
	// Audit discovers locations matching field-name patterns and samples data.
	// Returns map[location]sampledData where location identifies the source
	// (e.g. "schema.table.column" for SQL, "key:name" for Redis).
	Audit(patterns []string, limit int) (map[string]string, error)
}

var DefaultAuditPatterns []string
