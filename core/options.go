package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/fileutils"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"strings"
)

type Option struct {
	InputOptions  `group:"Input Options"`
	OutputOptions `group:"Output Options"`
	WordOptions   `group:"Word Options"`
	ActionOptions `group:"Post-Auth Actions"`
	MiscOptions   `group:"Misc Options"`
}

type InputOptions struct {
	IP            []string          `short:"i" long:"ip" alias:"ipp" description:"String, input ip"`
	IPFile        string            `short:"I" long:"IP" description:"File, input ip list filename"`
	CIDR          []string          `short:"c" long:"cidr" description:"String, input cidr"`
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
	Debug        bool   `long:"debug" description:"Bool, enable debug"`
	Quiet        bool   `short:"q" long:"quiet" description:"Bool, quiet mode"`
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
	Strict      bool   `long:"strict" description:"Bool, strict mode, when finger check pass will brute"`
	Threads     int    `short:"t" default:"100" description:"Int, threads"`
	Concurrency int    `long:"concurrency" default:"8" description:"Int, max concurrent connections per host, keep below service rate-limits e.g. sshd MaxStartups(default 10) to avoid random connection drops being misread as wrong password; 0=unlimited"`
	Timeout     int    `long:"timeout" default:"5" description:"Int, timeout"`
	Mod         string `short:"m" default:"clusterbomb" description:"String, clusterbomb/pitchfork/sniper"`
	ListService bool   `short:"l" long:"list" description:"Bool, list all service"`
	Bar         bool   `long:"bar" description:"Bool, enable bar"`
	Version     bool   `long:"version" description:"Bool, show version"`
}

type ActionOptions struct {
	Proton        bool     `long:"proton" description:"post-auth: collect info + run proton credential scan"`
	ScanTemplates []string `long:"scan-template" description:"proton template file or directory for --proton"`
	DBLimit       int      `long:"db-limit" default:"1000" description:"max rows per column in DB credential scan"`
}

func (opt *Option) Validate() error {
	if opt.Mod == "" {
		opt.Mod = ModBomb
	}
	switch opt.Mod {
	case ModBomb, ModPitchFork, ModSniper:
	default:
		return fmt.Errorf("unsupported mod %q, want clusterbomb, pitchfork, or sniper", opt.Mod)
	}
	if len(opt.IP) == 0 && opt.IPFile == "" && opt.JsonFile == "" && opt.GogoFile == "" && opt.CIDR == nil {
		return errors.New("please input ip or or file or json file or gogo file")
	}
	if opt.Mod == ModPitchFork && opt.Auth == nil && opt.AuthFile == "" {
		return errors.New("pitchfork mode requires auth, please set -a/-A")
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

	var file *fileutils.File
	var outfunc func(string)
	if opt.OutputFile != "" {
		file, err = fileutils.NewFile(opt.OutputFile, fileutils.ModeAppend, false, false)
		if err != nil {
			return nil, err
		}
		outfunc = func(s string) {
			if err := file.SyncWrite(s); err != nil {
				logs.Log.Warn(fmt.Sprintf("write output file failed: %v", err))
			}
		}
	}

	runnerOpt := &RunnerOption{
		Threads:         opt.Threads,
		Concurrency:     opt.Concurrency,
		Timeout:         opt.Timeout,
		Top:             opt.Top,
		Mod:             opt.Mod,
		FirstOnly:       !opt.ForceContinue,
		NoUnAuth:        opt.NoUnAuth,
		NoCheckHoneyPot: opt.NoCheckHoneyPot,
		Strict:          opt.Strict,
		Raw:             opt.Raw,
		Proton:          opt.Proton,
		ScanTemplates:   opt.ScanTemplates,
		DBLimit:         opt.DBLimit,
	}

	runner := NewRunner(runnerOpt)
	if err := runner.BuildPipeline(); err != nil {
		return nil, err
	}
	runner.File = file
	runner.OutFunc = outfunc
	runner.FileFormat = opt.FileFormat
	runner.OutputFormat = opt.OutputFormat

	if opt.Bar {
		pkg.InitBar()
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
		var ipg *Generator

		if opt.IP != nil {
			ipg = NewGeneratorWithInput(opt.IP)
		} else if opt.IPFile != "" {
			ipg, err = NewGeneratorWithFile(opt.IPFile)
			if err != nil {
				return nil, err
			}
		} else if opt.CIDR != nil {
			ipg = NewGeneratorWithChan(transformChan(utils.ParseCIDRs(opt.CIDR).Range()))
		}

		if ipg == nil {
			return nil, fmt.Errorf("not any ip input")
		}

		ipg.Run()

		// 处理输入参数
		for input := range ipg.C {
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

		if opt.FilterService != "" {
			var ok bool
			for _, s := range strings.Split(opt.FilterService, ",") {
				if s == t.Service {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		// 命令行中指定的 param 会覆盖原有的配置
		if len(opt.Param) > 0 {
			t.Param = opt.Param
		}
		runner.Targets = append(runner.Targets, t)
	}

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

	var users, pwds *Generator
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
	var auths *Generator
	if opt.Auth != nil {
		auths = NewGeneratorWithInput(opt.Auth)
	} else if opt.AuthFile != "" {
		auths, err = NewGeneratorWithFile(opt.AuthFile)
		if err != nil {
			return nil, err
		}
		logs.Log.Importantf("load auth from %s", opt.AuthFile)
	}
	if auths != nil {
		runner.Auths = auths
		runner.Mod = ModPitchFork
	}

	runner.bar = pkg.NewBar("targets", len(targets), runner.stat)

	return runner, nil
}
