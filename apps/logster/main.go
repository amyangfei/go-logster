package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amyangfei/go-logster/logster"
	"github.com/gofrs/flock"
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

	ParserOptions string `long:"parser-options" short:"P" description:"specific parser options"`

	OutputOptions string `long:"output-options" short:"O" description:"specific output options"`

	ParseInfo struct {
		ParserPlugin string
		LogFile      string
	} `positional-args:"yes" required:"yes" description:"parser plugin name and log file"`
}

func logErrorAndExit(logger zerolog.Logger, err error) {
	logger.Error().Msg(err.Error())
	os.Exit(1)
}

func process(logger zerolog.Logger) {
	baseName := strings.Replace(opts.ParseInfo.ParserPlugin+"-"+opts.ParseInfo.LogFile, "/", "-", -1)
	stateFile := filepath.Join(opts.StateDir, baseName+".state")
	lockFile := filepath.Join(opts.StateDir, baseName+".lock")

	logger.Debug().Msgf("State file %s, lock file %s", stateFile, lockFile)
	logger.Info().Msgf("Executing parser %s on logfile %s", opts.ParseInfo.ParserPlugin, opts.ParseInfo.LogFile)

	fileLock := flock.New(lockFile)
	locked, err := fileLock.TryLock()

	if err != nil {
		logErrorAndExit(logger, err)
	} else if !locked {
		logErrorAndExit(logger, fmt.Errorf("Failed to acquire file lock"))
	}
	defer fileLock.Unlock()

	tailer := logster.LogtailTailer{
		Logfile:   opts.ParseInfo.LogFile,
		Statefile: stateFile,
		Binary:    logster.DefaultLogtailPath,
	}
	var duration time.Duration
	info, err := os.Stat(stateFile)
	if err != nil {
		err2 := tailer.CreateStateFile()
		if err2 != nil {
			logErrorAndExit(logger, err2)
		}
		duration = time.Second
	} else {
		mtime := info.ModTime()
		duration = time.Since(mtime)
	}
	logger.Debug().Msgf("Setting duration to %s seconds", duration)

	parser, err := logster.LoadParserPlugin(opts.ParseInfo.ParserPlugin)
	if err != nil {
		logErrorAndExit(logger, err)
	}
	parser.Init(opts.ParserOptions)

	c := make(chan string)
	go tailer.ReadLines(c)
	for line := range c {
		parser.ParseLine(line)
	}

	if metrics, err := parser.GetState(duration.Seconds()); err != nil {
		logErrorAndExit(logger, err)
	} else {
		for _, pluginPath := range opts.Output {
			output, err := logster.LoadOutputPlugin(pluginPath)
			if err != nil {
				logErrorAndExit(logger, err)
			}
			if err := output.Init(opts.MetricPrefix, opts.MetricSuffix,
				opts.OutputOptions, opts.DryRun, logger); err != nil {
				logErrorAndExit(logger, err)
			}
			if err := output.Submit(metrics); err != nil {
				logErrorAndExit(logger, err)
			}
		}
	}
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
			fmt.Println(logster.ReleaseInfo(binName))
			return
		}
		fmt.Println(err)
		return
	}

	if opts.ShowVersion {
		fmt.Println(logster.ReleaseInfo(binName))
		return
	}

	logger := initLogger()

	process(logger)
}
