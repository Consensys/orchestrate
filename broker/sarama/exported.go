package sarama

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	component             = "broker.sarama"
	config                *sarama.Config
	client                sarama.Client
	initClientOnce        = &sync.Once{}
	producer              sarama.SyncProducer
	initProducerOnce      = &sync.Once{}
	group                 sarama.ConsumerGroup
	initConsumerGroupOnce = &sync.Once{}
)

// InitConfig initialize global Sarama configuration
func InitConfig() {
	// Init config
	config = sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
}

// GlobalConfig returns Sarama global configuration
func GlobalConfig() *sarama.Config {
	return config
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalConfig(cfg *sarama.Config) {
	config = cfg
}

// InitClient initilialize Sarama Client
// It bases on viper configuration to get Kafka address
func InitClient(ctx context.Context) {
	initClientOnce.Do(func() {
		if client != nil {
			return
		}

		// We need a config no create client
		if config == nil {
			InitConfig()
		}

		// Create sarama client
		var err error
		client, err = NewClient(viper.GetStringSlice(kafkaAddressViperKey), config)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not to start client")
		}

		// Retrieve and log connected brokers
		var brokers = make(map[int32]string)
		for _, v := range client.Brokers() {
			brokers[v.ID()] = v.Addr()
		}
		log.Infof("sarama: client ready (connected to brokers: %v)", brokers)

		// Close when context is canceled
		go func() {
			<-ctx.Done()
			client.Close()
		}()
	})
}

// GlobalClient returns Sarama global client
func GlobalClient() sarama.Client {
	return client
}

// SetGlobalClient sets Sarama global client
func SetGlobalClient(c sarama.Client) {
	client = c
}

// InitSyncProducer initialize Sarama SyncProducer
func InitSyncProducer(ctx context.Context) {
	initProducerOnce.Do(func() {
		if producer != nil {
			return
		}

		// Initialize client
		InitClient(ctx)

		// Create sarama sync producer
		var err error
		producer, err = NewSyncProducerFromClient(client)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not create producer")
		}
		log.Infof("sarama: producer ready")

		// Close when context is canceled
		go func() {
			<-ctx.Done()
			producer.Close()
		}()
	})
}

// GlobalSyncProducer returns Sarama global SyncProducer
func GlobalSyncProducer() sarama.SyncProducer {
	return producer
}

// SetGlobalSyncProducer sets Sarama global SyncProducer
func SetGlobalSyncProducer(p sarama.SyncProducer) {
	producer = p
}

// InitConsumerGroup initialize consumer group
func InitConsumerGroup(ctx context.Context) {
	initConsumerGroupOnce.Do(func() {
		if group != nil {
			return
		}

		// Initialize Client
		InitClient(ctx)

		// Create group
		var err error
		group, err = NewConsumerGroupFromClient(viper.GetString(kafkaGroupViperKey), client)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not create consumer group")
		}
		log.WithFields(log.Fields{
			"group": viper.GetString(kafkaGroupViperKey),
		}).Infof("sarama: consumer group ready")

		// Close when context is canceled
		go func() {
			<-ctx.Done()
			group.Close()
		}()
	})
}

// GlobalConsumerGroup returns Sarama global ConsumerGroup
func GlobalConsumerGroup() sarama.ConsumerGroup {
	return group
}

// SetGlobalConsumerGroup sets Sarama global ConsumerGroup
func SetGlobalConsumerGroup(g sarama.ConsumerGroup) {
	group = g
}

// Consume start consuming using global ConsumerGroup
func Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
consumeLoop:
	for {
		select {
		case <-ctx.Done():
			break consumeLoop
		default:
			err := group.Consume(ctx, topics, handler)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
