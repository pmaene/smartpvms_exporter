package collectors

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iancoleman/strcase"
	"github.com/pmaene/smartpvms_exporter/internal"
	"github.com/pmaene/smartpvms_exporter/internal/smartpvms"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/exp/maps"
)

const (
	plantsNamespace = "smartpvms"
	plantsSubsystem = "plant"
)

var (
	plantsUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "up"),
		"Whether collecting plant metrics was successful.",
		nil,
		nil,
	)

	plantsInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "info"),
		"Status of the plant.",
		[]string{"station_code", "name", "address", "capacity", "status"},
		nil,
	)

	plantsDayYieldDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "day_yield"),
		"Yield of the plant today.",
		[]string{"station_code"},
		nil,
	)

	plantsMonthYieldDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "month_yield"),
		"Yield of the plant this month.",
		[]string{"station_code"},
		nil,
	)

	plantsTotalYieldDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "total_yield"),
		"Total yield of the plant.",
		[]string{"station_code"},
		nil,
	)

	plantsDayIncomeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "day_income"),
		"Income of the plant today.",
		[]string{"station_code"},
		nil,
	)

	plantsTotalIncomeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(plantsNamespace, plantsSubsystem, "total_income"),
		"Total income of the plant.",
		[]string{"station_code"},
		nil,
	)
)

type Plant struct {
	smartpvms.Plant
	Data smartpvms.PlantData
}

type PlantsCollector struct {
	Cache *internal.Cache[Plant]
}

func (c *PlantsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- plantsUpDesc
	ch <- plantsInfoDesc
	ch <- plantsDayYieldDesc
	ch <- plantsMonthYieldDesc
	ch <- plantsTotalYieldDesc
	ch <- plantsDayIncomeDesc
	ch <- plantsTotalIncomeDesc
}

func (c *PlantsCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		plantsUpDesc,
		prometheus.GaugeValue,
		c.up(),
	)

	for _, v := range c.Cache.Data() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				plantsInfoDesc,
				prometheus.GaugeValue,
				1,
				v.StationCode,
				v.Name,
				v.Address,
				strconv.FormatFloat(1000*1000*v.Capacity, 'f', -1, 64),
				strcase.ToSnake(v.Data.Status.String()),
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				plantsDayYieldDesc,
				prometheus.GaugeValue,
				1000*v.Data.DayYield,
				v.Plant.StationCode,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				plantsMonthYieldDesc,
				prometheus.GaugeValue,
				1000*v.Data.MonthYield,
				v.Plant.StationCode,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				plantsTotalYieldDesc,
				prometheus.CounterValue,
				1000*v.Data.TotalYield,
				v.Plant.StationCode,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				plantsDayIncomeDesc,
				prometheus.GaugeValue,
				v.Data.DayIncome,
				v.Plant.StationCode,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				plantsTotalIncomeDesc,
				prometheus.CounterValue,
				v.Data.TotalIncome,
				v.Plant.StationCode,
			),
		)
	}
}

func (c *PlantsCollector) up() float64 {
	if c.Cache.IsValid() {
		return 1
	}

	return 0
}

func NewPlantsCollector(c *resty.Client, i time.Duration, l log.Logger) *PlantsCollector {
	r := &plantsRefresher{
		client:   c,
		interval: i,
	}

	return &PlantsCollector{
		Cache: internal.NewCache[Plant](l, r),
	}
}

type plantsRefresher struct {
	client   *resty.Client
	interval time.Duration
}

func (r *plantsRefresher) Interval() time.Duration {
	return r.interval
}

func (r *plantsRefresher) Refresh() ([]Plant, error) {
	res, err := smartpvms.GetPlantList(r.client)
	if err != nil {
		return nil, err
	}

	if res.Success {
		ps := make(map[string]Plant, 0)
		for _, v := range res.Data {
			ps[v.StationCode] = Plant{Plant: v}
		}

		res, err := smartpvms.GetRealtimePlantData(r.client, maps.Keys(ps)...)
		if err != nil {
			return nil, err
		}

		if res.Success {
			for _, v := range res.Data {
				if p, ok := ps[v.StationCode]; ok {
					ps[v.StationCode] = Plant{
						Plant: p.Plant,
						Data:  v.DataItemMap,
					}
				}
			}

			return maps.Values(ps), nil
		}
	}

	return nil, errors.New("collectors: failed to refresh plants")
}
