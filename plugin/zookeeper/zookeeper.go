package zookeeper

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

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
