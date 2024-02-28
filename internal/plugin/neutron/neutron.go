package neutron

import (
	"errors"
	"fmt"
	templates "github.com/chainreactors/neutron/templates_gogo"
	"github.com/chainreactors/zombie/pkg"
)

type NeutronPlugin struct {
	*pkg.Task
	Service pkg.Service
}

func (s *NeutronPlugin) Name() string {
	return s.Service.String()
}

func (s *NeutronPlugin) Unauth() (bool, error) {
	return false, nil
}

func (s *NeutronPlugin) Login() error {
	if template, ok := pkg.TemplateMap[s.Service.String()]; ok {
		err := NeutronScan(s.Scheme, s.Address(), s.Username, s.Password, template)
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
	return pkg.NilConnError{s.Service}
}

func NeutronScan(scheme, target, user, pwd string, template *templates.Template) error {
	if scheme == "" {
		if template.RequestsHTTP != nil {
			scheme = "http"
		} else if template.RequestsNetwork != nil {
			scheme = "tcp"
		}
	}

	res, err := template.Execute(fmt.Sprintf("%s://%s", scheme, target), map[string]interface{}{
		"username": user,
		"password": pwd,
	})
	if res == nil {
		return errors.New(fmt.Sprintf("%s failed", template.Id))
	}
	return err
}
