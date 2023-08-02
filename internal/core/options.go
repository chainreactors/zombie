package core

import (
	"encoding/json"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"strings"
)

type Option struct {
	InputOptions  `group:"Input Options"`
	OutputOptions `group:"Output Options"`
	WordOptions   `group:"Word Options"`
	MiscOptions   `group:"Misc Options"`
}

type InputOptions struct {
	IP            string   `short:"i" long:"ip" alias:"ipp" description:"String, input ip"`
	IPFile        string   `short:"I" long:"IP" description:"File, input ip list filename"`
	Username      []string `short:"u" long:"user" description:"Strings, input usernames"`
	UsernameFile  string   `short:"U" long:"USER" description:"File, input username list filename"`
	UsernameWord  string   `long:"userword" description:"String, input username generator dsl"`
	UsernameRule  string   `long:"userrule" description:"String, input username generator rule filename"`
	Password      []string `short:"p" long:"pwd" description:"String, input passwords"`
	PasswordFile  string   `short:"P" long:"PWD" description:"File, input password list filename"`
	PasswordWord  string   `long:"pwdword" description:"String, input password generator dsl"`
	PasswordRule  string   `long:"pwdrule" description:"String, input password generator rule filename"`
	GogoFile      string   `long:"go" description:"File, input gogo result filename"`
	ServiceName   string   `short:"s" long:"service" description:"String, input service name"`
	FilterService string   `short:"S" long:"filter-service" description:"String, filter service name"`
}

type OutputOptions struct {
	OutputFile   string `short:"f" long:"file" description:"File, output result filename"`
	FileFormat   string `short:"O" long:"file-format" default:"json" description:"String, output result file format"`
	OutputFormat string `short:"o" long:"format" default:"string" description:"String, output result format"`
}

type WordOptions struct {
	Top           int  `long:"top" default:"0" description:"Int, top n words"`
	ForceContinue bool `long:"force-continue" description:"Bool, force continue, not only stop when first success ever host"`
}

type MiscOptions struct {
	Threads int    `short:"t" default:"100" description:"Int, threads"`
	Timeout int    `short:"d" long:"timeout" default:"5" description:"Int, timeout"`
	Mod     string `short:"m" default:"clusterbomb" description:"String, mod"`
	Debug   bool   `long:"debug" description:"Bool, enable debug"`
}

func (opt *Option) Validate() error {
	return nil
}

func (opt *Option) Prepare() (*Runner, error) {
	var err error
	var targets []*Target
	var users, pwds *Generator
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
			ipg, err := NewGeneratorWithFile(opt.IPFile)
			if err != nil {
				return nil, err
			}
			ipslice = ipg.RunAsSlice()
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
	if opt.Username != nil {
		users = NewGeneratorWithInput(opt.Username)
	} else if opt.UsernameFile != "" {
		users, err = NewGeneratorWithFile(opt.UsernameFile)
		if err != nil {
			return nil, err
		}
	} else if opt.UsernameWord != "" {
		users, err = NewGeneratorWithWord(opt.UsernameWord, nil, nil)
		if err != nil {
			return nil, err
		}
	}

	if opt.UsernameRule != "" {
		err := users.SetRule(opt.UsernameRule)
		if err != nil {
			return nil, err
		}
	}
	// load password
	if opt.Password != nil {
		pwds = NewGeneratorWithInput(opt.Password)
	} else if opt.PasswordFile != "" {
		pwds, err = NewGeneratorWithFile(opt.PasswordFile)
		if err != nil {
			return nil, err
		}
	} else if opt.PasswordWord != "" {
		pwds, err = NewGeneratorWithWord(opt.PasswordWord, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	if opt.PasswordRule != "" {
		err := users.SetRule(opt.PasswordRule)
		if err != nil {
			return nil, err
		}
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
		Stat:      &pkg.Statistor{},
	}
	if opt.ForceContinue {
		runner.FirstOnly = false
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
