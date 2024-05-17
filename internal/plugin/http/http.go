package http

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/http"
)

type HttpPlugin struct {
	*pkg.Task
	Path   string `json:"path"`
	Host   string `json:"host"`
	Method string `json:"method"`
}

func (s *HttpPlugin) Name() string {
	return s.Service
}

func (s *HttpPlugin) Unauth() (bool, error) {
	return false, nil
}

func (s *HttpPlugin) Login() error {
	url := fmt.Sprintf("%s://%s:%s/%s", s.Service, s.IP, s.Port, s.Path)
	if s.Method == "" {
		s.Method = "GET"
	}
	req, err := http.NewRequest(s.Method, url, nil)
	if err != nil {
		return err
	}
	if s.Host != "" {
		req.Host = s.Host
	}
	req.Header.Set("User-Agent", pkg.RandomUA())
	req.SetBasicAuth(s.Username, s.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return pkg.ErrorWrongUserOrPwd
	}
	return nil
}

func (s *HttpPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *HttpPlugin) Close() error {
	return nil
}
