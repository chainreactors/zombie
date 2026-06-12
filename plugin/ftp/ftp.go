package ftp

import (
	"io"

	"github.com/chainreactors/zombie/pkg"
	"github.com/jlaffaye/ftp"
)

// ftpSession implements pkg.FileSession over an authenticated FTP connection.
type ftpSession struct {
	service string
	conn    *ftp.ServerConn
}

func (s *ftpSession) Service() string  { return s.service }
func (s *ftpSession) Raw() interface{} { return s.conn }

func (s *ftpSession) Close() error {
	if s.conn != nil {
		return s.conn.Quit()
	}
	return nil
}

func (s *ftpSession) List(path string) ([]string, error) {
	entries, err := s.conn.List(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names, nil
}

func (s *ftpSession) Read(path string) ([]byte, error) {
	resp, err := s.conn.Retr(path)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	return io.ReadAll(resp)
}

// FtpPlugin is stateless; all connection state lives in ftpSession.
type FtpPlugin struct{}

func (p *FtpPlugin) Name() string { return "ftp" }

// dial establishes an FTP control connection using the task's proxy-aware dialer.
func (p *FtpPlugin) dial(task *pkg.Task) (*ftp.ServerConn, error) {
	netConn, err := task.DialTimeout("tcp", task.Address(), task.Duration())
	if err != nil {
		return nil, err
	}
	conn, err := ftp.Dial(task.Address(), ftp.DialWithNetConn(netConn))
	if err != nil {
		netConn.Close()
		return nil, err
	}
	return conn, nil
}

func (p *FtpPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	conn, err := p.dial(task)
	if err != nil {
		return nil, err
	}
	if err := conn.Login(task.Username, task.Password); err != nil {
		conn.Quit()
		return nil, err
	}
	return &ftpSession{service: task.Service, conn: conn}, nil
}

func (p *FtpPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	conn, err := p.dial(task)
	if err != nil {
		return nil, err
	}
	if err := conn.Login("anonymous", ""); err != nil {
		conn.Quit()
		return nil, err
	}
	return &ftpSession{service: task.Service, conn: conn}, nil
}
