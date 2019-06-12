package sarama

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestPipeSuite struct {
	suite.Suite
}

func (s *TestPipeSuite) TestPipe() {
	in := make(chan *sarama.ConsumerMessage)

	// Initialize pipe
	piped := Pipe(context.Background(), in)

	// Feed pipe
	rounds := 100
	go feed(in, rounds)

	for i := 0; i < rounds; i++ {
		msg := <-piped
		assert.Equal(s.T(), []byte{byte(i)}, msg.(*Msg).Key, "Message should have correct Key")
	}
}

func feed(in chan<- *sarama.ConsumerMessage, rounds int) {
	for i := 0; i < rounds; i++ {
		in <- &sarama.ConsumerMessage{Key: []byte{byte(i)}}
	}
}

func (s *TestPipeSuite) TestPipeInterupted() {
	in := make(chan *sarama.ConsumerMessage)
	// Initialize pipe with a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	piped := Pipe(ctx, in)

	// Feed pipe
	rounds := 100
	go feed(in, rounds)

	// Eraly pipe cancellation
	time.Sleep(10 * time.Millisecond)
	cancel()

	count := 0
	for range piped {
		count++
	}

	assert.True(s.T(), count > 0, "At least one message should have been processed")
	assert.True(s.T(), count < rounds, "All message should not have been processed")
}

func TestPipe(t *testing.T) {
	suite.Run(t, new(TestPipeSuite))
}
