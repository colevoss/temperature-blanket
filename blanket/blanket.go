package blanket

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/colevoss/temperature-blanket/messenger"
	"github.com/colevoss/temperature-blanket/weather"
)

type TemperatureBlanket struct {
	weather   weather.Weather
	messenger messenger.Messenger
}

func NewTemperatureBlanket(weather weather.Weather, messenger messenger.Messenger) *TemperatureBlanket {
	return &TemperatureBlanket{
		weather,
		messenger,
	}
}

func (t *TemperatureBlanket) GetPhoneNumbers() ([]string, bool) {
	envNumbers, present := os.LookupEnv("TB_PHONE_NUMBERS")

	if !present {
		log.Printf("No numbers are set in the environment")
		return nil, false
	}

	numbers := strings.Split(envNumbers, ",")

	return numbers, true
}

func (t *TemperatureBlanket) DoIt() {
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

	numbers, present := t.GetPhoneNumbers()

	if !present {
		log.Printf("No numbers present. Nothing to send")
		return
	}

	for _, number := range numbers {
		err = t.messenger.SendMessage("+1"+number, message)

		if err != nil {
			log.Printf("Error sending message %s", err)
		}
	}
}
