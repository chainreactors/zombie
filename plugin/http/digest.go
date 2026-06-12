package http

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	digest_auth_client "github.com/xinsnake/go-http-digest-auth-client"
	"net/http"
)

// httpDigestSession implements pkg.Session for HTTP digest auth.
// HTTP is stateless, so Close is a no-op and Raw returns the http.Client.
type httpDigestSession struct {
	service string
	client  *http.Client
}

func (s *httpDigestSession) Service() string  { return s.service }
func (s *httpDigestSession) Raw() interface{} { return s.client }
func (s *httpDigestSession) Close() error     { return nil }

// HTTPDigestPlugin is stateless; all connection state lives in httpDigestSession.
type HTTPDigestPlugin struct{}

func (p *HTTPDigestPlugin) Name() string { return "digest" }

func (p *HTTPDigestPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	u := fmt.Sprintf("%s://%s:%s/", task.Service, task.IP, task.Port)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	digestClient := digest_auth_client.NewRequest(task.Username, task.Password, "GET", u, "")
	client := task.HTTPClient(true)
	digestClient.HTTPClient = client
	resp, err := digestClient.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to connect with digest auth, status code: %d", resp.StatusCode)
	}

	return &httpDigestSession{service: task.Service, client: client}, nil
}

func (p *HTTPDigestPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}
