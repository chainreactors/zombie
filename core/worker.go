package core

import (
	"errors"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

var ErrNoUnauth = errors.New("cannot unauth login")
var ErrNoPlugin = errors.New("no plugin for service")

func Execute(task *pkg.Task, plugins map[string]plugin.Plugin, pipeline []pkg.Action) *pkg.Result {
	p := resolvePlugin(task.Service, plugins)
	if p == nil {
		return pkg.NewResult(task, ErrNoPlugin)
	}

	session, err := p.Open(task)
	if err != nil {
		return pkg.NewResult(task, err)
	}
	defer session.Close()

	result := &pkg.Result{Task: task, OK: true}
	for _, action := range pipeline {
		ar, err := action.Run(session, task)
		if err != nil {
			continue
		}
		result.Merge(ar)
	}
	return result
}

func ExecuteUnauth(task *pkg.Task, plugins map[string]plugin.Plugin, pipeline []pkg.Action) *pkg.Result {
	p := resolvePlugin(task.Service, plugins)
	if p == nil {
		return pkg.NewResult(task, ErrNoPlugin)
	}

	session, err := p.Unauth(task)
	if err != nil {
		return pkg.NewResult(task, err)
	}
	if session == nil {
		return pkg.NewResult(task, ErrNoUnauth)
	}
	defer session.Close()

	result := &pkg.Result{Task: task, OK: true}
	for _, action := range pipeline {
		ar, _ := action.Run(session, task)
		result.Merge(ar)
	}
	return result
}

func resolvePlugin(service string, plugins map[string]plugin.Plugin) plugin.Plugin {
	if p, ok := plugins[service]; ok {
		return p
	}
	if p, ok := plugins["neutron"]; ok {
		return p
	}
	return nil
}
