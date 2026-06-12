package http

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/http"
)

// httpAuthSession implements pkg.Session for HTTP basic auth.
// HTTP is stateless, so Close is a no-op and Raw returns the http.Client.
type httpAuthSession struct {
	service string
	client  *http.Client
}

func (s *httpAuthSession) Service() string  { return s.service }
func (s *httpAuthSession) Raw() interface{} { return s.client }
func (s *httpAuthSession) Close() error     { return nil }

// HttpAuthPlugin is stateless; all connection state lives in httpAuthSession.
type HttpAuthPlugin struct{}

func (p *HttpAuthPlugin) Name() string { return "http" }

func (p *HttpAuthPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	path := task.Param["path"]
	host := task.Param["host"]
	method := task.Param["method"]

	url := fmt.Sprintf("%s://%s:%s/%s", task.Service, task.IP, task.Port, path)
	if method == "" {
		method = "GET"
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if host != "" {
		req.Host = host
	}
	req.Header.Set("User-Agent", pkg.RandomUA())
	req.SetBasicAuth(task.Username, task.Password)

	client := task.HTTPClient(true)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, pkg.ErrorWrongUserOrPwd
	}
	return &httpAuthSession{service: task.Service, client: client}, nil
}

func (p *HttpAuthPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}
