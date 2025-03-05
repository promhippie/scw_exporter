package command

import (
	"fmt"
	"os"
	"time"

	"github.com/promhippie/scw_exporter/pkg/action"
	"github.com/promhippie/scw_exporter/pkg/config"
	"github.com/promhippie/scw_exporter/pkg/version"
	"github.com/urfave/cli/v2"
)

// Run parses the command line arguments and executes the program.
func Run() error {
	cfg := config.Load()

	app := &cli.App{
		Name:    "scw_exporter",
		Version: version.String,
		Usage:   "Scaleway Exporter",
		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas@webhippie.de",
			},
		},
		Flags: RootFlags(cfg),
		Commands: []*cli.Command{
			Health(cfg),
		},
		Action: func(_ *cli.Context) error {
			logger := setupLogger(cfg)

			if cfg.Target.AccessKey == "" {
				logger.Error("Missing required scw.access-key")
				return fmt.Errorf("missing required scw.access-key")
			}

			if cfg.Target.SecretKey == "" {
				logger.Error("Missing required scw.secret-key")
				return fmt.Errorf("missing required scw.secret-key")
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

	return app.Run(os.Args)
}

// RootFlags defines the available root flags.
func RootFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
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
			Value:       "0.0.0.0:9503",
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
		&cli.BoolFlag{
			Name:        "web.debug",
			Value:       false,
			Usage:       "Enable pprof debugging for server",
			EnvVars:     []string{"SCW_EXPORTER_WEB_PPROF"},
			Destination: &cfg.Server.Pprof,
		},
		&cli.DurationFlag{
			Name:        "web.timeout",
			Value:       10 * time.Second,
			Usage:       "Server metrics endpoint timeout",
			EnvVars:     []string{"SCW_EXPORTER_WEB_TIMEOUT"},
			Destination: &cfg.Server.Timeout,
		},
		&cli.StringFlag{
			Name:        "web.config",
			Value:       "",
			Usage:       "Path to web-config file",
			EnvVars:     []string{"SCW_EXPORTER_WEB_CONFIG"},
			Destination: &cfg.Server.Web,
		},
		&cli.DurationFlag{
			Name:        "request.timeout",
			Value:       5 * time.Second,
			Usage:       "Request timeout as duration",
			EnvVars:     []string{"SCW_EXPORTER_REQUEST_TIMEOUT"},
			Destination: &cfg.Target.Timeout,
		},
		&cli.StringFlag{
			Name:        "scw.access-key",
			Value:       "",
			Usage:       "Access key for the Scaleway API",
			EnvVars:     []string{"SCW_EXPORTER_ACESS_KEY"},
			Destination: &cfg.Target.AccessKey,
		},
		&cli.StringFlag{
			Name:        "scw.secret-key",
			Value:       "",
			Usage:       "Secret key for the Scaleway API",
			EnvVars:     []string{"SCW_EXPORTER_SECRET_KEY"},
			Destination: &cfg.Target.SecretKey,
		},
		&cli.StringFlag{
			Name:        "scw.org",
			Value:       "",
			Usage:       "Organization for the Scaleway API",
			EnvVars:     []string{"SCW_EXPORTER_ORG"},
			Destination: &cfg.Target.Org,
		},
		&cli.StringFlag{
			Name:        "scw.project",
			Value:       "",
			Usage:       "Project for the Scaleway API",
			EnvVars:     []string{"SCW_EXPORTER_PROJECT"},
			Destination: &cfg.Target.Project,
		},
		&cli.StringFlag{
			Name:        "scw.region",
			Value:       "",
			Usage:       "Region for the Scaleway API",
			EnvVars:     []string{"SCW_EXPORTER_REGION"},
			Destination: &cfg.Target.Region,
		},
		&cli.StringFlag{
			Name:        "scw.zone",
			Value:       "",
			Usage:       "Zone for the Scaleway API",
			EnvVars:     []string{"SCW_EXPORTER_ZONE"},
			Destination: &cfg.Target.Zone,
		},
		&cli.BoolFlag{
			Name:        "collector.dashboard",
			Value:       true,
			Usage:       "Enable collector for dashboard",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_DASHBOARD"},
			Destination: &cfg.Collector.Dashboard,
		},
		&cli.BoolFlag{
			Name:        "collector.consumption",
			Value:       true,
			Usage:       "Enable collector for billing consumption",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_CONSUMPTION"},
			Destination: &cfg.Collector.Consumption,
		},
		&cli.BoolFlag{
			Name:        "collector.security-groups",
			Value:       true,
			Usage:       "Enable collector for security groups",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_SECURITY_GROUPS"},
			Destination: &cfg.Collector.SecurityGroups,
		},
		&cli.BoolFlag{
			Name:        "collector.servers",
			Value:       true,
			Usage:       "Enable collector for servers",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_SERVERS"},
			Destination: &cfg.Collector.Servers,
		},
		&cli.BoolFlag{
			Name:        "collector.baremetal",
			Value:       true,
			Usage:       "Enable collector for servers",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_BAREMETAL"},
			Destination: &cfg.Collector.Baremetal,
		},
		&cli.BoolFlag{
			Name:        "collector.snapshots",
			Value:       true,
			Usage:       "Enable collector for snapshots",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_SNAPSHOTS"},
			Destination: &cfg.Collector.Snapshots,
		},
		&cli.BoolFlag{
			Name:        "collector.volumes",
			Value:       true,
			Usage:       "Enable collector for volumes",
			EnvVars:     []string{"SCW_EXPORTER_COLLECTOR_VOLUMES"},
			Destination: &cfg.Collector.Volumes,
		},
	}
}
