package offset

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
)

var (
	s = mock.NewConsumerGroupSession(context.Background(), "test-group", make(map[string][]int32))
	c = mock.NewConsumerGroupClaim("test-topic", 0, 0)
)

func makeMarkerContext(i int) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare([]engine.HandlerFunc{}, log.NewEntry(log.StandardLogger()), nil)
	ctx := broker.WithConsumerGroupSessionAndClaim(context.Background(), s, c)
	txctx.WithContext(ctx)

	txctx.Msg = &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    int64(i),
	}

	return txctx
}

type MarkerTestSuite struct {
	testutils.HandlerTestSuite
}

func (suite *MarkerTestSuite) SetupSuite() {
	suite.Handler = Marker
}

func (suite *MarkerTestSuite) TestMarker() {
	rounds := 100
	txctxs := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeMarkerContext(i))
	}

	// Handle contexts
	suite.Handle(txctxs)

	assert.Equal(suite.T(), int64(rounds), s.LastMarkedOffset("test-topic", 0).Offset, "Expected message to have been marked")
}

func TestLoader(t *testing.T) {
	suite.Run(t, new(MarkerTestSuite))
}
