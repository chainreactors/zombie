package core

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/utils/parsers"
	"github.com/chainreactors/zombie/action"
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
	"github.com/chainreactors/zombie/service"
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

	bar  *pkg.Bar
	stat *pkg.Statistor
	wg   *sync.WaitGroup

	Plugins        map[string]plugin.Plugin
	FallbackPlugin plugin.Plugin
	Pipeline       []pkg.Action
	PostAction     *action.PostAction

	Users    *Generator
	Pwds     *Generator
	Auths    *Generator
	Addrs    utils.Addrs
	Targets  []*Target
	Services []string
	OnResult ResultHandler
	Pool     *ants.PoolWithFunc
	hostSem  *hostLimiter
}

// ResultHandler receives each executed attempt synchronously in its worker.
// Different workers may call the handler concurrently.
type ResultHandler func(*pkg.Result)

func NewRunner(opt *RunnerOption) *Runner {
	if opt == nil {
		opt = NewDefaultRunnerOption()
	}
	return &Runner{
		RunnerOption:   opt,
		Plugins:        defaultPlugins(),
		FallbackPlugin: defaultFallbackPlugin(),
		wg:             &sync.WaitGroup{},
		stat: &pkg.Statistor{
			Tasks: make(map[string]int),
		},
	}
}

// RegisterService adds a plugin to this Runner only.
func (r *Runner) RegisterService(service plugin.Service, p plugin.Plugin) error {
	service.Name = strings.ToLower(strings.TrimSpace(service.Name))
	if service.Name == "" {
		return fmt.Errorf("plugin service name is required")
	}
	if p == nil {
		return fmt.Errorf("plugin for service %q is nil", service.Name)
	}
	if _, exists := r.Plugins[service.Name]; exists {
		return fmt.Errorf("plugin service %q is already registered", service.Name)
	}
	if service.Source == "" {
		service.Source = pkg.PluginSource
	}

	aliases := make([]string, 0, len(service.Alias))
	seen := map[string]struct{}{service.Name: {}}
	for _, alias := range service.Alias {
		alias = strings.ToLower(strings.TrimSpace(alias))
		if alias == "" {
			continue
		}
		if _, duplicate := seen[alias]; duplicate {
			continue
		}
		if _, exists := r.Plugins[alias]; exists {
			return fmt.Errorf("plugin service alias %q is already registered", alias)
		}
		seen[alias] = struct{}{}
		aliases = append(aliases, alias)
	}
	service.Alias = aliases
	r.Plugins[service.Name] = p
	for _, alias := range aliases {
		r.Plugins[alias] = p
	}
	pkg.Services.Register(&service)
	return nil
}

