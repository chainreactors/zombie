package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func NewHTTPPlugin(method string, task *pkg.Task) *HTTPPlugin {
	plugin := &HTTPPlugin{
		Task:   task,
		Method: method,
		Path:   task.Param["path"],
		Host:   task.Param["host"],
		Type:   task.Param["type"],
		Header: make(map[string]string),
		Forms:  make(map[string]string),
		Params: make(map[string]string),
		//Keymap: make(map[string]string),
	}
	if task.Scheme == "" {
		plugin.Scheme = "http"
	}

	if task.Param["match_status"] == "" {
		plugin.MatchStatus = "200"
	}
	if method == "GET" {
		if userParam, ok := task.Param["username"]; ok {
			plugin.Params["username"] = userParam
		} else {
			plugin.Params["username"] = "username"
		}
		if passParam, ok := task.Param["password"]; ok {
			plugin.Params["password"] = passParam
		} else {
			plugin.Params["password"] = "password"
		}
	} else if method == "POST" {
		if userParam, ok := task.Param["username"]; ok {
			plugin.Forms["username"] = userParam
		} else {
			plugin.Forms["username"] = "username"
		}
		if passParam, ok := task.Param["password"]; ok {
			plugin.Forms["password"] = passParam
		} else {
			plugin.Forms["password"] = "password"
		}
	}
	return plugin
}

type HTTPPlugin struct {
	*pkg.Task
	Path        string            `json:"path"`
	Host        string            `json:"host"`
	Method      string            `json:"method"`
	Header      map[string]string `json:"header"`
	Forms       map[string]string `json:"forms"`
	Params      map[string]string `json:"params"` // map username/password param name to target param name
	Keymap      map[string]string `json:"keymap"`
	Type        string            `json:"type"`
	MatchStatus string            `json:"match_status"`
	MatchBody   string            `json:"match_body"`
	MatchHeader string            `json:"match_header"`
}

func (s *HTTPPlugin) Name() string {
	return s.Service
}

func (s *HTTPPlugin) Unauth() (bool, error) {
	return false, pkg.NotImplUnauthorized
}

func (s *HTTPPlugin) Login() error {
	u := fmt.Sprintf("%s://%s:%s/%s", s.Scheme, s.IP, s.Port, s.Path)
	if s.Method == "" {
		s.Method = "GET"
	}

	var reqBody []byte
	var err error

	if len(s.Params) > 0 {
		// 使用 Params
		query := url.Values{}
		for key, value := range s.Params {
			if key == "username" {
				query.Set(value, s.Task.Username)
			} else if key == "password" {
				query.Set(value, s.Task.Password)
			} else {
				query.Set(key, value)
			}
		}
		reqBody = []byte(query.Encode())
		req, err := http.NewRequest(s.Method, u+"?"+query.Encode(), nil)
		if err != nil {
			return err
		}
		s.setupRequestHeaders(req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return pkg.ErrorWrongUserOrPwd
		}
		return nil
	} else if len(s.Forms) > 0 {
		// 使用 Forms
		formData := url.Values{}
		for key, value := range s.Forms {
			if key == "username" {
				formData.Set(value, s.Task.Username)
			} else if key == "password" {
				formData.Set(value, s.Task.Password)
			} else {
				formData.Set(key, value)
			}
		}

		if s.Type == "json" {
			reqBody, err = json.Marshal(formData)
			if err != nil {
				return err
			}
		} else if s.Type == "xml" {
			reqBody, err = xml.Marshal(formData)
			if err != nil {
				return err
			}
		} else {
			reqBody = []byte(formData.Encode())
		}

		req, err := http.NewRequest(s.Method, u, bytes.NewBuffer(reqBody))
		if err != nil {
			return err
		}
		s.setupRequestHeaders(req)
		if s.Type == "json" {
			req.Header.Set("Content-Type", "application/json")
		} else if s.Type == "xml" {
			req.Header.Set("Content-Type", "application/xml")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		return s.matchResponse(resp)
	}

	return fmt.Errorf("no valid params or form data provided")
}

func (s *HTTPPlugin) setupRequestHeaders(req *http.Request) {
	if s.Host != "" {
		req.Host = s.Host
	}
	req.Header.Set("User-Agent", pkg.RandomUA())
	for key, value := range s.Header {
		req.Header.Set(key, value)
	}
}

func (s *HTTPPlugin) matchResponse(resp *http.Response) error {
	if iutils.ToString(resp.StatusCode) != s.MatchStatus {
		return pkg.ErrorWrongUserOrPwd
	}

	if s.MatchBody != "" {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, s.MatchBody) {
			return pkg.ErrorWrongUserOrPwd
		}
	}

	if s.MatchHeader != "" {
		matchFound := false
		for key, values := range resp.Header {
			for _, value := range values {
				if key == s.MatchHeader || value == s.MatchHeader {
					matchFound = true
					break
				}
			}
			if matchFound {
				break
			}
		}
		if !matchFound {
			return pkg.ErrorWrongUserOrPwd
		}
	}

	return nil
}

func (s *HTTPPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *HTTPPlugin) Close() error {
	return nil
}
