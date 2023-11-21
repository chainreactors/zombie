package telnet

import (
	"github.com/chainreactors/zombie/pkg"
	"strconv"
)

type TelnetPlugin struct {
	*pkg.Task
}

func (s *TelnetPlugin) Unauth() (bool, error) {
	return false, nil
}

func (s *TelnetPlugin) Login() error {

	port, _ := strconv.Atoi(s.Port)
	c := Init(s.IP, port)
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
	return pkg.NilConnError{Service: s.Service}
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
