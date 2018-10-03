package exporter

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	scw "github.com/scaleway/go-scaleway"
)

// VolumeCollector collects metrics about the volumes.
type VolumeCollector struct {
	client   *scw.ScalewayAPI
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	timeout  time.Duration

	Available *prometheus.Desc
	Size      *prometheus.Desc
	Created   *prometheus.Desc
	Modified  *prometheus.Desc
}

// NewVolumeCollector returns a new VolumeCollector.
func NewVolumeCollector(logger log.Logger, client *scw.ScalewayAPI, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *VolumeCollector {
	failures.WithLabelValues("volume").Add(0)

	labels := []string{"id", "name"}
	return &VolumeCollector{
		client:   client,
		logger:   log.With(logger, "collector", "volume"),
		failures: failures,
		duration: duration,
		timeout:  timeout,

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
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *VolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Available
	ch <- c.Size
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *VolumeCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	volumes, err := c.client.GetVolumes()
	c.duration.WithLabelValues("volume").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch volumes",
			"err", err,
		)

		c.failures.WithLabelValues("volume").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched volumes",
		"count", len(*volumes),
	)

	for _, volume := range *volumes {
		var (
			created float64
			updated float64
		)

		labels := []string{
			volume.Identifier,
			volume.Name,
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

		if num, err := time.Parse("2006-01-02T15:04:05.000000-07:00", volume.CreationDate); err == nil {
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

		if num, err := time.Parse("2006-01-02T15:04:05.000000-07:00", volume.ModificationDate); err == nil {
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
