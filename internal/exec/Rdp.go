package exec

import (
	"fmt"
	"github.com/chainreactors/grdp"
	"github.com/chainreactors/zombie/pkg/utils"
	"strings"
)

type RdpService struct {
	utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Input    string
}

func (s *RdpService) Query() bool {
	return false
}

func (s *RdpService) GetInfo() bool {
	return false
}

func (s *RdpService) Connect() bool {
	err, res := RdpConnectTest(s.Username, s.Password, s.IpInfo)
	if err == nil && res {
		return true
	}
	return false

}

func (s *RdpService) DisConnect() bool {
	return false
}

func (s *RdpService) SetQuery(query string) {
	s.Input = query
}

func (s *RdpService) Output(res interface{}) {

}

func RdpConnectTest(User string, Password string, info utils.IpInfo) (err error, result bool) {

	var UserName, DoaminName string

	if strings.Contains(User, "/") {
		UserName = strings.Split(User, "/")[1]
		DoaminName = strings.Split(User, "/")[0]
	} else {
		UserName = User
		DoaminName = ""
	}

	res, _ := RdpConnect(info.Ip, DoaminName, UserName, Password, info.Port)

	if res == true {
		result = true
	}

	return err, result
}

func RdpConnect(ip, domain, username, password string, port int) (bool, error) {
	target := fmt.Sprintf("%v:%v", ip, port)
	//SSL协议登录测试
	err := grdp.LoginForSSL(target, domain, username, password)
	if err == nil {
		//fmt.Println("Login Success")
		return true, nil
	}
	if err.Error() != "PROTOCOL_RDP" {
		//fmt.Println("Login Error:", err)
		return false, err
	}
	//RDP协议登录测试
	err = grdp.LoginForRDP(target, domain, username, password)
	if err == nil {
		//fmt.Println("Login Success")
		return true, err
	} else {
		//fmt.Println("Login Error:", err)
		return false, err
	}
}
