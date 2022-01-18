package main

import (
	"Zombie/src/Cmd"
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
	app.Commands = []*cli.Command{&Cmd.Brute, &Cmd.Query}
	app.Run(os.Args)
}
