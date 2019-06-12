package sarama

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

func TestContext(t *testing.T) {
	ctx := WithConsumerGroupSessionAndClaim(
		context.Background(),
		mock.NewConsumerGroupSession(context.TODO(), "test-group", make(map[string][]int32)),
		mock.NewConsumerGroupClaim("topic-test", 1, 0),
	)

	s, c := GetConsumerGroupSessionAndClaim(ctx)
	assert.True(t, s.GenerationID() >= 1, "Generation should be greater than 1")
	assert.Equal(t, "topic-test", c.Topic(), "Topic should match")
}

type CounterHandler struct {
	counter int32
}

func (h *CounterHandler) Handle(txctx *engine.TxContext) {
	txctx.Logger.WithFields(log.Fields{
		"topic":     txctx.Msg.(*Msg).Topic,
		"offset":    txctx.Msg.(*Msg).Offset,
		"partition": txctx.Msg.(*Msg).Partition,
	}).Infof("Handling message")
	atomic.AddInt32(&h.counter, 1)
}

func TestConsumerGroupHandler(t *testing.T) {
	conf := engine.NewConfig()
	e := engine.NewEngine(&conf)

	counter := CounterHandler{}
	e.Register(counter.Handle)

	cgHandler := NewEngineConsumerGroupHandler(e)

	msgs := make(map[string]map[int32][]*sarama.ConsumerMessage)
	for _, topic := range []string{"test-topic-1", "test-topic-2"} {
		msgs[topic] = make(map[int32][]*sarama.ConsumerMessage)
		for _, partition := range []int32{0, 1, 2} {
			msgs[topic][partition] = []*sarama.ConsumerMessage{}
			for i := range make([]int, 10) {
				msgs[topic][partition] = append(msgs[topic][partition], &sarama.ConsumerMessage{
					Topic:     topic,
					Partition: partition,
					Offset:    int64(i),
				})
			}
		}
	}

	// Create consumer group
	cg := mock.NewConsumerGroup("test-group", msgs)
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	go func() {
		err = cg.Consume(ctx, []string{"test-topic-1", "test-topic-2"}, cgHandler)
		assert.NoError(t, err, "No error expected")
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Equal(t, int32(60), counter.counter, "Count of processed message should be correct")
}
