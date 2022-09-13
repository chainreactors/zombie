package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg/utils"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type SshService struct {
	*utils.Task
	MysqlInf
	Cmd  string
	conn *ssh.Client
}

func (s *SshService) Connect() error {
	conn, err := SSHConnect(s.Task)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *SshService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return NilConnError{s.Service}
}

func (s *SshService) GetInfo() bool {

	if s.Cmd != "" {
		session, err := s.conn.NewSession()
		cmd := "ping -c 5 " + s.Cmd
		buf, err := session.Output(cmd)

		if err != nil {
			return false
		}

		re, _ := regexp.Compile(`\d received`)

		FindRes := string(re.Find([]byte(buf)))

		reslist := strings.Split(FindRes, " ")
		if reslist[1] == "received" {
			if reslist[0] != "0" {
				fmt.Printf("%v can reach %v\n", s.IP.String(), s.Cmd)
			}
		}
	} else {
		panic("Please input ip")
	}

	return true
}

func (s *SshService) SetQuery(cmd string) {
	s.Cmd = cmd
}

func (s *SshService) Query() bool {
	session, err := s.conn.NewSession()
	defer session.Close()
	defer s.conn.Close()
	buf, err := session.Output(s.Cmd)

	if err != nil {
		return false
	}
	res := fmt.Sprintf(s.IP.String() + ":\n" + string(buf) + "\n")
	s.Output(res)
	return true
}

func (s *SshService) Output(res interface{}) {
	finres := res.(string)
	utils.TDatach <- finres
}

func SSHConnect(task *utils.Task) (conn *ssh.Client, err error) {
	config := &ssh.ClientConfig{
		User:    task.Username,
		Timeout: time.Duration(task.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if strings.HasPrefix(task.Password, "pk:") {
		config.Auth = []ssh.AuthMethod{
			publicKeyAuthFunc(task.Password[3:]),
		}
	} else {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(task.Password),
		}
	}

	conn, err = ssh.Dial("tcp", task.Address(), config)
	if err != nil {
		return nil, err
	}
	//session, err := client.NewSession()
	//if err != nil {
	//
	//}
	//defer session.Close()
	//errRet := session.Run("whoami")
	//if err == nil && errRet == nil {
	//	result = true
	//}
	//connect = client
	return conn, nil
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {

	key, err := ioutil.ReadFile(kPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
