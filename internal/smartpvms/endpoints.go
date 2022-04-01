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
	Data []Plant `json:"data"`
}

type GetRealtimePlantDataBody struct {
	StationCodes []string `json:"stationCodes"`
}

func (b *GetRealtimePlantDataBody) MarshalJSON() ([]byte, error) {
	type Alias GetRealtimePlantDataBody

	return json.Marshal(&struct {
		StationCodes string `json:"stationCodes"`
		*Alias
	}{
		StationCodes: strings.Join(b.StationCodes, ","),
		Alias:        (*Alias)(b),
	})
}

type GetRealtimePlantDataResult struct {
	Result
	Data []struct {
		DataItemMap PlantData `json:"dataItemMap"`
		StationCode string    `json:"stationCode"`
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
	Data []Device `json:"data"`
}

type GetRealtimeDeviceDataBody struct {
	Type DeviceType `json:"devTypeId"`
	IDs  []int64    `json:"devIds"`
}

func (b *GetRealtimeDeviceDataBody) MarshalJSON() ([]byte, error) {
	type Alias GetRealtimeDeviceDataBody

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

type GetRealtimeDeviceDataResult[T any] struct {
	Result
	Data []struct {
		DataItemMap T     `json:"dataItemMap"`
		DeviceID    int64 `json:"devId"`
	} `json:"data"`
}
