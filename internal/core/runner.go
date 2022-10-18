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
			addrs = ipcs.NewAddrsWithDefaultPort(ipslice, utils.ServicePortMap[opt.ServiceName])
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
	outfunc := logs.Log.Console
	if opt.OutputFile != "" {
		file, err = files.NewFile(opt.OutputFile, false, true, true)
		if err != nil {
			return nil, err
		}
		outfunc = func(s string) {
			file.SafeWrite(s)
		}
	}

	runner := &Runner{
		Users:    users,
		Pwds:     pwds,
		Addrs:    addrs,
		Targets:  targets,
		Services: strings.Split(opt.ServiceName, ","),
		Threads:  opt.Threads,
		Timeout:  opt.Timeout,
		Mod:      opt.Mod,
		File:     file,
		OutFunc:  outfunc,
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
	Threads    int
	Timeout    int
	ExecString string
	Mod        string
	FirstOnly  bool
	OutputCh   chan *utils.Result
	File       *files.File
	OutFunc    func(string)
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
		result := Brute(task)
		if result.OK {
			if r.FirstOnly {
				task.Canceler()
			}
			r.OutputCh <- result
		} else {
			logs.Log.Debugf(" %s\t%s\t%s ,failed, %s", task.URI(), task.Username, task.Password, result.Err.Error())
		}
		wg.Done()
	})

	for target := range targets {
		ctx, canceler := context.WithCancel(rootContext)
		ch := r.clusterBombGenerate(target, canceler)
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
			case <-ctx.Done():
				break loop
			}
		}
	}
	wg.Wait()
}

func (r *Runner) clusterBombGenerate(target *Target, canceler context.CancelFunc) chan *utils.Task {
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
		for _, user := range users {
			for _, pwd := range pwds {
				ch <- &utils.Task{
					IP:         target.IP,
					Port:       target.Port,
					Service:    target.Service,
					Username:   user,
					Password:   pwd,
					ExecString: r.ExecString,
					Canceler:   canceler,
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
				if slice.Contains(r.Services, t.Service) {
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
	for {
		select {
		case result, ok := <-r.OutputCh:
			if ok {
				r.OutFunc(result.String())
			} else {
				logs.Log.Debug(result.Err.Error())
			}
		}
	}
}
