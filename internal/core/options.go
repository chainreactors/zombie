package core

import (
	"encoding/json"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/words"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"os"
	"strings"
)

type Option struct {
	InputOptions  `group:"Input Options"`
	OutputOptions `group:"Output Options"`
	WordOptions   `group:"Word Options"`
	MiscOptions   `group:"Misc Options"`
}

type InputOptions struct {
	IP            string `short:"i" long:"ip"`
	IPFile        string `short:"I" long:"IP"`
	Username      string `short:"u" long:"user"`
	UsernameFile  string `short:"U" long:"USER"`
	Password      string `short:"p" long:"pwd"`
	PasswordFile  string `short:"P" long:"PWD"`
	GogoFile      string `long:"go"`
	ServiceName   string `short:"s" long:"service"`
	FilterService string `short:"S" long:"filter-service"`
}

type OutputOptions struct {
	OutputFile   string `short:"f" long:"file"`
	FileFormat   string `short:"O" long:"file-format" default:"json"`
	OutputFormat string `short:"o" long:"format" default:"string"`
}

type WordOptions struct {
	Top int `long:"top" default:"10"`
}

type MiscOptions struct {
	Threads int    `short:"t" default:"100"`
	Timeout int    `short:"d" long:"timeout" default:"5"`
	Mod     string `short:"m" default:"clusterbomb"`
	Debug   bool   `long:"debug"`
}

func (opt *Option) Validate() error {
	return nil
}

func (opt *Option) Prepare() (*Runner, error) {
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
			addrs = ipcs.NewAddrsWithDefaultPort(ipslice, pkg.ServicePortMap[strings.ToUpper(opt.ServiceName)])
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
		OutputCh:  make(chan *pkg.Result),
	}
	if opt.ServiceName != "" {
		runner.Services = strings.Split(strings.ToUpper(opt.ServiceName), ",")
	}
	return runner, nil

}

type Target struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Service string `json:"service"`
}

func (t Target) Addr() *ipcs.Addr {
	return &ipcs.Addr{IP: ipcs.NewIP(t.IP), Port: t.Port}
}
