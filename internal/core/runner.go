package core

import (
	"context"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers/iutils"
	"github.com/chainreactors/zombie/pkg"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

type Runner struct {
	*Option
	wg         *sync.WaitGroup
	Stat       *pkg.Statistor
	Users      *Generator
	Pwds       *Generator
	Addrs      ipcs.Addrs
	Targets    []*Target
	Services   []string
	OutputCh   chan *pkg.Result
	File       *files.File
	OutFunc    func(string)
	ExecString string
	FirstOnly  bool
}

func (r *Runner) Run() {
	go r.Output()
	ch := r.targetGenerate()
	switch r.Mod {
	case "pitchfork":
	case "clusterbomb":
		r.RunWithClusterBomb(ch)
	}
}

func (r *Runner) RunWithPitchfork() {
	//rootContext, rootCancel := context.WithCancel(context.Background())
	//for _, addr := range r.Addrs{
	//	 for _, user := range r.Users{}
	//}
}

func (r *Runner) RunWithClusterBomb(targets chan *Target) {
	rootContext, _ := context.WithCancel(context.Background())
	pool, _ := ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		task := i.(*pkg.Task)
		ctx, cancel := context.WithCancel(task.Context) // current task context
		go func() {
			err := Brute(task)
			if err != nil {
				r.OutputCh <- &pkg.Result{
					Task: task,
					Err:  err,
				}
			} else {
				r.OutputCh <- &pkg.Result{
					Task: task,
					OK:   true,
				}
				if r.FirstOnly {
					cancel()        // 退出当前任务
					task.Canceler() // 取消正在执行的所有任务
				}
			}
			cancel()
		}()

		// 设置超时时间, 防止任务挂死
		select {
		case <-ctx.Done():
		case <-task.Context.Done():
		case <-time.After(time.Duration(task.Timeout+1) * time.Second):
			r.OutputCh <- &pkg.Result{
				Task: task,
				Err:  fmt.Errorf("timeout"),
			}
		}
	})

	// 执行

	for target := range targets {
		ch := r.clusterBombGenerate(rootContext, target)
	loop:
		for {
			select {
			case task, ok := <-ch:
				// 从生成器中取任务.
				if ok {
					r.wg.Add(1)
					r.Stat.Count++
					_ = pool.Invoke(task)
				} else {
					break loop
				}
			case <-rootContext.Done():
				// todo 为断点续传做准备
				break loop
			}
		}
	}
	r.wg.Wait()
}

func (r *Runner) clusterBombGenerate(ctx context.Context, target *Target) chan *pkg.Task {
	// 通过用户名与密码的笛卡尔积生成数据
	tctx, canceler := context.WithCancel(ctx)
	ch := make(chan *pkg.Task)
	var users, pwds []string
	// 自动选择默认的用户名与密码字典
	if r.Users == nil {
		users = pkg.UseDefaultUser(target.Service.String())
	} else {
		users = r.Users.RunAsSlice()
	}

	if r.Pwds == nil {
		pwds = pkg.UseDefaultPassword(target.Service.String(), r.Top)
	} else {
		pwds = r.Pwds.RunAsSlice()
	}

	// task生成器
	go func() {
	Loop:
		for _, user := range users {
			for _, pwd := range pwds {
				if target.Service == "" {
					logs.Log.Warn("unknown service " + target.Service.String())
					continue
				}
				select {
				case ch <- &pkg.Task{
					IP:         target.IP,
					Port:       target.Port,
					Service:    target.Service,
					Username:   user,
					Password:   pwd,
					Timeout:    r.Timeout,
					ExecString: r.ExecString,
					Context:    tctx,
					Canceler:   canceler,
				}:
				case <-tctx.Done():
					break Loop
				}
			}
		}
		close(ch)
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

var outed int

func (r *Runner) Output() {

loop:
	for {
		select {
		case result, ok := <-r.OutputCh:
			if !ok {
				break loop
			}
			outed++
			if result.OK {
				if r.File != nil {
					r.OutFunc(result.Format(r.Option.FileFormat))
				}
				logs.Log.Console(result.Format(r.Option.OutputFormat))
			} else {
				logs.Log.Debugf(" %s\t%s\t%s ,failed, %s", result.URI(), result.Username, result.Password, result.Err.Error())
			}
			r.wg.Done()
		}
	}
}
