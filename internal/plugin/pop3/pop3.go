package pop3

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/knadh/go-pop3"
	"strconv"
)

type Pop3Plugin struct {
	*pkg.Task
}

func (s *Pop3Plugin) Unauth() (bool, error) {
	return false, nil
}

func (s *Pop3Plugin) Login() error {
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

func (s *Pop3Plugin) Name() string {
	return s.Service
}

func (s *Pop3Plugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *Pop3Plugin) Close() error {
	return pkg.NilConnError{s.Service}
}

//
//func (s *Pop3Plugin) SetQuery(query string) {
//	//s.Input = query
//}
//
//func (s *Pop3Plugin) Output(res interface{}) {
//
//}
