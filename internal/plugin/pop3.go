package plugin

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/knadh/go-pop3"
	"strconv"
)

type Pop3Service struct {
	*pkg.Task
}

func (s *Pop3Service) Query() bool {
	return false
}

func (s *Pop3Service) GetInfo() bool {
	return false
}

func (s *Pop3Service) Connect() error {
	port, _ := strconv.Atoi(s.Port)

	p := pop3.New(pop3.Opt{
		Host:       s.IP,
		Port:       port,
		TLSEnabled: false,
	})

	c, err := p.NewConn()
	if err != nil {
		return err
	}
	defer c.Quit()

	// Authenticate.
	if err := c.Auth(s.Username, s.Password); err != nil {
		return err
	}

	return nil

}

func (s *Pop3Service) Close() error {
	return NilConnError{s.Service}
}

func (s *Pop3Service) SetQuery(query string) {
	//s.Input = query
}

func (s *Pop3Service) Output(res interface{}) {

}
