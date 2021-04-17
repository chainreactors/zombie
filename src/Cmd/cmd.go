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
		StringFlag("userdict", "U", "", ""),
		StringFlag("passdict", "P", "", ""),
		StringFlag("uppair", "UP", "", ""),
		SimpleStringFlag("ip", "", ""),
		SimpleStringFlag("IP", "", ""),
		StringFlag("file", "f", "", ""),
		StringFlag("server", "s", "", ""),
		BoolSimpleFlag("ssl", ""),
		IntSimpleFlag("timeout", 2, ""),
		IntFlag("thread", "t", "", 60),
		BoolFlag("simple", "e", false, ""),
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

func BoolSimpleFlag(name, usage string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func BoolFlag(name, alases string, value bool, usage string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    name,
		Aliases: []string{alases},
		Value:   value,
		Usage:   usage,
	}
}

func IntSimpleFlag(name string, value int, usage string) *cli.IntFlag {
	return &cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func IntFlag(name, alases, usage string, value int) *cli.IntFlag {
	return &cli.IntFlag{
		Name:    name,
		Aliases: []string{alases},
		Value:   value,
		Usage:   usage,
	}
}
