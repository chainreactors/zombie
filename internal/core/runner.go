package core

import (
	"context"
	"fmt"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/words"
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/panjf2000/ants/v2"
	"os"
	"strings"
	"sync"
)

func PrepareRunner(opt *Option) (*Runner, error) {
	var err error
	opt.ServiceName = strings.ToUpper(opt.ServiceName)
	if err = opt.Validate(); err != nil {
		return nil, err
	}

	var users, pwds []string
	var addrs ipcs.Addrs
	if opt.Debug {
		logs.Log.Level = logs.Debug
	}
	if opt.GogoFile != "" {

	} else {
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

	if opt.Username != "" {
		users = strings.Split(opt.Username, ",")
	} else if opt.UsernameFile != "" {
		userf, err := os.Open(opt.UsernameFile)
		if err != nil {
			return nil, err
		}
		users = words.NewWorderWithFile(userf).All()
	} else if tmp, ok := utils.DefaultUsernames[opt.ServiceName]; ok {
		users = tmp
	} else {
		users = []string{"admin"}
	}

	if opt.Password != "" {
		pwds = strings.Split(opt.Password, ",")
	} else if opt.PasswordFile != "" {
		pwdf, err := os.Open(opt.PasswordFile)
		if err != nil {
			return nil, err
		}
		pwds = words.NewWorderWithFile(pwdf).All()
	} else if tmp, ok := utils.DefaultPasswords[opt.ServiceName]; ok {
		pwds = tmp
	} else {
		pwds = []string{"admin"}
	}

	runner := &Runner{
		Users:      users,
		Pwds:       pwds,
		Addrs:      addrs,
		Service:    opt.ServiceName,
		Threads:    opt.Threads,
		Timeout:    opt.Timeout,
		ExecString: "",
		Mod:        opt.Mod,
	}
	return runner, nil
}

type Runner struct {
	Users      []string
	Pwds       []string
	Addrs      ipcs.Addrs
	Service    string
	Threads    int
	Timeout    int
	ExecString string
	Mod        string
	First      bool
	OutputCh   chan *utils.Result
}

func (r *Runner) Run() {
	go r.Outputting()

	switch r.Mod {
	case "pitchfork":
	case "clusterbomb":
		r.RunWithClusterBomb()
	}
}

func (r *Runner) RunWithPitchfork() {
	//rootContext, rootCancel := context.WithCancel(context.Background())
	//for _, addr := range r.Addrs{
	//	 for _, user := range r.Users{}
	//}

}

func (r *Runner) RunWithClusterBomb() {
	rootContext, _ := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(r.Threads, func(i interface{}) {
		req := i.(*utils.Task)
		result := Brute(req)
		if result.OK {
			if r.First {
				req.Canceler()
			}
			r.OutputCh <- result
		} else {
			logs.Log.Debugf(" %s\t%s\tfailed, %s", req.String(), req.Username, req.Password, result.Err.Error())
		}
		wg.Done()
	})

	for _, addr := range r.Addrs {
		ctx, canceler := context.WithCancel(rootContext)
		ch := r.clusterBombGenerate(rootContext, canceler, addr)
		select {
		case req := <-ch:
			wg.Add(1)
			err := pool.Invoke(req)
			if err != nil {
				logs.Log.Error(err.Error())
			}
		case <-ctx.Done():
			continue
		}
	}
	wg.Wait()
}

func (r *Runner) clusterBombGenerate(ctx context.Context, canceler context.CancelFunc, addr *ipcs.Addr) chan *utils.Task {
	ch := make(chan *utils.Task)
	go func() {
		for _, user := range r.Users {
			for _, pwd := range r.Pwds {
				ch <- &utils.Task{
					Addr:       addr,
					Service:    r.Service,
					Username:   user,
					Password:   pwd,
					ExecString: r.ExecString,
					Context:    ctx,
					Canceler:   canceler,
				}
			}
		}
		close(ch)
	}()
	return ch
}

func (r *Runner) Outputting() {
Loop:
	for {
		select {
		case result, ok := <-r.OutputCh:
			if ok {
				logs.Log.Console(result.String())
			} else {
				break Loop
			}

		}
	}
}
