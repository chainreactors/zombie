package core

import (
	"errors"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/action"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

var ErrNoUnauth = errors.New("cannot unauth login")
var ErrNoPlugin = errors.New("no plugin for service")

func Execute(task *pkg.Task, plugins map[string]plugin.Plugin, fallback plugin.Plugin, pipeline []pkg.Action, postAction *action.PostAction) *pkg.Result {
	p := resolvePlugin(task.Service, plugins, fallback)
	if p == nil {
		return pkg.NewResult(task, ErrNoPlugin)
	}

	session, err := p.Open(task)
	if err != nil {
		return pkg.NewResult(task, err)
	}
	if session == nil {
		return pkg.NewResult(task, errors.New("plugin returned nil session"))
	}
	defer session.Close()

	result := pkg.NewResult(task, nil)
	for _, a := range pipeline {
		ar, err := a.Run(session, task)
		if err != nil {
			logs.Log.Debugf("[%s] action %s failed on %s: %v", task.Service, a.Name(), task.URI(), err)
			continue
		}
		result.Merge(ar)
	}
	if postAction != nil {
		for label, data := range result.Loot {
			result.Extracteds = append(result.Extracteds, postAction.ScanData(data, label)...)
		}
	}
	return result
}

func ExecuteUnauth(task *pkg.Task, plugins map[string]plugin.Plugin, fallback plugin.Plugin, pipeline []pkg.Action, postAction *action.PostAction) *pkg.Result {
	p := resolvePlugin(task.Service, plugins, fallback)
	if p == nil {
		return pkg.NewResult(task, ErrNoPlugin)
	}

	unauth, ok := p.(plugin.UnauthPlugin)
	if !ok {
		return pkg.NewResult(task, pkg.NotImplUnauthorized)
	}
	session, err := unauth.Unauth(task)
	if err != nil {
		return pkg.NewResult(task, err)
	}
	if session == nil {
		return pkg.NewResult(task, ErrNoUnauth)
	}
	defer session.Close()

	result := pkg.NewResult(task, nil)
	for _, a := range pipeline {
		ar, err := a.Run(session, task)
		if err != nil {
			logs.Log.Debugf("[%s] action %s failed on %s: %v", task.Service, a.Name(), task.URI(), err)
		}
		result.Merge(ar)
	}
	if postAction != nil {
		for label, data := range result.Loot {
			result.Extracteds = append(result.Extracteds, postAction.ScanData(data, label)...)
		}
	}
	return result
}

func resolvePlugin(service string, plugins map[string]plugin.Plugin, fallback plugin.Plugin) plugin.Plugin {
	if p, ok := plugins[service]; ok {
		return p
	}
	if s, ok := pkg.Services.Get(service); ok {
		if p, ok := plugins[s.Name]; ok {
			return p
		}
	}
	return fallback
}
