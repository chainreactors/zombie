package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/jlaffaye/ftp"
)

type FtpService struct {
	*pkg.Task
	Input string
	conn  *ftp.ServerConn
}

func (s *FtpService) Query() bool {
	return false
}

func (s *FtpService) GetInfo() bool {
	return false
}

func (s *FtpService) Connect() error {
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

func (s *FtpService) SetQuery(query string) {
	s.Input = query
}

func (s *FtpService) Output(res interface{}) {

}

func (s *FtpService) Close() error {
	if s.conn != nil {
		return s.conn.Quit()
	}
	return NilConnError{s.Service}
}
