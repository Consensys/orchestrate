package integrationtests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	sarama2 "github.com/Shopify/sarama"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
	"gopkg.in/h2non/gock.v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/config"
	kafkaDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/kafka"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/zookeeper"
)

const kafkaContainerID = "Kafka-tx-signer"
const zookeeperContainerID = "zookeeper-tx-signer"
const txSchedulerURL = "http://transaction-scheduler:8081"
const keyManagerURL = "http://key-manager:8081"
const txSchedulerMetricsURL = "http://transaction-scheduler:8082"
const keyManagerMetricsURL = "http://key-manager:8082"
const networkName = "tx-signer"

var envKafkaHostPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	ctx            context.Context
	logger         log.Logger
	txSigner       *app.App
	client         *docker.Client
	consumer       *integrationtest.KafkaConsumer
	producer       sarama2.SyncProducer
	metricsURL     string
	txSignerConfig *txsigner.Config
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))
	envKafkaHostPort = strconv.Itoa(rand.IntnRange(20000, 29092))

	// Define external hostname
	kafkaExternalHostname := os.Getenv("KAFKA_HOST")
	if kafkaExternalHostname == "" {
		kafkaExternalHostname = "localhost"
	}
	kafkaExternalHostname = fmt.Sprintf("%s:%s", kafkaExternalHostname, envKafkaHostPort)

	// Initialize environment flags
	flgs := pflag.NewFlagSet("tx-signer-integration-test", pflag.ContinueOnError)
	postgres.DBPort(flgs)
	sarama.KafkaURL(flgs)
	httputils.MetricFlags(flgs)
	httputils.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--kafka-url=" + kafkaExternalHostname,
	}

	err := flgs.Parse(args)
	if err != nil {
		logger.WithError(err).Error("cannot parse environment flags")
		return nil, err
	}

	// Initialize environment container setup
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			zookeeperContainerID: {Zookeeper: zookeeper.NewDefault()},
			kafkaContainerID: {Kafka: kafkaDocker.NewDefault().
				SetHostPort(envKafkaHostPort).
				SetZookeeperHostname(zookeeperContainerID).
				SetKafkaInternalHostname(kafkaContainerID).
				SetKafkaExternalHostname(kafkaExternalHostname),
			},
		},
	}

	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	return &IntegrationEnvironment{
		ctx:        ctx,
		logger:     logger,
		client:     dockerClient,
		metricsURL: "http://localhost:" + envMetricsPort,
		producer:   sarama.GlobalSyncProducer(),
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	err := env.client.CreateNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not create network")
		return err
	}

	// Start Kafka + zookeeper
	err = env.client.Up(ctx, zookeeperContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up zookeeper")
		return err
	}

	err = env.client.Up(ctx, kafkaContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up Kafka")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, kafkaContainerID, 20*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start Kafka")
		return err
	}

	// Create app
	env.txSignerConfig = txsigner.NewConfig(viper.GetViper())
	env.txSigner, err = newTxSigner(ctx, env.txSignerConfig)
	if err != nil {
		env.logger.WithError(err).Error("could not initialize tx-signer")
		return err
	}

	// Start Kafka consumer
	env.consumer, err = integrationtest.NewKafkaTestConsumer(
		ctx,
		"tx-signer-integration-listener-group",
		sarama.GlobalClient(),
		[]string{env.txSignerConfig.SenderTopic, env.txSignerConfig.RecoverTopic},
	)
	if err != nil {
		env.logger.WithError(err).Error("could initialize Kafka")
		return err
	}
	err = env.consumer.Start(context.Background())
	if err != nil {
		env.logger.WithError(err).Error("could not run Kafka consumer")
		return err
	}

	// Set producer
	env.producer = sarama.GlobalSyncProducer()

	// Start tx-signer app
	err = env.txSigner.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start tx-signer")
		return err
	}

	integrationtest.WaitForServiceLive(ctx, fmt.Sprintf("%s/live", env.metricsURL), "tx-signer", 15*time.Second)
	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Infof("tearing test suite down")

	err := env.txSigner.Stop(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not stop tx-signer")
	}

	err = env.client.Down(ctx, kafkaContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down Kafka")
	}

	err = env.client.Down(ctx, zookeeperContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down zookeeper")
	}

	err = env.client.RemoveNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Errorf("could not remove network")
	}
}

func newTxSigner(ctx context.Context, txSignerConfig *txsigner.Config) (*app.App, error) {
	// Initialize dependencies
	sarama.InitSyncProducer(ctx)
	sarama.InitConsumerGroup(ctx, txSignerConfig.GroupName)

	httpClient := httputils.NewClient(httputils.NewDefaultConfig())
	gock.InterceptClient(httpClient)

	// We mock the calls to the key manager and tx-scheduler
	conf := keymanager.NewConfig(keyManagerURL, nil)
	conf.MetricsURL = keyManagerMetricsURL
	keyManagerClient := keymanager.NewHTTPClient(httpClient, conf)

	conf2 := txscheduler.NewConfig(txSchedulerURL, nil)
	conf2.MetricsURL = txSchedulerMetricsURL
	txSchedulerClient := txscheduler.NewHTTPClient(httpClient, conf2)

	return txsigner.NewTxSigner(txSignerConfig, sarama.GlobalConsumerGroup(), sarama.GlobalSyncProducer(), keyManagerClient, txSchedulerClient)
}
