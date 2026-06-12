package pkg

type Session interface {
	Service() string
	Close() error
	Raw() interface{}
}

type ShellSession interface {
	Session
	Exec(cmd string) ([]byte, error)
}

type SQLSession interface {
	Session
	Query(query string, args ...any) ([][]string, error)
	Databases() ([]string, error)
}

type KVSession interface {
	Session
	Get(key string) ([]byte, error)
	Keys(pattern string) ([]string, error)
}

type FileSession interface {
	Session
	List(path string) ([]string, error)
	Read(path string) ([]byte, error)
}

type DirectorySession interface {
	Session
	Search(baseDN, filter string, attrs []string) ([]map[string][]string, error)
}

func AsShell(s Session) (ShellSession, bool) {
	ss, ok := s.(ShellSession)
	return ss, ok
}

func AsSQL(s Session) (SQLSession, bool) {
	ss, ok := s.(SQLSession)
	return ss, ok
}

func AsKV(s Session) (KVSession, bool) {
	ss, ok := s.(KVSession)
	return ss, ok
}

func AsFile(s Session) (FileSession, bool) {
	ss, ok := s.(FileSession)
	return ss, ok
}

func AsDirectory(s Session) (DirectorySession, bool) {
	ss, ok := s.(DirectorySession)
	return ss, ok
}
