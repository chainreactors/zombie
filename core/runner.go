package core

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/fileutils"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/zombie/pkg"
	"github.com/panjf2000/ants/v2"
)

var (
	ModSniper    = "sniper"
	ModBomb      = "clusterbomb"
	ModPitchFork = "pitchfork"
)

type Runner struct {
	*RunnerOption

	bar     *pkg.Bar
	stat    *pkg.Statistor
	wg      *sync.WaitGroup
	outlock *sync.WaitGroup
	addlock *sync.Mutex

	Users        *Generator
	Pwds         *Generator
	Auths        *Generator
	Addrs        utils.Addrs
	Targets      []*Target
	Services     []string
	OutputCh     chan *pkg.Result
	File         *fileutils.File
	OutFunc      func(string)
	FileFormat   string
	OutputFormat string
	Pool         *ants.PoolWithFunc
}

func NewRunner(opt *RunnerOption) *Runner {
	if opt == nil {
		opt = NewDefaultRunnerOption()
	}
	return &Runner{
		RunnerOption: opt,
		OutputCh:     make(chan *pkg.Result),
		wg:           &sync.WaitGroup{},
		outlock:      &sync.WaitGroup{},
		addlock:      &sync.Mutex{},
		stat: &pkg.Statistor{
			Tasks: make(map[string]int),
		},
	}
}

func (r *Runner) SetTargets(targets []*Target) {
	r.Targets = targets
}

func (r *Runner) SetUsers(users []string) {
	if len(users) > 0 {
		r.Users = NewGeneratorWithInput(users)
	}
}

func (r *Runner) SetPasswords(passwords []string) {
	if len(passwords) > 0 {
		r.Pwds = NewGeneratorWithInput(passwords)
	}
}

func (r *Runner) SetAuths(pairs []string) {
	if len(pairs) > 0 {
		r.Auths = NewGeneratorWithInput(pairs)
		r.Mod = ModPitchFork
	}
}

func (r *Runner) Run() {
	_ = r.RunWithContext(context.Background())
}

