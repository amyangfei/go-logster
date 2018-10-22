package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amyangfei/go-logster/logster"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

const binName = "go-logster"

var opts struct {
	MetricPrefix string `long:"metric-prefix" short:"p" default:"" description:"Add prefix to all published metrics. This is for people that may multiple instances of same service on same host."`

	MetricSuffix string `long:"metric-suffix" short:"x" default:"" description:"Add suffix to all published metrics. This is for people that may add suffix at the end of their metrics."`

	LogDir string `long:"log-dir" short:"l" default:"/var/log/logster" description:"Where to store the logster logfile."`

	StateDir string `long:"state-dir" short:"s" default:"/var/run" description:"Where to store the tailer state file."`

	Output []string `long:"output" short:"o" description:"Where to send metrics (can specify multiple times)" required:"true"`

	ShowVersion bool `long:"version" short:"v" description:"print version"`

	DryRun bool `long:"dry-run" short:"d" description:"Parse the log file but send stats to standard output."`

	Debug bool `long:"debug" short:"D" description:"Provide more verbose logging for debugging."`

	Options string `long:"options" short:"O" description:"specific output options"`

	ParseInfo struct {
		ParserPlugin string
		LogFile      string
	} `positional-args:"yes" required:"yes" description:"parser plugin name and log file"`
}

func process(logger zerolog.Logger) {
	stateFile := filepath.Join(opts.StateDir, opts.ParseInfo.ParserPlugin, opts.ParseInfo.LogFile+".state")
	lockFile := filepath.Join(opts.StateDir, opts.ParseInfo.ParserPlugin, opts.ParseInfo.LogFile+".lock")

	logger.Debug().Msg(fmt.Sprintf("State file %s, lock file %s", stateFile, lockFile))
	logger.Info().Msg(fmt.Sprintf("Executing parser %s on logfile %s", opts.ParseInfo.ParserPlugin, opts.ParseInfo.LogFile))

	// TODO: load plugin
}

func initLogger() zerolog.Logger {
	if opts.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if _, err := os.Stat(opts.LogDir); os.IsNotExist(err) {
		err := os.Mkdir(opts.LogDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	logfileName := filepath.Join(opts.LogDir, "go-logster.log")
	logFile, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return zerolog.New(logFile).With().Timestamp().Logger()
}

func main() {
	args := make([]string, len(os.Args)-1)
	copy(args, os.Args[1:])

	args, err := flags.NewParser(&opts, flags.PassDoubleDash|flags.HelpFlag|flags.IgnoreUnknown).ParseArgs(args)
	if err != nil {
		if opts.ShowVersion {
			fmt.Println(logster.Version(binName))
			return
		}
		fmt.Println(err)
		return
	}

	if opts.ShowVersion {
		fmt.Println(logster.Version(binName))
		return
	}

	logger := initLogger()

	process(logger)
}
