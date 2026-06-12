package socks5

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
)

// socks5Session implements pkg.Session over a SOCKS5 proxy dialer.
type socks5Session struct {
	service string
	dialer  proxy.Dialer
}

func (s *socks5Session) Service() string  { return s.service }
func (s *socks5Session) Raw() interface{} { return s.dialer }
func (s *socks5Session) Close() error     { return nil }

// Socks5Plugin is stateless; all connection state lives in socks5Session.
type Socks5Plugin struct{}

func (p *Socks5Plugin) Name() string { return "socks5" }

func (p *Socks5Plugin) Open(task *pkg.Task) (pkg.Session, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("socks5://%s:%s@%s:%s", task.Username, task.Password, task.IP, task.Port))
	if err != nil {
		return nil, err
	}
	return dialAndTest(task, proxyURL)
}

func (p *Socks5Plugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	proxyURL, _ := url.Parse(fmt.Sprintf("socks5://%s:%s", task.IP, task.Port))
	return dialAndTest(task, proxyURL)
}

func dialAndTest(task *pkg.Task, proxyURL *url.URL) (pkg.Session, error) {
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	testURL := task.Param["url"]
	if testURL == "" {
		testURL = "http://baidu.com"
	}
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return nil, err
	}
	_, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return &socks5Session{service: task.Service, dialer: dialer}, nil
}