func (r *Runner) RunWithContext(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	pkg.RunOpt.Raw = r.Raw

	if r.Mod == "" {
		r.Mod = ModBomb
	}
	switch r.Mod {
	case ModSniper, ModBomb, ModPitchFork:
	default:
		return fmt.Errorf("unsupported mod %q, want clusterbomb, pitchfork, or sniper", r.Mod)
	}
	if r.Mod == ModPitchFork && r.Auths == nil {
		return fmt.Errorf("pitchfork mode requires auth, please set -a/-A")
	}

	if r.OutFunc != nil {
		go r.OutputHandler()
	}

	r.Pool, _ = ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		task := i.(*pkg.Task)
		defer func() {
			r.wg.Done()
			if task.Locker != nil {
				task.Locker.Unlock()
			}
		}()
		ctx, tcancel := context.WithCancel(task.Context)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logs.Log.Debugf("%s panic: %v", task.String(), r)
					tcancel()
				}
			}()
			var res *pkg.Result
			if task.Mod == parsers.ZombieModUnauth {
				res = Unauth(task)
			} else if task.Mod == parsers.ZombieModCheck {
				res = Brute(task)
			} else {
				res = Brute(task)
			}

			select {
			case <-ctx.Done():
				return
			case <-task.Context.Done():
				return
			default:
				r.Output(res)
			}

			if res.OK && r.FirstOnly && task.Mod != parsers.ZombieModSniper {
				tcancel()
				task.Canceler()
			}
			tcancel()
		}()

		select {
		case <-ctx.Done():
		case <-task.Context.Done():
			logs.Log.Debugf("all task %s cancel", task.URI())
		case <-time.After(time.Duration(task.Timeout*2) * time.Second):
			tcancel()
			r.Output(&pkg.Result{
				Task: task,
				Err:  fmt.Errorf("goroutine timeout, force cancel"),
			})
		}
	}, ants.WithPanicHandler(func(err interface{}) {
		debug.PrintStack()
		r.wg.Done()
	}))

	ch := r.targetGenerate()
	switch r.Mod {
	case ModSniper:
		r.RunWithSniper(ctx, ch)
	case ModBomb:
		r.RunWithClusterBomb(ctx, ch)
	case ModPitchFork:
		r.RunWithPitchfork(ctx, ch)
	default:
		return nil
	}
	if r.OutFunc != nil {
		r.outlock.Wait()
	}
	close(r.OutputCh)

	if !r.Quiet {
		logs.Log.Importantf("%s", r.stat.TaskString())
		logs.Log.Importantf("total: %d, success: %d", r.stat.Total, r.stat.Success)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (r *Runner) Stat() *pkg.Statistor {
	return r.stat
}

func (r *Runner) RunWithSniper(ctx context.Context, targets chan *Target) {
	for target := range targets {
		select {
		case <-ctx.Done():
			r.wg.Wait()
			return
		default:
		}
		targetCtx, cancel := context.WithCancel(ctx)
		r.add(&pkg.Task{
			ZombieResult: &parsers.ZombieResult{
				IP:       target.IP,
				Port:     target.Port,
				Service:  target.Service,
				Scheme:   target.Scheme,
				Username: target.Username,
				Password: target.Password,
				Param:    target.Param,
				Mod:      parsers.ZombieModSniper,
			},
			Context:  targetCtx,
			Canceler: cancel,
			Timeout:  r.Timeout,
		})
	}
	r.wg.Wait()
}

func (r *Runner) RunWithPitchfork(ctx context.Context, target chan *Target) {
	if r.Auths == nil {
		return
	}
	var pairs [][]string
	for _, auth := range r.Auths.RunAsSlice() {
		username, password := parseAuthPair(auth)
		pairs = append(pairs, []string{username, password})
	}

	for target := range target {
		select {
		case <-ctx.Done():
			r.wg.Wait()
			return
		default:
		}
		targetCtx, cancel := context.WithCancel(ctx)
		defer cancel()
	pairLoop:
		for _, pair := range pairs {
			select {
			case <-targetCtx.Done():
				break pairLoop
			default:
			}
			r.add(&pkg.Task{
				ZombieResult: &parsers.ZombieResult{
					IP:       target.IP,
					Port:     target.Port,
					Service:  target.Service,
					Scheme:   target.Scheme,
					Username: pair[0],
					Password: pair[1],
					Param:    target.Param,
					Mod:      parsers.ZombieModPitchfork,
				},
				Context:  targetCtx,
				Canceler: cancel,
				Timeout:  r.Timeout,
			})
		}
	}
	r.wg.Wait()
}

func (r *Runner) RunWithClusterBomb(ctx context.Context, targets chan *Target) {
	targetWG := &sync.WaitGroup{}
	for target := range targets {
		select {
		case <-ctx.Done():
			targetWG.Wait()
			r.wg.Wait()
			return
		default:
		}
		targetWG.Add(1)
		targetCtx, cancel := context.WithCancel(ctx)
		cur := target

		go func() {
			defer targetWG.Done()
			if r.Strict {
				if open := cur.CheckOpen(); !open {
					cancel()
					return
				}
				if matched := cur.CheckFinger(); !matched {
					cancel()
					return
				}
			}

			if !r.NoCheckHoneyPot {
				locker := &sync.Mutex{}
				locker.Lock()
				r.add(&pkg.Task{
					ZombieResult: &parsers.ZombieResult{
						IP:       cur.IP,
						Port:     cur.Port,
						Service:  cur.Service,
						Scheme:   cur.Scheme,
						Param:    cur.Param,
						Username: randomString(10),
						Password: randomString(10),
						Mod:      parsers.ZombieModCheck,
					},
					Context:  targetCtx,
					Canceler: cancel,
					Timeout:  r.Timeout,
					Locker:   locker,
				})
				locker.Lock()
				locker.Unlock()
			}

			ch := r.clusterBombGenerate(targetCtx, cancel, cur)
		loop:
			for {
				select {
				case task, ok := <-ch:
					if ok {
						r.add(task)
					} else {
						break loop
					}
				case <-targetCtx.Done():
					break loop
				}
			}
		}()
	}
	targetWG.Wait()
	r.wg.Wait()
}

func (r *Runner) clusterBombGenerate(ctx context.Context, canceler context.CancelFunc, target *Target) chan *pkg.Task {
	ch := make(chan *pkg.Task)
	var users, pwds []string
	if target.Username != "" {
		users = []string{target.Username}
	} else if r.Users == nil {
		users = pkg.UseDefaultUser(target.Service, r.Top)
	} else {
		users = r.Users.RunAsSlice()
	}

	if target.Password != "" {
		pwds = []string{target.Password}
	} else if r.Pwds == nil {
		pwds = pkg.UseDefaultPassword(target.Service, r.Top)
	} else {
		pwds = r.Pwds.RunAsSlice()
	}
	wg := &sync.WaitGroup{}

	go func() {
		defer close(ch)
		for _, user := range users {
			select {
			case <-ctx.Done():
				return
			default:
			}
			wg.Add(1)
			usr := user
			go func() {
				defer wg.Done()
				if !r.NoUnAuth {
					userLocker := &sync.Mutex{}
					userLocker.Lock()
					ch <- &pkg.Task{
						ZombieResult: &parsers.ZombieResult{
							IP:       target.IP,
							Port:     target.Port,
							Service:  target.Service,
							Scheme:   target.Scheme,
							Username: usr,
							Param:    target.Param,
							Mod:      parsers.ZombieModUnauth,
						},
						Timeout:  r.Timeout,
						Context:  ctx,
						Canceler: canceler,
						Locker:   userLocker,
					}
					userLocker.Lock()
					userLocker.Unlock()
				}

				for _, pwd := range pwds {
					if target.Service == "" {
						logs.Log.Warn("unknown service " + target.Service)
						continue
					}
					select {
					case ch <- &pkg.Task{
						ZombieResult: &parsers.ZombieResult{
							IP:       target.IP,
							Port:     target.Port,
							Service:  target.Service,
							Scheme:   target.Scheme,
							Username: usr,
							Password: pwd,
							Param:    target.Param,
							Mod:      parsers.ZombieModBrute,
						},
						Timeout:  r.Timeout,
						Context:  ctx,
						Canceler: canceler,
					}:
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		wg.Wait()
		if r.bar != nil {
			r.bar.Done()
		}
	}()
	return ch
}

func (r *Runner) targetGenerate() chan *Target {
	ch := make(chan *Target)
	go func() {
		for _, target := range r.Targets {
			if r.Services == nil || (r.Services != nil && iutils.StringsContains(r.Services, target.Service)) {
				ch <- target
			}
		}
		close(ch)
	}()

	return ch
}

func (r *Runner) add(task *pkg.Task) {
	task.ProxyDial = r.ProxyDial
	r.stat.Cur = task.String()
	r.addlock.Lock()
	r.stat.Tasks[task.Service]++
	r.wg.Add(1)
	r.stat.Total++
	r.addlock.Unlock()
	_ = r.Pool.Invoke(task)
}

func (r *Runner) Output(res *pkg.Result) {
	if r.OutFunc != nil {
		r.outlock.Add(1)
	}
	if res.OK {
		r.stat.Success++
	}
	r.OutputCh <- res
}

func (r *Runner) OutputHandler() {
loop:
	for {
		select {
		case result, ok := <-r.OutputCh:
			if !ok {
				break loop
			}
			if result.OK {
				if r.File != nil {
					r.OutFunc(result.Format(r.FileFormat))
				}
				logs.Log.Console(result.Format(r.OutputFormat))
			} else {
				logs.Log.Debugf("[%s] %s %s %s ,%s login failed, %s", result.Mod.String(), result.URI(), result.Username, result.Password, result.Service, result.Err.Error())
			}
			r.outlock.Done()
		}
	}
}
