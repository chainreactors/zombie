package core

import (
	"errors"
	"testing"

	"github.com/chainreactors/utils/parsers"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
	mongoplugin "github.com/chainreactors/zombie/plugin/mongo"
)

func TestMongoUnauthUsesRealProbe(t *testing.T) {
	task := &pkg.Task{
		ZombieResult: &parsers.ZombieResult{
			IP:      "127.0.0.1",
			Port:    "27017",
			Service: "mongo",
			Mod:     parsers.ZombieModUnauth,
		},
		Timeout: 2,
	}
	plugins := map[string]plugin.Plugin{"mongo": &mongoplugin.MongoPlugin{}}

	res := ExecuteUnauth(task, plugins, nil, nil, nil)

	if errors.Is(res.Err, pkg.NotImplUnauthorized) {
		t.Fatalf("mongo Unauth is still a stub (NotImplUnauthorized)")
	}
}
