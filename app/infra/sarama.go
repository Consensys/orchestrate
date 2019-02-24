package infra

import (
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
)

func initSarama(infra *Infra, wait *sync.WaitGroup) {
	initClient(infra, wait)
	initUnmarshaller(infra, wait)
	wait.Done()
}

func initClient(infra *Infra, wait *sync.WaitGroup) {
	// Init config
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient(viper.GetStringSlice("kafka.addresses"), config)
	if err != nil {
		log.WithError(err).Fatalf("infra-sarama: could not to start client")
		return
	}

	// Retrieve and log connectted brokers
	var brokers = make(map[int32]string)
	for _, v := range client.Brokers() {
		brokers[v.ID()] = v.Addr()
	}
	log.Infof("infra-sarama: client ready (connected to brokers: %v)", brokers)

	// Attach client
	infra.SaramaClient = client

	// Close when infra is cut
	go func() {
		<-infra.ctx.Done()
		client.Close()
	}()
}

func initUnmarshaller(infra *Infra, wait *sync.WaitGroup) {
	infra.Unmarshaller = infSarama.NewUnmarshaller()
}
