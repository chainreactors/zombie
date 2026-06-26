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

type Plugin interface {
	Name() string
	Open(task *Task) (Session, error)
	Unauth(task *Task) (Session, error)
}

var pluginRegistry = map[string]Plugin{}

func RegisterPlugin(name string, p Plugin) {
	pluginRegistry[name] = p
}

func GetPlugin(service string) (Plugin, bool) {
	p, ok := pluginRegistry[service]
	return p, ok
}

func DefaultPluginRegistry() map[string]Plugin {
	m := make(map[string]Plugin, len(pluginRegistry))
	for k, v := range pluginRegistry {
		m[k] = v
	}
	return m
}
