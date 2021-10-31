package internal

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "smartpvms"
)

var (
	upDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Whether collecting SmartPVMS metrics was successful.",
		nil,
		nil,
	)
)

type Collector struct {
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		upDesc,
		prometheus.GaugeValue,
		c.up(),
	)
}

func (c *Collector) up() float64 {
	return 1
}

func NewCollector() *Collector {
	return &Collector{}
}
