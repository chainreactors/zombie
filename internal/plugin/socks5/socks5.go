package socks5

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
)

type Socks5Plugin struct {
	*pkg.Task
	Socks5Inf
}

func (s *Socks5Plugin) Unauth() (bool, error) {
	Socks5Url := fmt.Sprintf("%s://%s:%s@%s:%s", s.Service, "", "", s.IP, s.Port)
	proxyURL, _ := url.Parse(Socks5Url)
	password, _ := proxyURL.User.Password()
	dialer, _ := proxy.SOCKS5("tcp", proxyURL.Host, &proxy.Auth{User: proxyURL.User.Username(), Password: password}, proxy.Direct)
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	// 使用HTTP Transport创建HTTP客户端
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	_, err = client.Do(req)
	if err != nil {
		return false, err
	}
	return true, nil
}

type Socks5Inf struct {
	Url string `json:"url"`
}

func (s *Socks5Plugin) Login() error {
	Socks5Url := fmt.Sprintf("%s://%s:%s@%s:%s", s.Service, s.Username, s.Password, s.IP, s.Port)
	proxyURL, err := url.Parse(Socks5Url)
	if err != nil {
		return err
	}
	password, _ := proxyURL.User.Password()
	dialer, _ := proxy.SOCKS5("tcp", proxyURL.Host, &proxy.Auth{User: proxyURL.User.Username(), Password: password}, proxy.Direct)
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	// 使用HTTP Transport创建HTTP客户端
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil

}

func (s *Socks5Plugin) Close() error {
	return pkg.NilConnError{s.Service}
}

func (s *Socks5Plugin) Name() string {
	return s.Service.String()
}

func (s *Socks5Plugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}
