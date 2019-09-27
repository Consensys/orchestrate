package steps

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/chain"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/envelope"
)

func TestGetChainCounts(t *testing.T) {
	envelopes := map[string]*envelope.Envelope{
		"1": &envelope.Envelope{
			Chain: chain.FromInt(888),
		},
		"2": &envelope.Envelope{
			Chain: chain.FromInt(777),
		},
		"3": &envelope.Envelope{
			Chain: chain.FromInt(888),
		},
	}

	counts := GetChainCounts(envelopes)
	expected := map[string]uint{"777": 1, "888": 2}

	assert.Equal(t, counts, expected, "Should be equal")
}

func TestChanTimeout(t *testing.T) {

	testChan := make(chan *envelope.Envelope)

	e, err := ReadChanWithTimeout(testChan, 1, 1)
	assert.Nil(t, e, "Should get an nil envelope slice")
	assert.Error(t, err, "Should get an error for not having received an envelope")

	testEnvelope := &envelope.Envelope{Chain: chain.FromInt(888)}
	go func() {
		testChan <- testEnvelope
	}()
	e, err = ReadChanWithTimeout(testChan, 1, 1)
	assert.NoError(t, err, "Should not get an error")
	assert.Equal(t, testEnvelope.GetChain().GetId(), e[0].GetChain().GetId(), "Should be the same envelope")
}

func TestSendEnvelope(t *testing.T) {

	producer := mocks.NewSyncProducer(t, nil)
	producer.ExpectSendMessageAndSucceed()
	producer.ExpectSendMessageAndFail(sarama.ErrOutOfBrokers)
	broker.SetGlobalSyncProducer(producer)

	e := &envelope.Envelope{
		Chain: chain.FromInt(888),
	}

	err := SendEnvelope(e, "topic-tx-crafter")
	assert.NoError(t, err, "Should not get an error")
	err = SendEnvelope(e, "topic-tx-crafter")
	assert.Error(t, err, "Should get an error")
}
