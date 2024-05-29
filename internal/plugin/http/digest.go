package http

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/xinsnake/go-http-digest-auth-client"
	"net/http"
)

type HTTPDigestPlugin struct {
	*pkg.Task
}

func (s *HTTPDigestPlugin) Name() string {
	return s.Service
}

func (s *HTTPDigestPlugin) Unauth() (bool, error) {
	return false, nil
}

func (s *HTTPDigestPlugin) Login() error {
	u := fmt.Sprintf("%s://%s:%s/", s.Service, s.IP, s.Port)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	digestClient := digest_auth_client.NewRequest(s.Username, s.Password, "GET", u, "")
	resp, err := digestClient.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to connect with digest auth, status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *HTTPDigestPlugin) GetResult() *pkg.Result {
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *HTTPDigestPlugin) Close() error {
	return nil
}
