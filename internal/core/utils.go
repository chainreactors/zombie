package core

import (
	"fmt"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/zombie/pkg"
	"net"
	"net/url"
	"strings"
)

type Target struct {
	IP      string            `json:"ip"`
	Port    string            `json:"port"`
	Service pkg.Service       `json:"service"`
	Param   map[string]string `json:"param"`
}

func (t *Target) String() string {
	return fmt.Sprintf("%s://%s:%s", t.Service, t.IP, t.Port)
}

func (t *Target) UpdateService(s string) {
	t.Service = pkg.Service(strings.ToLower(s))
	if t.Port == "" {
		t.Port = t.Service.DefaultPort()
	}
}

func (t *Target) Addr() *utils.Addr {
	return &utils.Addr{IP: utils.ParseIP(t.IP), Port: t.Port}
}

func ParseUrl(u string) (*Target, bool) {
	var t *Target
	parsed, err := url.Parse(u)
	if err != nil {
		if ip, port, err := net.SplitHostPort(u); err == nil {
			t = &Target{
				IP:   ip,
				Port: port,
			}
		} else {
			return nil, false
		}

	}

	if parsed.Host == "" {
		if utils.IsIp(u) {
			return &Target{
				IP: u,
			}, true
		}
		return nil, false
	}

	t = &Target{
		IP: parsed.Hostname(),
	}
	if parsed.Port() != "" {
		t.Port = parsed.Port()
	}

	if parsed.Scheme != "" {
		t.Service = pkg.Service(parsed.Scheme)
		if t.Port == "" {
			t.Port = t.Service.DefaultPort()
		}
	} else if t.Port != "" {
		t.Service = pkg.GetDefault(t.Port)
	}

	return t, true
}

func SimpleParseUrl(u string) *Target {
	result := strings.Split(u, ":")
	t := &Target{
		IP:   result[0],
		Port: result[1],
	}
	return t
}
