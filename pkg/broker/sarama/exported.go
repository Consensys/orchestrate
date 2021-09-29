package sarama

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/hashicorp/go-multierror"
	healthz "github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
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
	checker               healthz.Check
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

// GlobalConfig returns Sarama global configuration
func GlobalConfig() *sarama.Config {
	return config
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalConfig(cfg *sarama.Config) {
	config = cfg
}

// initialize Sarama Client
// It bases on viper configuration to get Kafka address
func InitClient(ctx context.Context) (err error) {
	initClientOnce.Do(func() {
		if client != nil {
			return
		}

		logger := log.FromContext(ctx)
		// We need a config no create client
		if config == nil {
			config, err = NewSaramaConfig()
			if err != nil {
				logger.WithError(err).Fatal("failed to initialize kafka client")
				return
			}
		}

		// Create sarama client
		hostnames := viper.GetStringSlice(KafkaURLViperKey)
		client, err = NewClient(hostnames, config)
		if err != nil {
			logger.WithField("hosts", hostnames).Fatal("could not to start client")
			return
		}

		// Retrieve and log connected brokers
		var brokers = make(map[int32]string)
		for _, v := range client.Brokers() {
			brokers[v.ID()] = v.Addr()
		}

		checker = func() error {
			gr := &multierror.Group{}
			for _, host := range hostnames {
				gr.Go(healthz.TCPDialCheck(host, time.Second*3))
			}

			return gr.Wait().ErrorOrNil()
		}

		sarama.Logger = newLogger(logger).SetLevel(logrus.DebugLevel)
		logger.WithField("host", hostnames).WithField("broker", brokers).Info("client ready")
	})

	return nil
}

// GlobalClient returns Sarama global client
func GlobalClient() sarama.Client {
	return client
}

func GlobalClientChecker() healthz.Check {
	return checker
}

// SetGlobalClient sets Sarama global client
func SetGlobalClient(c sarama.Client) {
	client = c
}

// InitSyncProducer initialize Sarama SyncProducer
func InitSyncProducer(ctx context.Context) {
	initProducerOnce.Do(func() {
		logger := log.NewLogger().WithContext(ctx).SetComponent(component + ".producer")
		ctx = log.With(ctx, logger)
		if producer != nil {
			return
		}

		// Initialize client
		err := InitClient(ctx)
		if err != nil {
			return
		}

		if client == nil {
			logger.WithError(err).Fatal("client is not initialize")
			return
		}

		// Create sarama sync producer
		producer, err = NewSyncProducerFromClient(client)
		if err != nil {
			logger.WithError(err).Fatal("could not create producer")
			return
		}
		logger.Info("producer ready")
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
func InitConsumerGroup(ctx context.Context, kafkaGroup string) {
	initConsumerGroupOnce.Do(func() {
		logger := log.NewLogger().SetComponent(component + ".consumer").WithContext(ctx)
		ctx = log.With(ctx, logger)
		if group != nil {
			return
		}

		err := InitClient(ctx)
		if err != nil {
			return
		}

		group, err = NewConsumerGroupFromClient(kafkaGroup, client)
		if err != nil {
			logger.WithError(err).Fatal("could not create consumer group")
			return
		}

		logger.WithField("group", kafkaGroup).Info("consumer group ready")
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
