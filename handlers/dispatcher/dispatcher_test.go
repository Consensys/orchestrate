package dispatcher

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

const testsNum = 3

type MockChanRegistry struct {
	MockChan chan *envelope.Envelope
}

func (r *MockChanRegistry) NewEnvelopeChan(scenarioID, topic string) chan *envelope.Envelope {
	return r.MockChan
}

func (r *MockChanRegistry) GetEnvelopeChan(scenarioID, topic string) chan *envelope.Envelope {
	return r.MockChan
}

func (r *MockChanRegistry) CloseEnvelopeChan(scenarioID, topic string) error {
	return fmt.Errorf("error")
}

func makeCrafterContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())

	switch i {
	case 0:
		// Input an empty envelope which return error
		ctx.Envelope.Metadata = &envelope.Metadata{}
		ctx.Set("errors", 1)
		ctx.Set("expectedErrorMessage", "message:\"invalid input message format\" ")
		ctx.Set("result", "")
	case 1:
		// Input a standard envelope with an extra data with ScenarioID
		ctx.Msg = &sarama.ConsumerMessage{
			Topic: "testTopic",
		}
		extra := make(map[string]string)
		extra["ScenarioID"] = "test"
		ctx.Envelope.Metadata = &envelope.Metadata{Id: "test", Extra: extra}
		ctx.Set("errors", 0)
		ctx.Set("result", "")
	case 2:
		// Input an envelope without ScenarioID in extra data
		extra := make(map[string]string)
		ctx.Msg = &sarama.ConsumerMessage{
			Topic: "testTopic",
		}
		ctx.Envelope.Metadata = &envelope.Metadata{Id: "test", Extra: extra}
		ctx.Set("errors", 1)
		ctx.Set("expectedErrorMessage", "message:\"no ScenarioID found, envelope not dispatched\" ")
		ctx.Set("result", "")
	}

	return ctx
}

type DispacherTestSuite struct {
	testutils.HandlerTestSuite
	MockChan chan *envelope.Envelope
}

func (s *DispacherTestSuite) SetupSuite() {
	mock := make(chan *envelope.Envelope)
	s.MockChan = mock
	s.Handler = Dispacher(&MockChanRegistry{
		MockChan: mock,
	})
}

func (s *DispacherTestSuite) TestDispatcher() {
	txctxs := []*engine.TxContext{}
	for i := 0; i < testsNum; i++ {
		txctxs = append(txctxs, makeCrafterContext(i))
	}

	go func() {
		for i := 0; i < testsNum; i++ {
			<-s.MockChan
		}
	}()

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors", txctx.Envelope.Call)
		if len(txctx.Envelope.Errors) == 1 {
			assert.Equal(s.T(), txctx.Get("expectedErrorMessage").(string), txctx.Envelope.Errors[0].String(), "Expected the right error message")

		}
	}
}

func TestDispatcher(t *testing.T) {
	suite.Run(t, new(DispacherTestSuite))
}
