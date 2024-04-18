package ssh

import (
	"github.com/chainreactors/zombie/pkg"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type SshPlugin struct {
	*pkg.Task
	//Cmd            string
	conn *ssh.Client
}

func (s *SshPlugin) Login() error {
	var auth []ssh.AuthMethod
	if method, pkfile := pkg.ParseMethod(s.Password); method == "pk" && pkfile != "" {
		auth = []ssh.AuthMethod{
			publicKeyAuthFunc(pkfile),
		}
	} else {
		auth = []ssh.AuthMethod{
			ssh.Password(s.Password),
		}
	}

	conn, err := SSHConnect(s.Task, auth)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *SshPlugin) Unauth() (bool, error) {
	conn, err := SSHConnect(s.Task, []ssh.AuthMethod{ssh.Password("")})
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

func (s *SshPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}

func (s *SshPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *SshPlugin) Name() string {
	return s.Service
}

func SSHConnect(task *pkg.Task, auth []ssh.AuthMethod) (conn *ssh.Client, err error) {
	config := &ssh.ClientConfig{
		User:    task.Username,
		Timeout: time.Duration(task.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	config.Auth = auth

	conn, err = ssh.Dial("tcp", task.Address(), config)
	if err != nil {
		return nil, err
	}

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
