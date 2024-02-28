package ftp

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/jlaffaye/ftp"
)

type FtpPlugin struct {
	*pkg.Task
	Input string
	conn  *ftp.ServerConn
}

func (s *FtpPlugin) Name() string {
	return s.Service
}

func (s *FtpPlugin) Unauth() (bool, error) {
	conn, err := ftp.DialTimeout(s.Address(), s.Duration())
	if err != nil {
		return false, err
	}
	err = conn.Login("anonymous", "")
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

func (s *FtpPlugin) Login() error {
	conn, err := ftp.DialTimeout(s.Address(), s.Duration())
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

func (s *FtpPlugin) GetResult() *pkg.Result {
	// todo list root dir
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *FtpPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Quit()
	}
	return pkg.NilConnError{Service: s.Service}
}
