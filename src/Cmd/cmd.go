package Cmd

import (
	"Zombie/src/Moudle"
	"github.com/urfave/cli/v2"
)

var Brute = cli.Command{
	Name:    "Brute",
	Action:  Moudle.Brute,
	Aliases: []string{"B"},
	Flags: []cli.Flag{
		StringFlag("username", "u", "", ""),
		StringFlag("password", "p", "", ""),
		SimpleStringFlag("ip", "", ""),
		StringFlag("server", "s", "", ""),
		IntFlag("t", 2, ""),
	},
}

var Exec = cli.Command{
	Name:    "Exec",
	Action:  Moudle.Exec,
	Aliases: []string{"E"},
	Flags: []cli.Flag{
		StringFlag("username", "u", "", ""),
		StringFlag("password", "p", "", ""),
		SimpleStringFlag("ip", "", ""),
		StringFlag("server", "s", "", ""),
		StringFlag("query", "q", "", ""),
		IntFlag("t", 2, ""),
	},
}

func StringFlag(name, alases, value, usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    name,
		Aliases: []string{alases},
		Value:   value,
		Usage:   usage,
	}
}

func SimpleStringFlag(name, value, usage string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func BoolFlag(name, usage string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func IntFlag(name string, value int, usage string) *cli.IntFlag {
	return &cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
