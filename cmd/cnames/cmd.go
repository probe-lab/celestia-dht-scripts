package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v3"

	"github.com/probe-lab/celestia-dht-scripts/dht"
)

const (
	flagCategoryLogging = "Logging Configuration:"
)

var rootConfig = &dht.RootConfig{
	LogLevel:  dht.DefaultLogLevel,
	LogFormat: dht.DefaultLogFormat,
}

var app = &cli.Command{
	Name:                  "cnames",
	Usage:                 "A Celestia's DHT namespace scrapper",
	EnableShellCompletion: true,
	Flags:                 rootFlags,
	Before:                rootBefore,
	Commands: []*cli.Command{
		cmdLookup,
		cmdCrawl,
		cmdDHTKeys,
	},
	After: rootAfter,
}

var rootFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "log.level",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_LOG_LEVEL")},
		},
		Usage:       "sets an explicity logging level: debug, info, warn, error. Takes precedence over the verbose flag.",
		Destination: &rootConfig.LogLevel,
		Value:       rootConfig.LogLevel,
		Category:    flagCategoryLogging,
	},
	&cli.StringFlag{
		Name: "log.format",
		Sources: cli.ValueSourceChain{
			Chain: []cli.ValueSource{cli.EnvVar("CNAMES_LOG_FORMAT")},
		},
		Usage:       "sets the format to output the log statements in: text, json",
		Destination: &rootConfig.LogFormat,
		Value:       rootConfig.LogFormat,
		Category:    flagCategoryLogging,
	},
}

func main() {
	sigs := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		defer cancel()
		defer signal.Stop(sigs)

		select {
		case <-ctx.Done():
		case sig := <-sigs:
			log.WithField("signal", sig.String()).Info("Received termination signal - Stopping...")
		}
	}()

	if err := app.Run(ctx, os.Args); err != nil && !errors.Is(err, context.Canceled) {
		log.Errorf("ookla terminated abnormally %s", err.Error())
		os.Exit(1)
	}
}

func rootBefore(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// don't set up anything if ookla is run without arguments
	if cmd.NArg() == 0 {
		return ctx, nil
	}

	// read CLI args and configure the global logger
	if err := configureLogger(ctx, cmd); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func rootAfter(c context.Context, cmd *cli.Command) error {
	log.Info("dht script successfully finished")
	return nil
}

// configureLogger configures the global logger based on the provided CLI
// context. It sets the log level based on the "--log-level" flag or the
// "--verbose" flag. The log format is determined by the "--log.format" flag.
// The function returns an error if the log level or log format is not supported.
// Possible log formats include "tint", "hlog", "text", and "json". The default
// logger is overwritten with the configured logger.
func configureLogger(_ context.Context, cmd *cli.Command) error {
	// log level
	logLevel := log.InfoLevel
	if cmd.IsSet("log.level") {
		switch strings.ToLower(rootConfig.LogLevel) {
		case "debug":
			logLevel = log.DebugLevel
		case "info":
			logLevel = log.InfoLevel
		case "warn":
			logLevel = log.WarnLevel
		case "error":
			logLevel = log.ErrorLevel
		default:
			return fmt.Errorf("unknown log level: %s", rootConfig.LogLevel)
		}
	}
	log.SetLevel(log.Level(logLevel))

	// log format
	switch strings.ToLower(rootConfig.LogFormat) {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("unknown log format: %q", rootConfig.LogFormat)
	}

	return nil
}
