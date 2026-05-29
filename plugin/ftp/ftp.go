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

// dial 通过 task 配置（含代理）建立 FTP 控制连接。
func (s *FtpPlugin) dial() (*ftp.ServerConn, error) {
	netConn, err := s.DialTimeout("tcp", s.Address(), s.Duration())
	if err != nil {
		return nil, err
	}
	conn, err := ftp.Dial(s.Address(), ftp.DialWithNetConn(netConn))
	if err != nil {
		netConn.Close()
		return nil, err
	}
	return conn, nil
}

func (s *FtpPlugin) Unauth() (bool, error) {
	conn, err := s.dial()
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
	conn, err := s.dial()
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
	return nil
}
