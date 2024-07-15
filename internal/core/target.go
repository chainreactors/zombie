package core

import (
	"fmt"
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/fingers/fingers"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/httputils"
	"github.com/chainreactors/zombie/pkg"
	"net"
	"net/http"
	"strings"
)

type Target struct {
	IP       string            `json:"ip"`
	Port     string            `json:"port"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Service  string            `json:"service"`
	Scheme   string            `json:"scheme"`
	Param    map[string]string `json:"param"`
	Network  string            `json:"network"`
	conn     net.Conn
}

func (t *Target) String() string {
	return fmt.Sprintf("%s://%s:%s", t.Service, t.IP, t.Port)
}

func (t *Target) Address() string {
	return fmt.Sprintf("%s:%s", t.IP, t.Port)
}

func (t *Target) URL() string {
	if t.Scheme != "" {
		return fmt.Sprintf("%s://%s", t.Scheme, t.Address())
	} else {
		return fmt.Sprintf("http://%s/", t.Address())
	}
}

func (t *Target) UpdateService(s string) {
	t.Service = strings.ToLower(s)
	if t.Port == "" {
		t.Port = pkg.Services.DefaultPort(t.Service)
	}
}

func (t *Target) Addr() *utils.Addr {
	return &utils.Addr{IP: utils.ParseIP(t.IP), Port: t.Port}
}

func (t *Target) Conn() (net.Conn, error) {
	if t.conn == nil {
		if t.Network == "" {
			t.Network = "tcp"
		}
		conn, err := net.Dial(t.Network, t.Address())
		if err != nil {
			logs.Log.Debugf("Dial %s error: %s", t.Address(), err.Error())
			return nil, err
		}
		t.conn = conn
		return t.conn, nil
	}
	return t.conn, nil

}
func (t *Target) CheckOpen() bool {
	_, err := t.Conn()
	if err != nil {
		logs.Log.Debugf("%s connect error, %s", t.String(), err.Error())
		return false
	}
	return true
}

func (t *Target) CheckFinger() bool {
	var frames common.Frameworks
	if group, ok := pkg.FingersEngine.SocketGroup[t.Port]; ok {
		frames, _ = group.Match(fingers.NewContent(nil, "", false), 1, t.socketSender, nil, true)
	} else {
		resp, err := http.Get(t.URL())
		if err != nil {
			return false
		}
		frames, _ = pkg.FingersEngine.HTTPMatch(httputils.ReadRaw(resp), "")
	}
	if _, ok := frames[t.Service]; ok {
		return true
	}
	logs.Log.Debugf("%s finger:%s not match, found: %v", t.String(), t.Service, frames.GetNames())
	return false
}

//func (t *Target) httpSender(sendData []byte) ([]byte, bool) {
//	url := t.URL() + string(sendData)
//	resp, err := http.Get(url)
//	if err == nil {
//		return httputils.ReadRaw(resp), true
//	} else {
//		return nil, false
//	}
//}

func (t *Target) socketSender(sendData []byte) ([]byte, bool) {
	conn, err := t.Conn()
	if err != nil {
		return nil, false
	}

	_, err = conn.Write(sendData)
	if err != nil {
		return nil, false
	}
	var data = make([]byte, 1024)
	_, err = conn.Read(data)
	if err != nil {
		return nil, false
	}
	return data, true
}
