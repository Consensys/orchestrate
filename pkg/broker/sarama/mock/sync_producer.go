package mock

import (
	"github.com/Shopify/sarama"
)

type MockSyncProducer struct {
	lastMessage *sarama.ProducerMessage
}

var _ sarama.SyncProducer = &MockSyncProducer{}

func NewMockSyncProducer() *MockSyncProducer {
	return &MockSyncProducer{}
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	m.lastMessage = msg
	return 0, 0, nil
}

func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	for _, msg := range msgs {
		_,_, _ = m.SendMessage(msg)
	}
	return nil
}

func (m *MockSyncProducer) Close() error {
	return nil
}

func (m *MockSyncProducer) LastMessage() *sarama.ProducerMessage {
	return m.lastMessage
}

func (m *MockSyncProducer) Clean() {
	m.lastMessage = nil
}
