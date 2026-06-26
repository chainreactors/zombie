package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/neutron/templates"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/service"
	"gopkg.in/yaml.v3"
)

type ServiceAction struct {
	index    map[string]*service.Template
	chain    *templates.ChainExecutor
	vars     map[string]interface{}
	payloads map[string]interface{}
	risk     string
	tags     []string
}

func NewServiceAction(loaded []*service.Template, vars map[string]interface{}, payloads ...map[string]interface{}) (*ServiceAction, error) {
	if len(loaded) == 0 {
		return nil, fmt.Errorf("no service templates loaded")
	}

	index := make(map[string]*service.Template, len(loaded))
	chain := templates.NewChainExecutor(templates.ChainConfig{
		DepthFirst:    true,
		PassVariables: true,
	})
	for _, t := range loaded {
		index[t.Id] = t
		chain.Add(t.Id, t.Chains)
	}

	var cliPayloads map[string]interface{}
	if len(payloads) > 0 {
		cliPayloads = payloads[0]
	}

	return &ServiceAction{
		index:    index,
		chain:    chain,
		vars:     vars,
		payloads: cliPayloads,
	}, nil
}

func (a *ServiceAction) SetRisk(risk string) { a.risk = risk }
func (a *ServiceAction) SetTags(tags []string) { a.tags = tags }

func (a *ServiceAction) Name() string { return "service" }

func (a *ServiceAction) Run(session pkg.Session, task *pkg.Task) (*pkg.ActionResult, error) {
	result := &pkg.ActionResult{}
	host := task.Address()
	svc := session.Service()

	a.chain.Execute(a.chain.Entrypoints(), func(id string, vars map[string]interface{}) *templates.ChainResult {
		tmpl, ok := a.index[id]
		if !ok || !tmpl.Match(svc) {
			return nil
		}
		if !tmpl.RiskAllowed(a.risk) {
			return nil
		}
		if len(a.tags) > 0 && !a.matchTags(tmpl) {
			return nil
		}

		mergedVars := copyVars(a.vars)
		for k, v := range vars {
			mergedVars[k] = v
		}

		opResult, err := tmpl.ExecuteWithOptions(session, host, mergedVars, a.payloads)
		if err != nil {
			logs.Log.Debugf("[service] template %s failed on %s: %v", id, host, err)
			return nil
		}
		if opResult == nil {
			return nil
		}

		if opResult.Matched || opResult.Extracted {
			for name, extracts := range opResult.Extracts {
				result.Extracteds = append(result.Extracteds, &parsers.Extracted{
					Name:          fmt.Sprintf("%s:%s", id, name),
					ExtractResult: extracts,
				})
			}
			for _, output := range opResult.OutputExtracts {
				result.Extracteds = append(result.Extracteds, &parsers.Extracted{
					Name:          id,
					ExtractResult: []string{output},
				})
			}
		}

		if opResult.Response != "" {
			if result.Loot == nil {
				result.Loot = make(map[string][]byte)
			}
			result.Loot[id] = []byte(opResult.Response)
		}

		chainVars := copyVars(mergedVars)
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
		return &templates.ChainResult{Vars: chainVars}
	})

	return result, nil
}

func (a *ServiceAction) matchTags(tmpl *service.Template) bool {
	for _, tag := range a.tags {
		if tmpl.HasTag(tag) {
			return true
		}
	}
	return false
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

func LoadServiceTemplatesFromPaths(paths []string) ([]*service.Template, error) {
	execOpts := &protocols.ExecuterOptions{Options: &protocols.Options{}}
	var all []*service.Template
	for _, p := range paths {
		tmpls, err := loadServiceTemplatesFromPath(p, execOpts)
		if err != nil {
			return nil, fmt.Errorf("load service templates from %s: %w", p, err)
		}
		all = append(all, tmpls...)
	}
	return all, nil
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
	return loadServiceTemplateBytes(data, execOpts)
}

func loadServiceTemplateBytes(data []byte, execOpts *protocols.ExecuterOptions) ([]*service.Template, error) {
	var tmpl service.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, err
	}
	if len(tmpl.Services) == 0 && len(tmpl.RequestsHTTP) == 0 && len(tmpl.RequestsNetwork) == 0 {
		return nil, nil
	}
	if err := tmpl.Compile(execOpts); err != nil {
		return nil, err
	}
	return []*service.Template{&tmpl}, nil
}

func LoadServiceTemplatesFromData(data []byte) ([]*service.Template, error) {
	if len(data) == 0 {
		return nil, nil
	}
	execOpts := &protocols.ExecuterOptions{Options: &protocols.Options{}}
	var list []service.Template
	if err := yaml.Unmarshal(data, &list); err != nil {
		return loadServiceTemplateBytes(data, execOpts)
	}
	var all []*service.Template
	for i := range list {
		tmpl := &list[i]
		if len(tmpl.Services) == 0 && len(tmpl.RequestsHTTP) == 0 && len(tmpl.RequestsNetwork) == 0 {
			continue
		}
		if err := tmpl.Compile(execOpts); err != nil {
			logs.Log.Debugf("[service] skip embedded template %s: %v", tmpl.Id, err)
			continue
		}
		all = append(all, tmpl)
	}
	return all, nil
}
