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
	residentialInvertersNamespace = "smartpvms"
	residentialInvertersSubsystem = "residential_inverter"
)

var (
	residentialInvertersUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "up"),
		"Whether collecting residential inverter metrics was successful.",
		nil,
		nil,
	)

	residentialInvertersInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "info"),
		"Status of the residential inverter.",
		[]string{
			"station_code",
			"serial",
			"model",
			"software_version",
			"latitude",
			"longitude",
			"run_status",
			"status",
			"startup_time",
			"shutdown_time",
		},
		nil,
	)

	residentialInvertersTemperatureDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "temperature"),
		"Internal temperature of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersEfficiencyDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "efficiency"),
		"Efficiency of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersPowerFactorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "power_factor"),
		"Power factor of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersActivePowerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "active_power"),
		"Active output power of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersReactivePowerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "reactive_power"),
		"Reactive output power of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersPVPowerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "pv_power"),
		"Solar input power of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersVoltageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "voltage"),
		"Output voltage of the residential inverter.",
		[]string{"station_code", "serial", "phase"},
		nil,
	)

	residentialInvertersCurrentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "current"),
		"Output voltage of the residential inverter.",
		[]string{"station_code", "serial", "phase"},
		nil,
	)

	residentialInvertersPVVoltageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "pv_voltage"),
		"Voltage of the solar panel string.",
		[]string{"station_code", "serial", "string"},
		nil,
	)

	residentialInvertersPVCurrentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "pv_current"),
		"Current of the solar panel string.",
		[]string{"station_code", "serial", "string"},
		nil,
	)

	residentialInvertersDayYieldDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "day_yield"),
		"Yield of the residential inverter today.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersTotalYieldDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "total_yield"),
		"Total yield of the residential inverter.",
		[]string{"station_code", "serial"},
		nil,
	)

	residentialInvertersMPPTTotalYieldDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "mppt_total_yield"),
		"Total yield of the MPP tracker.",
		[]string{"station_code", "serial", "tracker"},
		nil,
	)

	residentialInvertersGridVoltageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "grid_voltage"),
		"Grid voltage of the residential inverter.",
		[]string{"station_code", "serial", "phase"},
		nil,
	)

	residentialInvertersGridFrequencyDesc = prometheus.NewDesc(
		prometheus.BuildFQName(residentialInvertersNamespace, residentialInvertersSubsystem, "grid_frequency"),
		"Frequency of the grid.",
		[]string{"station_code", "serial"},
		nil,
	)
)

type ResidentialInverter struct {
	smartpvms.Device
	Data smartpvms.ResidentialInverterData
}

type ResidentialInvertersCollector struct {
	Cache *internal.Cache[ResidentialInverter]
}

func (c *ResidentialInvertersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- residentialInvertersUpDesc
	ch <- residentialInvertersInfoDesc
	ch <- residentialInvertersTemperatureDesc
	ch <- residentialInvertersEfficiencyDesc
	ch <- residentialInvertersPowerFactorDesc
	ch <- residentialInvertersActivePowerDesc
	ch <- residentialInvertersReactivePowerDesc
	ch <- residentialInvertersPVPowerDesc
	ch <- residentialInvertersVoltageDesc
	ch <- residentialInvertersCurrentDesc
	ch <- residentialInvertersPVVoltageDesc
	ch <- residentialInvertersPVCurrentDesc
	ch <- residentialInvertersDayYieldDesc
	ch <- residentialInvertersTotalYieldDesc
	ch <- residentialInvertersMPPTTotalYieldDesc
	ch <- residentialInvertersGridVoltageDesc
	ch <- residentialInvertersGridFrequencyDesc
}

