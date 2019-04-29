package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
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
		broker.GlobalSyncProducer().SendMessage(msg)
	}
}
