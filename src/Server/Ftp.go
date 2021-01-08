package Server

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/jlaffaye/ftp"
)

func FtpConnect(User string, Password string, info Utils.IpInfo) (err error, result bool) {
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", info.Ip, info.Port), Utils.Timeout)
	if err == nil {
		err = conn.Login(User, Password)
		if err == nil {
			defer conn.Logout()
			result = true
		}
	}
	return err, result
}
