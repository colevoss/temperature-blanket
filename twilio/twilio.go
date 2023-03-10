package twilio

import (
	"encoding/json"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"log"
	"os"
)

var TWILIO_ACCOUNT_SID string
var TWILIO_API_TOKEN string
var TWILIO_MESSAGE_SERVICE_ID string

type Twilio struct {
}

func New() *Twilio {
	return &Twilio{}
}

func (t *Twilio) SendMessage(to string, message string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_API_TOKEN,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetMessagingServiceSid(TWILIO_MESSAGE_SERVICE_ID)
	params.SetBody(message)

	log.Println("Sending message to:", to)
	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	response, _ := json.Marshal(*resp)
	log.Println("Twilio response:", string(response))
	return nil
}

func init() {
	TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_API_TOKEN = os.Getenv("TWILIO_API_TOKEN")
	TWILIO_MESSAGE_SERVICE_ID = os.Getenv("TWILIO_MESSAGE_SERVICE_ID")
}
