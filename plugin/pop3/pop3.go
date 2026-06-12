package pop3

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/knadh/go-pop3"
	"strconv"
)

// pop3Session implements pkg.Session over an authenticated POP3 connection.
type pop3Session struct {
	service string
	conn    *pop3.Conn
}

func (s *pop3Session) Service() string  { return s.service }
func (s *pop3Session) Raw() interface{} { return s.conn }

func (s *pop3Session) Close() error {
	if s.conn != nil {
		return s.conn.Quit()
	}
	return nil
}

// Pop3Plugin is stateless; all connection state lives in pop3Session.
type Pop3Plugin struct{}

func (p *Pop3Plugin) Name() string { return "pop3" }

func (p *Pop3Plugin) Open(task *pkg.Task) (pkg.Session, error) {
	port, _ := strconv.Atoi(task.Port)

	pp := pop3.New(pop3.Opt{
		Host:       task.IP,
		Port:       port,
		TLSEnabled: false,
	})

	c, err := pp.NewConn()
	if err != nil {
		return nil, err
	}

	if err := c.Auth(task.Username, task.Password); err != nil {
		c.Quit()
		return nil, err
	}

	return &pop3Session{service: task.Service, conn: c}, nil
}

func (p *Pop3Plugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}
