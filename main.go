package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/colevoss/temperature-blanket/blanket"
	"github.com/colevoss/temperature-blanket/synoptic"
	"github.com/colevoss/temperature-blanket/twilio"
)

func Handler(ctx context.Context) {
	synopticApi := synoptic.New()
	m := twilio.New()

	blanket := blanket.NewTemperatureBlanket(synopticApi, m)

	blanket.DoIt()
}

func main() {
	lambda.Start(Handler)
}
