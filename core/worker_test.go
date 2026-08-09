package core

import (
	"errors"
	"testing"

	"github.com/chainreactors/utils/parsers"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

type aliasPluginSession struct{}

func (aliasPluginSession) Service() string { return "mongo" }
func (aliasPluginSession) Close() error    { return nil }

type aliasPlugin struct{}

func (aliasPlugin) Open(*pkg.Task) (pkg.Session, error)   { return aliasPluginSession{}, nil }
func (aliasPlugin) Unauth(*pkg.Task) (pkg.Session, error) { return aliasPluginSession{}, nil }

func TestResolvePluginAcceptsServiceAliases(t *testing.T) {
	plugins := map[string]plugin.Plugin{"mongo": aliasPlugin{}}

	if p := resolvePlugin("mongodb", plugins, nil); p == nil {
		t.Fatalf("resolvePlugin(mongodb) = %#v, want mongo plugin", p)
	}
}

type openOnlyPlugin struct{}

func (openOnlyPlugin) Open(*pkg.Task) (pkg.Session, error) {
	return nil, errors.New("not used")
}

func TestExecuteUnauthRejectsOpenOnlyPlugin(t *testing.T) {
	task := &pkg.Task{ZombieResult: &parsers.ZombieResult{Service: "open-only"}}
	result := ExecuteUnauth(task, map[string]plugin.Plugin{"open-only": openOnlyPlugin{}}, nil, nil, nil)
	if !errors.Is(result.Err, pkg.NotImplUnauthorized) {
		t.Fatalf("error = %v, want NotImplUnauthorized", result.Err)
	}
	if result.OK || result.ErrString == "" {
		t.Fatalf("public result did not preserve unauth failure: %#v", result.ZombieResult)
	}
}
