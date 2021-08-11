package ExecAble

import (
	"Zombie/src/Utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type SshService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	MysqlInf
	Cmd    string
	SshCon *ssh.Client
}

func (s *SshService) Connect() bool {
	err, _, conn := SSHConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		s.SshCon = conn
		return true
	}
	return false
}

func (s *SshService) GetInfo() bool {
	return true
}

func (s *SshService) SetQuery(cmd string) {
	s.Cmd = cmd
}

func (s *SshService) Query() bool {

	session, err := s.SshCon.NewSession()
	defer session.Close()
	defer s.SshCon.Close()
	buf, err := session.Output("ping -c 3 114.114.114.114")

	if err != nil {
		return false
	}
	fmt.Println(string(buf))
	return true
}

func SSHConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, connect *ssh.Client) {
	config := &ssh.ClientConfig{
		User: User,
		Auth: []ssh.AuthMethod{
			ssh.Password(Password),
		},
		Timeout: time.Duration(Utils.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", info.Ip, info.Port), config)
	if err == nil {
		session, err := client.NewSession()
		defer session.Close()
		errRet := session.Run("whoami")
		if err == nil && errRet == nil {
			result = true
		}
		connect = client
	}
	return err, result, connect
}

func SSHConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, conn := SSHConnect(User, Password, info)

	if err == nil {
		result.Result = res
		conn.Close()
	}

	return err, result
}
