package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
)

func main() {
	// Initialize producer
	broker.InitSyncProducer(context.Background())

	rounds := 1000
	log.Infof("Sending %v messages", rounds)
	for i := 0; i < rounds; i++ {
		msg := &sarama.ProducerMessage{
			Topic: "topic-e2e",
			Key:   sarama.StringEncoder(string(i)),
		}

		partition, offset, err := broker.GlobalSyncProducer().SendMessage(msg)
		if err != nil {
			log.WithError(err).Errorf("Could not send message")
			return
		}
		log.WithFields(log.Fields{
			"partition": partition,
			"offset":    offset,
		}).Info("Message sent")
	}
}
