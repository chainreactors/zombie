package http

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/http"
	"net/url"
)

// httpProxySession implements pkg.Session for HTTP proxy test results.
// HTTP is stateless, so Close is a no-op.
type httpProxySession struct {
	service string
	client  *http.Client
}

func (s *httpProxySession) Service() string  { return s.service }
func (s *httpProxySession) Raw() interface{} { return s.client }
func (s *httpProxySession) Close() error     { return nil }

// HTTPProxyPlugin is stateless; all connection state lives in httpProxySession.
type HTTPProxyPlugin struct{}

func (p *HTTPProxyPlugin) Name() string { return "http_proxy" }

func (p *HTTPProxyPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", task.Scheme, task.IP, task.Port))
	if err != nil {
		return nil, err
	}

	// Set proxy authentication
	proxyURL.User = url.UserPassword(task.Username, task.Password)

	testURL := task.Param["url"]
	if testURL == "" {
		testURL = "http://baidu.com"
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusProxyAuthRequired {
		return nil, pkg.ErrorWrongUserOrPwd
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &httpProxySession{service: task.Service, client: client}, nil
}

func (p *HTTPProxyPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", task.Scheme, task.IP, task.Port))
	if err != nil {
		return nil, err
	}

	testURL := task.Param["url"]
	if testURL == "" {
		testURL = "http://baidu.com"
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusProxyAuthRequired {
		return nil, pkg.ErrorWrongUserOrPwd
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &httpProxySession{service: task.Service, client: client}, nil
}
