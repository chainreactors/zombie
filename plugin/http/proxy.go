package http

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/http"
	"net/url"
)

type HTTPProxyPlugin struct {
	*pkg.Task
	TestURL string `json:"url"`
}

func (s *HTTPProxyPlugin) Name() string {
	return s.Service
}

func (s *HTTPProxyPlugin) Unauth() (bool, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", s.Scheme, s.IP, s.Port))
	if err != nil {
		return false, err
	}

	if s.TestURL == "" {
		s.TestURL = "http://baidu.com"
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", s.TestURL, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 检查是否通过认证
	if resp.StatusCode == http.StatusProxyAuthRequired {
		return false, pkg.ErrorWrongUserOrPwd
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	s.Username = ""
	return true, nil
}

func (s *HTTPProxyPlugin) Login() error {
	proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", s.Scheme, s.IP, s.Port))
	if err != nil {
		return err
	}

	// 设置代理认证
	proxyURL.User = url.UserPassword(s.Username, s.Password)
	if s.TestURL == "" {
		s.TestURL = "http://baidu.com"
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", s.TestURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查是否通过认证
	if resp.StatusCode == http.StatusProxyAuthRequired {
		return pkg.ErrorWrongUserOrPwd
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *HTTPProxyPlugin) GetResult() *pkg.Result {
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *HTTPProxyPlugin) Close() error {
	return nil
}
