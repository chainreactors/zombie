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
	return false, nil
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
	return pkg.NilConnError{s.Service}
}

//func (s *VNCPlugin) SetQuery(query string) {
//	s.Input = query
//}
//
//func (s *VNCPlugin) Output(res interface{}) {
//
//}

func (s *VNCPlugin) Name() string {
	return s.Service.String()
}

func (s *VNCPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}
