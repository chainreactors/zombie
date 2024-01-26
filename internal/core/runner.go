package core

import (
	"context"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/zombie/pkg"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

var (
	ModSniper = "sniper"
	ModBomb   = "clusterbomb"
)

type Runner struct {
	*Option
	wg        *sync.WaitGroup
	Stat      *pkg.Statistor
	Users     *Generator
	Pwds      *Generator
	Addrs     utils.Addrs
	Targets   []*Target
	Services  []string
	OutputCh  chan *pkg.Result
	outlock   *sync.WaitGroup
	File      *files.File
	OutFunc   func(string)
	FirstOnly bool
	Pool      *ants.PoolWithFunc
}

func (r *Runner) Run() {
	go r.OutputHandler()
	r.Pool, _ = ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		defer r.wg.Done()
		task := i.(*pkg.Task)
		ctx, tcancel := context.WithCancel(task.Context) // current task context
		go func() {
			var res *pkg.Result
			if task.Mod == pkg.TaskModUnauth {
				res = Unauth(task)
				task.Locker.Unlock()
			} else if task.Mod == pkg.TaskModCheck {
				res = Brute(task)
				task.Locker.Unlock()
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

			if res.OK && r.FirstOnly && task.Mod != pkg.TaskModSniper {
				tcancel()       // 退出当前任务
				task.Canceler() // 取消正在执行的所有任务
			}
			tcancel()
		}()

		// 设置超时时间, 防止任务挂死
		select {
		case <-ctx.Done():
			//logs.Log.Debugf("current task %s %s %s cancel", task.URI(), task.Username, task.Password)
		case <-task.Context.Done():
			logs.Log.Debugf("all task %s cancel", task.URI())
		case <-time.After(time.Duration(task.Timeout+10) * time.Second):
			tcancel()
			r.Output(&pkg.Result{
				Task: task,
				Err:  fmt.Errorf("goroutine timeout, force cancel"),
			})
		}
	})

	ch := r.targetGenerate()
	switch r.Mod {
	case ModSniper:
		r.RunWithSniper(ch)
	case ModBomb:
		r.RunWithClusterBomb(ch)
	}
	r.outlock.Wait()
}

func (r *Runner) RunWithSniper(targets chan *Target) {
	for target := range targets {
		r.add(&pkg.Task{
			IP:       target.IP,
			Port:     target.Port,
			Service:  target.Service,
			Username: target.Username,
			Password: target.Password,
			Param:    target.Param,
			Context:  context.Background(),
			Timeout:  r.Timeout,
			Mod:      pkg.TaskModSniper,
		})
	}
	r.wg.Wait()
}

func (r *Runner) RunWithClusterBomb(targets chan *Target) {
	targetWG := &sync.WaitGroup{}
	// unauth
	for target := range targets {
		// check honey pot
		targetWG.Add(1)
		targetCtx, cancel := context.WithCancel(context.Background())
		cur := target
		go func() {
			defer targetWG.Done()

			if !r.NoCheckHoneyPot {
				locker := &sync.Mutex{}
				locker.Lock()
				r.add(&pkg.Task{
					IP:       cur.IP,
					Port:     cur.Port,
					Service:  cur.Service,
					Context:  targetCtx,
					Canceler: cancel,
					Timeout:  r.Timeout,
					Username: randomString(10),
					Password: randomString(10),
					Mod:      pkg.TaskModCheck,
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
					// 从生成器中取任务.
					if ok {
						r.add(task)
					} else {
						break loop
					}
				case <-targetCtx.Done():
					// todo 为断点续传做准备
					break loop
				}
			}
		}()
	}
	targetWG.Wait()
	r.wg.Wait()
}

func (r *Runner) clusterBombGenerate(ctx context.Context, canceler context.CancelFunc, target *Target) chan *pkg.Task {
	// 通过用户名与密码的笛卡尔积生成数据
	ch := make(chan *pkg.Task)
	var users, pwds []string
	// 自动选择默认的用户名与密码字典
	if target.Username != "" {
		users = []string{target.Username}
	} else if r.Users == nil {
		users = pkg.UseDefaultUser(target.Service.String())
	} else {
		users = r.Users.RunAsSlice()
	}

	if target.Password != "" {
		pwds = []string{target.Password}
	} else if r.Pwds == nil {
		pwds = pkg.UseDefaultPassword(target.Service.String(), r.Top)
	} else {
		pwds = r.Pwds.RunAsSlice()
	}
	wg := &sync.WaitGroup{}

	// task生成器
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
						IP:       target.IP,
						Port:     target.Port,
						Service:  target.Service,
						Username: usr,
						Param:    target.Param,
						Timeout:  r.Timeout,
						Mod:      pkg.TaskModUnauth,
						Context:  ctx,
						Canceler: canceler,
						Locker:   userLocker,
					}
					userLocker.Lock()
					userLocker.Unlock()
				}

				for _, pwd := range pwds {
					if target.Service == "" {
						logs.Log.Warn("unknown service " + target.Service.String())
						continue
					}
					select {
					case ch <- &pkg.Task{
						IP:       target.IP,
						Port:     target.Port,
						Service:  target.Service,
						Username: usr,
						Password: pwd,
						Param:    target.Param,
						Timeout:  r.Timeout,
						Mod:      pkg.TaskModBrute,
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
	}()
	return ch
}

func (r *Runner) targetGenerate() chan *Target {
	ch := make(chan *Target)
	go func() {
		// 通过targets生成目标
		for _, target := range r.Targets {
			if r.Services == nil || (r.Services != nil && iutils.StringsContains(r.Services, target.Service.String())) {
				// 如果从gogo中输入的目标, 可以通过-s过滤特定的服务进行扫描
				ch <- target
			}
		}
		close(ch)
	}()

	return ch
}

func (r *Runner) add(task *pkg.Task) {
	r.wg.Add(1)
	r.Stat.Count++
	_ = r.Pool.Invoke(task)
}

func (r *Runner) Output(res *pkg.Result) {
	r.outlock.Add(1)
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
					r.OutFunc(result.Format(r.Option.FileFormat))
				}
				logs.Log.Console(result.Format(r.Option.OutputFormat))
			} else {
				logs.Log.Debugf("[%s] %s %s %s failed, %s", result.Mod.String(), result.URI(), result.Username, result.Password, result.Err.Error())
			}
			r.outlock.Done()
		}
	}
}