func (c *ResidentialInvertersCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		residentialInvertersUpDesc,
		prometheus.GaugeValue,
		c.up(),
	)

	for _, v := range c.Cache.Data() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersInfoDesc,
				prometheus.GaugeValue,
				1,
				v.Device.StationCode,
				v.Device.Serial,
				v.Device.Model,
				v.Device.SoftwareVersion,
				strconv.FormatFloat(v.Device.Latitude, 'f', -1, 64),
				strconv.FormatFloat(v.Device.Longitude, 'f', -1, 64),
				strcase.ToSnake(v.Data.RunStatus.String()),
				strcase.ToSnake(v.Data.Status.String()),
				v.Data.StartupTime.String(),
				v.Data.ShutdownTime.String(),
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersTemperatureDesc,
				prometheus.GaugeValue,
				v.Data.Temperature,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersEfficiencyDesc,
				prometheus.GaugeValue,
				v.Data.Efficiency/100,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPowerFactorDesc,
				prometheus.GaugeValue,
				v.Data.PowerFactor,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersActivePowerDesc,
				prometheus.GaugeValue,
				1000*v.Data.ActivePower,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersReactivePowerDesc,
				prometheus.GaugeValue,
				1000*v.Data.ReactivePower,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVPowerDesc,
				prometheus.GaugeValue,
				1000*v.Data.PVPower,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersVoltageDesc,
				prometheus.GaugeValue,
				v.Data.L1Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"l1",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersVoltageDesc,
				prometheus.GaugeValue,
				v.Data.L2Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"l2",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersVoltageDesc,
				prometheus.GaugeValue,
				v.Data.L3Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"l3",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersCurrentDesc,
				prometheus.GaugeValue,
				v.Data.L1Current,
				v.Device.StationCode,
				v.Device.Serial,
				"l1",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersCurrentDesc,
				prometheus.GaugeValue,
				v.Data.L2Current,
				v.Device.StationCode,
				v.Device.Serial,
				"l2",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersCurrentDesc,
				prometheus.GaugeValue,
				v.Data.L3Current,
				v.Device.StationCode,
				v.Device.Serial,
				"l3",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV1Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv1",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV2Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv2",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV3Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv3",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV4Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv4",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV5Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv5",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV6Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv6",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV7Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv7",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVVoltageDesc,
				prometheus.GaugeValue,
				v.Data.PV8Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"pv8",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV1Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv1",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV2Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv2",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV3Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv3",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV4Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv4",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV5Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv5",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV6Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv6",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV7Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv7",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersPVCurrentDesc,
				prometheus.GaugeValue,
				v.Data.PV8Current,
				v.Device.StationCode,
				v.Device.Serial,
				"pv8",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersDayYieldDesc,
				prometheus.GaugeValue,
				1000*v.Data.DayYield,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersTotalYieldDesc,
				prometheus.CounterValue,
				1000*v.Data.TotalYield,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersMPPTTotalYieldDesc,
				prometheus.CounterValue,
				1000*v.Data.MPPT1TotalYield,
				v.Device.StationCode,
				v.Device.Serial,
				"mppt1",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersMPPTTotalYieldDesc,
				prometheus.CounterValue,
				1000*v.Data.MPPT2TotalYield,
				v.Device.StationCode,
				v.Device.Serial,
				"mppt2",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersMPPTTotalYieldDesc,
				prometheus.CounterValue,
				1000*v.Data.MPPT3TotalYield,
				v.Device.StationCode,
				v.Device.Serial,
				"mppt3",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersMPPTTotalYieldDesc,
				prometheus.CounterValue,
				1000*v.Data.MPPT4TotalYield,
				v.Device.StationCode,
				v.Device.Serial,
				"mppt4",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersGridVoltageDesc,
				prometheus.GaugeValue,
				v.Data.GridL1L2Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"l1l2",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersGridVoltageDesc,
				prometheus.GaugeValue,
				v.Data.GridL2L3Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"l2l3",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersGridVoltageDesc,
				prometheus.GaugeValue,
				v.Data.GridL3L1Voltage,
				v.Device.StationCode,
				v.Device.Serial,
				"l3l1",
			),
		)

		ch <- prometheus.NewMetricWithTimestamp(
			c.Cache.Timestamp(),
			prometheus.MustNewConstMetric(
				residentialInvertersGridFrequencyDesc,
				prometheus.GaugeValue,
				v.Data.GridFrequency,
				v.Device.StationCode,
				v.Device.Serial,
			),
		)
	}
}

func (c *ResidentialInvertersCollector) up() float64 {
	if c.Cache.IsValid() {
		return 1
	}

	return 0
}

func NewResidentialInvertersCollector(c *resty.Client, i time.Duration, l log.Logger) *ResidentialInvertersCollector {
	r := &residentialInvertersRefresher{
		client:   c,
		interval: i,
	}

	return &ResidentialInvertersCollector{
		Cache: internal.NewCache[ResidentialInverter](l, r),
	}
}

type residentialInvertersRefresher struct {
	client   *resty.Client
	interval time.Duration
}

func (r *residentialInvertersRefresher) Interval() time.Duration {
	return r.interval
}

func (r *residentialInvertersRefresher) Refresh() ([]ResidentialInverter, error) {
	res, err := smartpvms.GetPlantList(r.client)
	if err != nil {
		return nil, err
	}

	if res.Success {
		var cs []string
		for _, v := range res.Data {
			cs = append(cs, v.StationCode)
		}

		res, err := smartpvms.GetDeviceList(r.client, cs...)
		if err != nil {
			return nil, err
		}

		if res.Success {
			ds := make(map[int64]ResidentialInverter, 0)
			for _, v := range res.Data {
				if v.Type != smartpvms.DeviceTypeResidentialInverter {
					continue
				}

				ds[v.ID] = ResidentialInverter{Device: v}
			}

			res, err := smartpvms.GetRealtimeDeviceData[smartpvms.ResidentialInverterData](
				r.client,
				smartpvms.DeviceTypeResidentialInverter,
				maps.Keys(ds)...,
			)

			if err != nil {
				return nil, err
			}

			if res.Success {
				for _, v := range res.Data {
					if d, ok := ds[v.DeviceID]; ok {
						ds[v.DeviceID] = ResidentialInverter{
							Device: d.Device,
							Data:   v.DataItemMap,
						}
					}
				}

				return maps.Values(ds), nil
			}
		}
	}

	return nil, errors.New("collectors: failed to refresh residential inverters")
}
