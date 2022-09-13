package plugin

import (
	"github.com/chainreactors/grdp"
	"github.com/chainreactors/grdp/glog"
	"github.com/chainreactors/zombie/pkg/utils"
	"strings"
)

type RdpService struct {
	*utils.Task
	Input string
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
	s.Input = query
}

func (s *RdpService) Output(res interface{}) {

}

func RdpConnect(task *utils.Task) error {
	var username, domain string
	if strings.Contains(task.Username, "/") {
		username = strings.Split(task.Username, "/")[1]
		domain = strings.Split(task.Username, "/")[0]
	}

	//SSL协议登录测试
	client := grdp.NewClient(task.Address(), glog.NONE)
	err := client.LoginForSSL(domain, username, task.Password)
	if err == nil {
		//fmt.Println("Login Success")
		return nil
	}

	//RDP协议登录测试
	err = client.LoginForRDP(domain, username, task.Password)
	if err != nil {
		//fmt.Println("Login Success")
		return err
	}
	return nil
}
