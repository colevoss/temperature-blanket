# Temperature Blanket

Small Go lambda that sends the previous day's High, Low, and Average temperature
to configured phone numbers via text for the use of crocheting a temperature blanket.

Currently only works for the Lincoln, NE.

It will send the text message every morning at 10 AM (CST)

## Integrations

### [Synoptic Weather API](https://synopticdata.com/)

The [Synoptic Mesonet Timeseries API](https://developers.synopticdata.com/mesonet/) is
used for gathering the historical air temperature for the previous day.

### Twilio

Text messages are sent using Twilio's SMS service

## Build

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
```

## Zip
```bash
zip main.zip main
```

## Test

Phone numbers should be in the `1112223333` format where 1's are the area code, 2's are the
prefix and 3's are the final four. No need to include `+1`

```bash
TB_PHONE_NUMBERS="<COMMA DELIMTED LIST OF PHONE NUMBERS>" go test ./blanket -v
```

## Lambda

Currently this lambda is manually deployed using the AWS Lambda console

## Env Variables

The following environment variables are necessary for the lambda to function properly

* `SYNOPTIC_API_TOKEN` - Private API token for access to the Synoptic weather API
* `TWILIO_ACCOUNT_SID` - Twilio account ID
* `TWILIO_API_TOKEN` - Private Twilio API Token
* `TWILIO_MESSAGE_SERVICE_ID` - Twilio Temperature Blanket Message Service ID
* `TB_PHONE_NUMBERS` - Comma delimited list of phone numbers to send the text to
