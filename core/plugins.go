package core

import (
	internalplugins "github.com/chainreactors/zombie/internal/plugins"
	"github.com/chainreactors/zombie/plugin"
)

func defaultPlugins() map[string]plugin.Plugin {
	return internalplugins.Default()
}

func defaultFallbackPlugin() plugin.Plugin {
	return internalplugins.Fallback()
}

func registerBuiltinServices() {
	internalplugins.RegisterServices()
}
