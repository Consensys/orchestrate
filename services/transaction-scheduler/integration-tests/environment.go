package integrationtests

import (
	"context"
	"net/http"
	"time"

	kafkaDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/kafka"
	postgresDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/zookeeper"
	transactionscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
)

const postgresContainerID = "postgres-transaction-scheduler"
const kafkaContainerID = "kafka-transaction-scheduler"
const zookeeperContainerID = "zookeeper-transaction-scheduler"

type IntegrationEnvironment struct {
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
				ZookeeperHostname: zookeeperContainerID,
				KafkaHostname:     kafkaContainerID,
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
	ctx := context.Background()

	// Start kafka
	networkID, err := env.client.CreateNetwork(ctx, "transaction-scheduler")
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not create network")
		return
	}
	env.networkID = networkID

	err = env.client.Up(ctx, zookeeperContainerID, networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up kafka")
		return
	}
	time.Sleep(10 * time.Second)
	err = env.client.Up(ctx, kafkaContainerID, networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up kafka")
		return
	}
	time.Sleep(20 * time.Second)

	// Start postgres database
	err = env.client.Up(ctx, postgresContainerID, networkID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up postgres")
		return
	}
	time.Sleep(10 * time.Second)

	// Migrate database
	err = env.migrate()
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not migrate postgres")
		return
	}

	// Start transaction scheduler
	err = transactionscheduler.Start(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not start transaction-scheduler")
		return
	}

	env.waitForService()
	log.WithoutContext().Infof("transaction-scheduler ready")
}

func (env *IntegrationEnvironment) Teardown() {
	log.WithoutContext().Infof("tearing test suite down")
	err := transactionscheduler.Stop(env.ctx)
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
