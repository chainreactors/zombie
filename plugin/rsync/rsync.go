package rsync

import (
	"github.com/chainreactors/zombie/pkg"
)

// rsyncSession implements pkg.Session. Rsync uses short-lived socket
// connections per operation, so there is no persistent conn to wrap.
type rsyncSession struct {
	service string
}

func (s *rsyncSession) Service() string  { return s.service }
func (s *rsyncSession) Close() error     { return nil }
func (s *rsyncSession) Raw() interface{} { return nil }

// RsyncPlugin is stateless; all connection state lives in rsyncSession.
type RsyncPlugin struct{}

func (p *RsyncPlugin) Name() string { return "rsync" }

func (p *RsyncPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	ver, modules, err := RsyncDetect(task.Address(), task.Timeout, task.DialTimeout)
	if err != nil {
		return nil, err
	}

	err = RsyncLogin(task.Address(), task.Username, task.Password, ver, modules, task.Timeout, task.DialTimeout)
	if err != nil {
		return nil, err
	}

	return &rsyncSession{service: task.Service}, nil
}

func (p *RsyncPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	ver, modules, err := RsyncDetect(task.Address(), task.Timeout, task.DialTimeout)
	if err != nil {
		return nil, err
	}
	err = RsyncUnauth(task.Address(), ver, modules, task.Timeout, task.DialTimeout)
	if err != nil {
		return nil, err
	}
	return &rsyncSession{service: task.Service}, nil
}
