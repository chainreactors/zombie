package ssh

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/chainreactors/zombie/pkg"
	"golang.org/x/crypto/ssh"
)

type SshPlugin struct {
	*pkg.Task
	conn *ssh.Client
}

func (s *SshPlugin) Login() error {
	var auth []ssh.AuthMethod
	if method, pkdata := pkg.ParseMethod(s.Password); method == "pk" && pkdata != "" {
		am, err := publicKeyAuth(pkdata)
		if err != nil {
			return err
		}
		auth = []ssh.AuthMethod{am}
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
	return nil
}

func (s *SshPlugin) GetResult() *pkg.Result {
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
		Auth: auth,
	}

	netConn, err := task.DialTimeout("tcp", task.Address(), config.Timeout)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(netConn, task.Address(), config)
	if err != nil {
		netConn.Close()
		return nil, err
	}
	conn = ssh.NewClient(c, chans, reqs)

	return conn, nil
}

// publicKeyAuth resolves a private key from either base64-encoded PEM
// data or a file path, and returns an SSH auth method.
func publicKeyAuth(data string) (ssh.AuthMethod, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		keyBytes, err = os.ReadFile(data)
		if err != nil {
			return nil, fmt.Errorf("ssh key read failed: %w", err)
		}
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("ssh key parse failed: %w", err)
	}
	return ssh.PublicKeys(signer), nil
}
