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
	qkmclient "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client"
	"github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	httputils "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	redis2 "github.com/ConsenSys/orchestrate/pkg/toolkit/database/redis"
	ethclient "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/rpc"
	integrationtest "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test"
	hashicorpDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/hashicorp"
	quorumkeymanagerDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/quorum-key-manager"
	txsender "github.com/ConsenSys/orchestrate/services/tx-sender"
	"github.com/ConsenSys/orchestrate/services/tx-sender/store"
	"github.com/ConsenSys/orchestrate/services/tx-sender/store/redis"
	sarama2 "github.com/Shopify/sarama"
	"github.com/alicebob/miniredis"
	"github.com/cenkalti/backoff/v4"
	"gopkg.in/h2non/gock.v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/config"
	kafkaDocker "github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/kafka"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/integration-test/docker/container/zookeeper"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const kafkaContainerID = "Kafka-tx-sender"
const zookeeperContainerID = "zookeeper-tx-sender"
const apiURL = "http://api:8081"
const apiMetricsURL = "http://api:8082"
const networkName = "tx-sender"
const qkmStoreName = "orchestrate-eth1"
const hashicorpContainerID = "hashicorp"
const qkmContainerID = "quorum-key-manager"
const hashicorpMountPath = "orchestrate"
const maxRecoveryDefault = 1

var envQKMHostPort string
var envVaultHostPort string
var envKafkaHostPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	ctx        context.Context
	logger     log.Logger
	txSender   *app.App
	client     *docker.Client
	consumer   *integrationtest.KafkaConsumer
	producer   sarama2.SyncProducer
	metricsURL string
	ns         store.NonceSender
	redis      *redis2.Client
	srvConfig  *txsender.Config
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))
	envKafkaHostPort = strconv.Itoa(rand.IntnRange(20000, 29092))
	envQKMHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envVaultHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))

	// Define external hostname
	kafkaExternalHostname := os.Getenv("KAFKA_HOST")
	if kafkaExternalHostname == "" {
		kafkaExternalHostname = "localhost"
	}

	kafkaExternalHostname = fmt.Sprintf("%s:%s", kafkaExternalHostname, envKafkaHostPort)

	quorumKeyManagerURL := fmt.Sprintf("http://localhost:%s", envQKMHostPort)

	// Initialize environment flags
	flgs := pflag.NewFlagSet("tx-sender-integration-test", pflag.ContinueOnError)
	txsender.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--kafka-url=" + kafkaExternalHostname,
		"--nonce-manager-type=" + txsender.NonceManagerTypeRedis,
		"--key-manager-url=" + quorumKeyManagerURL,
		"--key-manager-store-name=" + qkmStoreName,
		"--api-url=" + apiURL,
		"--log-level=panic",
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
			zookeeperContainerID: {Zookeeper: zookeeper.NewDefault()},
			kafkaContainerID: {Kafka: kafkaDocker.NewDefault().
				SetHostPort(envKafkaHostPort).
				SetZookeeperHostname(zookeeperContainerID).
				SetKafkaInternalHostname(kafkaContainerID).
				SetKafkaExternalHostname(kafkaExternalHostname),
			},
			hashicorpContainerID: {HashicorpVault: vaultContainer},
			qkmContainerID:       {QuorumKeyManager: qkmContainer},
		},
	}

	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	mredis, _ := miniredis.Run()
	conf := &redis2.Config{
		Expiration: 100000,
		Host:       mredis.Host(),
		Port:       mredis.Port(),
	}

	pool, _ := redis2.NewPool(conf)
	redisCli := redis2.NewClient(pool, conf)
	return &IntegrationEnvironment{
		ctx:        ctx,
		logger:     logger,
		client:     dockerClient,
		metricsURL: "http://localhost:" + envMetricsPort,
		producer:   sarama.GlobalSyncProducer(),
		redis:      redisCli,
		ns:         redis.NewNonceSender(redisCli),
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
	env.srvConfig = txsender.NewConfig(viper.GetViper())
	env.srvConfig.BckOff = testBackOff()
	env.txSender, err = newTxSender(ctx, env.srvConfig, env.redis)
	if err != nil {
		env.logger.WithError(err).Error("could not initialize tx-sender")
		return err
	}

	// Start Kafka consumer
	env.consumer, err = integrationtest.NewKafkaTestConsumer(
		ctx,
		"tx-sender-integration-listener-group",
		sarama.GlobalClient(),
		[]string{env.srvConfig.SenderTopic, env.srvConfig.RecoverTopic},
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

	// Start tx-sender app
	err = env.txSender.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start tx-sender")
		return err
	}

	integrationtest.WaitForServiceLive(ctx, fmt.Sprintf("%s/live", env.metricsURL), "tx-sender", 15*time.Second)
	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Info("tearing test suite down")

	err := env.txSender.Stop(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not stop tx-sender")
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

func newTxSender(ctx context.Context, txSenderConfig *txsender.Config, redisCli *redis2.Client) (*app.App, error) {
	// Initialize dependencies
	sarama.InitSyncProducer(ctx)
	sarama.InitConsumerGroup(ctx, txSenderConfig.GroupName)

	httpClient := httputils.NewClient(httputils.NewDefaultConfig())
	gock.InterceptClient(httpClient)

	ec := ethclient.NewClient(testBackOff, httpClient)
	qkmClient := qkmclient.NewHTTPClient(httputils.NewClient(httputils.NewDefaultConfig()), &qkmclient.Config{
		URL: fmt.Sprintf("http://localhost:%s", envQKMHostPort),
	})

	conf2 := client.NewConfig(apiURL, nil)
	conf2.MetricsURL = apiMetricsURL
	apiClient := client.NewHTTPClient(httpClient, conf2)

	txSenderConfig.NonceMaxRecovery = maxRecoveryDefault
	return txsender.NewTxSender(txSenderConfig, []sarama2.ConsumerGroup{sarama.GlobalConsumerGroup()}, sarama.GlobalSyncProducer(),
		qkmClient, apiClient, ec, redisCli)
}

func testBackOff() backoff.BackOff {
	return backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), maxRecoveryDefault)
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
