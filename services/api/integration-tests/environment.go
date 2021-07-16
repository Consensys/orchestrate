package integrationtests

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	authjwt "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/jwt"
	authkey "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/key"
	httputils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	ethclient "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/rpc"
	integrationtest "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test"
	ganacheDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/ganache"
	hashicorpDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/hashicorp"
	quorumkeymanagerDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/services/api"
	"gopkg.in/h2non/gock.v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/config"
	kafkaDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/kafka"
	postgresDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/postgres"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/zookeeper"
	"github.com/ConsenSys/orchestrate/services/api/store/postgres/migrations"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const postgresContainerID = "postgres-api"
const kafkaContainerID = "Kafka-api"
const zookeeperContainerID = "zookeeper-api"
const ganacheContainerID = "ganache-api"
const qkmContainerID = "quorum-key-manager"
const hashicorpContainerID = "hashicorp"
const networkName = "api"
const localhost = "localhost"
const qkmStoreName = "orchestrate-eth1"
const hashicorpMountPath = "orchestrate"

var envPGHostPort string
var envKafkaHostPort string
var envHTTPPort string
var envMetricsPort string
var envGanacheHostPort string
var envQKMHostPort string
var envVaultHostPort string

type IntegrationEnvironment struct {
	ctx               context.Context
	logger            log.Logger
	api               *app.App
	client            *docker.Client
	consumer          *integrationtest.KafkaConsumer
	pgmngr            postgres.Manager
	baseURL           string
	metricsURL        string
	kafkaTopicConfig  *sarama.KafkaTopicConfig
	blockchainNodeURL string
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envPGHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envHTTPPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))
	envKafkaHostPort = strconv.Itoa(rand.IntnRange(20000, 29092))
	envGanacheHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envQKMHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envVaultHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))

	// Define external hostname
	kafkaExternalHostname := os.Getenv("KAFKA_HOST")
	if kafkaExternalHostname == "" {
		kafkaExternalHostname = localhost
	}
	kafkaExternalHostname = fmt.Sprintf("%s:%s", kafkaExternalHostname, envKafkaHostPort)

	ganacheExternalHostname := os.Getenv("GANACHE_HOST")
	if ganacheExternalHostname == "" {
		ganacheExternalHostname = localhost
	}

	quorumKeyManagerURL := fmt.Sprintf("http://localhost:%s", envQKMHostPort)

	// Initialize environment flags
	flgs := pflag.NewFlagSet("api-integration-test", pflag.ContinueOnError)
	api.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--rest-port=" + envHTTPPort,
		"--db-port=" + envPGHostPort,
		"--kafka-url=" + kafkaExternalHostname,
		"--key-manager-url=" + quorumKeyManagerURL,
		"--key-manager-store-name=" + qkmStoreName,
		"--log-level=error",
	}

	err := flgs.Parse(args)
	if err != nil {
		logger.WithError(err).Error("cannot parse environment flags")
		return nil, err
	}

	rootToken := fmt.Sprintf("root_token_%v", strconv.Itoa(rand.IntnRange(0, 10000)))
	pluginsPath, err := getPluginsPath()
	if err != nil {
		return nil, err
	}

	vaultContainer := hashicorpDocker.NewDefault().
		SetHostPort(envVaultHostPort).
		SetRootToken(rootToken).
		SetPluginSourceDirectory(pluginsPath).
		SetMountPath(hashicorpMountPath)

	err = vaultContainer.DownloadPlugin()
	if err != nil {
		return nil, err
	}

	manifestsPath, err := getManifestsPath()
	if err != nil {
		return nil, err
	}

	qkmContainer := quorumkeymanagerDocker.NewDefault().
		SetHostPort(envQKMHostPort).
		SetManifestDirectory(manifestsPath)

	err = qkmContainer.CreateManifest("manifest.yml", &qkm.Manifest{
		Kind:    "Eth1Account",
		Version: "0.0.1",
		Name:    qkmStoreName,
		Specs: map[string]interface{}{
			"keystore": "HashicorpKeys",
			"specs": map[string]string{
				"mountPoint": hashicorpMountPath,
				"address":    "http://" + hashicorpContainerID + ":8200",
				"token":      rootToken,
				"namespace":  "",
			},
		},
	})

	if err != nil {
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
			hashicorpContainerID: {HashicorpVault: vaultContainer},
			qkmContainerID:       {QuorumKeyManager: qkmContainer},
			ganacheContainerID:   {Ganache: ganacheDocker.NewDefault().SetHostPort(envGanacheHostPort).SetHost(ganacheExternalHostname)},
		},
	}

	// Docker client
	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	return &IntegrationEnvironment{
		ctx:               ctx,
		logger:            logger,
		client:            dockerClient,
		pgmngr:            postgres.NewManager(),
		baseURL:           "http://localhost:" + envHTTPPort,
		metricsURL:        "http://localhost:" + envMetricsPort,
		blockchainNodeURL: fmt.Sprintf("http://%s:%s", ganacheExternalHostname, envGanacheHostPort),
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	err := env.client.CreateNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not create network")
		return err
	}

	// Start Hashicorp Vault
	err = env.client.Up(ctx, hashicorpContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up vault container")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, hashicorpContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start vault")
		return err
	}

	// Start quorum key manager
	err = env.client.Up(ctx, qkmContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up quorum-key-manager")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, qkmContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start quorum-key-manager")
		return err
	}

	// Start postgres Database
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

	// Start ganache
	err = env.client.Up(ctx, ganacheContainerID, "")
	if err != nil {
		env.logger.WithError(err).Error("could not up ganache")
		return err
	}
	err = env.client.WaitTillIsReady(ctx, ganacheContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start ganache")
		return err
	}

	// Run postgres migrations
	err = env.migrate(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not migrate postgres")
		return err
	}

	env.kafkaTopicConfig = sarama.NewKafkaTopicConfig(viper.GetViper())

	env.api, err = newAPI(ctx, env.kafkaTopicConfig)
	if err != nil {
		env.logger.WithError(err).Error("could initialize API")
		return err
	}

	// Start Kafka consumer
	env.consumer, err = integrationtest.NewKafkaTestConsumer(ctx, "api-group", sarama.GlobalClient(),
		[]string{env.kafkaTopicConfig.Sender})
	if err != nil {
		env.logger.WithError(err).Error("could initialize Kafka")
		return err
	}
	err = env.consumer.Start(context.Background())
	if err != nil {
		env.logger.WithError(err).Error("could not run Kafka consumer")
		return err
	}

	// Start API
	err = env.api.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start API")
		return err
	}

	integrationtest.WaitForServiceLive(
		ctx,
		fmt.Sprintf("%s/live", env.metricsURL),
		"api",
		15*time.Second,
	)

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Info("tearing test suite down")

	if env.api != nil {
		err := env.api.Stop(ctx)
		if err != nil {
			env.logger.WithError(err).Error("could not stop API")
		}
	}

	err := env.client.Down(ctx, qkmContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down zookeeper")
	}

	err = env.client.Down(ctx, hashicorpContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down zookeeper")
	}

	err = env.client.Down(ctx, ganacheContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down ganache")
	}

	err = env.client.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down postgres")
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

