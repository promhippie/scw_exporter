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

// SecurityGroupCollector collects metrics about the security groups.
type SecurityGroupCollector struct {
	client   *scw.Client
	instance *instance.API
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target
	org      *string
	project  *string

	Defined         *prometheus.Desc
	EnableDefault   *prometheus.Desc
	ProjectDefault  *prometheus.Desc
	Stateful        *prometheus.Desc
	InboundDefault  *prometheus.Desc
	OutboundDefault *prometheus.Desc
	Servers         *prometheus.Desc
	Created         *prometheus.Desc
	Modified        *prometheus.Desc
}

// NewSecurityGroupCollector returns a new SecurityGroupCollector.
func NewSecurityGroupCollector(logger log.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *SecurityGroupCollector {
	if failures != nil {
		failures.WithLabelValues("security_group").Add(0)
	}

	labels := []string{"id", "name", "zone", "org", "project"}
	collector := &SecurityGroupCollector{
		client:   client,
		instance: instance.NewAPI(client),
		logger:   log.With(logger, "collector", "security_group"),
		failures: failures,
		duration: duration,
		config:   cfg,

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
		ProjectDefault: prometheus.NewDesc(
			"scw_security_group_project_default",
			"1 if the security group is an project default, 0 otherwise",
			labels,
			nil,
		),
		Stateful: prometheus.NewDesc(
			"scw_security_group_stateful",
			"1 if the security group is stateful by default, 0 otherwise",
			labels,
			nil,
		),
		InboundDefault: prometheus.NewDesc(
			"scw_security_group_inbound_default_policy",
			"1 if the security group inbound default policy is accept, 0 otherwise",
			labels,
			nil,
		),
		OutboundDefault: prometheus.NewDesc(
			"scw_security_group_outbound_default_policy",
			"1 if the security group outbound default policy is accept, 0 otherwise",
			labels,
			nil,
		),
		Servers: prometheus.NewDesc(
			"scw_security_group_servers_count",
			"Number of servers attached to the security group",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"scw_security_group_created_timestamp",
			"Timestamp when the security group have been created",
			labels,
			nil,
		),
		Modified: prometheus.NewDesc(
			"scw_security_group_modified_timestamp",
			"Timestamp when the security group have been modified",
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
func (c *SecurityGroupCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Defined,
		c.EnableDefault,
		c.ProjectDefault,
		c.Stateful,
		c.InboundDefault,
		c.OutboundDefault,
		c.Servers,
		c.Created,
		c.Modified,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *SecurityGroupCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Defined
	ch <- c.EnableDefault
	ch <- c.ProjectDefault
	ch <- c.Stateful
	ch <- c.InboundDefault
	ch <- c.OutboundDefault
	ch <- c.Servers
	ch <- c.Created
	ch <- c.Modified
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *SecurityGroupCollector) Collect(ch chan<- prometheus.Metric) {
	for _, zone := range affectedZones(c.client) {
		now := time.Now()
		resp, err := c.instance.ListSecurityGroups(&instance.ListSecurityGroupsRequest{
			Zone:         zone,
			Organization: c.org,
			Project:      c.project,
		}, scw.WithAllPages())
		c.duration.WithLabelValues("security_group").Observe(time.Since(now).Seconds())

		if err != nil {
			level.Error(c.logger).Log(
				"msg", "Failed to fetch security groups",
				"zone", zone,
				"err", err,
			)

			c.failures.WithLabelValues("security_group").Inc()
			return
		}

		level.Debug(c.logger).Log(
			"msg", "Fetched security groups",
			"zone", zone,
			"count", resp.TotalCount,
		)

		for _, securityGroup := range resp.SecurityGroups {
			var (
				enableDefault   float64
				projectDefault  float64
				stateful        float64
				inboundDefault  float64
				outboundDefault float64
			)

			labels := []string{
				securityGroup.ID,
				securityGroup.Name,
				securityGroup.Zone.String(),
				securityGroup.Organization,
				securityGroup.Project,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Defined,
				prometheus.GaugeValue,
				1.0,
				labels...,
			)

			if securityGroup.EnableDefaultSecurity {
				enableDefault = 1.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.EnableDefault,
				prometheus.GaugeValue,
				enableDefault,
				labels...,
			)

			if securityGroup.ProjectDefault {
				projectDefault = 1.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.ProjectDefault,
				prometheus.GaugeValue,
				projectDefault,
				labels...,
			)

			if securityGroup.Stateful {
				stateful = 1.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.Stateful,
				prometheus.GaugeValue,
				stateful,
				labels...,
			)

			switch val := securityGroup.InboundDefaultPolicy; val {
			case instance.SecurityGroupPolicyAccept:
				inboundDefault = 1.0
			case instance.SecurityGroupPolicyDrop:
				inboundDefault = 0.0
			default:
				inboundDefault = 0.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.InboundDefault,
				prometheus.GaugeValue,
				inboundDefault,
				labels...,
			)

			switch val := securityGroup.OutboundDefaultPolicy; val {
			case instance.SecurityGroupPolicyAccept:
				outboundDefault = 1.0
			case instance.SecurityGroupPolicyDrop:
				outboundDefault = 0.0
			default:
				outboundDefault = 0.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.OutboundDefault,
				prometheus.GaugeValue,
				outboundDefault,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Servers,
				prometheus.GaugeValue,
				float64(len(securityGroup.Servers)),
				labels...,
			)

			if securityGroup.CreationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Created,
					prometheus.GaugeValue,
					float64(securityGroup.CreationDate.Unix()),
					labels...,
				)
			}

			if securityGroup.ModificationDate != nil {
				ch <- prometheus.MustNewConstMetric(
					c.Modified,
					prometheus.GaugeValue,
					float64(securityGroup.ModificationDate.Unix()),
					labels...,
				)
			}
		}
	}
}
