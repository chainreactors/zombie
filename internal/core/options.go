package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chainreactors/files"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/pkg"
	"github.com/vbauerster/mpb/v8"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Option struct {
	InputOptions  `group:"Input Options"`
	OutputOptions `group:"Output Options"`
	WordOptions   `group:"Word Options"`
	MiscOptions   `group:"Misc Options"`
}

type InputOptions struct {
	IP            []string          `short:"i" long:"ip" alias:"ipp" description:"String, input ip"`
	IPFile        string            `short:"I" long:"IP" description:"File, input ip list filename"`
	Username      []string          `short:"u" long:"user" description:"Strings, input usernames"`
	UsernameFile  string            `short:"U" long:"USER" description:"File, input username list filename"`
	Auth          []string          `short:"a" long:"auth" description:"Strings, input auth, username::password"`
	AuthFile      string            `short:"A" long:"AUTH" description:"File, input auth list filename"`
	UsernameRule  string            `long:"userrule" description:"String, input username generator rule filename"`
	Password      []string          `short:"p" long:"pwd" description:"String, input passwords"`
	PasswordFile  string            `short:"P" long:"PWD" description:"File, input password list filename"`
	PasswordRule  string            `long:"pwdrule" description:"String, input password generator rule filename"`
	Dictionaries  []string          `short:"d" long:"dict" description:"Strings, input dictionaries"`
	JsonFile      string            `short:"j" long:"json" description:"File, input json result filename"`
	GogoFile      string            `short:"g" long:"gogo" description:"File, input gogo result filename"`
	ServiceName   string            `short:"s" long:"service" description:"String, input service name"`
	FilterService string            `short:"S" long:"filter-service" description:"String, filter service when input json/gogo file"`
	Param         map[string]string `long:"param" description:"params"`
}

type OutputOptions struct {
	OutputFile   string `short:"f" long:"file" description:"File, output result filename"`
	FileFormat   string `short:"O" long:"file-format" default:"json" description:"String, output result file format"`
	OutputFormat string `short:"o" long:"format" default:"string" description:"String, output result format"`
}

type WordOptions struct {
	Top             int  `long:"top" default:"0" description:"Int, top n words"`
	ForceContinue   bool `long:"force-continue" description:"Bool, force continue, not only stop when first success ever host"`
	WeakPassWord    bool `long:"weakpass" description:"Bool, common weak password rule"`
	NoUnAuth        bool `long:"no-unauth" description:"Bool, skip check unauth"`
	NoCheckHoneyPot bool `long:"no-honeypot" description:"Bool, skip check honeypot"`
}

type MiscOptions struct {
	Raw         bool   `long:"raw" description:"Bool, parser raw username/password"`
	Threads     int    `short:"t" default:"100" description:"Int, threads"`
	Timeout     int    `long:"timeout" default:"5" description:"Int, timeout"`
	Mod         string `short:"m" default:"clusterbomb" description:"String, clusterbomb/sniper"`
	Debug       bool   `long:"debug" description:"Bool, enable debug"`
	ListService bool   `short:"l" long:"list" description:"Bool, list all service"`
	Bar         bool   `long:"bar" description:"Bool, enable bar"`
	Version     bool   `long:"version" description:"Bool, show version"`
}

func (opt *Option) Validate() error {
	if len(opt.IP) == 0 && opt.IPFile == "" && opt.JsonFile == "" && opt.GogoFile == "" {
		return errors.New("please input ip or or file or json file or gogo file")
	}
	if opt.WeakPassWord && (opt.Password == nil && opt.PasswordFile == "") {
		return errors.New("use weak-password rule must set password, please set -p/-P")
	}
	if opt.PasswordRule != "" && (opt.Password == nil && opt.PasswordFile == "") {
		return errors.New("use custom password rule must set password, please set -p/-P")
	}
	if opt.UsernameRule != "" && (opt.Username == nil && opt.UsernameFile == "") {
		return errors.New("use custom username rule must set username, please set -u/-U")
	}
	return nil
}

