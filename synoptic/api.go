package synoptic

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/colevoss/temperature-blanket/weather"
)

var SYNOPTIC_API_TOKEN string
var SYNOPTIC_API_URL *url.URL

type SynopticApi struct {
}

func New() *SynopticApi {
	return &SynopticApi{}
}

func (s *SynopticApi) GetPreviousDaysWeatherInfo(day time.Time) (*weather.WeatherInfo, error) {
	start, end := s.GetPreviousDay(day)

	timeseriesData, err := s.GetTemparatureData(start.UTC(), end.UTC())

	if err != nil {
		return nil, err
	}

	station := timeseriesData.Station[0]
	temps := station.Observations.AirTemp

	total := 0.0
	high := 0.0
	low := 0.0

	for i, temp := range temps {
		fTemp := weather.CelciusToFahrenheit(temp)

		if i == 0 {
			high = fTemp
			low = fTemp
		}

		if fTemp > high {
			high = fTemp
		}

		if fTemp < low {
			low = fTemp
		}

		total += fTemp
	}

	avg := total / float64(len(temps))

	weatherInfo := &weather.WeatherInfo{
		Date:    start,
		High:    high,
		Low:     low,
		Average: avg,
	}

	return weatherInfo, nil
}

/**
 * Given a day, return 00:00 and 23:59 of the previous day
 */
func (s *SynopticApi) GetPreviousDay(day time.Time) (time.Time, time.Time) {
	tz, err := time.LoadLocation("America/Chicago")

	if err != nil {
		log.Fatalf("Cannot load timezone %s", err)
	}

	// now := time.Now()
	yesterday := day.Add(time.Hour * -24)

	startOfYesterday := time.Date(
		yesterday.Year(),
		yesterday.Month(),
		yesterday.Day(),
		0,
		0,
		0,
		0,
		tz,
	)

	endOfYesterday := startOfYesterday.Add(time.Hour * 23).Add(time.Minute * 59)

	return startOfYesterday, endOfYesterday
}

/**
 * Makes request to get temperature data
 * @see https://developers.synopticdata.com/mesonet/v2/stations/timeseries/
 */
