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

// httpSession implements pkg.Session for HTTP GET/POST login.
// HTTP is stateless, so Close is a no-op and Raw returns the http.Client.
type httpSession struct {
	service string
	client  *http.Client
}

func (s *httpSession) Service() string  { return s.service }
func (s *httpSession) Raw() interface{} { return s.client }
func (s *httpSession) Close() error     { return nil }

// HTTPPlugin is stateless; all per-request state is derived from the task.
type HTTPPlugin struct {
	Method string
}

func NewHTTPPlugin(method string) *HTTPPlugin {
	return &HTTPPlugin{Method: method}
}

func (p *HTTPPlugin) Name() string { return strings.ToLower(p.Method) }

func (p *HTTPPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	path := task.Param["path"]
	host := task.Param["host"]
	contentType := task.Param["type"]
	matchStatus := task.Param["match_status"]
	matchBody := task.Param["match_body"]
	matchHeader := task.Param["match_header"]

	scheme := task.Scheme
	if scheme == "" {
		scheme = "http"
	}
	if matchStatus == "" {
		matchStatus = "200"
	}

	u := fmt.Sprintf("%s://%s:%s/%s", scheme, task.IP, task.Port, path)
	method := p.Method
	if method == "" {
		method = "GET"
	}

	// Build params / forms from task.Param
	params := make(map[string]string)
	forms := make(map[string]string)
	headers := make(map[string]string)

	if method == "GET" {
		if userParam, ok := task.Param["username"]; ok {
			params["username"] = userParam
		} else {
			params["username"] = "username"
		}
		if passParam, ok := task.Param["password"]; ok {
			params["password"] = passParam
		} else {
			params["password"] = "password"
		}
	} else if method == "POST" {
		if userParam, ok := task.Param["username"]; ok {
			forms["username"] = userParam
		} else {
			forms["username"] = "username"
		}
		if passParam, ok := task.Param["password"]; ok {
			forms["password"] = passParam
		} else {
			forms["password"] = "password"
		}
	}

	client := task.HTTPClient(true)

	if len(params) > 0 {
		query := url.Values{}
		for key, value := range params {
			if key == "username" {
				query.Set(value, task.Username)
			} else if key == "password" {
				query.Set(value, task.Password)
			} else {
				query.Set(key, value)
			}
		}
		req, err := http.NewRequest(method, u+"?"+query.Encode(), nil)
		if err != nil {
			return nil, err
		}
		setupRequestHeaders(req, host, headers)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, pkg.ErrorWrongUserOrPwd
		}
		return &httpSession{service: task.Service, client: client}, nil
	} else if len(forms) > 0 {
		formData := url.Values{}
		for key, value := range forms {
			if key == "username" {
				formData.Set(value, task.Username)
			} else if key == "password" {
				formData.Set(value, task.Password)
			} else {
				formData.Set(key, value)
			}
		}

		var reqBody []byte
		var err error
		if contentType == "json" {
			reqBody, err = json.Marshal(formData)
			if err != nil {
				return nil, err
			}
		} else if contentType == "xml" {
			reqBody, err = xml.Marshal(formData)
			if err != nil {
				return nil, err
			}
		} else {
			reqBody = []byte(formData.Encode())
		}

		req, err := http.NewRequest(method, u, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
		setupRequestHeaders(req, host, headers)
		if contentType == "json" {
			req.Header.Set("Content-Type", "application/json")
		} else if contentType == "xml" {
			req.Header.Set("Content-Type", "application/xml")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		err = matchResponse(resp, matchStatus, matchBody, matchHeader)
		if err != nil {
			return nil, err
		}
		return &httpSession{service: task.Service, client: client}, nil
	}

	return nil, fmt.Errorf("no valid params or form data provided")
}

func (p *HTTPPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}

func setupRequestHeaders(req *http.Request, host string, headers map[string]string) {
	if host != "" {
		req.Host = host
	}
	req.Header.Set("User-Agent", pkg.RandomUA())
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func matchResponse(resp *http.Response, matchStatus, matchBody, matchHeader string) error {
	if iutils.ToString(resp.StatusCode) != matchStatus {
		return pkg.ErrorWrongUserOrPwd
	}

	if matchBody != "" {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, matchBody) {
			return pkg.ErrorWrongUserOrPwd
		}
	}

	if matchHeader != "" {
		matchFound := false
		for key, values := range resp.Header {
			for _, value := range values {
				if key == matchHeader || value == matchHeader {
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
