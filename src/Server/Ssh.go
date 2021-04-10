package Server

import (
	"Zombie/src/Utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func SSHConnect(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	config := &ssh.ClientConfig{
		User: User,
		Auth: []ssh.AuthMethod{
			ssh.Password(Password),
		},
		Timeout: time.Duration(Utils.Timeout+10) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", info.Ip, info.Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("whoami")
		if err == nil && errRet == nil {
			defer session.Close()
			result.Result = true
		}
	}
	return err, result
}