func (r *Runner) BuildPipeline() error {
	if r.Proton {
		if len(r.ScanTemplates) == 0 {
			return fmt.Errorf("--proton requires --scan-template to specify proton template path")
		}
		pa, err := action.NewPostAction(r.ScanTemplates)
		if err != nil {
			return fmt.Errorf("failed to init post action: %w", err)
		}
		r.PostAction = pa
	}

	var serviceTemplates []*service.Template
	if r.Gather {
		embedded, err := action.LoadServiceTemplatesFromData(pkg.ServiceTemplateData)
		if err != nil {
			return fmt.Errorf("failed to load embedded service templates: %w", err)
		}
		serviceTemplates = append(serviceTemplates, embedded...)
	}
	if len(r.ServiceTemplates) > 0 {
		fromPaths, err := action.LoadServiceTemplatesFromPaths(r.ServiceTemplates)
		if err != nil {
			return fmt.Errorf("failed to load service templates: %w", err)
		}
		serviceTemplates = append(serviceTemplates, fromPaths...)
	}
	if len(serviceTemplates) > 0 {
		serviceAction, err := action.NewServiceAction(serviceTemplates, r.ServiceVars, r.ServicePayloads)
		if err != nil {
			return fmt.Errorf("failed to init service action: %w", err)
		}
		if r.Risk != "" {
			serviceAction.SetRisk(r.Risk)
		}
		if r.Gather && len(r.Tags) == 0 {
			serviceAction.SetTags([]string{"gather"})
		} else if len(r.Tags) > 0 {
			serviceAction.SetTags(r.Tags)
		}
		r.Pipeline = append(r.Pipeline, serviceAction)
	}

	if r.Gather {
		r.Pipeline = append(r.Pipeline, action.NewAuditAction())
	}

	if r.Gather && r.PostAction == nil {
		if data := pkg.LootTemplateData; len(data) > 0 {
			pa, err := action.NewPostActionFromData(data)
			if err != nil {
				logs.Log.Debugf("loot scanner disabled: %v", err)
			} else {
				r.PostAction = pa
			}
		}
	}

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
	if r.Threads <= 0 {
		return fmt.Errorf("threads must be greater than zero")
	}
	if r.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}
	if r.Concurrency < 0 {
		return fmt.Errorf("concurrency must not be negative")
	}
	if r.Top < 0 {
		return fmt.Errorf("top must not be negative")
	}

	r.hostSem = newHostLimiter(r.Concurrency)
	pool, err := ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		task := i.(*pkg.Task)
		defer func() {
			r.wg.Done()
			if task.Completed != nil {
				close(task.Completed)
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
		taskCtx, cancel := context.WithTimeout(task.Context, task.Duration())
		defer cancel()
		task.Context = taskCtx

		var res *pkg.Result
		if task.Mod == parsers.ZombieModUnauth {
			res = ExecuteUnauth(task, r.Plugins, r.FallbackPlugin, r.Pipeline, r.PostAction)
		} else {
			res = Execute(task, r.Plugins, r.FallbackPlugin, r.Pipeline, r.PostAction)
		}
		r.Output(res)

		if res.OK && r.FirstOnly && task.Mod != parsers.ZombieModSniper {
			task.Cancel()
		}
	}, ants.WithPanicHandler(func(err interface{}) {
		debug.PrintStack()
	}))
	if err != nil {
		return fmt.Errorf("create worker pool: %w", err)
	}
	r.Pool = pool
	defer pool.Release()

	ch := r.targetGenerate(ctx)
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
			Context: targetCtx,
			Cancel:  cancel,
			Timeout: r.Timeout,
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
				Context: targetCtx,
				Cancel:  cancel,
				Timeout: r.Timeout,
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
				completed := make(chan struct{})
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
					Context:   targetCtx,
					Cancel:    cancel,
					Timeout:   r.Timeout,
					Completed: completed,
				})
				<-completed
			}

			for task := range r.clusterBombGenerate(targetCtx, cancel, cur) {
				r.add(task)
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
	genLoop:
		for _, user := range users {
			select {
			case <-ctx.Done():
				break genLoop
			default:
			}
			wg.Add(1)
			usr := user
			go func() {
				defer wg.Done()
				if !r.NoUnAuth {
					userCompleted := make(chan struct{})
					select {
					case ch <- &pkg.Task{
						ZombieResult: &parsers.ZombieResult{
							IP:       target.IP,
							Port:     target.Port,
							Service:  target.Service,
							Scheme:   target.Scheme,
							Username: usr,
							Param:    target.Param,
							Mod:      parsers.ZombieModUnauth,
						},
						Timeout:   r.Timeout,
						Context:   ctx,
						Cancel:    canceler,
						Completed: userCompleted,
					}:
					case <-ctx.Done():
						return
					}
					<-userCompleted
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
						Timeout: r.Timeout,
						Context: ctx,
						Cancel:  canceler,
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

func (r *Runner) targetGenerate(ctx context.Context) chan *Target {
	ch := make(chan *Target)
	go func() {
		defer close(ch)
		for _, target := range r.Targets {
			if r.Services == nil || (r.Services != nil && iutils.StringsContains(r.Services, target.Service)) {
				select {
				case ch <- target:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch
}

func (r *Runner) add(task *pkg.Task) {
	task.ProxyDial = r.ProxyDial
	task.Raw = r.Raw
	r.wg.Add(1)
	r.stat.RecordTask(task.Service, task.String())
	_ = r.Pool.Invoke(task)
}

func (r *Runner) Output(res *pkg.Result) {
	r.stat.RecordResult(res)
	if r.OnResult != nil {
		r.OnResult(res)
	}
}
