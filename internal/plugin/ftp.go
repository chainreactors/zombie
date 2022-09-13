package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/jlaffaye/ftp"
	"time"
)

type FtpService struct {
	*utils.Task
	Input  string
	Ftpcon *ftp.ServerConn
}

func FtpConnect(info *utils.Task) (conn *ftp.ServerConn, err error) {
	conn, err = ftp.DialTimeout(fmt.Sprintf(info.Address()), time.Duration(utils.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	err = conn.Login(info.Username, info.Password)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func (s *FtpService) Query() bool {
	return false
}

func (s *FtpService) GetInfo() bool {
	return false
}

func (s *FtpService) Connect() error {
	conn, err := FtpConnect(s.Task)
	if err != nil {
		return err
	}
	s.Ftpcon = conn
	return nil
}

func (s *FtpService) SetQuery(query string) {
	s.Input = query
}

func (s *FtpService) Output(res interface{}) {

}

func (s *FtpService) Close() error {
	if s.Ftpcon != nil {
		return s.Ftpcon.Quit()
	}
	return NilConnError{s.Service}
}
