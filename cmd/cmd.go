package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chainreactors/zombie/core"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var ver = "dev"

func Zombie() {
	os.Exit(Run(os.Args[1:], os.Stdout))
}

func Run(args []string, output io.Writer) int {
	if output == nil {
		output = os.Stdout
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := core.RunWithArgs(ctx, args, core.RunOptions{Output: output, Version: ver})
	if err != nil {
		fmt.Fprintln(output, err.Error())
		return 1
	}
	return 0
}
