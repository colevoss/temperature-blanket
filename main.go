package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/colevoss/temperature-blanket/messenger"
	"github.com/colevoss/temperature-blanket/synoptic"

	// "github.com/colevoss/temperature-blanket/twilio"
	"github.com/colevoss/temperature-blanket/weather"
)

type temperatureBlanket struct {
	weather   weather.Weather
	messenger messenger.Messenger
}

func NewTemperatureBlanket(weather weather.Weather, messenger messenger.Messenger) *temperatureBlanket {
	return &temperatureBlanket{
		weather,
		messenger,
	}
}

func (t *temperatureBlanket) DoIt() {
	weatherInfo, err := t.weather.GetPreviousDaysWeatherInfo(time.Now())

	if err != nil {
		log.Printf("Bad thigns: %s", err)
	}

	formattedDate := weatherInfo.Date.Format("Jan 2 2006")

	message := fmt.Sprintf(
		"\nWeather for %s:\n\u2600\ufe0f High: %.0f°\n\u2744\ufe0f Low: %.0f°\n\U0001f600 Avg: %.0f°",
		formattedDate,
		math.Ceil(weatherInfo.High),
		math.Ceil(weatherInfo.Low),
		math.Ceil(weatherInfo.Average),
	)

	err = t.messenger.SendMessage("+14027209808", message)

	if err != nil {
		log.Printf("Error sending message %s", err)
	}
}

func main() {
	synopticApi := synoptic.New()
	// m := twilio.New()
	m := messenger.NewMockMessenger()

	blanket := NewTemperatureBlanket(synopticApi, m)

	blanket.DoIt()
}
