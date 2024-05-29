package zookeeper

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type ZookeeperPlugin struct {
	*pkg.Task
	conn *zk.Conn
}

func (s *ZookeeperPlugin) Name() string {
	return s.Service
}

func (s *ZookeeperPlugin) Unauth() (bool, error) {
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", s.IP, s.Port)}, time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

func (s *ZookeeperPlugin) Login() error {
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", s.IP, s.Port)}, time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return err
	}
	err = conn.AddAuth("digest", []byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *ZookeeperPlugin) GetResult() *pkg.Result {
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *ZookeeperPlugin) Close() error {
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}
