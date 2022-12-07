package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/words"
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/chainreactors/zombie/pkg/utils/slice"
	"github.com/panjf2000/ants/v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

func PrepareRunner(opt *Option) (*Runner, error) {
	var err error

	if err = opt.Validate(); err != nil {
		return nil, err
	}

	if opt.Debug {
		logs.Log.Level = logs.Debug
	}

	var targets []*Target
	var users, pwds []string
	var addrs ipcs.Addrs
	if opt.GogoFile != "" {
		// load gogo result
		content, err := ioutil.ReadFile(opt.GogoFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &targets)
		if err != nil {
			return nil, err
		}
	} else {
		// load target
		var ipslice []string
		if opt.IP != "" {
			ipslice = strings.Split(opt.IP, ",")
		} else if opt.IPFile != "" {
			ipf, err := os.Open(opt.IPFile)
			if err != nil {
				return nil, err
			}
			ipslice = words.NewWorderWithFile(ipf).All()
		}

		if len(ipslice) == 0 {
			return nil, fmt.Errorf("not any ip input")
		}

		if strings.Contains(ipslice[0], ":") {
			addrs = ipcs.NewAddrs(ipslice)
		} else {
			addrs = ipcs.NewAddrsWithDefaultPort(ipslice, utils.ServicePortMap[strings.ToUpper(opt.ServiceName)])
		}
	}

	// load username
	if opt.Username != "" {
		users = strings.Split(opt.Username, ",")
	} else if opt.UsernameFile != "" {
		userf, err := os.Open(opt.UsernameFile)
		if err != nil {
			return nil, err
		}
		users = words.NewWorderWithFile(userf).All()
	}

	// load password
	if opt.Password != "" {
		pwds = strings.Split(opt.Password, ",")
	} else if opt.PasswordFile != "" {
		pwdf, err := os.Open(opt.PasswordFile)
		if err != nil {
			return nil, err
		}
		pwds = words.NewWorderWithFile(pwdf).All()
	}

	var file *files.File
	var outfunc func(string)
	if opt.OutputFile != "" {
		file, err = files.NewFile(opt.OutputFile, false, true, true)
		if err != nil {
			return nil, err
		}
		outfunc = func(s string) {
			file.SafeWrite(s)
			file.SafeSync()
		}
	}

	runner := &Runner{
		Users:     users,
		Pwds:      pwds,
		Addrs:     addrs,
		Targets:   targets,
		Option:    opt,
		File:      file,
		FirstOnly: true,
		OutFunc:   outfunc,
		OutputCh:  make(chan *utils.Result),
	}
	if opt.ServiceName != "" {
		runner.Services = strings.Split(strings.ToUpper(opt.ServiceName), ",")
	}
	return runner, nil
}

type Runner struct {
	Users      []string
	Pwds       []string
	Addrs      ipcs.Addrs
	Targets    []*Target
	Services   []string
	Generator  chan *Target
	OutputCh   chan *utils.Result
	File       *files.File
	OutFunc    func(string)
	Done       bool
	ExecString string
	FirstOnly  bool
	*Option
}

func (r *Runner) Run() {
	go r.Outputting()
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
		task := i.(*utils.Task)
		ctx, cancel := context.WithTimeout(task.Context, time.Duration(task.Timeout)*time.Second)
		go func() {
			err := Brute(task)
			if err != nil {
				r.OutputCh <- &utils.Result{
					Task: task,
					Err:  err,
				}
			} else {
				r.OutputCh <- &utils.Result{
					Task: task,
					OK:   true,
				}
				if r.FirstOnly {
					task.Canceler()
				}
			}
			cancel()
		}()

		select {
		case <-ctx.Done():
		case <-time.After(time.Duration(task.Timeout+1) * time.Second):
			r.OutputCh <- &utils.Result{
				Task: task,
				Err:  fmt.Errorf("timeout"),
			}
			cancel()
		}
		wg.Done()
	})

	for target := range targets {
		ctx, canceler := context.WithCancel(rootContext)
		ch := r.clusterBombGenerate(ctx, target, canceler)
	loop:
		for {
			select {
			case task, ok := <-ch:
				if ok {
					wg.Add(1)
					_ = pool.Invoke(task)
				} else {
					break loop
				}
			}
		}
	}
	wg.Wait()

	for len(r.OutputCh) == 0 {
		close(r.OutputCh)
		break
	}
	time.Sleep(100)
}

func (r *Runner) clusterBombGenerate(ctx context.Context, target *Target, canceler context.CancelFunc) chan *utils.Task {
	ch := make(chan *utils.Task)
	var users, pwds []string
	if r.Users == nil {
		users = utils.UseDefaultUser(target.Service)
	} else {
		users = r.Users
	}

	if r.Pwds == nil {
		pwds = utils.UseDefaultPassword(target.Service)
	} else {
		pwds = r.Pwds
	}

	go func() {
	Loop:
		for _, user := range users {
			for _, pwd := range pwds {
				if _, ok := utils.ServicePortMap[target.Service]; !ok {
					logs.Log.Debugf("unknown service " + target.Service)
					continue
				}
				select {
				case ch <- &utils.Task{
					IP:         target.IP,
					Port:       target.Port,
					Service:    target.Service,
					Username:   user,
					Password:   pwd,
					Timeout:    r.Timeout,
					ExecString: r.ExecString,
					Context:    ctx,
					Canceler:   canceler,
				}:
					continue
				case <-ctx.Done():
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
		if r.Targets != nil {
			for _, t := range r.Targets {
				t.Service = strings.ToUpper(t.Service)
				if r.Services != nil {
					if slice.Contains(r.Services, t.Service) {
						ch <- t
					}
				} else {
					ch <- t
				}
			}
		} else {
			for _, addr := range r.Addrs {
				for _, service := range r.Services {
					ch <- &Target{
						IP:      addr.IP.String(),
						Port:    addr.Port,
						Service: service,
					}
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (r *Runner) Outputting() {
loop:
	for {
		select {
		case result, ok := <-r.OutputCh:
			if ok {

				if result.OK {
					if r.File != nil {
						r.OutFunc(result.Format(r.Option.FileFormat))
					}
					logs.Log.Console(result.Format(r.Option.OutputFormat))
				} else {
					logs.Log.Debugf(" %s\t%s\t%s ,failed, %s", result.URI(), result.Username, result.Password, result.Err.Error())
				}

			} else {
				break loop
			}
		}
	}
}
