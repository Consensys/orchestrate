package integrationtests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"

	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test/mocks"
	chainClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	transactionscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler"
	"gopkg.in/h2non/gock.v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	kafkaDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/kafka"
	postgresDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/zookeeper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
)

const postgresContainerID = "postgres-transaction-scheduler"
const kafkaContainerID = "kafka-transaction-scheduler"
const zookeeperContainerID = "zookeeper-transaction-scheduler"
const ChainRegistryURL = "http://chain-registry:8081"
const networkName = "transaction-scheduler"

var envPGHostPort string
var envKafkaHostPort string
var envHTTPPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	ctx                           context.Context
	logger                        log.Logger
	txScheduler                   *app.App
	client                        *docker.Client
	consumer                      *integrationtest.KafkaConsumer
	pgmngr                        postgres.Manager
	baseURL                       string
	contractRegistryResponseFaker *mocks.ContractRegistryFaker
	kafkaTopicConfig              *sarama.KafkaTopicConfig
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envPGHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envHTTPPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))

	// Define external hostname
	var kafkaExternalHostname = os.Getenv(sarama.KafkaURLEnv)
	if kafkaExternalHostname != "" {
		if strings.Contains(sarama.KafkaURLEnv, ":") {
			envKafkaHostPort = strings.Split(sarama.KafkaURLEnv, ":")[1]
		} else {
			envKafkaHostPort = "9092"
		}
	} else {
		envKafkaHostPort = strconv.Itoa(rand.IntnRange(20000, 29092))
		kafkaExternalHostname = fmt.Sprintf("localhost:%s", envKafkaHostPort)
	}

	// Initialize environment flags
	flgs := pflag.NewFlagSet("transaction-scheduler-integration-test", pflag.ContinueOnError)
	postgres.DBPort(flgs)
	sarama.KafkaURL(flgs)
	httputils.MetricFlags(flgs)
	httputils.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--rest-port=" + envHTTPPort,
		"--db-port=" + envPGHostPort,
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
			postgresContainerID:  {Postgres: postgresDocker.NewDefault().SetHostPort(envPGHostPort)},
			zookeeperContainerID: {Zookeeper: zookeeper.NewDefault()},
			kafkaContainerID: {Kafka: kafkaDocker.NewDefault().
				SetHostPort(envKafkaHostPort).
				SetZookeeperHostname(zookeeperContainerID).
				SetKafkaInternalHostname(kafkaContainerID).
				SetKafkaExternalHostname(kafkaExternalHostname),
			},
		},
	}

	// Docker client
	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	return &IntegrationEnvironment{
		ctx:     ctx,
		logger:  logger,
		client:  dockerClient,
		pgmngr:  postgres.NewManager(),
		baseURL: "http://localhost:" + envHTTPPort,
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	err := env.client.CreateNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not create network")
		return err
	}

	// Start postgres database
	err = env.client.Up(ctx, postgresContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up postgres")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, postgresContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start postgres")
		return err
	}

	// Start kafka + zookeeper
	err = env.client.Up(ctx, zookeeperContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up zookeeper")
		return err
	}

	err = env.client.Up(ctx, kafkaContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up kafka")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, kafkaContainerID, 20*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start kafka")
		return err
	}

	// Run postgres migrations
	err = env.migrate(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not migrate postgres")
		return err
	}

	env.kafkaTopicConfig = sarama.NewKafkaTopicConfig(viper.GetViper())
	env.contractRegistryResponseFaker = &mocks.ContractRegistryFaker{}

	env.txScheduler, err = newTransactionScheduler(ctx,
		mocks.NewContractRegistryClientMock(env.contractRegistryResponseFaker),
		env.kafkaTopicConfig)
	if err != nil {
		env.logger.WithError(err).Error("could initialize transaction scheduler")
		return err
	}

	// Start kafka consumer
	env.consumer, err = integrationtest.NewKafkaTestConsumer(ctx, "tx-scheduler-group", sarama.GlobalClient(),
		[]string{env.kafkaTopicConfig.Crafter, env.kafkaTopicConfig.Sender})
	if err != nil {
		env.logger.WithError(err).Error("could initialize kafka")
		return err
	}
	err = env.consumer.Start(context.Background())
	if err != nil {
		env.logger.WithError(err).Error("could not run kafka consumer")
		return err
	}

	// Start tx-scheduler app
	err = env.txScheduler.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start transaction-scheduler")
		return err
	}

	integrationtest.WaitForServiceReady(
		ctx,
		fmt.Sprintf("http://localhost:%s/ready", envMetricsPort),
		"transaction-scheduler",
		10*time.Second,
	)

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Infof("tearing test suite down")

	err := env.txScheduler.Stop(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not stop transaction-scheduler")
	}

	err = env.client.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down postgres")
	}

	err = env.client.Down(ctx, kafkaContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down kafka")
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

func (env *IntegrationEnvironment) migrate(ctx context.Context) error {
	// Set database connection
	opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
	if err != nil {
		return err
	}

	db := env.pgmngr.Connect(ctx, opts)
	env.logger.Debugf("initializing database migrations...")
	_, _, err = migrations.Run(db, "init")
	if err != nil {
		return err
	}

	env.logger.Debugf("running database migrations...")
	_, _, err = migrations.Run(db, "up")
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func newTransactionScheduler(
	ctx context.Context,
	contractRegistryClient contractregistry.ContractRegistryClient,
	topicCfg *sarama.KafkaTopicConfig,
) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	sarama.InitSyncProducer(ctx)

	// We mock the calls to the chain registry
	conf := chainClient.NewConfig(ChainRegistryURL)
	httpClient := httputils.NewClient()
	gock.InterceptClient(httpClient)
	chainRegistryClient := chainClient.NewHTTPClient(httpClient, conf)

	pgmngr := postgres.GetManager()
	txSchedulerConfig := transactionscheduler.NewConfig(viper.GetViper())

	return transactionscheduler.NewTxScheduler(
		txSchedulerConfig,
		pgmngr,
		authjwt.GlobalChecker(), authkey.GlobalChecker(),
		chainRegistryClient,
		contractRegistryClient,
		sarama.GlobalSyncProducer(),
		topicCfg,
	)
}
