package sarama

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const component = "broker.sarama"

var (
	config                *sarama.Config
	client                sarama.Client
	initClientOnce        = &sync.Once{}
	producer              sarama.SyncProducer
	initProducerOnce      = &sync.Once{}
	group                 sarama.ConsumerGroup
	initConsumerGroupOnce = &sync.Once{}
)

// NewTLSConfig inspired by https://medium.com/processone/using-tls-authentication-for-your-go-kafka-client-3c5841f2a625
func NewTLSConfig(clientCertFilePath, clientKeyFilePath, caCertFilePath string) (*tls.Config, error) {
	tlsConfig := tls.Config{}
	var err error

	if clientCertFilePath != "" && clientKeyFilePath != "" {
		// Load client cert
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(clientCertFilePath, clientKeyFilePath)
		if err != nil {
			return &tls.Config{}, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if clientCertFilePath != "" {
		// Load CA cert
		var caCert []byte
		caCert, err = ioutil.ReadFile(caCertFilePath)
		if err != nil {
			return &tls.Config{}, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.InsecureSkipVerify = viper.GetBool(kafkaTLSInsecureSkipVerifyViperKey)

	return &tlsConfig, err
}

// InitConfig initialize global Sarama configuration
func InitConfig() {
	// Init config
	config = sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.MaxWaitTime = time.Duration(viper.GetInt64(kafkaConsumerMaxWaitTimeViperKey)) * time.Millisecond

	config.Net.SASL.Enable = viper.GetBool(kafkaSASLEnabledViperKey)
	config.Net.SASL.Mechanism = sarama.SASLMechanism(viper.GetString(kafkaSASLMechanismViperKey))
	config.Net.SASL.Handshake = viper.GetBool(kafkaSASLHandshakeViperKey)
	config.Net.SASL.User = viper.GetString(kafkaSASLUserViperKey)
	config.Net.SASL.Password = viper.GetString(kafkaSASLPasswordViperKey)
	config.Net.SASL.SCRAMAuthzID = viper.GetString(kafkaSASLSCRAMAuthzIDViperKey)

	config.Net.TLS.Enable = viper.GetBool(kafkaTLSEnableViperKey)
	if config.Net.TLS.Enable {
		tlsConfig, err := NewTLSConfig(
			viper.GetString(kafkaTLSClientCertFilePathViperKey),
			viper.GetString(kafkaTLSClientKeyFilePathViperKey),
			viper.GetString(kafkaTLSCACertFilePathViperKey),
		)
		// Fatal if get error from NewTLSConfig
		if err != nil {
			log.Fatalf("sarama: cannot init TLS configuration for Kafka - got error: %q)", err)
		}
		config.Net.TLS.Config = tlsConfig
	}
}

// GlobalConfig returns Sarama global configuration
func GlobalConfig() *sarama.Config {
	return config
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalConfig(cfg *sarama.Config) {
	config = cfg
}

// InitClient initialize Sarama Client
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
		client, err = NewClient(viper.GetStringSlice(KafkaURLViperKey), config)
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
			closeErr := client.Close()
			if closeErr != nil {
				log.WithError(closeErr).Warn("could not close client")
			}
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
			closeErr := producer.Close()
			if closeErr != nil {
				log.WithError(closeErr).Warn("could not close client")
			}
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
		group, err = NewConsumerGroupFromClient(viper.GetString(KafkaGroupViperKey), client)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not create consumer group")
		}
		log.WithFields(log.Fields{
			"group": viper.GetString(KafkaGroupViperKey),
		}).Infof("sarama: consumer group ready")

		// Close when context is canceled
		go func() {
			<-ctx.Done()
			closeErr := group.Close()
			if closeErr != nil {
				log.WithError(closeErr).Warn("could not close client")
			}
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
