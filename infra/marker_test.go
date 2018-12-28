package infra

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

type mockConsumerGroupSession struct {
	mux      *sync.Mutex
	lastMark int64
}

func (s *mockConsumerGroupSession) Claims() map[string][]int32 { return make(map[string][]int32) }
func (s *mockConsumerGroupSession) MemberID() string           { return "" }
func (s *mockConsumerGroupSession) GenerationID() int32        { return 0 }

func (s *mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	// Simulate some io time
	r := rand.Intn(100)
	time.Sleep(time.Duration(r) * time.Millisecond)
	s.mux.Lock()
	defer s.mux.Unlock()
	if offset > s.lastMark {
		s.lastMark = offset
	}
}

func (s *mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (s *mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	s.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, metadata)
}

func (s *mockConsumerGroupSession) Context() context.Context {
	return context.Background()
}

func newOffsetMarkerTestMessage(i int64) *sarama.ConsumerMessage {
	msg := sarama.ConsumerMessage{}
	msg.Offset = i
	return &msg
}

func TestSimpleSaramaOffsetMarkerConcurrent(t *testing.T) {
	mockS := &mockConsumerGroupSession{mux: &sync.Mutex{}}
	offset := NewSimpleSaramaOffsetMarker(mockS)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			offset.Mark(newOffsetMarkerTestMessage(i))
		}(int64(i))
	}
	wg.Wait()

	if mockS.lastMark != int64(rounds) {
		t.Errorf("SimpleSaramaOffsetMarker: expected last mark to be %v but got %v", rounds, mockS.lastMark)
	}
}
