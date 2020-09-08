package sarama

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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
	var err error
	config, err = NewSaramaConfig()
	if err != nil {
		log.Fatalf("sarama: cannot init TLS configuration for Kafka - got error: %q)", err)
		return
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
func InitClient(_ context.Context) (err error) {
	initClientOnce.Do(func() {
		if client != nil {
			return
		}

		// We need a config no create client
		if config == nil {
			InitConfig()
		}

		// Create sarama client
		hostnames := viper.GetStringSlice(KafkaURLViperKey)
		client, err = NewClient(hostnames, config)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not to start client at host %v", hostnames)
			return
		}

		// Retrieve and log connected brokers
		var brokers = make(map[int32]string)
		for _, v := range client.Brokers() {
			brokers[v.ID()] = v.Addr()
		}
		log.Infof("sarama: client ready (connected to brokers: %v) at host %v", brokers, hostnames)
	})

	return nil
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
		err := InitClient(ctx)
		if err != nil {
			return
		}

		if client == nil {
			log.WithError(err).Fatalf("sarama: client is not initialize")
			return
		}

		// Create sarama sync producer
		producer, err = NewSyncProducerFromClient(client)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not create producer")
			return
		}
		log.Infof("sarama: producer ready")
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
		err := InitClient(ctx)
		if err != nil {
			return
		}

		// Create group
		group, err = NewConsumerGroupFromClient(viper.GetString(KafkaGroupViperKey), client)
		if err != nil {
			log.WithError(err).Fatalf("sarama: could not create consumer group")
			return
		}
		log.WithFields(log.Fields{
			"group": viper.GetString(KafkaGroupViperKey),
		}).Infof("sarama: consumer group ready")
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

func Stop(_ context.Context) error {
	var closeProducerErr, closeConsumerGroupErr error
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		closeProducerErr = producer.Close()
		wg.Done()
	}()
	go func() {
		closeConsumerGroupErr = group.Close()
		wg.Done()
	}()

	wg.Wait()

	closeClientErr := client.Close()
	return errors.CombineErrors(closeProducerErr, closeConsumerGroupErr, closeClientErr)
}
