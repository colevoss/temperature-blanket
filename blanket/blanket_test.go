package blanket

import (
	"testing"

	"github.com/colevoss/temperature-blanket/messenger"
	"github.com/colevoss/temperature-blanket/synoptic"
	// "github.com/colevoss/temperature-blanket/twilio"
)

func TestBlanket(t *testing.T) {
	synopticApi := synoptic.New()
	// m := twilio.New()
	m := messenger.NewMockMessenger()

	blanket := NewTemperatureBlanket(synopticApi, m)

	blanket.DoIt()
}
