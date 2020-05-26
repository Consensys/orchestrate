package integrationtests

import (
	"context"
	"net/http"
	"os"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	transactionscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler"
	"gopkg.in/h2non/gock.v1"

	"github.com/containous/traefik/v2/pkg/log"
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
const kafkaHostEnvName = "KAFKA_HOST"
const txCrafterTopic = "transaction-scheduler-integration-crafter-topic"
const ChainRegistryURL = "http://chain-registry:8081"

type IntegrationEnvironment struct {
	app       *app.App
	client    *docker.Client
	pgmngr    postgres.Manager
	ctx       context.Context
	networkID string
}

func NewIntegrationEnvironment(ctx context.Context) *IntegrationEnvironment {
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			postgresContainerID:  {Postgres: (&postgresDocker.Config{}).SetDefault()},
			zookeeperContainerID: {Zookeeper: (&zookeeper.Config{}).SetDefault()},
			kafkaContainerID: {Kafka: (&kafkaDocker.Config{
				ZookeeperHostname:     zookeeperContainerID,
				KafkaInternalHostname: kafkaContainerID,
				KafkaExternalHostname: os.Getenv(kafkaHostEnvName),
			}).SetDefault()},
		},
	}

	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		panic(err)
	}

	return &IntegrationEnvironment{
		client: dockerClient,
		pgmngr: postgres.NewManager(),
		ctx:    ctx,
	}
}

func (env *IntegrationEnvironment) Start() {
	// Start kafka
	networkID, err := env.client.CreateNetwork(env.ctx, "transaction-scheduler")
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not create network")
		return
	}
	env.networkID = networkID

	err = env.client.Up(env.ctx, zookeeperContainerID, networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up kafka")
		return
	}
	time.Sleep(2 * time.Second)

	err = env.client.Up(env.ctx, kafkaContainerID, networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up kafka")
		return
	}
	time.Sleep(2 * time.Second)

	// Start postgres database
	err = env.client.Up(env.ctx, postgresContainerID, networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up postgres")
		return
	}
	time.Sleep(2 * time.Second)

	// Migrate database
	err = env.migrate()
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not migrate postgres")
		return
	}

	// Start transaction scheduler
	txSchedulerApp, err := initService(env.ctx)
	if err != nil {
		panic(err)
	}
	env.app = txSchedulerApp
	err = env.app.Start(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not start transaction-scheduler")
		return
	}

	env.waitForService()
}

func (env *IntegrationEnvironment) Teardown() {
	log.WithoutContext().Infof("tearing test suite down")
	err := env.app.Stop(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not stop transaction-scheduler")
		return
	}

	err = env.client.Down(env.ctx, postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not down postgres")
		return
	}

	err = env.client.Down(env.ctx, kafkaContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not down kafka")
		return
	}

	err = env.client.Down(env.ctx, zookeeperContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not down zookeeper")
		return
	}

	err = env.client.RemoveNetwork(env.ctx, env.networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not remove network")
		return
	}
}

func (env *IntegrationEnvironment) migrate() error {
	// Set database connection
	opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
	if err != nil {
		return err
	}

	db := env.pgmngr.Connect(env.ctx, opts)

	_, _, err = migrations.Run(db, "init")
	if err != nil {
		return err
	}

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

func (env *IntegrationEnvironment) waitForService() {
	for {
		resp, _ := http.Get("http://localhost:8082/ready")
		if resp != nil && resp.StatusCode == 200 {
			return
		}

		time.Sleep(2 * time.Second)
	}
}

func initService(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	sarama.InitSyncProducer(ctx)

	// We mock the calls to the chain registry
	conf := client.NewConfig("http://chain-registry:8081")
	httpClient := httputils.NewClient()
	gock.InterceptClient(httpClient)
	chainRegistryClient := client.NewHTTPClient(httpClient, conf)

	return transactionscheduler.New(
		transactionscheduler.NewConfig(viper.GetViper()),
		postgres.GetManager(),
		authjwt.GlobalChecker(), authkey.GlobalChecker(),
		chainRegistryClient,
		sarama.GlobalSyncProducer(),
		txCrafterTopic,
	)
}