func (env *IntegrationEnvironment) migrate(ctx context.Context) error {
	// Set Database connection
	opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
	if err != nil {
		return err
	}

	db := env.pgmngr.Connect(ctx, opts)
	env.logger.Debugf("initializing Database migrations...")
	_, _, err = migrations.Run(db, "init")
	if err != nil {
		return err
	}

	env.logger.Debugf("running Database migrations...")
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

func newAPI(ctx context.Context, topicCfg *sarama.KafkaTopicConfig) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	sarama.InitSyncProducer(ctx)
	ethclient.Init(ctx)
	qkm.Init()

	interceptedHTTPClient := httputils.NewClient(httputils.NewDefaultConfig())
	gock.InterceptClient(interceptedHTTPClient)

	pgmngr := postgres.GetManager()
	txSchedulerConfig := api.NewConfig(viper.GetViper())

	return api.NewAPI(
		txSchedulerConfig,
		pgmngr,
		authjwt.GlobalChecker(), authkey.GlobalChecker(),
		qkm.GlobalClient(),
		ethclient.GlobalClient(),
		sarama.GlobalSyncProducer(),
		topicCfg,
	)
}

func getManifestsPath() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(currDir, "manifests"), nil
}

func getPluginsPath() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(currDir, "plugins"), nil
}
