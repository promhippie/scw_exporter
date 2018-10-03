package exporter

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	scw "github.com/scaleway/go-scaleway"
)

// SnapshotCollector collects metrics about the snapshots.
type SnapshotCollector struct {
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

// NewSnapshotCollector returns a new SnapshotCollector.
func NewSnapshotCollector(logger log.Logger, client *scw.ScalewayAPI, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *SnapshotCollector {
	failures.WithLabelValues("snapshot").Add(0)

	labels := []string{"id", "name"}
	return &SnapshotCollector{
		client:   client,
		logger:   log.With(logger, "collector", "snapshot"),
		failures: failures,
		duration: duration,
		timeout:  timeout,

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
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *SnapshotCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Available
	ch <- c.Size
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *SnapshotCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	snapshots, err := c.client.GetSnapshots()
	c.duration.WithLabelValues("snapshot").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch snapshots",
			"err", err,
		)

		c.failures.WithLabelValues("snapshot").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched snapshots",
		"count", len(*snapshots),
	)

	for _, snapshot := range *snapshots {
		var (
			created float64
			updated float64
		)

		labels := []string{
			snapshot.Identifier,
			snapshot.Name,
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

		if num, err := time.Parse("2006-01-02T15:04:05.000000-07:00", snapshot.CreationDate); err == nil {
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

		if num, err := time.Parse("2006-01-02T15:04:05.000000-07:00", snapshot.ModificationDate); err == nil {
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
