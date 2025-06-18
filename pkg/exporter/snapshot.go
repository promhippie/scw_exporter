package exporter

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/scw_exporter/pkg/config"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// SnapshotCollector collects metrics about the snapshots.
type SnapshotCollector struct {
	client   *scw.Client
	instance *instance.API
	logger   *slog.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target
	org      *string
	project  *string

	Available  *prometheus.Desc
	Size       *prometheus.Desc
	VolumeType *prometheus.Desc
	State      *prometheus.Desc
	Created    *prometheus.Desc
	Modified   *prometheus.Desc
}

// NewSnapshotCollector returns a new SnapshotCollector.
func NewSnapshotCollector(logger *slog.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *SnapshotCollector {
	if failures != nil {
		failures.WithLabelValues("snapshot").Add(0)
	}

	labels := []string{"id", "name", "zone", "org", "project_id", "project_name"}
	collector := &SnapshotCollector{
		client:   client,
		instance: instance.NewAPI(client),
		logger:   logger.With("collector", "snapshot"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Available: prometheus.NewDesc(
			"scw_snapshot_available",
			"Constant value of 1 that this snapshot is available",
			labels,
			nil,
		),
		Size: prometheus.NewDesc(
			"scw_snapshot_size_bytes",
			"Size of the snapshot in bytes",
			labels,
			nil,
		),
		VolumeType: prometheus.NewDesc(
			"scw_snapshot_type",
			"Type of the snapshot",
			labels,
			nil,
		),
		State: prometheus.NewDesc(
			"scw_snapshot_state",
			"State of the snapshot",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"scw_snapshot_created_timestamp",
			"Timestamp when the snapshot have been created",
			labels,
			nil,
		),
		Modified: prometheus.NewDesc(
			"scw_snapshot_modified_timestamp",
			"Timestamp when the snapshot have been modified",
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
func (c *SnapshotCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Available,
		c.Size,
		c.VolumeType,
		c.State,
		c.Created,
		c.Modified,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *SnapshotCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Available
	ch <- c.Size
	ch <- c.VolumeType
	ch <- c.State
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *SnapshotCollector) Collect(ch chan<- prometheus.Metric) {
	for _, zone := range affectedZones(c.client) {
		now := time.Now()
		resp, err := c.instance.ListSnapshots(&instance.ListSnapshotsRequest{
			Zone:         zone,
			Organization: c.org,
			Project:      c.project,
		}, scw.WithAllPages())
		c.duration.WithLabelValues("snapshot").Observe(time.Since(now).Seconds())

		if err != nil {
			c.logger.Error("Failed to fetch snapshots",
				"zone", zone,
				"err", err,
			)

			c.failures.WithLabelValues("snapshot").Inc()
			return
		}

		c.logger.Debug("Fetched snapshots",
			"zone", zone,
			"count", resp.TotalCount,
		)

		for _, snapshot := range resp.Snapshots {
			var (
				volumeType float64
				state      float64
			)

			projectName, err := retrieveProject(c.logger, c.client, snapshot.Project)
			if err != nil {
				c.logger.Error("Failed to retrieve project",
					"project", snapshot.Project,
					"err", err,
				)

				continue
			}

			labels := []string{
				snapshot.ID,
				snapshot.Name,
				snapshot.Zone.String(),
				snapshot.Organization,
				snapshot.Project,
				projectName,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Available,
				prometheus.GaugeValue,
				1.0,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Size,
				prometheus.GaugeValue,
				float64(snapshot.Size),
				labels...,
			)

			switch val := snapshot.VolumeType; val {
			case instance.VolumeVolumeTypeLSSD:
				volumeType = 1.0
			case instance.VolumeVolumeTypeBSSD:
				volumeType = 2.0
			case instance.VolumeVolumeTypeUnified:
				volumeType = 3.0
			default:
				volumeType = 0.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.VolumeType,
				prometheus.GaugeValue,
				volumeType,
				labels...,
			)

			switch val := snapshot.State; val {
			case instance.SnapshotStateAvailable:
				state = 1.0
			case instance.SnapshotStateSnapshotting:
				state = 2.0
			case instance.SnapshotStateInvalidData:
				state = 3.0
			case instance.SnapshotStateImporting:
				state = 4.0
			case instance.SnapshotStateExporting:
				state = 5.0
			default:
				state = 0.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.State,
				prometheus.GaugeValue,
				state,
				labels...,
			)

			if snapshot.CreationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Created,
					prometheus.GaugeValue,
					float64(snapshot.CreationDate.Unix()),
					labels...,
				)
			}

			if snapshot.ModificationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Modified,
					prometheus.GaugeValue,
					float64(snapshot.ModificationDate.Unix()),
					labels...,
				)
			}
		}
	}
}
