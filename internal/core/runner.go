package core

import (
	"context"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/pkg"
	"github.com/panjf2000/ants/v2"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	*Option
	Users      []string
	Pwds       []string
	Addrs      ipcs.Addrs
	Targets    []*Target
	Services   []string
	Generator  chan *Target
	OutputCh   chan *pkg.Result
	File       *files.File
	OutFunc    func(string)
	Done       bool
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
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		task := i.(*pkg.Task)
		ctx, cancel := context.WithTimeout(task.Context, time.Duration(task.Timeout)*time.Second)
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
					task.Canceler()
				}
			}
			cancel()
		}()

		// 设置超时时间, 防止任务挂死
		select {
		case <-ctx.Done():
		case <-time.After(time.Duration(task.Timeout+1) * time.Second):
			r.OutputCh <- &pkg.Result{
				Task: task,
				Err:  fmt.Errorf("timeout"),
			}
			cancel()
		}
		wg.Done()
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
					wg.Add(1)
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
	wg.Wait()

	// 某些情况下, 任务执行完毕后, 还在处理输出结果, 等待结果输出完毕
	for len(r.OutputCh) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	close(r.OutputCh)
	time.Sleep(100 * time.Millisecond)
}

func (r *Runner) clusterBombGenerate(ctx context.Context, target *Target) chan *pkg.Task {
	// 通过用户名与密码的笛卡尔积生成数据
	tctx, canceler := context.WithCancel(ctx)
	ch := make(chan *pkg.Task)
	var users, pwds []string
	// 自动选择默认的用户名与密码字典
	if r.Users == nil {
		users = pkg.UseDefaultUser(target.Service)
	} else {
		users = r.Users
	}

	if r.Pwds == nil {
		pwds = pkg.UseDefaultPassword(target.Service, r.Top)
	} else {
		pwds = r.Pwds
	}

	// task生成器
	go func() {
	Loop:
		for _, user := range users {
			for _, pwd := range pwds {
				if _, ok := pkg.ServicePortMap[target.Service]; !ok {
					logs.Log.Warn("unknown service " + target.Service)
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
			target.Service = strings.ToUpper(target.Service)
			if r.Services == nil || (r.Services != nil && pkg.SliceContains(r.Services, target.Service)) {
				// 如果从gogo中输入的目标, 可以通过-s过滤特定的服务进行扫描
				ch <- target
			}
		}

		// 通过addrs生成目标
		for _, addr := range r.Addrs {
			for _, service := range r.Services {
				ch <- &Target{
					IP:      addr.IP.String(),
					Port:    addr.Port,
					Service: service,
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (r *Runner) Output() {
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
				logs.Log.Debugf(" %s\t%s\t%s ,failed, %s", result.URI(), result.Username, result.Password, result.Err.Error())
			}
		}
	}
}
