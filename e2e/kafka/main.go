package main

import (
	//	"encoding/hex"
	"fmt"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type LogProcessor struct {
	Name string
}

func (p *LogProcessor) ProcessMessage(message *sarama.ConsumerMessage) {
	fmt.Printf("%v: Process message %q (offset=%v)\n", p.Name, string(message.Value), message.Offset)
}

func (p *LogProcessor) ProcessError(err error) {
	fmt.Printf("%v: Process error %v\n", p.Name, err)
}

func produceMessages(p sarama.AsyncProducer) {
	for i := 0; i < 20; i++ {
		var message sarama.ProducerMessage
		if i%2 == 0 {

			message = sarama.ProducerMessage{
				Topic:     "test-A",
				Partition: -1,
				Value:     sarama.StringEncoder(fmt.Sprintf("testing %v", i)),
			}
		} else {
			message = sarama.ProducerMessage{
				Topic:     "test-B",
				Partition: -1,
				Value:     sarama.StringEncoder(fmt.Sprintf("testing %v", i)),
			}
		}

		p.Input() <- &message
	}
}

func main() {
	// Create client
	client, err := sarama.NewClient([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Client ready")

	// Create consumer
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Consumer ready")
	defer consumer.Close()

	// Create producer
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Producer ready")
	defer producer.Close()

	// Create worker
	worker := infra.NewSaramaWorker(consumer)
	fmt.Printf("Worker created %v\n", worker)

	// Subscribe
	processorA := LogProcessor{"Processor A"}
	worker.Subscribe("test-A", 0, sarama.OffsetOldest, &processorA)

	processorB := LogProcessor{"Processor B"}
	worker.Subscribe("test-B", 0, sarama.OffsetOldest, &processorB)

	produceMessages(producer)

	// Start worker
	worker.Run()
}
