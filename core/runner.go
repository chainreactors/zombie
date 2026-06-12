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
	"github.com/chainreactors/zombie/action"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
	"github.com/panjf2000/ants/v2"
)

var (
	ModSniper    = "sniper"
	ModBomb      = "clusterbomb"
	ModPitchFork = "pitchfork"
)

// hostLimiter 给每个 host 一个并发闸,把单 host 在飞连接数限制在服务端
// 限速安全区内(如 sshd 默认 MaxStartups 10:30:100),从源头避免连接被拒。
type hostLimiter struct {
	mu    sync.Mutex
	sems  map[string]chan struct{}
	limit int
}

func newHostLimiter(limit int) *hostLimiter {
	return &hostLimiter{sems: make(map[string]chan struct{}), limit: limit}
}

// acquire 阻塞直到该 host 有空闲额度,返回 (释放函数, 是否取得)。limit<=0 时不限。
// 等额度期间若 ctx 被取消(如该目标已命中口令),立即返回 ok=false 且不建连,
// 否则被取消的目标会被排空队列高速循环建连,瞬间冲破服务端限速。
func (h *hostLimiter) acquire(ctx context.Context, key string) (func(), bool) {
	if h == nil || h.limit <= 0 {
		return func() {}, true
	}
	h.mu.Lock()
	sem, ok := h.sems[key]
	if !ok {
		sem = make(chan struct{}, h.limit)
		h.sems[key] = sem
	}
	h.mu.Unlock()
	select {
	case sem <- struct{}{}:
		return func() { <-sem }, true
	case <-ctx.Done():
		return func() {}, false
	}
}

type Runner struct {
	*RunnerOption

	bar     *pkg.Bar
	stat    *pkg.Statistor
	wg      *sync.WaitGroup
	outlock *sync.WaitGroup
	addlock *sync.Mutex

	Plugins  map[string]plugin.Plugin
	Pipeline []pkg.Action

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
	hostSem      *hostLimiter
}

func NewRunner(opt *RunnerOption) *Runner {
	if opt == nil {
		opt = NewDefaultRunnerOption()
	}
	return &Runner{
		RunnerOption: opt,
		Plugins:      plugin.DefaultRegistry(),
		OutputCh:     make(chan *pkg.Result),
		wg:           &sync.WaitGroup{},
		outlock:      &sync.WaitGroup{},
		addlock:      &sync.Mutex{},
		stat: &pkg.Statistor{
			Tasks: make(map[string]int),
		},
	}
}

func (r *Runner) BuildPipeline() error {
	if !r.Proton {
		return nil
	}
	if len(r.ScanTemplates) == 0 {
		return fmt.Errorf("--proton requires --scan-template to specify proton template path")
	}
	postAction, err := action.NewPostAction(r.ScanTemplates, r.DBLimit)
	if err != nil {
		return fmt.Errorf("failed to init post action: %w", err)
	}
	r.Pipeline = append(r.Pipeline, postAction)
	return nil
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

	r.hostSem = newHostLimiter(r.Concurrency)
	r.Pool, _ = ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		task := i.(*pkg.Task)
		defer func() {
			r.wg.Done()
			if task.Locker != nil {
				task.Locker.Unlock()
			}
		}()
		// 该目标已命中/被取消,无需再发起连接,直接跳过。避免 first-success 之后
		// 队列里剩余任务被排空时,每个都瞬间建连又立刻取消,冲破服务端限速。
		select {
		case <-task.Context.Done():
			return
		default:
		}
		// per-host 闸放在外层、且在超时计时之前:既能严格把单 host 在飞连接限制
		// 在服务端安全区(如 sshd MaxStartups 10),排队等额度的时间又不计入
		// 任务超时(否则重竞争下会误判 goroutine timeout)。外层 select 会等
		// ctx.Done(inner 跑完后 tcancel),所以 sem 一直持有到连接结束。
		releaseHost, ok := r.hostSem.acquire(task.Context, task.Address())
		if !ok {
			return // 等额度期间目标已命中/取消,不再建连
		}
		defer releaseHost()
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
				res = ExecuteUnauth(task, r.Plugins, r.Pipeline)
			} else {
				res = Execute(task, r.Plugins, r.Pipeline)
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
	defer r.Pool.Release()

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

	select {
	case <-ctx.Done():
		if !r.Quiet {
			logs.Log.Warnf("interrupted, printing partial results")
		}
	default:
	}

	if !r.Quiet {
		logs.Log.Importantf("%s", r.stat.TaskString())
		logs.Log.Importantf("%s", r.stat.SummaryString())
	}

	if r.File != nil {
		r.File.Close()
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
	r.stat.RecordResult(res)
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
				errMsg := "unknown error"
				if result.Err != nil {
					errMsg = result.Err.Error()
				}
				logs.Log.Debugf("[%s] %s %s %s ,%s login failed, %s", result.Mod.String(), result.URI(), result.Username, result.Password, result.Service, errMsg)
			}
			r.outlock.Done()
		}
	}
}