func (opt *Option) Prepare() (*Runner, error) {
	var err error
	var targets []*Target

	var file *files.File
	var outfunc func(string)
	if opt.OutputFile != "" {
		file, err = files.NewFile(opt.OutputFile, false, false, true)
		if err != nil {
			return nil, err
		}
		outfunc = func(s string) {
			file.SafeWrite(s)
			file.SafeSync()
		}
	}

	runner := &Runner{
		File:      file,
		OutFunc:   outfunc,
		FirstOnly: !opt.ForceContinue,
		Option:    opt,
		wg:        &sync.WaitGroup{},
		outlock:   &sync.WaitGroup{},
		addlock:   &sync.Mutex{},
		OutputCh:  make(chan *pkg.Result),
		stat: &pkg.Statistor{
			Tasks: make(map[string]int),
		},
	}

	if opt.Bar {
		runner.progress = mpb.New(
			mpb.WithRefreshRate(200*time.Millisecond),
			mpb.WithOutput(os.Stdout),
		)
		logs.Log.SetOutput(runner.progress)
	}

	logs.Log.Importantf("mod: %s, check-unauth: %t, check-honeypot: %t", runner.Mod, !runner.NoUnAuth, !runner.NoCheckHoneyPot)

	if opt.ServiceName != "" {
		runner.Services = strings.Split(strings.ToLower(opt.ServiceName), ",")
	}

	if opt.JsonFile != "" {
		// load json file
		content, err := ioutil.ReadFile(opt.JsonFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &targets)
		if err != nil {
			return nil, err
		}
		logs.Log.Importantf("load %d targets from json: %s ", len(targets), opt.JsonFile)
	} else if opt.GogoFile != "" {
		targets, err = LoadGogoFile(opt.GogoFile)
		if err != nil {
			return nil, err
		}
		logs.Log.Importantf("load %d targets from gogo: %s ", len(targets), opt.GogoFile)
	} else {
		var ipslice []string
		if opt.IP != nil {
			ipslice = opt.IP
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

		// 处理输入参数
		for _, input := range ipslice {
			t, ok := ParseUrl(input)
			if !ok {
				t = SimpleParseUrl(input)
			}

			targets = append(targets, t)
		}
		if opt.IPFile != "" {
			logs.Log.Importantf("load %d targets from file: %s", len(targets), opt.IPFile)
		}
	}

	for _, t := range targets {
		// 如果指定了service, 将会覆盖json或gogo中的字段
		if opt.ServiceName != "" {
			t.UpdateService(opt.ServiceName)
		}
		if t.Service == "" {
			logs.Log.Warn(t.String() + " null service")
			continue
		}

		// 命令行中指定的 param 会覆盖原有的配置
		if len(opt.Param) > 0 {
			t.Param = opt.Param
		}
	}
	runner.Targets = targets

	var dicts [][]string
	if opt.Dictionaries != nil {
		var s strings.Builder
		dicts = make([][]string, len(opt.Dictionaries))
		for i, f := range opt.Dictionaries {
			dicts[i], err = loadFileToSlice(f)
			if err != nil {
				return nil, err
			}
			s.WriteString(fmt.Sprintf("%s: %ditems", f, len(dicts[i])))
		}

		logs.Log.Importantf("load dictionaries: %s", s.String())
	}

	var users, pwds, auths *Generator
	// load username
	if opt.Username != nil {
		if len(opt.Username) == 1 && dicts != nil {
			users, err = NewGeneratorWithWord(opt.Username[0], dicts, nil)
			if err != nil {
				return nil, err
			}
			logs.Log.Importantf("parse username from %s", opt.Username[0])
		} else {
			users = NewGeneratorWithInput(opt.Username)
		}
	} else if opt.UsernameFile != "" {
		users, err = NewGeneratorWithFile(opt.UsernameFile)
		if err != nil {
			return nil, err
		}
		logs.Log.Importantf("load username from %s", opt.UsernameFile)
	}
	if opt.UsernameRule != "" {
		err := users.SetRuleFile(opt.UsernameRule)
		if err != nil {
			return nil, err
		}
	}
	runner.Users = users

	// load password
	if opt.Password != nil {
		if len(opt.Password) == 1 && dicts != nil {
			pwds, err = NewGeneratorWithWord(opt.Password[0], dicts, nil)
			if err != nil {
				return nil, err
			}
			logs.Log.Importantf("parse password from %s ", opt.Password[0])
		} else {
			pwds = NewGeneratorWithInput(opt.Password)
		}
	} else if opt.PasswordFile != "" {
		pwds, err = NewGeneratorWithFile(opt.PasswordFile)
		if err != nil {
			return nil, err
		}
		logs.Log.Importantf("load password from %s", opt.PasswordFile)
	}
	if opt.PasswordRule != "" {
		err := pwds.SetRuleFile(opt.PasswordRule)
		if err != nil {
			return nil, err
		}
	} else if opt.WeakPassWord {
		err := pwds.SetInternalRule("weakpass")
		if err != nil {
			return nil, err
		}
	}
	runner.Pwds = pwds

	// load auth pair
	if opt.Auth != nil {
		auths = NewGeneratorWithInput(opt.Auth)
	} else if opt.AuthFile != "" {
		auths, err = NewGeneratorWithFile(opt.AuthFile)
		if err != nil {
			return nil, err
		}
		logs.Log.Importantf("load auth from %s", opt.AuthFile)
	}
	runner.Auths = auths

	if runner.progress != nil {
		runner.bar = pkg.NewBar("targets", len(targets), runner.stat, runner.progress)
	}
	return runner, nil
}
