package action

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/promhippie/scw_exporter/pkg/config"
	"github.com/promhippie/scw_exporter/pkg/exporter"
	"github.com/promhippie/scw_exporter/pkg/middleware"
	"github.com/promhippie/scw_exporter/pkg/version"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Server handles the server sub-command.
func Server(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("Launching Scaleway Exporter",
		"version", version.String,
		"revision", version.Revision,
		"date", version.Date,
		"go", version.Go,
	)

	accessKey, err := config.Value(cfg.Target.AccessKey)

	if err != nil {
		logger.Error("Failed to load access key from file",
			"err", err,
		)

		return err
	}

	secretKey, err := config.Value(cfg.Target.SecretKey)

	if err != nil {
		logger.Error("Failed to load secret key from file",
			"err", err,
		)

		return err
	}

	opts := []scw.ClientOption{
		scw.WithAuth(
			accessKey,
			secretKey,
		),
		scw.WithDefaultPageSize(
			100,
		),
		scw.WithUserAgent(
			fmt.Sprintf(
				"scw_exporter/%s",
				version.String,
			),
		),
	}

	if cfg.Target.Org != "" {
		opts = append(opts, scw.WithDefaultOrganizationID(
			cfg.Target.Org,
		))
	}

	if cfg.Target.Project != "" {
		opts = append(opts, scw.WithDefaultProjectID(
			cfg.Target.Project,
		))
	}

	if cfg.Target.Region != "" {
		region, err := scw.ParseRegion(
			cfg.Target.Region,
		)

		if err != nil {
			logger.Error("Failed to parse region",
				"err", err,
			)

			return err
		}

		opts = append(opts, scw.WithDefaultRegion(
			region,
		))
	}

	if cfg.Target.Zone != "" {
		zone, err := scw.ParseZone(
			cfg.Target.Zone,
		)

		if err != nil {
			logger.Error("Failed to parse zone",
				"err", err,
			)

			return err
		}

		opts = append(opts, scw.WithDefaultZone(
			zone,
		))
	}

	client, err := scw.NewClient(
		opts...,
	)

	if err != nil {
		logger.Error("Failed to initialize Scaleway client",
			"err", err,
		)

		return err
	}

	var gr run.Group

	{
		server := &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      handler(cfg, logger, client),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: cfg.Server.Timeout,
		}

		gr.Add(func() error {
			logger.Info("Starting metrics server",
				"address", cfg.Server.Addr,
			)

			return web.ListenAndServe(
				server,
				&web.FlagConfig{
					WebListenAddresses: sliceP([]string{cfg.Server.Addr}),
					WebSystemdSocket:   boolP(false),
					WebConfigFile:      stringP(cfg.Server.Web),
				},
				logger,
			)
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown metrics gracefully",
					"err", err,
				)

				return
			}

			logger.Info("Metrics shutdown gracefully",
				"reason", reason,
			)
		})
	}

	{
		stop := make(chan os.Signal, 1)

		gr.Add(func() error {
			signal.Notify(stop, os.Interrupt)

			<-stop

			return nil
		}, func(_ error) {
			close(stop)
		})
	}

	return gr.Run()
}

func handler(cfg *config.Config, logger *slog.Logger, client *scw.Client) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer(logger))
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Timeout)
	mux.Use(middleware.Cache)

	if cfg.Server.Pprof {
		mux.Mount("/debug", middleware.Profiler())
	}

	if cfg.Collector.Consumption {
		logger.Debug("Consumption collector registered")

		registry.MustRegister(exporter.NewConsumptionCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Dashboard {
		logger.Debug("Dashboard collector registered")

		registry.MustRegister(exporter.NewDashboardCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.SecurityGroups {
		logger.Debug("Security group collector registered")

		registry.MustRegister(exporter.NewSecurityGroupCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Servers {
		logger.Debug("Server collector registered")

		registry.MustRegister(exporter.NewServerCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Snapshots {
		logger.Debug("Snaptshot collector registered")

		registry.MustRegister(exporter.NewSnapshotCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Volumes {
		logger.Debug("Volume collector registered")

		registry.MustRegister(exporter.NewVolumeCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	reg := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			ErrorLog: promLogger{logger},
		},
	)

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.Server.Path, http.StatusMovedPermanently)
	})

	mux.Route("/", func(root chi.Router) {
		root.Get(cfg.Server.Path, func(w http.ResponseWriter, r *http.Request) {
			reg.ServeHTTP(w, r)
		})

		root.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})

		root.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})
	})

	return mux
}
