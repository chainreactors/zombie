package cmd

import (
	"fmt"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/internal/core"
	"github.com/chainreactors/zombie/pkg"
	"github.com/jessevdk/go-flags"
	"strings"
)

var ver = "dev"

func Zombie() {
	var opt core.Option
	parser := flags.NewParser(&opt, flags.Default)
	parser.Usage = `

	WIKI: https://chainreactors.github.io/wiki/zombie

	QUICKSTART:
		simple example:
			zombie -i 1.1.1.1 -u root -s ssh
	
		brute multiple ssh targets(ip list):
			zombie -I targets.txt -u root -p password -s ssh

		brute from file and auto parse
			zombie -I targets.txt
	
			targets.txt:
			mysql://user:pass@1.1.1.1:3307  
			ssh://user@2.2.2.2             
			mssql://3.3.3.3:1433            
	
	
		rude brute:
			zombie -I targets.txt -U user.txt -P pass.txt

	
		brute from gogo dat:
			zombie --gogo 1.dat
	
		brute from json file:
			zombie -j 1.json

		weak password generate:
			zombie -l 1.txt -p google --weakpass
`
	_, err := parser.Parse()
	if err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Println(err.Error())
		}
		return
	}

	if opt.Version {
		fmt.Println(ver)
		return
	}

	err = pkg.Load()
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}

	if opt.ListService {
		fmt.Println("support service list:\n    service\t\tsource\taliases\n	---------------\t\t------")
		for k, s := range pkg.Services.Plugins {
			fmt.Printf("    %15s\t\t%s\t%v\n", k, s.Source, strings.Join(s.Alias, ","))
		}
		return
	}

	if err = opt.Validate(); err != nil {
		logs.Log.Error(err.Error())
		return
	}

	if opt.Debug {
		logs.Log.SetLevel(logs.Debug)
	}

	runner, err := opt.Prepare()
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}

	runner.Run()
}
