package vnc

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/mitchellh/go-vnc"
	"time"
)

// vncSession implements pkg.Session over an authenticated VNC connection.
type vncSession struct {
	service string
	conn    *vnc.ClientConn
}

func (s *vncSession) Service() string  { return s.service }
func (s *vncSession) Raw() interface{} { return s.conn }

func (s *vncSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// VNCPlugin is stateless; all connection state lives in vncSession.
type VNCPlugin struct{}

func (p *VNCPlugin) Name() string { return "vnc" }

func (p *VNCPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	return dial(task, task.Password)
}

func (p *VNCPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return dial(task, "")
}

func dial(task *pkg.Task, password string) (pkg.Session, error) {
	tcpconn, err := task.DialTimeout("tcp", task.Address(), time.Duration(task.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: password},
		},
	}
	conn, err := vnc.Client(tcpconn, &config)
	if err != nil {
		return nil, err
	}
	return &vncSession{service: task.Service, conn: conn}, nil
}
