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
	IP       string            `json:"ip"`
	Port     string            `json:"port"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Service  string            `json:"service"`
	Scheme   string            `json:"scheme"`
	Param    map[string]string `json:"param"`
}

func (t *Target) String() string {
	return fmt.Sprintf("%s://%s:%s", t.Service, t.IP, t.Port)
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

func ParseUrl(u string) (*Target, bool) {
	var t *Target
	parsed, err := url.Parse(u)
	if err != nil && parsed != nil {
		if ip, port, err := net.SplitHostPort(u); err == nil {
			t = &Target{
				IP:   ip,
				Port: port,
			}
		} else {
			return nil, false
		}

	} else if parsed == nil {
		return nil, false
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
		t.Service = parsed.Scheme
		if t.Port == "" {
			t.Port = pkg.Services.DefaultPort(t.Service)
		}
		t.Scheme = parsed.Scheme
	} else if t.Port != "" {
		t.Service = pkg.GetDefault(t.Port)
	}
	if parsed.User != nil {
		if parsed.User.Username() != "" {
			t.Username = parsed.User.Username()
		}
		if pwd, _ := parsed.User.Password(); pwd != "" {
			t.Password, _ = parsed.User.Password()
		}
	}

	return t, true
}

func SimpleParseUrl(u string) *Target {
	result := strings.Split(u, ":")
	if len(result) == 1 {
		return &Target{
			IP: result[0],
		}
	} else {
		return &Target{
			IP:   result[0],
			Port: result[1],
		}
	}
}
