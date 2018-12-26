package handlers

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type testConsumerGroupSession struct {
	mux      *sync.Mutex
	lastMark int64
}

func (s *testConsumerGroupSession) Claims() map[string][]int32 { return make(map[string][]int32) }
func (s *testConsumerGroupSession) MemberID() string           { return "" }
func (s *testConsumerGroupSession) GenerationID() int32        { return 0 }

func (s *testConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	// Simulate some io time
	r := rand.Intn(100)
	time.Sleep(time.Duration(r) * time.Millisecond)
	s.mux.Lock()
	defer s.mux.Unlock()
	if offset > s.lastMark {
		s.lastMark = offset
	}
}

func (s *testConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (s *testConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	s.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, metadata)
}

func (s *testConsumerGroupSession) Context() context.Context {
	return context.Background()
}

func newMessage(i int64) *sarama.ConsumerMessage {
	msg := &sarama.ConsumerMessage{}
	msg.Offset = i
	return msg
}

func TestMarker(t *testing.T) {
	s := testConsumerGroupSession{&sync.Mutex{}, -1}
	offset := NewSimpleSaramaOffsetMarker(&s)
	h := Marker(offset)

	w := infra.NewWorker([]infra.HandlerFunc{h}, 100)

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newMessage(int64(i))
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if s.lastMark != int64(rounds)+1 {
		t.Errorf("Marker: expected lastMark to be %v but got %v", rounds+1, s.lastMark)
	}
}
