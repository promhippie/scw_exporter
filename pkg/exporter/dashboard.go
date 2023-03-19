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

// DashboardCollector collects metrics about the dashboard.
type DashboardCollector struct {
	client   *scw.Client
	instance *instance.API
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target
	org      *string
	project  *string

	Volumes         *prometheus.Desc
	Running         *prometheus.Desc
	Images          *prometheus.Desc
	Snapshots       *prometheus.Desc
	Servers         *prometheus.Desc
	IPs             *prometheus.Desc
	SecurityGroups  *prometheus.Desc
	UnusedIPs       *prometheus.Desc
	VolumesLSSD     *prometheus.Desc
	VolumesLSSDSize *prometheus.Desc
	VolumesBSSD     *prometheus.Desc
	VolumesBSSDSize *prometheus.Desc
	PrivateNics     *prometheus.Desc
	PlacementGroups *prometheus.Desc
	Types           *prometheus.Desc
}

// NewDashboardCollector returns a new DashboardCollector.
func NewDashboardCollector(logger log.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *DashboardCollector {
	if failures != nil {
		failures.WithLabelValues("dashboard").Add(0)
	}

	labels := []string{"zone"}
	collector := &DashboardCollector{
		client:   client,
		instance: instance.NewAPI(client),
		logger:   log.With(logger, "collector", "dashboard"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Volumes: prometheus.NewDesc(
			"scw_dashboard_volumes_count",
			"Count of used volumes",
			labels,
			nil,
		),
		Running: prometheus.NewDesc(
			"scw_dashboard_running_servers",
			"Count of running servers",
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
		Servers: prometheus.NewDesc(
			"scw_dashboard_servers_count",
			"Count of owned servers",
			labels,
			nil,
		),
		IPs: prometheus.NewDesc(
			"scw_dashboard_ips_count",
			"Count of used IPs",
			labels,
			nil,
		),
		SecurityGroups: prometheus.NewDesc(
			"scw_dashboard_security_groups_count",
			"Count of security groups",
			labels,
			nil,
		),
		UnusedIPs: prometheus.NewDesc(
			"scw_dashboard_unused_ips_count",
			"Count of unused IPs",
			labels,
			nil,
		),
		VolumesLSSD: prometheus.NewDesc(
			"scw_dashboard_volumes_lssd_count",
			"Count of unused IPs",
			labels,
			nil,
		),
		VolumesLSSDSize: prometheus.NewDesc(
			"scw_dashboard_volumes_lssd_total_size",
			"Count of unused IPs",
			labels,
			nil,
		),
		VolumesBSSD: prometheus.NewDesc(
			"scw_dashboard_volumes_bssd_count",
			"Count of unused IPs",
			labels,
			nil,
		),
		VolumesBSSDSize: prometheus.NewDesc(
			"scw_dashboard_volumes_bssd_total_size",
			"Count of unused IPs",
			labels,
			nil,
		),
		PrivateNics: prometheus.NewDesc(
			"scw_dashboard_private_nics_count",
			"Count of private nics",
			labels,
			nil,
		),
		PlacementGroups: prometheus.NewDesc(
			"scw_dashboard_placement_groups_count",
			"Count of placement groups",
			labels,
			nil,
		),
		Types: prometheus.NewDesc(
			"scw_dashboard_server_types_count",
			"Count of servers by type",
			append(labels, "type"),
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
func (c *DashboardCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Volumes,
		c.Running,
		c.Images,
		c.Snapshots,
		c.Servers,
		c.IPs,
		c.SecurityGroups,
		c.UnusedIPs,
		c.VolumesLSSD,
		c.VolumesLSSDSize,
		c.VolumesBSSD,
		c.VolumesBSSDSize,
		c.PrivateNics,
		c.PlacementGroups,
		c.Types,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *DashboardCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Volumes
	ch <- c.Running
	ch <- c.Images
	ch <- c.Snapshots
	ch <- c.Servers
	ch <- c.IPs
	ch <- c.SecurityGroups
	ch <- c.UnusedIPs
	ch <- c.VolumesLSSD
	ch <- c.VolumesLSSDSize
	ch <- c.VolumesBSSD
	ch <- c.VolumesBSSDSize
	ch <- c.PrivateNics
	ch <- c.PlacementGroups
	ch <- c.Types
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *DashboardCollector) Collect(ch chan<- prometheus.Metric) {
	for _, zone := range affectedZones(c.client) {
		now := time.Now()
		resp, err := c.instance.GetDashboard(&instance.GetDashboardRequest{
			Zone:         zone,
			Organization: c.org,
			Project:      c.project,
		})
		c.duration.WithLabelValues("dashboard").Observe(time.Since(now).Seconds())

		if err != nil {
			level.Error(c.logger).Log(
				"msg", "Failed to fetch dashboard",
				"zone", zone,
				"err", err,
			)

			c.failures.WithLabelValues("dashboard").Inc()
			return
		}

		dashboard := resp.Dashboard
		labels := []string{
			zone.String(),
		}

		ch <- prometheus.MustNewConstMetric(
			c.Volumes,
			prometheus.GaugeValue,
			float64(dashboard.VolumesCount),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Running,
			prometheus.GaugeValue,
			float64(dashboard.RunningServersCount),
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
			c.Servers,
			prometheus.GaugeValue,
			float64(dashboard.ServersCount),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IPs,
			prometheus.GaugeValue,
			float64(dashboard.IPsCount),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SecurityGroups,
			prometheus.GaugeValue,
			float64(dashboard.SecurityGroupsCount),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.UnusedIPs,
			prometheus.GaugeValue,
			float64(dashboard.IPsUnused),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VolumesLSSD,
			prometheus.GaugeValue,
			float64(dashboard.VolumesLSSDCount),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.VolumesLSSDSize,
			prometheus.GaugeValue,
			float64(dashboard.VolumesLSSDTotalSize),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.VolumesBSSD,
			prometheus.GaugeValue,
			float64(dashboard.VolumesBSSDCount),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.VolumesBSSDSize,
			prometheus.GaugeValue,
			float64(dashboard.VolumesBSSDTotalSize),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PrivateNics,
			prometheus.GaugeValue,
			float64(dashboard.PrivateNicsCount),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PlacementGroups,
			prometheus.GaugeValue,
			float64(dashboard.PlacementGroupsCount),
			labels...,
		)

		for kind, count := range dashboard.ServersByTypes {
			ch <- prometheus.MustNewConstMetric(
				c.Types,
				prometheus.GaugeValue,
				float64(count),
				append(labels, kind)...,
			)
		}
	}
}
