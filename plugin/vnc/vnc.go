package vnc

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/mitchellh/go-vnc"
	"net"
	"time"
)

type VNCPlugin struct {
	*pkg.Task
	conn  *vnc.ClientConn
	Input string
}

func (s *VNCPlugin) Unauth() (bool, error) {
	target := s.Address()

	tcpconn, err := net.DialTimeout("tcp", target, time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return false, err
	}

	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: ""},
		},
	}
	conn, err := vnc.Client(tcpconn, &config)
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

func (s *VNCPlugin) Login() error {
	target := s.Address()

	tcpconn, err := net.DialTimeout("tcp", target, time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return err
	}

	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: s.Password},
		},
	}
	conn, err := vnc.Client(tcpconn, &config)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *VNCPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

//func (s *VNCPlugin) SetQuery(query string) {
//	s.Input = query
//}
//
//func (s *VNCPlugin) Output(res interface{}) {
//
//}

func (s *VNCPlugin) Name() string {
	return s.Service
}

func (s *VNCPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}
