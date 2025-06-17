package exporter

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/scw_exporter/pkg/config"

	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ConsumptionCollector collects metrics about resources consumption.
type ConsumptionCollector struct {
	client      *scw.Client
	consumption *billing.API
	logger      *slog.Logger
	failures    *prometheus.CounterVec
	duration    *prometheus.HistogramVec
	config      config.Target
	org         *string
	project     *string

	Value    *prometheus.Desc
	Quantity *prometheus.Desc
}

// NewConsumptionCollector returns a new ServerCollector.
func NewConsumptionCollector(logger *slog.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *ConsumptionCollector {
	if failures != nil {
		failures.WithLabelValues("consumption").Add(0)
	}

	labels := cfg.Consumption.Labels
	collector := &ConsumptionCollector{
		client:      client,
		consumption: billing.NewAPI(client),
		logger:      logger.With("collector", "consumption"),
		failures:    failures,
		duration:    duration,
		config:      cfg,

		Value: prometheus.NewDesc(
			"scw_consumption_value",
			"sdasdas",
			labels,
			nil,
		),

		Quantity: prometheus.NewDesc(
			"scw_consumption_billed_quantity",
			"sdasdsadasd",
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
func (c *ConsumptionCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Value,
		c.Quantity,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ConsumptionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Value
	ch <- c.Quantity
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ConsumptionCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	resp, err := c.consumption.ListConsumptions(&billing.ListConsumptionsRequest{
		OrganizationID: c.org,
		ProjectID:      c.project,
	}, scw.WithAllPages())
	c.duration.WithLabelValues("consumption").Observe(time.Since(now).Seconds())

	if err != nil {
		c.logger.Error("Failed to fetch consumptions",
			"err", err,
		)

		c.failures.WithLabelValues("consumption").Inc()
		return
	}

	c.logger.Debug("Fetched consumptions",
		"count", resp.TotalCount,
	)

	for _, consumption := range resp.Consumptions {
		var (
			value    float64
			quantity float64
		)

		if consumption.Value != nil {
			value = float64(consumption.Value.Units)
		}

		quantity, err = strconv.ParseFloat(
			consumption.BilledQuantity,
			64,
		)

		if err != nil {
			c.logger.Error("Failed to parse consumptions",
				"err", err,
			)

			continue
		}

		labels := []string{}

		for _, label := range c.config.Consumption.Labels {
			labels = append(
				labels,
				consumptionLabel(consumption, label),
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.Value,
			prometheus.CounterValue,
			value,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Quantity,
			prometheus.CounterValue,
			quantity,
			labels...,
		)
	}
}

func consumptionLabel(consumption *billing.ListConsumptionsResponseConsumption, label string) string {
	switch label {
	case "category_name":
		return consumption.CategoryName
	case "product_name":
		return consumption.ProductName
	case "project_id":
		return consumption.ProjectID
	case "resource_name":
		return consumption.ResourceName
	case "sku":
		return consumption.Sku
	case "unit":
		return consumption.Unit
	case "currency":
		if consumption.Value != nil {
			return consumption.Value.CurrencyCode
		}
	}

	return ""
}
