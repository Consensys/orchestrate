package integrationtests

import (
	"context"
	"net/http"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/postgres/migrations"
)

const postgresContainerID = "postgres-chain-registry"

type IntegrationEnvironment struct {
	client *docker.Client
	pgmngr postgres.Manager
	ctx    context.Context
}

func NewIntegrationEnvironment(ctx context.Context) *IntegrationEnvironment {
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			postgresContainerID: {Postgres: (&config.Postgres{}).SetDefault()},
		},
	}

	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		panic(err)
	}

	initChains := []string{`{"name":"geth","urls":["http://geth:8545"],"listenerStartingBlock":"0"}`,
		`{"name":"besu","urls":["http://validator2:8545"],"listenerStartingBlock":"0"}`,
		`{"name":"quorum","urls":["http://172.16.239.11:8545"],"listenerStartingBlock":"0","privateTxManagers":[{"url":"http://tessera1:9080","type":"Tessera"}]}`}
	viper.SetDefault(chainregistry.InitViperKey, initChains)

	return &IntegrationEnvironment{
		client: dockerClient,
		pgmngr: postgres.NewManager(),
		ctx:    ctx,
	}
}

func (env *IntegrationEnvironment) Start() error {
	// Start postgres database
	err := env.client.Up(context.Background(), postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up postgres")
		return err
	}
	// Wait 10 seconds for postgres database to be up
	time.Sleep(10 * time.Second)

	// Migrate database
	err = env.migrate()
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not migrate postgres")
		return err
	}

	// Start chain registry API
	err = chainregistry.Start(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not start chain-registry")
		return err
	}

	env.waitForService()
	log.WithoutContext().Infof("chain-registry ready")
	return nil
}

func (env *IntegrationEnvironment) Teardown() {
	log.WithoutContext().Infof("tearing test suite down")

	err := env.client.Down(env.ctx, postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not down postgres")
		return
	}

	err = chainregistry.Stop(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not stop chain-registry")
		return
	}
}

func (env *IntegrationEnvironment) migrate() error {
	db := env.pgmngr.Connect(env.ctx, postgres.NewOptions(viper.GetViper()))

	_, _, err := migrations.Run(db, "init")
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
