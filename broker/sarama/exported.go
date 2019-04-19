package sarama

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config   *sarama.Config
	client   sarama.Client
	producer sarama.SyncProducer
	group    sarama.ConsumerGroup
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

// Config returns Sarama global configuration
func Config() *sarama.Config {
	return config
}

// SetConfig sets Sarama global configuration
func SetConfig(cfg *sarama.Config) {
	config = cfg
}

// InitClient initilialize Sarama Client
// It bases on viper configuration to get Kafka address
func InitClient(ctx context.Context) {
	// We need a config no create client
	if config == nil {
		InitConfig()
	}

	// Create sarama client
	var err error
	client, err = sarama.NewClient(viper.GetStringSlice(kafkaAddressViperKey), config)
	if err != nil {
		log.WithError(err).Fatalf("sarama: could not to start client")
		return
	}

	// Retrieve and log connected brokers
	var brokers = make(map[int32]string)
	for _, v := range client.Brokers() {
		brokers[v.ID()] = v.Addr()
	}
	log.Infof("sarama: client ready (connected to brokers: %v)", brokers)

	// Close when context is cancelled
	go func() {
		<-ctx.Done()
		client.Close()
	}()
}

// Client returns Sarama global client
func Client() sarama.Client {
	return client
}

// SetClient sets Sarama global client
func SetClient(c sarama.Client) {
	client = c
}

// InitSyncProducer initialize Sarama SyncProducer
func InitSyncProducer(ctx context.Context) {
	// We need a client no create SyncProcucer
	if client == nil {
		InitClient(ctx)
	}

	// Create sarama sync producer
	var err error
	producer, err = sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.WithError(err).Fatalf("sarama: could not create producer")
	}
	log.Infof("sarama: producer ready")

	// Close when context is cancelled
	go func() {
		<-ctx.Done()
		producer.Close()
	}()
}

// SyncProducer returns Sarama global SyncProducer
func SyncProducer() sarama.SyncProducer {
	return producer
}

// SetSyncProducer sets Sarama global SyncProducer
func SetSyncProducer(p sarama.SyncProducer) {
	producer = p
}

// InitConsumerGroup initialize consumer group
func InitConsumerGroup(ctx context.Context) {
	// We need a client no create ConsumerGroup
	if client == nil {
		InitClient(ctx)
	}

	// Create group
	var err error
	group, err = sarama.NewConsumerGroupFromClient(viper.GetString("worker.group"), client)
	if err != nil {
		log.WithError(err).Fatalf("sarama: could not create consumer group")
	}
	log.WithFields(log.Fields{
		"group": viper.GetString("worker.group"),
	}).Infof("sarama: consumer group ready")

	// Close when context is cancelled
	go func() {
		<-ctx.Done()
		group.Close()
	}()
}

// ConsumerGroup returns Sarama global ConsumerGroup
func ConsumerGroup() sarama.ConsumerGroup {
	return group
}

// SetConsumerGroup sets Sarama global ConsumerGroup
func SetConsumerGroup(g sarama.ConsumerGroup) {
	group = g
}
