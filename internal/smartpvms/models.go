package smartpvms

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Result struct {
	Success  bool `json:"success"`
	FailCode int  `json:"failCode"`
	Params   *struct {
		CurrentTime int64 `json:"currentTime"`
	}
	Message *string `json:"message"`
}

type LoginBody struct {
	Username string `json:"userName"`
	Password string `json:"systemCode"`
}

type LoginResult struct {
	Result
	Data interface{} `json:"data"`
}

type LogoutBody struct {
	XSRFToken string `json:"xsrfToken"`
}

type LogoutResult struct {
	Result
	Data interface{} `json:"data"`
}

type GetPlantListResult struct {
	Result
	Data []struct {
		StationCode        string  `json:"stationCode"`
		Name               string  `json:"stationName"`
		Address            string  `json:"stationAddr"`
		Capacity           float64 `json:"capacity"`
		BuildStatus        string  `json:"buildState"`
		GridConnectionType string  `json:"combineType"`
		AIDType            int     `json:"aidType"`
		ContactPerson      string  `json:"stationLinkman"`
		ContactPersonPhone string  `json:"linkmanPho"`
	} `json:"data"`
}

type GetRealTimePlantDataBody struct {
	StationCodes []string `json:"stationCodes"`
}

func (b *GetRealTimePlantDataBody) MarshalJSON() ([]byte, error) {
	type Alias GetRealTimePlantDataBody

	return json.Marshal(&struct {
		StationCodes string `json:"stationCodes"`
		*Alias
	}{
		StationCodes: strings.Join(b.StationCodes, ","),
		Alias:        (*Alias)(b),
	})
}

type GetRealTimePlantDataResult struct {
	Result
	Data []struct {
		DataItemMap struct {
			DayPower    float64 `json:"day_power"`
			MonthPower  float64 `json:"month_power"`
			TotalPower  float64 `json:"total_power"`
			DayIncome   float64 `json:"day_income"`
			TotalIncome float64 `json:"total_income"`
			Status      int     `json:"real_health_state"`
		} `json:"dataItemMap"`
		StationCode string `json:"stationCode"`
	} `json:"data"`
}

type GetDeviceListBody struct {
	StationCodes []string `json:"stationCodes"`
}

func (b *GetDeviceListBody) MarshalJSON() ([]byte, error) {
	type Alias GetDeviceListBody

	return json.Marshal(&struct {
		StationCodes string `json:"stationCodes"`
		*Alias
	}{
		StationCodes: strings.Join(b.StationCodes, ","),
		Alias:        (*Alias)(b),
	})
}

type GetDeviceListResult struct {
	Result
	Data []struct {
		ID              int64   `json:"id"`
		Name            string  `json:"devName"`
		StationCode     string  `json:"stationCode"`
		Serial          string  `json:"esnCode"`
		Type            int     `json:"devTypeId"`
		SoftwareVersion string  `json:"softwareVersion"`
		IverterType     string  `json:"invType"`
		Longitude       float64 `json:"longitude"`
		Latitude        float64 `json:"latitude"`
	} `json:"data"`
}

type GetDeviceDataBody struct {
	IDs  []int64 `json:"devIds"`
	Type int     `json:"devTypeId"`
}

func (b *GetDeviceDataBody) MarshalJSON() ([]byte, error) {
	type Alias GetDeviceDataBody

	var ids []string
	for _, v := range b.IDs {
		ids = append(ids, strconv.Itoa(int(v)))
	}

	return json.Marshal(&struct {
		IDs string `json:"devIds"`
		*Alias
	}{
		IDs:   strings.Join(ids, ","),
		Alias: (*Alias)(b),
	})
}

type GetDeviceDataResult struct {
	Result
	Data []struct {
		DataItemMap struct {
			InverterStatus float64 `json:"inverter_state"`
			GridABVoltage  float64 `json:"ab_u"`
			GridBCVoltage  float64 `json:"bc_u"`
			GridCAVoltage  float64 `json:"ca_u"`
			PhaseAVoltage  float64 `json:"a_u"`
			PhaseBVoltage  float64 `json:"b_u"`
			PhaseCVoltage  float64 `json:"c_u"`
			PhaseACurrent  float64 `json:"a_i"`
			PhaseBCurrent  float64 `json:"b_i"`
			PhaseCCurrent  float64 `json:"c_i"`
			Efficiency     float64 `json:"efficiency"`
			Temperature    float64 `json:"temperature"`
			PowerFactor    float64 `json:"power_factor"`
			GridFrequency  float64 `json:"elec_freq"`
			ActivePower    float64 `json:"active_power"`
			ReactivePower  float64 `json:"reactive_power"`
			DayPower       float64 `json:"day_cap"`
			MPPTPower      float64 `json:"mppt_power"`
			PV1Voltage     float64 `json:"pv1_u"`
			PV2Voltage     float64 `json:"pv2_u"`
			PV3Voltage     float64 `json:"pv3_u"`
			PV4Voltage     float64 `json:"pv4_u"`
			PV5Voltage     float64 `json:"pv5_u"`
			PV6Voltage     float64 `json:"pv6_u"`
			PV7Voltage     float64 `json:"pv7_u"`
			PV8Voltage     float64 `json:"pv8_u"`
			PV1Current     float64 `json:"pv1_i"`
			PV2Current     float64 `json:"pv2_i"`
			PV3Current     float64 `json:"pv3_i"`
			PV4Current     float64 `json:"pv4_i"`
			PV5Current     float64 `json:"pv5_i"`
			PV6Current     float64 `json:"pv6_i"`
			PV7Current     float64 `json:"pv7_i"`
			PV8Current     float64 `json:"pv8_i"`
			TotalPower     float64 `json:"total_cap"`
			StartupTime    int64   `json:"open_time"`
			ShutdownTime   int64   `json:"close_time"`
			MPPT1Power     float64 `json:"mppt_1_cap"`
			MPPT2Power     float64 `json:"mppt_2_cap"`
			MPPT3Power     float64 `json:"mppt_3_cap"`
			MPPT4Power     float64 `json:"mppt_4_cap"`
			Status         int     `json:"run_state"`
		} `json:"dataItemMap"`
		DeviceID int64 `json:"devId"`
	} `json:"data"`
}
