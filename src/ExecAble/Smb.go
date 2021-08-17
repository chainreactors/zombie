package ExecAble

import (
	"Zombie/src/Utils"
	"encoding/hex"
	"fmt"

	"github.com/hirochachacha/go-smb2"
	"net"
	"strings"
	"time"
)

type SmbService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string
	Input    string
}

func (s *SmbService) Query() bool {
	return false
}

func (s *SmbService) GetInfo() bool {
	return false
}

func (s *SmbService) Connect() bool {
	err, verison, res := SMBConnect(s.Username, s.Password, s.IpInfo)
	s.Version = verison
	if err == nil && res {
		return true
	}

	return false
}

func (s *SmbService) DisConnect() bool {
	return false
}

func (s *SmbService) SetQuery(query string) {
	s.Input = query
}

func (s *SmbService) Output(res interface{}) {

}

func SMBConnect(User string, Password string, info Utils.IpInfo) (err error, version string, result bool) {

	var UserName, DoaminName string

	if strings.Contains(User, "/") {
		UserName = strings.Split(User, "/")[1]
		DoaminName = strings.Split(User, "/")[0]
	} else {
		UserName = User
		DoaminName = ""
	}

	target := fmt.Sprintf("%v:%v", info.Ip, info.Port)

	conn, err := net.DialTimeout("tcp", target, time.Duration(Utils.Timeout)*time.Second)
	if err == nil {
		defer conn.Close()

		//hash := "11e7993210372b9634119676e7401289"
		//buf := make([]byte, len(hash)/2)
		//hex.Decode(buf, []byte(hash))

		d := &smb2.Dialer{}

		if strings.HasPrefix(Password, "hash:") {
			hash := Password[5:]
			buf := make([]byte, len(hash)/2)
			hex.Decode(buf, []byte(hash))
			d.Initiator = &smb2.NTLMInitiator{
				User:   UserName,
				Domain: DoaminName,
				Hash:   buf,
			}
		} else {
			d.Initiator = &smb2.NTLMInitiator{
				User:   UserName,
				Domain: DoaminName,
				//Hash: buf,
				Password: Password,
			}
		}

		_ = conn.SetDeadline(time.Now().Add(time.Duration(Utils.Timeout) * time.Second))

		s, curversion, err2 := d.Dial(conn)

		if err2 == nil {
			defer s.Logoff()
			result = true
			return err2, curversion, result
		}
		return err2, curversion, result
	}
	return err, "", result

}
