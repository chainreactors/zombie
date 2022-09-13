package plugin

import (
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/mitchellh/go-vnc"
	"net"
	"time"
)

type VNCService struct {
	*utils.Task
	conn  *vnc.ClientConn
	Input string
}

func (s *VNCService) Query() bool {
	return false
}

func (s *VNCService) GetInfo() bool {
	return false
}

func (s *VNCService) Connect() error {
	conn, err := VNCConnect(s.Task)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil

}

func (s *VNCService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return NilConnError{s.Service}
}

func (s *VNCService) SetQuery(query string) {
	s.Input = query
}

func (s *VNCService) Output(res interface{}) {

}

func VNCConnect(task *utils.Task) (conn *vnc.ClientConn, err error) {
	target := task.Address()

	tcpconn, err := net.DialTimeout("tcp", target, time.Duration(task.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: task.Password},
		},
	}
	conn, err = vnc.Client(tcpconn, &config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
