package utils

import (
	"fmt"
	"time"

	"github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/tests/utils"
	"github.com/consensys/orchestrate/tests/utils/chanregistry"
	"github.com/consensys/orchestrate/tests/utils/tracker"
)

var Topics = map[string]string{
	utils.TxDecodedTopicKey: sarama.TxDecodedViperKey,
	utils.TxRecoverTopicKey: sarama.TxRecoverViperKey,
}

func NewEnvelopeTracker(chanReg *chanregistry.ChanRegistry, e *tx.Envelope, testID string) *tracker.Tracker {
	logger := log.NewLogger().SetComponent("stress-test.tracker")
	// Prepare envelope metadata
	if testID != "" {
		_ = e.SetContextLabelsValue("id", testID)
	}
	// Set envelope metadata so it can be tracked

	// Create tracker and attach envelope
	t := tracker.NewTracker()
	t.Current = e

	// Initialize output channels on tracker and register channels on channel registry
	for topic := range Topics {
		ckey := utils.LongKeyOf(topic, testID)
		var ch = make(chan *tx.Envelope, 10)
		// Register channel on channel registry
		logger.WithField("id", ckey).WithField("testId", testID).
			WithField("topic", topic).Debug("registered new envelope channel")
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
