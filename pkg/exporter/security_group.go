package exporter

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	scw "github.com/scaleway/go-scaleway"
)

// SecurityGroupCollector collects metrics about the security groups.
type SecurityGroupCollector struct {
	client   *scw.ScalewayAPI
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	timeout  time.Duration

	Defined       *prometheus.Desc
	EnableDefault *prometheus.Desc
	OrgDefault    *prometheus.Desc
}

// NewSecurityGroupCollector returns a new SecurityGroupCollector.
func NewSecurityGroupCollector(logger log.Logger, client *scw.ScalewayAPI, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *SecurityGroupCollector {
	failures.WithLabelValues("security_group").Add(0)

	labels := []string{"id", "name"}
	return &SecurityGroupCollector{
		client:   client,
		logger:   log.With(logger, "collector", "security_group"),
		failures: failures,
		duration: duration,
		timeout:  timeout,

		Defined: prometheus.NewDesc(
			"scw_security_group_defined",
			"Constant value of 1 that this security group is defined",
			labels,
			nil,
		),
		EnableDefault: prometheus.NewDesc(
			"scw_security_group_enable_default",
			"1 if the security group is enabled by default, 0 otherwise",
			labels,
			nil,
		),
		OrgDefault: prometheus.NewDesc(
			"scw_security_group_organization_default",
			"1 if the security group is an organization default, 0 otherwise",
			labels,
			nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *SecurityGroupCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Defined
	ch <- c.EnableDefault
	ch <- c.OrgDefault
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *SecurityGroupCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	sgs, err := c.client.GetSecurityGroups()
	c.duration.WithLabelValues("security_group").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch security groups",
			"err", err,
		)

		c.failures.WithLabelValues("security_group").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched security groups",
		"count", len(sgs.SecurityGroups),
	)

	for _, sg := range sgs.SecurityGroups {
		var (
			enableDefault float64
			orgDefault    float64
		)

		labels := []string{
			sg.ID,
			sg.Name,
		}

		ch <- prometheus.MustNewConstMetric(
			c.Defined,
			prometheus.GaugeValue,
			1.0,
			labels...,
		)

		if sg.EnableDefaultSecurity {
			enableDefault = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.EnableDefault,
			prometheus.GaugeValue,
			enableDefault,
			labels...,
		)

		if sg.OrganizationDefault {
			orgDefault = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.OrgDefault,
			prometheus.GaugeValue,
			orgDefault,
			labels...,
		)
	}
}
