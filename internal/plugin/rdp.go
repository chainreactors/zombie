package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	rdpclient "github.com/tomatome/grdp/client"
	"sync"
)

type RdpService struct {
	*pkg.Task
	conn *rdpclient.Client
}

func (s *RdpService) Query() bool {
	return false
}

func (s *RdpService) GetInfo() bool {
	return false
}

func (s *RdpService) Connect() error {
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

func (s *RdpService) Close() error {
	return nil
}

func (s *RdpService) SetQuery(query string) {
	//s.Input = query
}

func (s *RdpService) Output(res interface{}) {

}
