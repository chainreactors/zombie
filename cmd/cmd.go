package cmd

import (
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/internal/core"
	"github.com/jessevdk/go-flags"
)

func Zombie() {
	var opt core.Option
	_, err := flags.Parse(&opt)
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
