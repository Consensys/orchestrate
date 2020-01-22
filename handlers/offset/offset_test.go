package offset

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
)

var (
	session = mocks.NewConsumerGroupSession(context.Background(), "test-group", make(map[string][]int32))
	c       = mocks.NewConsumerGroupClaim("test-topic", 0, 0)
)

func makeMarkerContext(i int) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare(log.NewEntry(log.StandardLogger()), nil)
	ctx := broker.WithConsumerGroupSessionAndClaim(context.Background(), session, c)
	txctx.WithContext(ctx)

	txctx.In = &broker.Msg{
		ConsumerMessage: sarama.ConsumerMessage{
			Topic:     "test-topic",
			Partition: 0,
			Offset:    int64(i),
		},
	}

	return txctx
}

type MarkerTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *MarkerTestSuite) SetupSuite() {
	s.Handler = Marker
}

func (s *MarkerTestSuite) TestMarker() {
	rounds := 100
	var txctxs []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeMarkerContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	assert.Equal(s.T(), int64(rounds), session.LastMarkedOffset("test-topic", 0).Offset, "Expected message to have been marked")
}

func TestLoader(t *testing.T) {
	suite.Run(t, new(MarkerTestSuite))
}
