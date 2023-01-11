package synoptic

import "time"

type SynopticTimeSeriesResponse struct {
	Units   *Units     `json:"UNITS"`
	Station []*Station `json:"STATION"`
	Summary *Summary   `json:"SUMMARY"`
}

type Units struct {
	AirTemp string `json:"air_temp"`
}

type Station struct {
	Status         string                 `json:"STATUS"`
	MnetId         string                 `json:"MNET_ID"`
	Longitude      string                 `json:"LONGITUDE"`
	Latitude       string                 `json:"LATITUDE"`
	Timezone       string                 `json:"TIMEZONE"`
	Id             string                 `json:"ID"`
	State          string                 `json:"STATE"`
	PeriodOfRecord map[string]interface{} `json:"PERIOD_OF_RECORD"`
	// Try setting to int
	Elevation       int32            `json:"ELEVATION,string"`
	Name            string           `json:"NAME"`
	QcFlagged       bool             `json:"QC_FLAGGED"`
	SensorVariables *SensorVariables `json:"SENSOR_VARIABLES"`
	Observations    *Observations    `json:"OBSERVATIONS"`
}

type SensorVariables struct {
	DateTime map[string]interface{} `json:"date_time"`
	AirTemp  map[string]interface{} `json:"air_temp"`
}

type Observations struct {
	DateTime []time.Time `json:"date_time"`
	AirTemp  []float64   `json:"air_temp_set_1"`
}

type Summary struct {
	NumberOfObjects int    `json:"NUMBER_OF_OBJECTS"`
	ResponseCode    int    `json:"RESPONSE_CODE"`
	ResponseMessage string `json:"RESPONSE_MESSAGE"`
	TotalDataTime   string `json:"TOTAL_DATA_TIME"`
}
