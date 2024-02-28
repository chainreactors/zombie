package rdp

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	rdpclient "github.com/tomatome/grdp/client"
	"sync"
)

type RdpPlugin struct {
	*pkg.Task
	conn *rdpclient.Client
}

func (s *RdpPlugin) Unauth() (bool, error) {
	return false, nil
}

//func (s *RdpPlugin) Query() bool {
//	return false
//}
//
//func (s *RdpPlugin) GetInfo() bool {
//	return false
//}

func (s *RdpPlugin) Login() error {
	client := rdpclient.NewClient(s.Address(), s.Username, s.Password, rdpclient.TC_RDP, nil)
	if client == nil {
		return fmt.Errorf("init error")
	}
	err := client.Login()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	client.OnReady(func() {
		wg.Done()
	})
	wg.Wait()
	s.conn = client
	return nil
}

func (s *RdpPlugin) Close() error {
	return nil
}

func (s *RdpPlugin) Name() string {
	return s.Service.String()
}

func (s *RdpPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

//func (s *RdpPlugin) SetQuery(query string) {
//	//s.Input = query
//}
//
//func (s *RdpPlugin) Output(res interface{}) {
//
//}
