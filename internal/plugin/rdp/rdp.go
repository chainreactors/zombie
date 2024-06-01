package rdp

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/lcvvvv/kscan/grdp"
)

type RdpPlugin struct {
	*pkg.Task
	conn *grdp.Client
}

func (s *RdpPlugin) Unauth() (bool, error) {
	return false, pkg.NotImplUnauthorized
}

func (s *RdpPlugin) Login() error {
	user, domain := pkg.SplitUserDomain(s.Username)
	err := grdp.Login(s.Address(), user, domain, s.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *RdpPlugin) Close() error {
	return nil
}

func (s *RdpPlugin) Name() string {
	return s.Service
}

func (s *RdpPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}
