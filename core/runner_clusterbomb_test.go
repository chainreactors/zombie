package core

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

type clusterBombSession struct{}

func (clusterBombSession) Service() string  { return "faketest" }
func (clusterBombSession) Close() error     { return nil }
func (clusterBombSession) Raw() interface{} { return nil }

type clusterBombPlugin struct{}

func (clusterBombPlugin) Name() string                        { return "faketest" }
func (clusterBombPlugin) Open(*pkg.Task) (pkg.Session, error) { return nil, errors.New("auth failed") }
func (clusterBombPlugin) Unauth(*pkg.Task) (pkg.Session, error) {
	return clusterBombSession{}, nil
}

func TestClusterBombStopsSendersBeforeClosingTaskChannel(t *testing.T) {
	r := NewRunner(NewDefaultRunnerOption())
	r.Plugins = map[string]plugin.Plugin{"faketest": clusterBombPlugin{}}
	r.Quiet = true

	drained := make(chan struct{})
	go func() {
		for range r.OutputCh {
		}
		close(drained)
	}()

	const nUsers = 400
	users := make([]string, nUsers)
	for i := range users {
		users[i] = fmt.Sprintf("u%d", i)
	}
	r.SetUsers(users)
	r.SetPasswords([]string{"x"})
	r.SetTargets([]*Target{{IP: "127.0.0.1", Port: "1", Service: "faketest"}})

	_ = r.RunWithContext(context.Background())
	<-drained
}
