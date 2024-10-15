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

func NewConsumptionCollector(logger *slog.Logger, client *scw.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *ConsumptionCollector {
	if failures != nil {
		failures.WithLabelValues("consumption").Add(0)
	}

	labels := []string{"category_name", "product_name", "project_id", "resource_name", "sku", "unit"}
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

func (c *ConsumptionCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Value,
		c.Quantity,
	}
}

func (c *ConsumptionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Value
	ch <- c.Quantity
}

func (c *ConsumptionCollector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()

	resp, err := c.consumption.ListConsumptions(
		&billing.ListConsumptionsRequest{},
		scw.WithAllPages(),
	)
	c.duration.WithLabelValues("consumption").Observe(time.Since(now).Seconds())

	if err != nil {
		c.logger.Error("failed to fetch consumptions",
			"organization", c.org,
			"err", err,
		)
		c.failures.WithLabelValues("consumption").Inc()
		return
	}

	for _, consumption := range resp.Consumptions {

		var (
			value    float64
			quantity float64
		)

		quantity, _ = strconv.ParseFloat(consumption.BilledQuantity, 64)

		labels := []string{
			consumption.CategoryName,
			consumption.ProductName,
			consumption.ProjectID,
			consumption.ResourceName,
			consumption.Sku,
			consumption.Unit,
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
