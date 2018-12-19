package infra

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
)

type TestProcessor struct {
	messages chan *sarama.ConsumerMessage
	errors   chan error
}

func (p *TestProcessor) ProcessMessage(message *sarama.ConsumerMessage) {
	p.messages <- message
}

func (p *TestProcessor) ProcessError(err error) {
	p.errors <- err
}

var (
	fooMsg = &sarama.ConsumerMessage{Value: []byte("hello Foo")}
	barMsg = &sarama.ConsumerMessage{Value: []byte("hello Bar")}
)

func TestListener(t *testing.T) {
	consumer := mocks.NewConsumer(t, nil)
	defer consumer.Close()

	// Prefill Partition test
	consumer.ExpectConsumePartition("test", 0, sarama.OffsetOldest).YieldMessage(fooMsg)
	consumer.ExpectConsumePartition("test", 0, sarama.OffsetOldest).YieldError(sarama.ErrOutOfBrokers)
	consumer.ExpectConsumePartition("test", 0, sarama.OffsetOldest).YieldMessage(barMsg)

	// Start consuming partition
	pc, err := consumer.ConsumePartition("test", 0, sarama.OffsetOldest)
	if err != nil {
		t.Errorf("Could not consume partition: %v", err)
	}

	// Create a new listener
	processor := TestProcessor{
		make(chan *sarama.ConsumerMessage),
		make(chan error),
	}
	l := NewListener(nil, pc, &processor)

	// Start listening
	go l.Listen()

	messages, errors := make([]*sarama.ConsumerMessage, 0), make([]error, 0)
	for i := 0; i < 3; i++ {
		select {
		case msg := <-processor.messages:
			messages = append(messages, msg)
		case err := <-processor.errors:
			errors = append(errors, err)
		}
	}
	l.Stop()

	if len(messages) != 2 {
		t.Errorf("Listener: expected 2 messages but got %v", len(messages))
	}

	if len(errors) != 1 {
		t.Errorf("Listener: expected 1 error but got %v", len(errors))
	}
}

func TestWorker(t *testing.T) {
	consumer := mocks.NewConsumer(t, nil)
	defer consumer.Close()

	// Prefill Partition test
	consumer.ExpectConsumePartition("testA", 0, sarama.OffsetOldest).YieldMessage(fooMsg)
	consumer.ExpectConsumePartition("testA", 1, sarama.OffsetOldest).YieldMessage(barMsg)
	consumer.ExpectConsumePartition("testB", 0, sarama.OffsetOldest).YieldMessage(barMsg)

	// Create worker
	w := NewSaramaWorker(consumer)

	// Create a new listener
	processor := TestProcessor{
		make(chan *sarama.ConsumerMessage),
		make(chan error),
	}

	err := w.Subscribe("testA", 0, -2, &processor)
	if err != nil {
		t.Errorf("Subscription: %v", err)
	}

	err = w.Subscribe("testA", 1, -2, &processor)
	if err != nil {
		t.Errorf("Subscription: %v", err)
	}

	err = w.Subscribe("testB", 0, -2, &processor)
	if err != nil {
		t.Errorf("Subscription: %v", err)
	}

	if len(w.new) != 3 {
		t.Errorf("Worker: expected 3 listeners but got %v", len(w.new))
	}

	for i := 0; i < 3; i++ {
		w.addListener(<-w.new)
	}

	if len(w.listeners) != 3 {
		t.Errorf("Worker: expected 3 listeners but got %v", len(w.listeners))
	}

	w.Stop()

	if len(w.listeners) != 3 {
		t.Errorf("Worker: expected 3 listeners but got %v", len(w.listeners))
	}

	for l := range w.listeners {
		if len(l.stop) != 1 {
			t.Errorf("Worker: expected listener to have 1 stop bug got %v", len(l.stop))
		}
		l.Close()
	}

	if len(w.dead) != 3 {
		t.Errorf("Worker: expected 3 dead listeners but got %v", len(w.dead))
	}

	for i := 0; i < 3; i++ {
		w.removeListener(<-w.dead)
	}

	if len(w.listeners) != 0 {
		t.Errorf("Worker: expected 0 listeners but got %v", len(w.listeners))
	}
}
