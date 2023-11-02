package plugin

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/mitchellh/go-vnc"
	"net"
	"time"
)

type VNCService struct {
	*pkg.Task
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

func (s *VNCService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}

func (s *VNCService) SetQuery(query string) {
	s.Input = query
}

func (s *VNCService) Output(res interface{}) {

}
