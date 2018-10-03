package main

import (
	"errors"
	"os"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
	"github.com/promhippie/scw_exporter/pkg/action"
	"github.com/promhippie/scw_exporter/pkg/config"
	"github.com/promhippie/scw_exporter/pkg/version"
	"gopkg.in/urfave/cli.v2"
)

var (
	// ErrMissingScalewayToken defines the error if scw.token is empty.
	ErrMissingScalewayToken = errors.New("Missing required scw.token")

	// ErrMissingScalewayOrg defines the error if scw.org is empty.
	ErrMissingScalewayOrg = errors.New("Missing required scw.org")

	// ErrMissingScalewayRegion defines the error if scw.region is empty.
	ErrMissingScalewayRegion = errors.New("Missing required scw.region")
)

func main() {
	cfg := config.Load()

	if env := os.Getenv("SCW_EXPORTER_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:    "scw_exporter",
		Version: version.Version,
		Usage:   "Scaleway Exporter",
		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas@webhippie.de",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log.level",
				Value:       "info",
				Usage:       "Only log messages with given severity",
				EnvVars:     []string{"SCW_EXPORTER_LOG_LEVEL"},
				Destination: &cfg.Logs.Level,
			},
			&cli.BoolFlag{
				Name:        "log.pretty",
				Value:       false,
				Usage:       "Enable pretty messages for logging",
				EnvVars:     []string{"SCW_EXPORTER_LOG_PRETTY"},
				Destination: &cfg.Logs.Pretty,
			},
			&cli.StringFlag{
				Name:        "web.address",
				Value:       "0.0.0.0:9109",
				Usage:       "Address to bind the metrics server",
				EnvVars:     []string{"SCW_EXPORTER_WEB_ADDRESS"},
				Destination: &cfg.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "web.path",
				Value:       "/metrics",
				Usage:       "Path to bind the metrics server",
				EnvVars:     []string{"SCW_EXPORTER_WEB_PATH"},
				Destination: &cfg.Server.Path,
			},
			&cli.DurationFlag{
				Name:        "request.timeout",
				Value:       5 * time.Second,
				Usage:       "Request timeout as duration",
				EnvVars:     []string{"SCW_EXPORTER_REQUEST_TIMEOUT"},
				Destination: &cfg.Target.Timeout,
			},
			&cli.StringFlag{
				Name:        "scw.token",
				Value:       "",
				Usage:       "Access token for the Scaleway API",
				EnvVars:     []string{"SCW_EXPORTER_TOKEN"},
				Destination: &cfg.Target.Token,
			},
			&cli.StringFlag{
				Name:        "scw.org",
				Value:       "",
				Usage:       "Organization for the Scaleway API",
				EnvVars:     []string{"SCW_EXPORTER_ORG"},
				Destination: &cfg.Target.Org,
			},
			&cli.StringFlag{
				Name:        "scw.region",
				Value:       "",
				Usage:       "Region for the Scaleway API",
				EnvVars:     []string{"SCW_EXPORTER_REGION"},
				Destination: &cfg.Target.Region,
			},
		},
		Action: func(c *cli.Context) error {
			logger := setupLogger(cfg)

			if cfg.Target.Token == "" {
				level.Error(logger).Log(
					"msg", ErrMissingScalewayToken,
				)

				return ErrMissingScalewayToken
			}

			if cfg.Target.Org == "" {
				level.Error(logger).Log(
					"msg", ErrMissingScalewayOrg,
				)

				return ErrMissingScalewayOrg
			}

			if cfg.Target.Region == "" {
				level.Error(logger).Log(
					"msg", ErrMissingScalewayRegion,
				)

				return ErrMissingScalewayRegion
			}

			return action.Server(cfg, logger)
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
