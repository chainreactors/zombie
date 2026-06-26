package zookeeper

import (
	"fmt"
	"strings"
	"time"

	"github.com/chainreactors/zombie/pkg"
	"github.com/samuel/go-zookeeper/zk"
)

func init() {
	pkg.RegisterPlugin("zookeeper", &ZookeeperPlugin{})
	pkg.Services.Register(&pkg.Service{Name: "zookeeper", DefaultPort: "2181", Source: pkg.PluginSource})
}

// zkSession implements pkg.Session over a ZooKeeper connection.
type zkSession struct {
	service string
	conn    *zk.Conn
}

func (s *zkSession) Service() string  { return s.service }
func (s *zkSession) Raw() interface{} { return s.conn }

func (s *zkSession) Close() error {
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}

func (s *zkSession) Get(key string) ([]byte, error) {
	data, _, err := s.conn.Get(key)
	return data, err
}

func (s *zkSession) Keys(pattern string) ([]string, error) {
	path := pattern
	if path == "*" || path == "" {
		path = "/"
	}
	children, _, err := s.conn.Children(path)
	return children, err
}

func (s *zkSession) Command(name string, args ...string) (interface{}, error) {
	switch strings.ToUpper(name) {
	case "SET":
		if len(args) < 2 {
			return nil, fmt.Errorf("SET requires path and data")
		}
		_, err := s.conn.Set(args[0], []byte(args[1]), -1)
		return "OK", err
	case "CREATE":
		if len(args) < 2 {
			return nil, fmt.Errorf("CREATE requires path and data")
		}
		path, err := s.conn.Create(args[0], []byte(args[1]), 0, zk.WorldACL(zk.PermAll))
		return path, err
	case "DELETE":
		if len(args) < 1 {
			return nil, fmt.Errorf("DELETE requires path")
		}
		return "OK", s.conn.Delete(args[0], -1)
	default:
		return nil, fmt.Errorf("unsupported zookeeper command: %s", name)
	}
}

// ZookeeperPlugin is stateless; all connection state lives in zkSession.
type ZookeeperPlugin struct{}

func (p *ZookeeperPlugin) Name() string { return "zookeeper" }

func (p *ZookeeperPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", task.IP, task.Port)}, time.Duration(task.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	err = conn.AddAuth("digest", []byte(fmt.Sprintf("%s:%s", task.Username, task.Password)))
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &zkSession{service: task.Service, conn: conn}, nil
}

func (p *ZookeeperPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", task.IP, task.Port)}, time.Duration(task.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	return &zkSession{service: task.Service, conn: conn}, nil
}
