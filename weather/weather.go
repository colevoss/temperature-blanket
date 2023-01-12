package weather

import "time"

type WeatherInfo struct {
	Date    time.Time
	High    float64
	Low     float64
	Average float64
}

func CelciusToFahrenheit(celcius float64) float64 {
	return (celcius * 1.8) + 32
}

type Weather interface {
	GetPreviousDaysWeatherInfo(day time.Time) (*WeatherInfo, error)
}
