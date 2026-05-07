package cmd

import (
	"context"
	"fmt"
	"github.com/chainreactors/zombie/core"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var ver = "dev"

func Zombie() {
	defer os.Exit(0)
	err := core.RunWithArgs(context.Background(), os.Args[1:], core.RunOptions{Version: ver})
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); !ok || flagsErr.Type != flags.ErrHelp {
			fmt.Println(err.Error())
		}
		return
	}
}
