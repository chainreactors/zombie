package cmd

import (
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/internal/core"
	"github.com/jessevdk/go-flags"
)

func Zombie() {
	var opts core.Option
	_, err := flags.Parse(&opts)
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}
	runner, err := core.PrepareRunner(&opts)
	if err != nil {
		logs.Log.Error(err.Error())
		return
	}

	runner.Run()
}
