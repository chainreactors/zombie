package rdp

import (
	"github.com/chainreactors/zombie/external/grdp"
	"github.com/chainreactors/zombie/pkg"
)

// rdpSession implements pkg.Session. RDP has no persistent connection,
// so Close is a no-op and Raw returns nil.
type rdpSession struct {
	service string
}

func (s *rdpSession) Service() string  { return s.service }
func (s *rdpSession) Close() error     { return nil }
func (s *rdpSession) Raw() interface{} { return nil }

// RdpPlugin is stateless; all connection state lives in rdpSession.
type RdpPlugin struct{}

func (p *RdpPlugin) Name() string { return "rdp" }

func (p *RdpPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	user, domain := pkg.SplitUserDomain(task.Username)
	err := grdp.Login(task.Address(), domain, user, task.Password)
	if err != nil {
		return nil, err
	}
	return &rdpSession{service: task.Service}, nil
}

func (p *RdpPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}
