package Cmd

import (
	"Zombie/src/Module"
	"github.com/urfave/cli/v2"
)

var Brute = cli.Command{
	Name:    "Brute",
	Action:  Module.Brute,
	Aliases: []string{"B"},
	Flags: []cli.Flag{
		StringFlag("username", "u", "", ""),
		StringFlag("password", "p", "", ""),
		StringFlag("userdict", "U", "", ""),
		StringFlag("passdict", "P", "", ""),
		StringFlag("uppair", "UP", "", ""),
		StringFlag("instance", "i", "orcl", ""),
		SimpleStringFlag("ip", "", ""),
		SimpleStringFlag("IP", "", ""),
		SimpleStringFlag("gt", "", ""),
		StringFlag("file", "f", "./.res.log", ""),
		StringFlag("server", "s", "", ""),
		IntSimpleFlag("timeout", 2, ""),
		IntFlag("thread", "t", "", 60),
		BoolFlag("simple", "e", true, ""),
		IntSimpleFlag("proc", 0, ""),
		SimpleStringFlag("type", "raw", ""),
	},
}

var Query = cli.Command{
	Name:    "Query",
	Action:  Module.Exec,
	Aliases: []string{"Q"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Value:   "",
			Usage:   "",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Value:   "",
			Usage:   "",
		},
		&cli.StringFlag{
			Name:  "ip",
			Value: "",
			Usage: "",
		},
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Value:   "",
			Usage:   "",
		},
		StringFlag("InputFile", "F", "", ""),
		StringFlag("OutputFile", "f", "./ExecRes.log", ""),
		StringFlag("server", "s", "", ""),
		BoolFlag("auto", "a", false, ""),
		SimpleStringFlag("type", "raw", ""),
		BoolSimpleFlag("more", false, ""),
		IntFlag("thread", "t", "", 60),
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

func BoolSimpleFlag(name string, value bool, usage string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  name,
		Value: value,
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
