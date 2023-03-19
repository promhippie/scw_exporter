package exporter

import (
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/scw_exporter/pkg/config"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ServerCollector collects metrics about the servers.
type ServerCollector struct {
	client   *scw.Client
	instance *instance.API
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target
	org      *string
	project  *string

	State           *prometheus.Desc
	VolumeCount     *prometheus.Desc
	PrivateNicCount *prometheus.Desc
	Created         *prometheus.Desc
	Modified        *prometheus.Desc
}

// NewServerCollector returns a new ServerCollector.
func NewServerCollector(logger log.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *ServerCollector {
	if failures != nil {
		failures.WithLabelValues("server").Add(0)
	}

	labels := []string{"id", "name", "zone", "org", "project", "type", "private_ip", "public_ip", "arch"}
	collector := &ServerCollector{
		client:   client,
		instance: instance.NewAPI(client),
		logger:   log.With(logger, "collector", "server"),
		failures: failures,
		duration: duration,
		config:   cfg,

		State: prometheus.NewDesc(
			"scw_server_state",
			"If 1 the server is running, depending on the state otherwise",
			labels,
			nil,
		),
		VolumeCount: prometheus.NewDesc(
			"scw_server_volume_count",
			"Number of volumes attached",
			labels,
			nil,
		),
		PrivateNicCount: prometheus.NewDesc(
			"scw_server_private_nic_count",
			"Number of private nics attached",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"scw_server_created_timestamp",
			"Timestamp when the server have been created",
			labels,
			nil,
		),
		Modified: prometheus.NewDesc(
			"scw_server_modified_timestamp",
			"Timestamp when the server have been modified",
			labels,
			nil,
		),
	}

	if cfg.Org != "" {
		collector.org = scw.StringPtr(cfg.Org)
	}

	if cfg.Project != "" {
		collector.project = scw.StringPtr(cfg.Project)
	}

	return collector
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *ServerCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.State,
		c.VolumeCount,
		c.PrivateNicCount,
		c.Created,
		c.Modified,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.State
	ch <- c.VolumeCount
	ch <- c.PrivateNicCount
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	for _, zone := range affectedZones(c.client) {
		now := time.Now()
		resp, err := c.instance.ListServers(&instance.ListServersRequest{
			Zone:         zone,
			Organization: c.org,
			Project:      c.project,
		}, scw.WithAllPages())
		c.duration.WithLabelValues("server").Observe(time.Since(now).Seconds())

		if err != nil {
			level.Error(c.logger).Log(
				"msg", "Failed to fetch servers",
				"zone", zone,
				"err", err,
			)

			c.failures.WithLabelValues("server").Inc()
			return
		}

		level.Debug(c.logger).Log(
			"msg", "Fetched servers",
			"zone", zone,
			"count", resp.TotalCount,
		)

		for _, server := range resp.Servers {
			var (
				privateIP string
				publicIP  string

				state float64
			)

			if server.PrivateIP != nil {
				privateIP = *server.PrivateIP
			}

			if server.PublicIP != nil {
				publicIP = server.PublicIP.Address.String()
			}

			labels := []string{
				server.ID,
				server.Name,
				server.Zone.String(),
				server.Organization,
				server.Project,
				server.CommercialType,
				privateIP,
				publicIP,
				server.Arch.String(),
			}

			switch val := server.State; val {
			case instance.ServerStateRunning:
				state = 1.0
			case instance.ServerStateStopped:
				state = 2.0
			case instance.ServerStateStoppedInPlace:
				state = 3.0
			case instance.ServerStateStarting:
				state = 4.0
			case instance.ServerStateStopping:
				state = 5.0
			case instance.ServerStateLocked:
				state = 6.0
			default:
				state = 0.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.State,
				prometheus.GaugeValue,
				state,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.VolumeCount,
				prometheus.GaugeValue,
				float64(len(server.Volumes)),
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.PrivateNicCount,
				prometheus.GaugeValue,
				float64(len(server.PrivateNics)),
				labels...,
			)

			if server.CreationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Created,
					prometheus.GaugeValue,
					float64(server.CreationDate.Unix()),
					labels...,
				)
			}

			if server.ModificationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Modified,
					prometheus.GaugeValue,
					float64(server.ModificationDate.Unix()),
					labels...,
				)
			}
		}
	}
}
