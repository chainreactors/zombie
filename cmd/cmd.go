package cmd

import (
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/internal/core"
	"github.com/chainreactors/zombie/pkg"
	"github.com/jessevdk/go-flags"
)

func Zombie() {
	var opt core.Option
	_, err := flags.Parse(&opt)
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}

	if err = opt.Validate(); err != nil {
		logs.Log.Error(err.Error())
		return
	}

	if opt.Debug {
		logs.Log.Level = logs.Debug
	}
	err = pkg.LoadKeyword()
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}
	err = pkg.LoadRules()
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}

	runner, err := opt.Prepare()
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}

	runner.Run()
}
