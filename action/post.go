package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/proton/proton/file"
	"github.com/chainreactors/proton/template"
	"github.com/chainreactors/zombie/pkg"
	"gopkg.in/yaml.v3"
)

type PostAction struct {
	scanner *file.Scanner
}

func NewPostAction(templatePaths []string) (*PostAction, error) {
	execOpts := &protocols.ExecuterOptions{Options: &protocols.Options{}}
	var rules []file.Rule
	for _, p := range templatePaths {
		tmpls, err := loadTemplatesFromPath(p, execOpts)
		if err != nil {
			return nil, fmt.Errorf("load templates from %s: %w", p, err)
		}
		for _, tmpl := range tmpls {
			if len(tmpl.RequestsFile) > 0 {
				rules = append(rules, file.Rule{
					ID: tmpl.Id, Name: tmpl.Info.Name,
					Severity: tmpl.Info.Severity, Requests: tmpl.RequestsFile,
				})
			}
		}
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("no file rules in loaded templates")
	}
	return &PostAction{scanner: file.NewScanner(rules, execOpts)}, nil
}

func (a *PostAction) Name() string { return "post" }

func (a *PostAction) Run(session pkg.Session, task *pkg.Task) (*pkg.ActionResult, error) {
	return nil, nil
}

func (a *PostAction) ScanData(data []byte, label string) []*parsers.Extracted {
	if a.scanner == nil || len(data) == 0 {
		return nil
	}
	var results []*parsers.Extracted
	for _, group := range a.scanner.Groups {
		findings := a.scanner.ScanData(data, label, group)
		for _, f := range findings {
			var extracts []string
			for _, e := range f.Extracts {
				extracts = append(extracts, e.Value)
			}
			for _, events := range f.Matches {
				for _, e := range events {
					extracts = append(extracts, e.Value)
				}
			}
			if len(extracts) > 0 {
				results = append(results, &parsers.Extracted{
					Name:          fmt.Sprintf("%s:%s", f.TemplateID, label),
					Severity:      f.Severity,
					ExtractResult: extracts,
				})
			}
		}
	}
	return results
}

func loadTemplatesFromPath(path string, execOpts *protocols.ExecuterOptions) ([]*template.Template, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return loadTemplateFile(path, execOpts)
	}
	var tmpls []*template.Template
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".yaml") && !strings.HasSuffix(p, ".yml") {
			return nil
		}
		loaded, err := loadTemplateFile(p, execOpts)
		if err != nil {
			return nil
		}
		tmpls = append(tmpls, loaded...)
		return nil
	})
	return tmpls, nil
}

func loadTemplateFile(path string, execOpts *protocols.ExecuterOptions) ([]*template.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tmpl template.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, err
	}
	if len(tmpl.RequestsFile) == 0 {
		return nil, nil
	}
	if err := tmpl.Compile(execOpts); err != nil {
		return nil, err
	}
	return []*template.Template{&tmpl}, nil
}
