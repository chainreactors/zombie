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
	Url string `json:"url"`
}

func (s *Socks5Plugin) Unauth() (bool, error) {
	proxyURL, _ := url.Parse(fmt.Sprintf("socks5://%s:%s", s.IP, s.Port))
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return false, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	if s.Url == "" {
		s.Url = "http://baidu.com"
	}
	req, err := http.NewRequest("GET", s.Url, nil)
	_, err = client.Do(req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Socks5Plugin) Login() error {
	proxyURL, err := url.Parse(fmt.Sprintf("socks5://%s:%s@%s:%s", s.Username, s.Password, s.IP, s.Port))
	if err != nil {
		return err
	}
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	if s.Url == "" {
		s.Url = "http://baidu.com"
	}
	req, err := http.NewRequest("GET", s.Url, nil)
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
