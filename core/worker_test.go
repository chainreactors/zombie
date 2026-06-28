package core

import (
	"testing"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

type aliasPluginSession struct{}

func (aliasPluginSession) Service() string  { return "mongo" }
func (aliasPluginSession) Close() error     { return nil }
func (aliasPluginSession) Raw() interface{} { return nil }

type aliasPlugin struct{}

func (aliasPlugin) Name() string                          { return "mongo" }
func (aliasPlugin) Open(*pkg.Task) (pkg.Session, error)   { return aliasPluginSession{}, nil }
func (aliasPlugin) Unauth(*pkg.Task) (pkg.Session, error) { return aliasPluginSession{}, nil }

func TestResolvePluginAcceptsServiceAliases(t *testing.T) {
	plugins := map[string]plugin.Plugin{"mongo": aliasPlugin{}}

	if p := resolvePlugin("mongodb", plugins); p == nil || p.Name() != "mongo" {
		t.Fatalf("resolvePlugin(mongodb) = %#v, want mongo plugin", p)
	}
}
