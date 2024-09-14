package exporter

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/scw_exporter/pkg/config"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// VolumeCollector collects metrics about the volumes.
type VolumeCollector struct {
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

// NewVolumeCollector returns a new VolumeCollector.
func NewVolumeCollector(logger *slog.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *VolumeCollector {
	if failures != nil {
		failures.WithLabelValues("volume").Add(0)
	}

	labels := []string{"id", "name", "zone", "org", "project"}
	collector := &VolumeCollector{
		client:   client,
		instance: instance.NewAPI(client),
		logger:   logger.With("collector", "volume"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Available: prometheus.NewDesc(
			"scw_volume_available",
			"Constant value of 1 that this volume is available",
			labels,
			nil,
		),
		Size: prometheus.NewDesc(
			"scw_volume_size_bytes",
			"Size of the volume in bytes",
			labels,
			nil,
		),
		VolumeType: prometheus.NewDesc(
			"scw_volume_type",
			"Type of the snapshot",
			labels,
			nil,
		),
		State: prometheus.NewDesc(
			"scw_volume_state",
			"State of the snapshot",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"scw_volume_created_timestamp",
			"Timestamp when the volume have been created",
			labels,
			nil,
		),
		Modified: prometheus.NewDesc(
			"scw_volume_modified_timestamp",
			"Timestamp when the volume have been modified",
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
func (c *VolumeCollector) Metrics() []*prometheus.Desc {
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
func (c *VolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Available
	ch <- c.Size
	ch <- c.VolumeType
	ch <- c.State
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *VolumeCollector) Collect(ch chan<- prometheus.Metric) {
	for _, zone := range affectedZones(c.client) {
		now := time.Now()
		resp, err := c.instance.ListVolumes(&instance.ListVolumesRequest{
			Zone:         zone,
			Organization: c.org,
			Project:      c.project,
		}, scw.WithAllPages())
		c.duration.WithLabelValues("volume").Observe(time.Since(now).Seconds())

		if err != nil {
			c.logger.Error("Failed to fetch volumes",
				"zone", zone,
				"err", err,
			)

			c.failures.WithLabelValues("volume").Inc()
			return
		}

		c.logger.Debug("Fetched volumes",
			"zone", zone,
			"count", resp.TotalCount,
		)

		for _, volume := range resp.Volumes {
			var (
				volumeType float64
				state      float64
			)

			labels := []string{
				volume.ID,
				volume.Name,
				volume.Zone.String(),
				volume.Organization,
				volume.Project,
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
				float64(volume.Size),
				labels...,
			)

			switch val := volume.VolumeType; val {
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

			switch val := volume.State; val {
			case instance.VolumeStateAvailable:
				state = 1.0
			case instance.VolumeStateSnapshotting:
				state = 2.0
			case instance.VolumeStateFetching:
				state = 3.0
			case instance.VolumeStateResizing:
				state = 4.0
			case instance.VolumeStateSaving:
				state = 5.0
			case instance.VolumeStateHotsyncing:
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

			if volume.CreationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Created,
					prometheus.GaugeValue,
					float64(volume.CreationDate.Unix()),
					labels...,
				)
			}

			if volume.ModificationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Modified,
					prometheus.GaugeValue,
					float64(volume.ModificationDate.Unix()),
					labels...,
				)
			}
		}
	}
}
