package main

import (
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

var kafkaUrls = []string{"localhost:9092"}
var inTopic = "topic-tx-sender"

func main() {
	consumer, err := sarama.NewConsumer(kafkaUrls, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if e := consumer.Close(); err != nil {
			log.Fatalln(e)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(inTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Infof("Consumed message offset %d", msg.Offset)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Infof("Consumed: %d", consumed)
}
