package plugin

import (
	"github.com/chainreactors/zombie/internal/plugin/lib"
	"github.com/chainreactors/zombie/pkg"
)

type TelnetService struct {
	*pkg.Task
}

func (s *TelnetService) Query() bool {
	return false
}

func (s *TelnetService) GetInfo() bool {
	return false
}

func (s *TelnetService) Connect() error {

	c := lib.New(s.IP, 23)
	err := c.Connect()
	if err != nil {
		return err
	}
	c.UserName = s.Username
	c.Password = s.Password
	c.ServerType = c.MakeServerType()
	err = c.Login()
	if err != nil {
		return err
	}

	return nil

}

func (s *TelnetService) Close() error {
	return NilConnError{s.Service}
}

func (s *TelnetService) SetQuery(query string) {
	//s.Input = query
}

func (s *TelnetService) Output(res interface{}) {

}
