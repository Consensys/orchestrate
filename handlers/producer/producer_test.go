package producer

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine/testutils"
)

const (
	caseCount = 2
	rounds    = 10
)

func makeProducerContext(i int, p *mocks.SyncProducer) *engine.TxContext {
	txctx := engine.NewTxContext().Prepare(log.NewEntry(log.StandardLogger()), nil)

	switch i % caseCount {
	case 0:
		txctx.Set("topic", "test-topic")
		p.ExpectSendMessageAndSucceed()
	default:
		txctx.Set("topic", "")
	}
	return txctx
}

func PrepareMessage(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	msg.Topic = txctx.Get("topic").(string)
	return nil
}

type ProducerTestSuite struct {
	testutils.HandlerTestSuite
	p *mocks.SyncProducer
}

func (s *ProducerTestSuite) SetupSuite() {
	s.p = mocks.NewSyncProducer(s.T(), nil)
	s.Handler = Producer(s.p, PrepareMessage)
}

func (s *ProducerTestSuite) TestProducer() {
	var txctxSlice []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxSlice = append(txctxSlice, makeProducerContext(i, s.p))
	}

	s.Handle(txctxSlice)
	err := s.p.Close() // This will error if messages have not been sent properly
	assert.Nil(s.T(), err, "Message should have been sent")
}

func TestProducer(t *testing.T) {
	suite.Run(t, new(ProducerTestSuite))
}
