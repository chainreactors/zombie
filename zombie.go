package main

import (
	"Zombie/v1/cmd"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Zombie"
	app.Authors = []*cli.Author{
		{Name: "U2"},
	}
	app.Version = "1.2.0-gtBeta"
	app.Usage = "None"
	app.Commands = []*cli.Command{&cmd.BruteCli, &cmd.ExecCli}
	app.Run(os.Args)
}
