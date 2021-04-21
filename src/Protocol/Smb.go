package Protocol

import (
	"Zombie/pkg/github.com/stacktitan/smb/smb"
	"Zombie/src/Utils"
	"strings"
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

	options := smb.Options{
		Host:        info.Ip,
		Port:        info.Port,
		User:        UserName,
		Password:    Password,
		Domain:      DoaminName,
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			result.Result = true
		}
	}
	return err, result
}
