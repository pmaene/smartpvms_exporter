//go:generate go-enum --file $GOFILE --lower --marshal --template ../../.go-enum/enum.tmpl
package smartpvms

import (
	"encoding/json"
	"time"
)

const (
	XSRFTokenRefreshInterval = 30*time.Minute - 15*time.Second
)

/*
ENUM(
StringInverter = 1
SmartLogger = 2
Transformer = 8
EMI = 10
ProtocolConverter = 13
GeneralDevice = 16
GridMeter = 17
PID = 22
PinnetDataLogger = 37
ResidentialInverter = 38
Battery = 39
BackupBox = 40
PLC = 45
Optimizer = 46
PowerSensor = 47
Dongle = 62
DistributedSmartLogger = 63
SafetyBox = 70
)
*/
type DeviceType int

/*
ENUM(
Disconnected
Connected
)
*/
type InverterRunStatus int

/*
ENUM(
StandbyInitializing = 0
StandbyInsulationResistanceDetection = 1
StandbySunlightDetection = 2
StandyPowerGridDetection = 3
Start = 256
GridConnection = 512
GridConnectionLimitedPower = 513
GridConnectionSelfDerating = 514
ShutdownUnexpected = 768
ShutdownCommandedShutdown = 769
ShutdownOVGR = 770
ShutdownCommunicationDisconnection = 771
ShutdownLimitedPower = 772
ShutdownManualStartupRequired = 773
ShutdownDCSwitchDisconnected = 774
GridSchedulingCosPhiPCurve = 1025
GridSchedulingQUCurve = 1026
SpotCheckReady = 1280
SpotChecking = 1281
Inspecting = 1536
AFCISelfCheck = 1792
IVScanning = 2048
DCInputDetection = 2304
StandbyNoSunlight = 40960
CommunicationDisconnection = 45056
Loading = 49152
)
*/
type InverterStatus float64

/*
ENUM(
NotConstructed
UnderConstruction
GridConnected
)
*/
type PlantBuildStatus int

/*
ENUM(
Utility = 1
CommercialAndIndustrial
Residential
)
*/
type PlantGridConnectionType int

/*
ENUM(
PovertyAlleviationPlant
NonPovertyAlleviationPlant
)
*/
type PlantAIDType int

/*
ENUM(
Disconnected = 1
Faulty
Healthy
)
*/
type PlantStatus int

type Device struct {
	Type            DeviceType `json:"devTypeId"`
	ID              int64      `json:"id"`
	Name            string     `json:"devName"`
	Serial          string     `json:"esnCode"`
	Model           string     `json:"invType"`
	SoftwareVersion string     `json:"softwareVersion"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	StationCode     string     `json:"stationCode"`
}

type Plant struct {
	StationCode        string                  `json:"stationCode"`
	Name               string                  `json:"stationName"`
	Address            string                  `json:"stationAddr"`
	Capacity           float64                 `json:"capacity"`
	BuildStatus        PlantBuildStatus        `json:"buildState"`
	GridConnectionType PlantGridConnectionType `json:"combineType"`
	AIDType            PlantAIDType            `json:"aidType"`
	ContactPerson      string                  `json:"stationLinkman"`
	ContactPersonPhone string                  `json:"linkmanPho"`
}

type PlantData struct {
	DayYield    float64     `json:"day_power"`
	MonthYield  float64     `json:"month_power"`
	TotalYield  float64     `json:"total_power"`
	DayIncome   float64     `json:"day_income"`
	TotalIncome float64     `json:"total_income"`
	Status      PlantStatus `json:"real_health_state"`
}

type ResidentialInverterData struct {
	RunStatus       InverterRunStatus `json:"run_state"`
	Status          InverterStatus    `json:"inverter_state"`
	StartupTime     time.Time         `json:"open_time"`
	ShutdownTime    time.Time         `json:"close_time"`
	Temperature     float64           `json:"temperature"`
	Efficiency      float64           `json:"efficiency"`
	PowerFactor     float64           `json:"power_factor"`
	ActivePower     float64           `json:"active_power"`
	ReactivePower   float64           `json:"reactive_power"`
	PVPower         float64           `json:"mppt_power"`
	L1Voltage       float64           `json:"a_u"`
	L2Voltage       float64           `json:"b_u"`
	L3Voltage       float64           `json:"c_u"`
	L1Current       float64           `json:"a_i"`
	L2Current       float64           `json:"b_i"`
	L3Current       float64           `json:"c_i"`
	PV1Voltage      float64           `json:"pv1_u"`
	PV2Voltage      float64           `json:"pv2_u"`
	PV3Voltage      float64           `json:"pv3_u"`
	PV4Voltage      float64           `json:"pv4_u"`
	PV5Voltage      float64           `json:"pv5_u"`
	PV6Voltage      float64           `json:"pv6_u"`
	PV7Voltage      float64           `json:"pv7_u"`
	PV8Voltage      float64           `json:"pv8_u"`
	PV1Current      float64           `json:"pv1_i"`
	PV2Current      float64           `json:"pv2_i"`
	PV3Current      float64           `json:"pv3_i"`
	PV4Current      float64           `json:"pv4_i"`
	PV5Current      float64           `json:"pv5_i"`
	PV6Current      float64           `json:"pv6_i"`
	PV7Current      float64           `json:"pv7_i"`
	PV8Current      float64           `json:"pv8_i"`
	DayYield        float64           `json:"day_cap"`
	TotalYield      float64           `json:"total_cap"`
	MPPT1TotalYield float64           `json:"mppt_1_cap"`
	MPPT2TotalYield float64           `json:"mppt_2_cap"`
	MPPT3TotalYield float64           `json:"mppt_3_cap"`
	MPPT4TotalYield float64           `json:"mppt_4_cap"`
	GridL1L2Voltage float64           `json:"ab_u"`
	GridL2L3Voltage float64           `json:"bc_u"`
	GridL3L1Voltage float64           `json:"ca_u"`
	GridFrequency   float64           `json:"elec_freq"`
}

func (i *ResidentialInverterData) UnmarshalJSON(data []byte) error {
	type Alias ResidentialInverterData

	a := &struct {
		StartupTime  int64 `json:"open_time"`
		ShutdownTime int64 `json:"close_time"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	i.StartupTime = time.Unix(0, a.StartupTime*int64(time.Millisecond))
	i.ShutdownTime = time.Unix(0, a.ShutdownTime*int64(time.Millisecond))

	return nil
}

type XSRFToken struct {
	XSRFToken string
	ExpiresAt time.Time
}

func (t *XSRFToken) IsValid() bool {
	if t == nil {
		return false
	}

	return t.ExpiresAt.After(time.Now())
}

func NewXSRFToken(t string) *XSRFToken {
	return &XSRFToken{
		XSRFToken: t,
		ExpiresAt: time.Now().Add(XSRFTokenRefreshInterval),
	}
}
