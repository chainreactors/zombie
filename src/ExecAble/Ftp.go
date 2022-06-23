package ExecAble

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)

type FtpService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Input    string
	Ftpcon   *ftp.ServerConn
}

func FtpConnect(User string, Password string, info Utils.IpInfo) (err error, conn *ftp.ServerConn, result bool) {
	conn, err = ftp.DialTimeout(fmt.Sprintf("%v:%v", info.Ip, info.Port), time.Duration(Utils.Timeout)*time.Second)
	if err == nil {
		err = conn.Login(User, Password)
		if err == nil {
			result = true
		}
	}
	return err, conn, result
}

func (s *FtpService) Query() bool {
	return false
}

func (s *FtpService) GetInfo() bool {
	return false
}

func (s *FtpService) Connect() bool {
	err, conn, res := FtpConnect(s.Username, s.Password, s.IpInfo)
	if conn != nil {
		s.Ftpcon = conn
	}
	if err == nil && res {
		return true
	}
	return false
}

func (s *FtpService) SetQuery(query string) {
	s.Input = query
}

func (s *FtpService) Output(res interface{}) {

}

func (s *FtpService) DisConnect() bool {
	if s.Ftpcon != nil {
		s.Ftpcon.Quit()
	}

	return false
}
