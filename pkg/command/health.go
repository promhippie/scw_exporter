package command

import (
	"context"
	"fmt"
	"net/http"

	"github.com/promhippie/scw_exporter/pkg/config"
	"github.com/urfave/cli/v3"
)

// Health provides the sub-command to perform a health check.
func Health(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "health",
		Usage: "Perform health checks",
		Flags: HealthFlags(cfg),
		Action: func(_ context.Context, _ *cli.Command) error {
			logger := setupLogger(cfg)

			resp, err := http.Get(
				fmt.Sprintf(
					"http://%s/healthz",
					cfg.Server.Addr,
				),
			)

			if err != nil {
				logger.Error("Failed to request health check",
					"err", err,
				)

				return err
			}

			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != 200 {
				logger.Error("Health check seems to be in bad state",
					"err", err,
					"code", resp.StatusCode,
				)

				return err
			}

			logger.Debug("Health check seems to be fine",
				"code", resp.StatusCode,
			)

			return nil
		},
	}
}

// HealthFlags defines the available health flags.
func HealthFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "web.address",
			Value:       "0.0.0.0:9503",
			Usage:       "Address to bind the metrics server",
			Sources:     cli.EnvVars("SCW_EXPORTER_WEB_ADDRESS"),
			Destination: &cfg.Server.Addr,
		},
	}
}
