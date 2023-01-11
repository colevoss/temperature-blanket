package main

import (
	"log"

	"github.com/colevoss/temperature-blanket/synoptic"
)

func main() {
	log.Printf("HELLO")
	tempData, err := synoptic.GetTemparatureData()

	if err != nil {
		log.Printf("Error getting temperature data %s", err)
		return
	}

	if tempData.Station == nil || len(tempData.Station) == 0 {
		log.Printf("No station data returned")
		return
	}

	station := tempData.Station[0]
	temps := station.Observations.AirTemp

	total := 0.0
	high := 0.0
	low := 0.0

	for i, temp := range temps {
		fTemp := (temp * 1.8) + 32

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

	log.Printf("Average temp: %f degrees", avg)
	log.Printf("Low temp: %f degrees", low)
	log.Printf("High temp: %f degrees", high)
}
