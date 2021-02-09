// +build unit

package offset

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"github.com/stretchr/testify/suite"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine/testutils"
)

func makeMarkerContext(session sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim, i int) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare(log.NewLogger(), nil)
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

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	mockConsumerGroupSession := mock.NewMockConsumerGroupSession(ctrl)
	mockConsumerGroupSession.EXPECT().Commit().Times(rounds)
	mockConsumerGroupSession.EXPECT().MarkMessage(gomock.Any(), "").Times(rounds)
	mockConsumerGroupClaim := mock.NewMockConsumerGroupClaim(ctrl)

	var txctxs []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeMarkerContext(mockConsumerGroupSession, mockConsumerGroupClaim, i))
	}

	// Handle contexts
	s.Handle(txctxs)
}

func TestLoader(t *testing.T) {
	suite.Run(t, new(MarkerTestSuite))
}
