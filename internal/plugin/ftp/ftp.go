package ftp

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/jlaffaye/ftp"
)

type FtpPlugin struct {
	*pkg.Task
	Input string
	conn  *ftp.ServerConn
}

func (s *FtpPlugin) Name() string {
	return s.Service.String()
}

func (s *FtpPlugin) Unauth() (bool, error) {
	// todo anoy login
	return false, nil
}

func (s *FtpPlugin) Login() error {
	conn, err := ftp.DialTimeout(fmt.Sprintf(s.Address()), s.Duration())
	if err != nil {
		return err
	}
	err = conn.Login(s.Username, s.Password)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *FtpPlugin) GetBasic() *pkg.Basic {
	// todo list root dir
	return &pkg.Basic{}
}

func (s *FtpPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Quit()
	}
	return pkg.NilConnError{Service: s.Service}
}
