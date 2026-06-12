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

// sshSession implements pkg.ShellSession over an authenticated SSH connection.
type sshSession struct {
	service string
	conn    *ssh.Client
}

func (s *sshSession) Service() string  { return s.service }
func (s *sshSession) Raw() interface{} { return s.conn }

func (s *sshSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *sshSession) Exec(cmd string) ([]byte, error) {
	sess, err := s.conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("ssh session: %w", err)
	}
	defer sess.Close()
	return sess.CombinedOutput(cmd)
}

// SshPlugin is stateless; all connection state lives in sshSession.
type SshPlugin struct{}

func (p *SshPlugin) Name() string { return "ssh" }

func (p *SshPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	var auth []ssh.AuthMethod
	if method, pkdata := pkg.ParseMethod(task.Password); method == "pk" && pkdata != "" {
		am, err := publicKeyAuth(pkdata)
		if err != nil {
			return nil, err
		}
		auth = []ssh.AuthMethod{am}
	} else {
		auth = []ssh.AuthMethod{
			ssh.Password(task.Password),
		}
	}

	conn, err := SSHConnect(task, auth)
	if err != nil {
		return nil, err
	}
	return &sshSession{service: task.Service, conn: conn}, nil
}

func (p *SshPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	conn, err := SSHConnect(task, []ssh.AuthMethod{ssh.Password("")})
	if err != nil {
		return nil, err
	}
	return &sshSession{service: task.Service, conn: conn}, nil
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
