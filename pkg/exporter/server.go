package exporter

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	scw "github.com/scaleway/go-scaleway"
)

// ServerCollector collects metrics about the servers.
type ServerCollector struct {
	client   *scw.ScalewayAPI
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	timeout  time.Duration

	Running  *prometheus.Desc
	Created  *prometheus.Desc
	Modified *prometheus.Desc
}

// NewServerCollector returns a new ServerCollector.
func NewServerCollector(logger log.Logger, client *scw.ScalewayAPI, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *ServerCollector {
	failures.WithLabelValues("server").Add(0)

	labels := []string{"id", "name", "datacenter", "type"}
	return &ServerCollector{
		client:   client,
		logger:   log.With(logger, "collector", "server"),
		failures: failures,
		duration: duration,
		timeout:  timeout,

		Running: prometheus.NewDesc(
			"scw_server_running",
			"If 1 the server is running, 0 otherwise",
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
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Running
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	servers, err := c.client.GetServers(true, 0)
	c.duration.WithLabelValues("server").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch servers",
			"err", err,
		)

		c.failures.WithLabelValues("server").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched servers",
		"count", len(*servers),
	)

	for _, server := range *servers {
		var (
			running float64
			created float64
			updated float64
		)

		labels := []string{
			server.Identifier,
			server.Name,
			server.Location.ZoneID,
			server.CommercialType,
		}

		if server.State == "running" {
			running = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Running,
			prometheus.GaugeValue,
			running,
			labels...,
		)

		if num, err := time.Parse("2006-01-02T15:04:05.000000-07:00", server.CreationDate); err == nil {
			created = float64(num.Unix())
		} else {
			level.Error(c.logger).Log(
				"msg", "Failed to parse creation time",
				"err", err,
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			created,
			labels...,
		)

		if num, err := time.Parse("2006-01-02T15:04:05.000000-07:00", server.ModificationDate); err == nil {
			updated = float64(num.Unix())
		} else {
			level.Error(c.logger).Log(
				"msg", "Failed to parse modification time",
				"err", err,
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.Modified,
			prometheus.GaugeValue,
			updated,
			labels...,
		)
	}
}
