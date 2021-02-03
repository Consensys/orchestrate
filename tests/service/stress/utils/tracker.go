package utils

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/tracker"
)

var Topics = [...]string{
	"tx.decoded",
	"tx.recover",
}

func NewEnvelopeTracker(chanReg *chanregistry.ChanRegistry, e *tx.Envelope, testID string) *tracker.Tracker {
	// Prepare envelope metadata
	if testID != "" {
		_ = e.SetContextLabelsValue("id", testID)
	}
	// Set envelope metadata so it can be tracked

	// Create tracker and attach envelope
	t := tracker.NewTracker()
	t.Current = e

	// Initialize output channels on tracker and register channels on channel registry
	for _, topic := range Topics {
		ckey := utils.LongKeyOf(topic, testID)
		var ch = make(chan *tx.Envelope, 10)
		// Register channel on channel registry
		log.WithFields(log.Fields{
			"id":     ckey,
			"testId": testID,
			"topic":  topic,
		}).Debugf("tracker: registered new envelope channel")
		chanReg.Register(ckey, ch)

		// Add channel as a tracker output
		t.AddOutput(topic, ch)
	}

	return t
}

func WaitForEnvelope(t *tracker.Tracker, d time.Duration) error {
	cerr := make(chan error, 1)

	go func() {
		cerr <- t.Load("tx.decoded", d)
	}()
	go func() {
		err := t.Load("tx.recover", d)
		if err != nil {
			cerr <- err
		}
		cerr <- fmt.Errorf("tx.recover: %s", t.Current.Error())
	}()

	return <-cerr
}
