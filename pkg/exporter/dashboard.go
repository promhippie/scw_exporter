package exporter

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	scw "github.com/scaleway/go-scaleway"
)

// DashboardCollector collects metrics about the dashboard.
type DashboardCollector struct {
	client   *scw.ScalewayAPI
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	timeout  time.Duration

	Running   *prometheus.Desc
	Servers   *prometheus.Desc
	Volumes   *prometheus.Desc
	Images    *prometheus.Desc
	Snapshots *prometheus.Desc
	IPs       *prometheus.Desc
}

// NewDashboardCollector returns a new DashboardCollector.
func NewDashboardCollector(logger log.Logger, client *scw.ScalewayAPI, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *DashboardCollector {
	failures.WithLabelValues("dashboard").Add(0)

	labels := []string{}
	return &DashboardCollector{
		client:   client,
		logger:   log.With(logger, "collector", "dashboard"),
		failures: failures,
		duration: duration,
		timeout:  timeout,

		Running: prometheus.NewDesc(
			"scw_dashboard_running_servers",
			"Count of running servers",
			labels,
			nil,
		),
		Servers: prometheus.NewDesc(
			"scw_dashboard_servers_count",
			"Count of owned servers",
			labels,
			nil,
		),
		Volumes: prometheus.NewDesc(
			"scw_dashboard_volumes_count",
			"Count of used volumes",
			labels,
			nil,
		),
		Images: prometheus.NewDesc(
			"scw_dashboard_images_count",
			"Count of used images",
			labels,
			nil,
		),
		Snapshots: prometheus.NewDesc(
			"scw_dashboard_snapshots_count",
			"Count of used snapshots",
			labels,
			nil,
		),
		IPs: prometheus.NewDesc(
			"scw_dashboard_ips_count",
			"Count of used IPs",
			labels,
			nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *DashboardCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Running
	ch <- c.Servers
	ch <- c.Volumes
	ch <- c.Images
	ch <- c.Snapshots
	ch <- c.IPs
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *DashboardCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	dashboard, err := c.client.GetDashboard()
	c.duration.WithLabelValues("dashboard").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch dashboard",
			"err", err,
		)

		c.failures.WithLabelValues("dashboard").Inc()
		return
	}

	labels := []string{}

	ch <- prometheus.MustNewConstMetric(
		c.Running,
		prometheus.GaugeValue,
		float64(dashboard.RunningServersCount),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Servers,
		prometheus.GaugeValue,
		float64(dashboard.ServersCount),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Volumes,
		prometheus.GaugeValue,
		float64(dashboard.VolumesCount),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Images,
		prometheus.GaugeValue,
		float64(dashboard.ImagesCount),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Snapshots,
		prometheus.GaugeValue,
		float64(dashboard.SnapshotsCount),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.IPs,
		prometheus.GaugeValue,
		float64(dashboard.IPsCount),
		labels...,
	)
}
