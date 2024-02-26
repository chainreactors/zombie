package telnet

import (
	"github.com/chainreactors/zombie/pkg"
)

type TelnetPlugin struct {
	*pkg.Task
}

func (s *TelnetPlugin) Unauth() (bool, error) {
	c, err := NewClient(s.Address(), "", "", s.Duration())
	if err != nil {
		return false, err
	}
	err = c.Login()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TelnetPlugin) Login() error {
	c, err := NewClient(s.Address(), s.Username, s.Password, s.Duration())
	if err != nil {
		return err
	}
	err = c.Login()
	if err != nil {
		return err
	}

	return nil

}

func (s *TelnetPlugin) Close() error {
	return pkg.NilConnError{Service: s.Service}
}

func (s *TelnetPlugin) Name() string {
	return s.Service.String()
}

func (s *TelnetPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}
