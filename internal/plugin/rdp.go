package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	rdpclient "github.com/tomatome/grdp/client"
	"sync"
)

type RdpService struct {
	*pkg.Task
}

func (s *RdpService) Query() bool {
	return false
}

func (s *RdpService) GetInfo() bool {
	return false
}

func (s *RdpService) Connect() error {
	err := RdpConnect(s.Task)
	if err != nil {
		return err
	}
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

func RdpConnect(task *pkg.Task) error {
	client := rdpclient.NewClient(task.Address(), task.Username, task.Password, rdpclient.TC_RDP, nil)
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
	return nil
}
