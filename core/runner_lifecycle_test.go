package core

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/chainreactors/utils/parsers"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
)

type contextProbePlugin struct {
	canceled chan struct{}
	release  chan struct{}
}

type failedAttemptPlugin struct{}

func (failedAttemptPlugin) Open(*pkg.Task) (pkg.Session, error) {
	return nil, errors.New("authentication failed")
}

func (failedAttemptPlugin) Unauth(*pkg.Task) (pkg.Session, error) {
	return nil, errors.New("unauthentication failed")
}

func TestRunnerEmitsEveryAttemptAndCompletesTask(t *testing.T) {
	opt := NewDefaultRunnerOption()
	opt.Mod = ModBomb
	opt.Threads = 2
	opt.Timeout = 1
	opt.Quiet = true
	opt.NoCheckHoneyPot = true
	opt.FirstOnly = false

	runner := NewRunner(opt)
	runner.Plugins = map[string]plugin.Plugin{"failed-attempt": failedAttemptPlugin{}}
	runner.SetUsers([]string{"root"})
	runner.SetPasswords([]string{"wrong"})
	runner.SetTargets([]*Target{{IP: "127.0.0.1", Port: "1", Service: "failed-attempt"}})

	var mu sync.Mutex
	var results []*pkg.Result
	var unauthCompleted chan struct{}
	completedDuringHandler := false
	runner.OnResult = func(result *pkg.Result) {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, result)
		if result.Mod == parsers.ZombieModUnauth {
			unauthCompleted = result.Completed
			select {
			case <-result.Completed:
				completedDuringHandler = true
			default:
			}
		}
	}

	if err := runner.RunWithContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d, want unauth + brute attempts", len(results))
	}
	for _, result := range results {
		if result.OK || result.Err == nil || result.ErrString == "" {
			t.Fatalf("incomplete failed result: %#v", result)
		}
	}
	if completedDuringHandler {
		t.Fatal("task completed before its result handler returned")
	}
	select {
	case <-unauthCompleted:
	default:
		t.Fatal("task completion was not signaled after execution")
	}
}

func (p *contextProbePlugin) Open(task *pkg.Task) (pkg.Session, error) {
	select {
	case <-task.Context.Done():
		close(p.canceled)
		return nil, task.Context.Err()
	case <-p.release:
		return nil, errors.New("released")
	}
}

func (*contextProbePlugin) Unauth(*pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}

func TestRunnerTaskTimeoutReachesPlugin(t *testing.T) {
	probe := &contextProbePlugin{
		canceled: make(chan struct{}),
		release:  make(chan struct{}),
	}
	t.Cleanup(func() { close(probe.release) })

	opt := NewDefaultRunnerOption()
	opt.Mod = ModSniper
	opt.Threads = 1
	opt.Timeout = 1
	opt.Quiet = true
	r := NewRunner(opt)
	r.Plugins = map[string]plugin.Plugin{"context-probe": probe}
	r.SetTargets([]*Target{{IP: "127.0.0.1", Port: "1", Service: "context-probe"}})

	started := time.Now()
	if err := r.RunWithContext(context.Background()); err != nil {
		t.Fatalf("RunWithContext: %v", err)
	}
	if elapsed := time.Since(started); elapsed > 1500*time.Millisecond {
		t.Fatalf("task timeout took %s, want close to configured 1s", elapsed)
	}

	select {
	case <-probe.canceled:
	case <-time.After(250 * time.Millisecond):
		t.Fatal("plugin did not observe task context cancellation")
	}
}

func TestRunnerDoesNotEmitResultsOnValidationError(t *testing.T) {
	opt := NewDefaultRunnerOption()
	opt.Threads = 0
	r := NewRunner(opt)
	called := false
	r.OnResult = func(*pkg.Result) { called = true }

	if err := r.RunWithContext(context.Background()); err == nil {
		t.Fatal("expected invalid thread count to fail")
	}
	if called {
		t.Fatal("result handler called before execution")
	}
}
