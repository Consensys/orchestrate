package utils

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/e2e/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/tracker"
)

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
	for _, topic := range utils2.TOPICS {
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
