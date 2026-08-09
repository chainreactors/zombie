package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils/fileutils"
	"github.com/chainreactors/zombie/pkg"
	"github.com/jessevdk/go-flags"
)

// RunOptions configures the reusable, no-exit zombie entrypoint.
type RunOptions struct {
	Output    io.Writer
	Version   string
	ProxyDial pkg.DialFunc
	OnResult  ResultHandler
}

func Help() string {
	var opt Option
	parser := flags.NewParser(&opt, flags.Default&^flags.PrintErrors)
	parser.Usage = Usage()
	var buf bytes.Buffer
	parser.WriteHelp(&buf)
	return buf.String()
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
	registerBuiltinServices()

	if opt.ListService {
		fmt.Fprintln(output, "support service list:\n    service\t\tsource\taliases\n\t---------------\t\t------")
		for k, s := range pkg.Services.All() {
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

	var outputFile *fileutils.File
	if opt.OutputFile != "" {
		outputFile, err = fileutils.NewFile(opt.OutputFile, fileutils.ModeAppend, false, false)
		if err != nil {
			return err
		}
		defer outputFile.Close()
	}
	runner.OnResult = cliResultHandler(opt.OutputFormat, opt.FileFormat, outputFile, opts.OnResult)
	if opts.ProxyDial != nil {
		runner.ProxyDial = opts.ProxyDial
	}
	return runner.RunWithContext(ctx)
}

func cliResultHandler(outputFormat, fileFormat string, outputFile *fileutils.File, next ResultHandler) ResultHandler {
	return func(result *pkg.Result) {
		if result.OK {
			if outputFile != nil {
				if err := outputFile.SyncWrite(result.Format(fileFormat)); err != nil {
					logs.Log.Warnf("write output file failed: %v", err)
				}
			}
			logs.Log.Console(result.Format(outputFormat))
		} else {
			errMsg := "unknown error"
			if result.Err != nil {
				errMsg = result.Err.Error()
			}
			logs.Log.Debugf("[%s] %s %s %s ,%s login failed, %s", result.Mod.String(), result.URI(), result.Username, result.Password, result.Service, errMsg)
		}
		if next != nil {
			next(result)
		}
	}
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
