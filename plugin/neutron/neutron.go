package neutron

import (
	"errors"
	"fmt"

	"github.com/chainreactors/logs"
	neutroncommon "github.com/chainreactors/neutron/common"
	templates "github.com/chainreactors/neutron/templates"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/zombie/pkg"
)

func init() {
	if neutroncommon.NeutronLog == nil {
		neutroncommon.NeutronLog = logs.Log
	}
	if neutroncommon.NeutronLog != nil {
		neutroncommon.NeutronLog.SetLevel(logs.ErrorLevel)
	}
}

type neutronSession struct {
	service string
}

func (s *neutronSession) Service() string  { return s.service }
func (s *neutronSession) Close() error     { return nil }
func (s *neutronSession) Raw() interface{} { return nil }

type NeutronPlugin struct{}

func (p *NeutronPlugin) Name() string { return "neutron" }

func (p *NeutronPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	template, ok := pkg.TemplateMap[task.Service]
	if !ok {
		return nil, errors.New("no template found")
	}
	_, _, err := NeutronScan(task.Scheme,
		task.Address(),
		map[string]interface{}{
			"username": task.Username,
			"password": task.Password,
		},
		template)
	if err != nil {
		return nil, err
	}
	return &neutronSession{service: task.Service}, nil
}

func (p *NeutronPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	template, ok := pkg.TemplateMap[task.Service]
	if !ok {
		return nil, errors.New("no template found")
	}
	usr, pwd, err := NeutronScan(task.Scheme, task.Address(), nil, template)
	if err != nil {
		return nil, err
	}
	task.Username = usr
	task.Password = pwd
	return &neutronSession{service: task.Service}, nil
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
		return "", "", fmt.Errorf("nil result, %s", template.Id)
	}
	if !res.Matched {
		return "", "", fmt.Errorf("not matched, %s", template.Id)
	}
	return iutils.ToString(res.PayloadValues["username"]), iutils.ToString(res.PayloadValues["password"]), nil
}
