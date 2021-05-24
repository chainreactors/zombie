package Protocol

import (
	"Zombie/src/Utils"
	"encoding/hex"
	"fmt"

	"github.com/hirochachacha/go-smb2"
	"net"
	"strings"
	"time"
)

func SMBConnect(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {

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

		//_ = conn.SetDeadline(time.Now().Add(time.Duration(Utils.Timeout) * time.Second))

		s, version, err := d.Dial(conn)
		result.Additional += version
		if err == nil {
			defer s.Logoff()
			result.Result = true
		}
		return err, result
	}

	return err, result
}