func (s *SynopticApi) GetTemparatureData(start time.Time, end time.Time) (*SynopticTimeSeriesResponse, error) {
	url, err := url.Parse("https://api.synopticdata.com/v2/stations/timeseries")

	if err != nil {
		log.Printf("Could not parse url %v", err)
		return nil, err
	}

	query := SYNOPTIC_API_URL.Query()

	query.Add("token", SYNOPTIC_API_TOKEN)
	query.Add("stid", "klnk")
	query.Add("vars", "air_temp")

	log.Printf("Date: %v - %v", start, end)

	// See format docs here https://pkg.go.dev/time#Time.Format
	formattedStart := start.Format("200601021504")
	formattedEnd := end.Format("200601021504")

	query.Add("start", formattedStart)
	query.Add("end", formattedEnd)

	url.RawQuery = query.Encode()

	log.Printf("Making Request to %s", url.String())

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		log.Printf("Could not create request %s", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Error making request: %s", err)
		return nil, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	var timeSeriesResponse SynopticTimeSeriesResponse
	err = json.Unmarshal(resBody, &timeSeriesResponse)

	if err != nil {
		log.Printf("Could not parse response body %s", err)
		return nil, err
	}

	// log.Printf("%v", timeSeriesResponse.Station[0].Observations.DateTime[0])
	// log.Printf("%v", timeSeriesResponse.Station[0].Name)

	return &timeSeriesResponse, nil
}

func Test() {
	var res SynopticTimeSeriesResponse
	err := json.Unmarshal([]byte(TEST_DATA), &res)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("%v", res.Station[0].Elevation)
}

func init() {
	SYNOPTIC_API_TOKEN = os.Getenv("SYNOPTIC_API_TOKEN")
	synopticUrl, err := url.Parse("https://api.synopticdata.com/v2/stations/timeseries")

	if err != nil {
		log.Fatalf("Cannot parse url %v", err)
		return
	}

	SYNOPTIC_API_URL = synopticUrl
}

const TEST_DATA = `
{
    "UNITS": {
        "position": "m",
        "elevation": "ft",
        "air_temp": "Celsius"
    },
    "QC_SUMMARY": {
        "QC_CHECKS_APPLIED": [
            "sl_range_check"
        ],
        "TOTAL_OBSERVATIONS_FLAGGED": 0.0,
        "PERCENT_OF_TOTAL_OBSERVATIONS_FLAGGED": 0.0
    },
    "STATION": [
        {
            "STATUS": "ACTIVE",
            "MNET_ID": "1",
            "PERIOD_OF_RECORD": {
                "start": "2002-04-29T00:00:00Z",
                "end": "2023-01-11T19:55:00Z"
            },
            "ELEVATION": "1214",
            "NAME": "Lincoln, Lincoln Municipal Airport",
            "STID": "KLNK",
            "SENSOR_VARIABLES": {
                "date_time": {
                    "date_time": {}
                },
                "air_temp": {
                    "air_temp_set_1": {
                        "position": "2.0"
                    }
                }
            },
            "ELEV_DEM": "1145",
            "LONGITUDE": "-96.76444",
            "UNITS": {
                "position": "m",
                "elevation": "ft"
            },
            "STATE": "NE",
            "OBSERVATIONS": {
                "date_time": [
                    "2023-01-10T06:00:00Z",
                    "2023-01-10T06:05:00Z",
                    "2023-01-10T06:10:00Z",
                    "2023-01-10T06:15:00Z",
                    "2023-01-10T06:20:00Z",
                    "2023-01-10T06:25:00Z",
                    "2023-01-10T06:30:00Z",
                    "2023-01-10T06:35:00Z",
                    "2023-01-10T06:40:00Z",
                    "2023-01-10T06:45:00Z",
                    "2023-01-10T06:50:00Z",
                    "2023-01-10T06:54:00Z",
                    "2023-01-10T06:55:00Z",
                    "2023-01-10T07:00:00Z",
                    "2023-01-10T07:10:00Z",
                    "2023-01-10T07:15:00Z",
                    "2023-01-10T07:20:00Z",
                    "2023-01-10T07:25:00Z",
                    "2023-01-10T07:30:00Z",
                    "2023-01-10T07:35:00Z",
                    "2023-01-10T07:40:00Z",
                    "2023-01-10T07:45:00Z",
                    "2023-01-10T07:50:00Z",
                    "2023-01-10T07:54:00Z",
                    "2023-01-10T07:55:00Z",
                    "2023-01-10T08:05:00Z",
                    "2023-01-10T08:10:00Z",
                    "2023-01-10T08:15:00Z",
                    "2023-01-10T08:20:00Z",
                    "2023-01-10T08:25:00Z",
                    "2023-01-10T08:30:00Z",
                    "2023-01-10T08:40:00Z",
                    "2023-01-10T08:45:00Z",
                    "2023-01-10T08:50:00Z",
                    "2023-01-10T08:54:00Z",
                    "2023-01-10T08:55:00Z",
                    "2023-01-10T09:00:00Z",
                    "2023-01-10T09:05:00Z",
                    "2023-01-10T09:10:00Z",
                    "2023-01-10T09:15:00Z",
                    "2023-01-10T09:20:00Z",
                    "2023-01-10T09:25:00Z",
                    "2023-01-10T09:30:00Z",
                    "2023-01-10T09:35:00Z",
                    "2023-01-10T09:40:00Z",
                    "2023-01-10T09:45:00Z",
                    "2023-01-10T09:50:00Z",
                    "2023-01-10T09:54:00Z",
                    "2023-01-10T10:00:00Z",
                    "2023-01-10T10:05:00Z",
                    "2023-01-10T10:10:00Z",
                    "2023-01-10T10:15:00Z",
                    "2023-01-10T10:20:00Z",
                    "2023-01-10T10:25:00Z",
                    "2023-01-10T10:30:00Z",
                    "2023-01-10T10:35:00Z",
                    "2023-01-10T10:40:00Z",
                    "2023-01-10T10:45:00Z",
                    "2023-01-10T10:54:00Z",
                    "2023-01-10T10:55:00Z",
                    "2023-01-10T11:05:00Z",
                    "2023-01-10T11:10:00Z",
                    "2023-01-10T11:15:00Z",
                    "2023-01-10T11:20:00Z",
                    "2023-01-10T11:25:00Z",
                    "2023-01-10T11:30:00Z",
                    "2023-01-10T11:35:00Z",
                    "2023-01-10T11:40:00Z",
                    "2023-01-10T11:50:00Z",
                    "2023-01-10T11:54:00Z",
                    "2023-01-10T11:55:00Z",
                    "2023-01-10T12:00:00Z",
                    "2023-01-10T12:05:00Z",
                    "2023-01-10T12:10:00Z",
                    "2023-01-10T12:15:00Z",
                    "2023-01-10T12:25:00Z",
                    "2023-01-10T12:30:00Z",
                    "2023-01-10T12:35:00Z",
                    "2023-01-10T12:40:00Z",
                    "2023-01-10T12:45:00Z",
                    "2023-01-10T12:50:00Z",
                    "2023-01-10T12:54:00Z",
                    "2023-01-10T13:00:00Z",
                    "2023-01-10T13:05:00Z",
                    "2023-01-10T13:10:00Z",
                    "2023-01-10T13:15:00Z",
                    "2023-01-10T13:20:00Z",
                    "2023-01-10T13:25:00Z",
                    "2023-01-10T13:30:00Z",
                    "2023-01-10T13:35:00Z",
                    "2023-01-10T13:40:00Z",
                    "2023-01-10T13:45:00Z",
                    "2023-01-10T13:50:00Z",
                    "2023-01-10T13:54:00Z",
                    "2023-01-10T13:55:00Z",
                    "2023-01-10T14:05:00Z",
                    "2023-01-10T14:10:00Z",
                    "2023-01-10T14:15:00Z",
                    "2023-01-10T14:20:00Z",
                    "2023-01-10T14:30:00Z",
                    "2023-01-10T14:40:00Z",
                    "2023-01-10T14:45:00Z",
                    "2023-01-10T14:50:00Z",
                    "2023-01-10T14:54:00Z",
                    "2023-01-10T14:55:00Z",
                    "2023-01-10T15:00:00Z",
                    "2023-01-10T15:05:00Z",
                    "2023-01-10T15:10:00Z",
                    "2023-01-10T15:15:00Z",
                    "2023-01-10T15:20:00Z",
                    "2023-01-10T15:25:00Z",
                    "2023-01-10T15:35:00Z",
                    "2023-01-10T15:45:00Z",
                    "2023-01-10T15:50:00Z",
                    "2023-01-10T15:54:00Z",
                    "2023-01-10T15:55:00Z",
                    "2023-01-10T16:00:00Z",
                    "2023-01-10T16:05:00Z",
                    "2023-01-10T16:10:00Z",
                    "2023-01-10T16:20:00Z",
                    "2023-01-10T16:25:00Z",
                    "2023-01-10T16:30:00Z",
                    "2023-01-10T16:35:00Z",
                    "2023-01-10T16:40:00Z",
                    "2023-01-10T16:45:00Z",
                    "2023-01-10T16:54:00Z",
                    "2023-01-10T16:55:00Z",
                    "2023-01-10T17:00:00Z",
                    "2023-01-10T17:05:00Z",
                    "2023-01-10T17:15:00Z",
                    "2023-01-10T17:20:00Z",
                    "2023-01-10T17:25:00Z",
                    "2023-01-10T17:35:00Z",
                    "2023-01-10T17:40:00Z",
                    "2023-01-10T17:45:00Z",
                    "2023-01-10T17:50:00Z",
                    "2023-01-10T17:54:00Z",
                    "2023-01-10T17:55:00Z",
                    "2023-01-10T18:00:00Z",
                    "2023-01-10T18:05:00Z",
                    "2023-01-10T18:10:00Z",
                    "2023-01-10T18:15:00Z",
                    "2023-01-10T18:20:00Z",
                    "2023-01-10T18:25:00Z",
                    "2023-01-10T18:30:00Z",
                    "2023-01-10T18:35:00Z",
                    "2023-01-10T18:40:00Z",
                    "2023-01-10T18:45:00Z",
                    "2023-01-10T18:50:00Z",
                    "2023-01-10T18:54:00Z",
                    "2023-01-10T18:55:00Z",
                    "2023-01-10T19:00:00Z",
                    "2023-01-10T19:10:00Z",
                    "2023-01-10T19:15:00Z",
                    "2023-01-10T19:20:00Z",
                    "2023-01-10T19:25:00Z",
                    "2023-01-10T19:30:00Z",
                    "2023-01-10T19:35:00Z",
                    "2023-01-10T19:40:00Z",
                    "2023-01-10T19:50:00Z",
                    "2023-01-10T19:54:00Z",
                    "2023-01-10T19:55:00Z",
                    "2023-01-10T20:00:00Z",
                    "2023-01-10T20:05:00Z",
                    "2023-01-10T20:10:00Z",
                    "2023-01-10T20:15:00Z",
                    "2023-01-10T20:20:00Z",
                    "2023-01-10T20:25:00Z",
                    "2023-01-10T20:30:00Z",
                    "2023-01-10T20:35:00Z",
                    "2023-01-10T20:40:00Z",
                    "2023-01-10T20:45:00Z",
                    "2023-01-10T20:50:00Z",
                    "2023-01-10T20:54:00Z",
                    "2023-01-10T20:55:00Z",
                    "2023-01-10T21:00:00Z",
                    "2023-01-10T21:10:00Z",
                    "2023-01-10T21:15:00Z",
                    "2023-01-10T21:20:00Z",
                    "2023-01-10T21:25:00Z",
                    "2023-01-10T21:30:00Z",
                    "2023-01-10T21:35:00Z",
                    "2023-01-10T21:40:00Z",
                    "2023-01-10T21:45:00Z",
                    "2023-01-10T21:50:00Z",
                    "2023-01-10T21:54:00Z",
                    "2023-01-10T21:55:00Z",
                    "2023-01-10T22:00:00Z",
                    "2023-01-10T22:05:00Z",
                    "2023-01-10T22:10:00Z",
                    "2023-01-10T22:15:00Z",
                    "2023-01-10T22:20:00Z",
                    "2023-01-10T22:25:00Z",
                    "2023-01-10T22:30:00Z",
                    "2023-01-10T22:35:00Z",
                    "2023-01-10T22:40:00Z",
                    "2023-01-10T22:50:00Z",
                    "2023-01-10T22:54:00Z",
                    "2023-01-10T23:00:00Z",
                    "2023-01-10T23:10:00Z",
                    "2023-01-10T23:15:00Z",
                    "2023-01-10T23:20:00Z",
                    "2023-01-10T23:25:00Z",
                    "2023-01-10T23:30:00Z",
                    "2023-01-10T23:35:00Z",
                    "2023-01-10T23:40:00Z",
                    "2023-01-10T23:45:00Z",
                    "2023-01-10T23:50:00Z",
                    "2023-01-10T23:54:00Z",
                    "2023-01-10T23:55:00Z",
                    "2023-01-11T00:00:00Z",
                    "2023-01-11T00:05:00Z",
                    "2023-01-11T00:10:00Z",
                    "2023-01-11T00:15:00Z",
                    "2023-01-11T00:20:00Z",
                    "2023-01-11T00:30:00Z",
                    "2023-01-11T00:35:00Z",
                    "2023-01-11T00:40:00Z",
                    "2023-01-11T00:45:00Z",
                    "2023-01-11T00:50:00Z",
                    "2023-01-11T00:54:00Z",
                    "2023-01-11T01:00:00Z",
                    "2023-01-11T01:05:00Z",
                    "2023-01-11T01:10:00Z",
                    "2023-01-11T01:15:00Z",
                    "2023-01-11T01:20:00Z",
                    "2023-01-11T01:25:00Z",
                    "2023-01-11T01:30:00Z",
                    "2023-01-11T01:35:00Z",
                    "2023-01-11T01:40:00Z",
                    "2023-01-11T01:45:00Z",
                    "2023-01-11T01:50:00Z",
                    "2023-01-11T01:54:00Z",
                    "2023-01-11T02:00:00Z",
                    "2023-01-11T02:05:00Z",
                    "2023-01-11T02:10:00Z",
                    "2023-01-11T02:15:00Z",
                    "2023-01-11T02:20:00Z",
                    "2023-01-11T02:25:00Z",
                    "2023-01-11T02:30:00Z",
                    "2023-01-11T02:35:00Z",
                    "2023-01-11T02:40:00Z",
                    "2023-01-11T02:45:00Z",
                    "2023-01-11T02:50:00Z",
                    "2023-01-11T02:54:00Z",
                    "2023-01-11T02:55:00Z",
                    "2023-01-11T03:00:00Z",
                    "2023-01-11T03:05:00Z",
                    "2023-01-11T03:10:00Z",
                    "2023-01-11T03:15:00Z",
                    "2023-01-11T03:20:00Z",
                    "2023-01-11T03:30:00Z",
                    "2023-01-11T03:35:00Z",
                    "2023-01-11T03:40:00Z",
                    "2023-01-11T03:45:00Z",
                    "2023-01-11T03:50:00Z",
                    "2023-01-11T03:54:00Z",
                    "2023-01-11T03:55:00Z",
                    "2023-01-11T04:00:00Z",
                    "2023-01-11T04:05:00Z",
                    "2023-01-11T04:15:00Z",
                    "2023-01-11T04:20:00Z",
                    "2023-01-11T04:25:00Z",
                    "2023-01-11T04:30:00Z",
                    "2023-01-11T04:35:00Z",
                    "2023-01-11T04:40:00Z",
                    "2023-01-11T04:50:00Z",
                    "2023-01-11T04:54:00Z",
                    "2023-01-11T04:55:00Z",
                    "2023-01-11T05:00:00Z",
                    "2023-01-11T05:10:00Z",
                    "2023-01-11T05:15:00Z",
                    "2023-01-11T05:20:00Z",
                    "2023-01-11T05:25:00Z",
                    "2023-01-11T05:30:00Z",
                    "2023-01-11T05:35:00Z",
                    "2023-01-11T05:45:00Z",
                    "2023-01-11T05:50:00Z",
                    "2023-01-11T05:54:00Z",
                    "2023-01-11T05:55:00Z"
                ],
                "air_temp_set_1": [
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -3.0,
                    -2.0,
                    -2.0,
                    -1.7,
                    -2.0,
                    -2.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.3,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -4.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -3.9,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -3.0,
                    -2.8,
                    -3.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -4.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -3.9,
                    -4.0,
                    -4.0,
                    -4.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -3.3,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -4.0,
                    -4.0,
                    -3.0,
                    -4.0,
                    -3.0,
                    -3.0,
                    -3.0,
                    -3.3,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -5.0,
                    -5.0,
                    -5.0,
                    -4.0,
                    -4.0,
                    -4.0,
                    -5.0,
                    -5.0,
                    -4.0,
                    -3.0,
                    -3.0,
                    -2.2,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -1.0,
                    -1.0,
                    -1.0,
                    -1.0,
                    -0.6,
                    -1.0,
                    0.0,
                    0.0,
                    0.0,
                    1.0,
                    1.0,
                    1.0,
                    1.0,
                    2.0,
                    2.0,
                    2.2,
                    2.0,
                    2.0,
                    2.0,
                    3.0,
                    3.0,
                    4.0,
                    4.0,
                    4.0,
                    4.0,
                    4.0,
                    4.4,
                    5.0,
                    6.0,
                    6.0,
                    6.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    8.0,
                    8.0,
                    8.0,
                    8.9,
                    9.0,
                    9.0,
                    9.0,
                    10.0,
                    10.0,
                    11.0,
                    10.0,
                    11.0,
                    11.0,
                    12.0,
                    11.7,
                    12.0,
                    12.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    13.0,
                    12.8,
                    13.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    11.0,
                    11.1,
                    11.0,
                    11.0,
                    11.0,
                    11.0,
                    10.0,
                    9.0,
                    9.0,
                    9.0,
                    9.0,
                    9.0,
                    9.0,
                    8.3,
                    8.0,
                    8.0,
                    8.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    6.7,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    7.0,
                    6.0,
                    5.0,
                    5.0,
                    6.0,
                    6.0,
                    6.1,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.1,
                    6.0,
                    7.0,
                    6.0,
                    7.0,
                    7.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    6.0,
                    5.0,
                    5.0,
                    5.0,
                    5.0,
                    5.0,
                    6.0,
                    6.0,
                    6.0,
                    4.0,
                    3.0,
                    5.0,
                    4.0,
                    4.0,
                    4.4,
                    4.0,
                    3.0,
                    1.0,
                    1.0,
                    -1.0,
                    1.0,
                    -1.0,
                    0.0,
                    2.0,
                    1.0,
                    1.1,
                    1.0,
                    1.0,
                    -1.0,
                    0.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.0,
                    -2.8,
                    -3.0
                ]
            },
            "RESTRICTED": false,
            "QC_FLAGGED": false,
            "LATITUDE": "40.83111",
            "TIMEZONE": "America/Chicago",
            "ID": "3912"
        }
    ],
    "SUMMARY": {
        "DATA_QUERY_TIME": "3.10397148132 ms",
        "RESPONSE_CODE": 1,
        "RESPONSE_MESSAGE": "OK",
        "METADATA_RESPONSE_TIME": "93.9450263977 ms",
        "DATA_PARSING_TIME": "6.46209716797 ms",
        "VERSION": "v2.17.0",
        "TOTAL_DATA_TIME": "9.56606864929 ms",
        "NUMBER_OF_OBJECTS": 1,
        "FUNCTION_USED": "time_data_parser"
    }
}
`
