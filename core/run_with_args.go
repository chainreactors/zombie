package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/pkg"
	"github.com/jessevdk/go-flags"
)

// RunOptions configures the reusable, no-exit zombie entrypoint.
type RunOptions struct {
	Output  io.Writer
	Version string
}

func RunWithArgs(ctx context.Context, args []string, opts RunOptions) error {
	if ctx == nil {
		ctx = context.Background()
	}
	var opt Option
	output := opts.Output
	if output == nil {
		output = os.Stdout
	}
	if opts.Output != nil {
		oldLog := logs.Log
		logs.Log = logs.NewLogger(oldLog.Level)
		logs.Log.SetOutput(output)
		defer func() {
			logs.Log = oldLog
		}()
	}

	parser := flags.NewParser(&opt, flags.Default&^flags.PrintErrors)
	parser.Usage = Usage()
	if _, err := parser.ParseArgs(args); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Fprintln(output, err.Error())
			return nil
		}
		return err
	}

	if opt.Version {
		version := opts.Version
		if version == "" {
			version = "dev"
		}
		fmt.Fprintln(output, version)
		return nil
	}

	if err := pkg.Load(); err != nil {
		return err
	}

	if opt.ListService {
		fmt.Fprintln(output, "support service list:\n    service\t\tsource\taliases\n\t---------------\t\t------")
		for k, s := range pkg.Services.Plugins {
			fmt.Fprintf(output, "    %15s\t\t%s\t%v\n", k, s.Source, strings.Join(s.Alias, ","))
		}
		return nil
	}

	if err := opt.Validate(); err != nil {
		return err
	}

	if opt.Debug {
		logs.Log.SetLevel(logs.Level(10))
	} else if opt.Quiet {
		logs.Log.SetLevel(logs.Level(51))
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	runner, err := opt.Prepare()
	if err != nil {
		return err
	}
	return runner.RunWithContext(ctx)
}

func Usage() string {
	return `

    WIKI: https://chainreactors.github.io/wiki/zombie

    QUICKSTART:
        simple example:
            zombie -i 1.1.1.1 -u root -s ssh

        brute multiple ssh targets(ip list):
            zombie -I targets.txt -u root -p password -s ssh

        brute from file and auto parse:
            zombie -I targets.txt
`
}
