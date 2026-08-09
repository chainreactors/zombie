package plugin_test

import (
	"errors"
	"testing"

	"github.com/chainreactors/zombie/core"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

type sdkPlugin struct{}

func (sdkPlugin) Open(*pkg.Task) (pkg.Session, error) {
	return nil, errors.New("not implemented")
}

func TestRunnerOwnsPluginRegistry(t *testing.T) {
	first := core.NewRunner(core.NewDefaultRunnerOption())
	second := core.NewRunner(core.NewDefaultRunnerOption())
	if _, ok := first.Plugins["redis"]; !ok {
		t.Fatal("built-in redis plugin is missing")
	}

	custom := sdkPlugin{}
	if err := first.RegisterService(plugin.Service{
		Name:        "SDK-Custom",
		Alias:       []string{"SDK-Alias"},
		DefaultPort: "4242",
	}, custom); err != nil {
		t.Fatalf("RegisterService: %v", err)
	}
	if _, ok := first.Plugins["sdk-custom"]; !ok {
		t.Fatal("custom plugin is missing from its runner")
	}
	if _, ok := first.Plugins["sdk-alias"]; !ok {
		t.Fatal("custom plugin alias is missing from its runner")
	}
	if _, ok := second.Plugins["sdk-custom"]; ok {
		t.Fatal("custom plugin leaked into another runner")
	}
	if err := first.RegisterService(plugin.Service{Name: "sdk-custom"}, custom); err == nil {
		t.Fatal("duplicate custom service registration succeeded")
	}
}

func TestUnauthPluginIsOptional(t *testing.T) {
	var p plugin.Plugin = sdkPlugin{}
	if _, ok := p.(plugin.UnauthPlugin); ok {
		t.Fatal("Open-only plugin unexpectedly implements UnauthPlugin")
	}
}
