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
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
)

const postgresContainerID = "postgres-contract-registry"

type IntegrationEnvironment struct {
	client *docker.Client
	pgmngr postgres.Manager
}

func NewIntegrationEnvironment() *IntegrationEnvironment {
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			postgresContainerID: {Postgres: (&config.Postgres{}).SetDefault()},
		},
	}

	client, err := docker.NewClient(composition)
	if err != nil {
		panic(err)
	}

	return &IntegrationEnvironment{
		client: client,
		pgmngr: postgres.NewManager(),
	}
}

func (env *IntegrationEnvironment) Start() {
	// Start postgres database
	err := env.client.Up(context.Background(), postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up postgres")
	}
	// Wait 10 seconds for postgres database to be up
	time.Sleep(10 * time.Second)

	// Migrate database
	err = env.migrate()
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not migrate postgres")
	}

	// Start contract registry
	err = contractregistry.Start(context.Background())
	if err != nil {
		// TODO: we should probably not panic here
		log.WithoutContext().WithError(err).Fatalf("could not start contract-registry")
	}

	env.waitForService()
	log.WithoutContext().Infof("contract-registry ready")
}

func (env *IntegrationEnvironment) Teardown() {
	log.WithoutContext().Infof("tearing test suite down")
	err := contractregistry.Stop(context.Background())
	if err != nil {
		// TODO: we should probably not panic here
		log.WithoutContext().WithError(err).Errorf("could not stop contract-registry")
	}

	err = env.client.Down(context.Background(), postgresContainerID)
	if err != nil {
		// TODO: we should probably not panic here
		log.WithoutContext().WithError(err).Errorf("could not down postgres")
	}
}

func (env *IntegrationEnvironment) migrate() error {
	db := env.pgmngr.Connect(context.Background(), postgres.NewOptions(viper.GetViper()))

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
