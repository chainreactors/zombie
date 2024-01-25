package telnet

import (
	"github.com/chainreactors/zombie/pkg"
	"strconv"
)

type TelnetPlugin struct {
	*pkg.Task
}

func (s *TelnetPlugin) Unauth() (bool, error) {
	port, _ := strconv.Atoi(s.Port)
	c := &Client{
		IPAddr:       s.IP,
		Port:         port,
		UserName:     "",
		Password:     "",
		conn:         nil,
		LastResponse: "",
		ServerType:   0,
	}
	err := c.Login()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TelnetPlugin) Login() error {
	port, _ := strconv.Atoi(s.Port)
	c := &Client{
		IPAddr:       s.IP,
		Port:         port,
		UserName:     s.Username,
		Password:     s.Password,
		conn:         nil,
		LastResponse: "",
		ServerType:   3,
	}
	err := c.Login()
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
