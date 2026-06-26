package plugin

import "github.com/chainreactors/zombie/pkg"

type Plugin = pkg.Plugin

func Register(name string, p Plugin) { pkg.RegisterPlugin(name, p) }
func Get(service string) (Plugin, bool) { return pkg.GetPlugin(service) }

func DefaultRegistry() map[string]Plugin {
	return pkg.DefaultPluginRegistry()
}
