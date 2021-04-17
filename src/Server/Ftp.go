package Server

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)

func FtpConnect(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", info.Ip, info.Port), time.Duration(Utils.Timeout)*time.Second)
	if err == nil {
		err = conn.Login(User, Password)
		if err == nil {
			defer conn.Logout()
			result.Result = true
		}
	}
	return err, result
}
