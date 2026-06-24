package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/service"
	"gopkg.in/yaml.v3"
)

type ServiceAction struct {
	templates []*service.Template
	index     map[string]*service.Template
	vars      map[string]interface{}
	payloads  map[string]interface{}
}

func NewServiceAction(templatePaths []string, vars map[string]interface{}, payloads ...map[string]interface{}) (*ServiceAction, error) {
	execOpts := &protocols.ExecuterOptions{Options: &protocols.Options{}}
	var templates []*service.Template
	for _, p := range templatePaths {
		tmpls, err := loadServiceTemplatesFromPath(p, execOpts)
		if err != nil {
			return nil, fmt.Errorf("load service templates from %s: %w", p, err)
		}
		templates = append(templates, tmpls...)
	}
	if len(templates) == 0 {
		return nil, fmt.Errorf("no service templates loaded")
	}

	index := make(map[string]*service.Template, len(templates))
	for _, t := range templates {
		index[t.Id] = t
	}

	var cliPayloads map[string]interface{}
	if len(payloads) > 0 {
		cliPayloads = payloads[0]
	}

	return &ServiceAction{templates: templates, index: index, vars: vars, payloads: cliPayloads}, nil
}

func (a *ServiceAction) Name() string { return "service" }

func (a *ServiceAction) Run(session pkg.Session, task *pkg.Task) (*pkg.ActionResult, error) {
	result := &pkg.ActionResult{}
	host := task.Address()
	executed := make(map[string]bool)

	for _, tmpl := range a.templates {
		if len(tmpl.Chains) > 0 {
			// templates with chains are entry points only — skip if chained from elsewhere
			continue
		}
		a.executeTemplate(tmpl, session, host, nil, result, executed)
	}

	// now run entry-point templates (those with chains)
	for _, tmpl := range a.templates {
		if len(tmpl.Chains) == 0 {
			continue
		}
		a.executeTemplate(tmpl, session, host, nil, result, executed)
	}

	return result, nil
}

func (a *ServiceAction) executeTemplate(tmpl *service.Template, session pkg.Session, host string, extraVars map[string]interface{}, result *pkg.ActionResult, executed map[string]bool) {
	if executed[tmpl.Id] {
		return
	}
	if !tmpl.Match(session.Service()) {
		return
	}
	executed[tmpl.Id] = true

	vars := copyVars(a.vars)
	for k, v := range extraVars {
		vars[k] = v
	}

	opResult, err := tmpl.ExecuteWithOptions(session, host, vars, a.payloads)
	if err != nil {
		logs.Log.Debugf("[service] template %s failed on %s: %v", tmpl.Id, host, err)
		return
	}
	if opResult == nil {
		return
	}

	// collect extractions into result
	if opResult.Matched || opResult.Extracted {
		for name, extracts := range opResult.Extracts {
			result.Extracteds = append(result.Extracteds, &parsers.Extracted{
				Name:          fmt.Sprintf("%s:%s", tmpl.Id, name),
				ExtractResult: extracts,
			})
		}
		for _, output := range opResult.OutputExtracts {
			result.Extracteds = append(result.Extracteds, &parsers.Extracted{
				Name:          tmpl.Id,
				ExtractResult: []string{output},
			})
		}
	}

	// execute chains — pass dynamic values from this template's result
	if len(tmpl.Chains) > 0 {
		chainVars := copyVars(vars)
		for k, v := range opResult.DynamicValues {
			if len(v) > 0 {
				chainVars[k] = v[0]
			}
		}
		for k, v := range opResult.Extracts {
			if len(v) > 0 {
				chainVars[k] = v[0]
			}
		}

		for _, chainID := range tmpl.Chains {
			target, ok := a.index[chainID]
			if !ok {
				logs.Log.Debugf("[service] chain target %q not found (from %s)", chainID, tmpl.Id)
				continue
			}
			a.executeTemplate(target, session, host, chainVars, result, executed)
		}
	}
}

func copyVars(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return make(map[string]interface{})
	}
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func loadServiceTemplatesFromPath(path string, execOpts *protocols.ExecuterOptions) ([]*service.Template, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return loadServiceTemplateFile(path, execOpts)
	}
	var templates []*service.Template
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".yaml") && !strings.HasSuffix(p, ".yml") {
			return nil
		}
		loaded, err := loadServiceTemplateFile(p, execOpts)
		if err != nil {
			logs.Log.Debugf("[service] skip %s: %v", p, err)
			return nil
		}
		templates = append(templates, loaded...)
		return nil
	})
	return templates, nil
}

func loadServiceTemplateFile(path string, execOpts *protocols.ExecuterOptions) ([]*service.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tmpl service.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, err
	}
	if len(tmpl.Services) == 0 {
		return nil, nil
	}
	if err := tmpl.Compile(execOpts); err != nil {
		return nil, err
	}
	return []*service.Template{&tmpl}, nil
}
