# Temperature Blanket

Small Go lambda that sends the previous day's High, Low, and Average temperature
to configured phone numbers via text for the use of crocheting a temperature blanket.

Currently only works for the Lincoln, NE.

It will send the text message every morning at 10 AM (CST)

# Build

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
```

# Zip
```bash
zip main.zip main
```

# Test

Phone numbers should be in the `1112223333` format where 1's are the area code, 2's are the
prefix and 3's are the final four. No need to include `+1`

```bash
TB_PHONE_NUMBERS="<COMMA DELIMTED LIST OF PHONE NUMBERS>" go test ./blanket -v
```
