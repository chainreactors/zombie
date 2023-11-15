package telnet

import (
	"github.com/chainreactors/zombie/internal/plugin/telnet/lib"
	"github.com/chainreactors/zombie/pkg"
	"strconv"
)

type TelnetPlugin struct {
	*pkg.Task
}

func (s *TelnetPlugin) Unauth() (bool, error) {
	//TODO implement me
	panic("implement me")
}

//func (s *TelnetService) Query() bool {
//	return false
//}
//
//func (s *TelnetService) GetInfo() bool {
//	return false
//}

func (s *TelnetPlugin) Login() error {

	port, _ := strconv.Atoi(s.Port)
	c := lib.Init(s.IP, port)
	err := c.Connect()
	if err != nil {
		return err
	}

	c.UserName = s.Username
	c.Password = s.Password
	c.ServerType = 3

	err = c.Login()
	if err != nil {
		return err
	}

	return nil

}

func (s *TelnetPlugin) Close() error {
	return pkg.NilConnError{s.Service}
}

func (s *TelnetPlugin) Name() string {
	return s.Service.String()
}

func (s *TelnetPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}

//func (s *TelnetService) SetQuery(query string) {
//	//s.Input = query
//}
//
//func (s *TelnetService) Output(res interface{}) {
//
//}
