package service

import (
	"fmt"
	"net"
	"strings"

	"github.com/chainreactors/neutron/common"
	"github.com/chainreactors/neutron/operators"
	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/zombie/pkg"
)

type Template struct {
	Id        string                 `json:"id" yaml:"id"`
	Service   []string               `json:"service" yaml:"service"`
	Chains    []string               `json:"chain,omitempty" yaml:"chain,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty"`
	Info      struct {
		Name        string `json:"name" yaml:"name"`
		Severity    string `json:"severity" yaml:"severity"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
		Tags        string `json:"tags,omitempty" yaml:"tags,omitempty"`
	} `json:"info" yaml:"info"`

	Services      []*Request `json:"services" yaml:"services"`
	TotalRequests int        `json:"-" yaml:"-"`
}

func (t *Template) Match(serviceName string) bool {
	if len(t.Service) == 0 {
		return true
	}
	for _, s := range t.Service {
		if strings.EqualFold(s, serviceName) {
			return true
		}
	}
	return false
}

func (t *Template) Compile(options *protocols.ExecuterOptions) error {
	if options == nil {
		options = &protocols.ExecuterOptions{Options: &protocols.Options{}}
	}
	for _, req := range t.Services {
		if len(req.Payloads) > 0 {
			attack := req.AttackType
			if attack == "" {
				attack = "pitchfork"
			}
			if _, ok := protocols.StringToType[strings.ToLower(attack)]; !ok {
				return fmt.Errorf("unsupported attack type %q in template %s", attack, t.Id)
			}
		}
		if err := req.Compile(options); err != nil {
			return err
		}
	}
	t.TotalRequests = len(t.Services)
	if t.TotalRequests == 0 {
		return fmt.Errorf("no service requests defined in template %s", t.Id)
	}
	return nil
}

func (t *Template) Execute(session pkg.Session, host string) (*operators.Result, error) {
	return t.ExecuteWithOptions(session, host, nil, nil)
}

func (t *Template) ExecuteWithVariables(session pkg.Session, host string, cliVars map[string]interface{}) (*operators.Result, error) {
	return t.ExecuteWithOptions(session, host, cliVars, nil)
}

func (t *Template) ExecuteWithOptions(session pkg.Session, host string, cliVars, cliPayloads map[string]interface{}) (*operators.Result, error) {
	if !t.Match(session.Service()) {
		return nil, nil
	}

	payloads := map[string]interface{}{
		"_session":              session,
		"_service_cli_vars":     cliVars,
		"_service_cli_payloads": cliPayloads,
	}
	scanCtx := protocols.NewScanContext(host, payloads)
	scanCtx.GlobalVars = t.executionVariables(host, cliVars)

	var merged *operators.Result
	dynamicValues := copyMap(scanCtx.GlobalVars)

	for _, req := range t.Services {
		err := req.ExecuteWithResults(scanCtx, dynamicValues, nil, func(event *protocols.InternalWrappedEvent) {
			if event.OperatorsResult == nil {
				return
			}
			if len(event.OperatorsResult.DynamicValues) > 0 {
				if dynamicValues == nil {
					dynamicValues = make(map[string]interface{})
				}
				for k, v := range event.OperatorsResult.DynamicValues {
					if len(v) > 0 {
						dynamicValues[k] = v[0]
					}
				}
			}
			if !event.OperatorsResult.Matched && !event.OperatorsResult.Extracted {
				return
			}
			if merged == nil {
				merged = event.OperatorsResult
			} else {
				mergeResult(merged, event.OperatorsResult)
			}
		})
		if err != nil {
			return nil, err
		}
	}

	if merged == nil {
		return &operators.Result{}, nil
	}
	return merged, nil
}

func (t *Template) executionVariables(host string, cliVars map[string]interface{}) map[string]interface{} {
	vars := serviceTargetVariables(host)
	for k, v := range t.Variables {
		vars[k] = v
	}
	for k, v := range cliVars {
		vars[k] = v
	}
	evaluated := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		value := common.ToString(v)
		if strings.Contains(value, "{{") {
			if got, err := common.Evaluate(value, vars); err == nil {
				evaluated[k] = got
				continue
			}
		}
		evaluated[k] = v
	}
	return evaluated
}

func serviceTargetVariables(host string) map[string]interface{} {
	vars := map[string]interface{}{
		"Hostname": host,
		"host":     host,
	}
	if h, p, err := net.SplitHostPort(host); err == nil {
		vars["Host"] = h
		vars["Port"] = p
		vars["hostname"] = host
		return vars
	}
	vars["Host"] = host
	return vars
}

func copyMap(values map[string]interface{}) map[string]interface{} {
	copied := make(map[string]interface{}, len(values))
	for k, v := range values {
		copied[k] = v
	}
	return copied
}

func mergeResult(dst, src *operators.Result) {
	if src.Matched {
		dst.Matched = true
	}
	if src.Extracted {
		dst.Extracted = true
	}
	for k, v := range src.Matches {
		if dst.Matches == nil {
			dst.Matches = make(map[string][]string)
		}
		dst.Matches[k] = append(dst.Matches[k], v...)
	}
	for k, v := range src.Extracts {
		if dst.Extracts == nil {
			dst.Extracts = make(map[string][]string)
		}
		dst.Extracts[k] = append(dst.Extracts[k], v...)
	}
	dst.OutputExtracts = append(dst.OutputExtracts, src.OutputExtracts...)
}
