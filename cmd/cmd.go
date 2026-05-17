package cmd

import (
	"context"
	"fmt"
	"github.com/chainreactors/zombie/core"
	"io"
	"io/ioutil"
	"log"
	"os"
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
	err := core.RunWithArgs(context.Background(), args, core.RunOptions{Output: output, Version: ver})
	if err != nil {
		fmt.Fprintln(output, err.Error())
		return 1
	}
	return 0
}
