package neutron

import (
	"errors"
	"fmt"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/neutron/common"
	templates "github.com/chainreactors/neutron/templates"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/zombie/pkg"
)

func init() {
	common.NeutronLog = logs.NewLogger(logs.Debug - 1)
}

type NeutronPlugin struct {
	*pkg.Task
	Service string
	Host    string
}

func (s *NeutronPlugin) Name() string {
	return s.Service
}

func (s *NeutronPlugin) Unauth() (bool, error) {
	if template, ok := pkg.TemplateMap[s.Service]; ok {
		var err error
		var usr, pwd string
		usr, pwd, err = NeutronScan(s.Scheme, s.Address(), nil, template)
		if err != nil {
			return false, err
		}

		s.Task.Username = usr
		s.Task.Password = pwd
		return true, nil
	}
	return false, errors.New("no template found")
}

func (s *NeutronPlugin) Login() error {
	if template, ok := pkg.TemplateMap[s.Service]; ok {
		_, _, err := NeutronScan(s.Scheme,
			s.Address(),
			map[string]interface{}{
				"username": s.Username,
				"password": s.Password,
			},
			template)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("no template found")
}

func (s *NeutronPlugin) GetResult() *pkg.Result {
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *NeutronPlugin) Close() error {
	return nil
}

func NeutronScan(scheme, target string, payload map[string]interface{}, template *templates.Template) (string, string, error) {
	if scheme == "" {
		if template.RequestsHTTP != nil {
			scheme = "http"
		} else if template.RequestsNetwork != nil {
			scheme = "tcp"
		}
	} else if scheme != "http" && scheme != "https" && scheme != "tcp" {
		scheme = "http"
	}

	res, err := template.Execute(fmt.Sprintf("%s://%s", scheme, target), payload)
	if err != nil {
		return "", "", err
	}
	if res == nil {
		return "", "", errors.New(fmt.Sprintf("nil result, %s", template.Id))
	}
	if !res.Matched {
		return "", "", errors.New(fmt.Sprintf("not matched, %s", template.Id))
	}
	return iutils.ToString(res.PayloadValues["username"]), iutils.ToString(res.PayloadValues["password"]), nil
}
