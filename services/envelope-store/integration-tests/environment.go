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
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/postgres/migrations"
)

const postgresContainerID = "postgres-envelope-store"

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

	return &IntegrationEnvironment{
		client: dockerClient,
		pgmngr: postgres.NewManager(),
		ctx:    ctx,
	}
}

func (env *IntegrationEnvironment) Start() {
	// Start postgres database
	err := env.client.Up(context.Background(), postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not up postgres")
		return
	}
	// Wait 10 seconds for postgres database to be up
	time.Sleep(10 * time.Second)

	// Migrate database
	err = env.migrate()
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not migrate postgres")
		return
	}

	// Start envelope store
	err = envelopestore.Start(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Fatalf("could not start envelope-store")
		return
	}

	env.waitForService()
	log.WithoutContext().Infof("envelope-store ready")
}

func (env *IntegrationEnvironment) Teardown() {
	log.WithoutContext().Infof("tearing test suite down")
	err := envelopestore.Stop(env.ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not stop envelope-store")
		return
	}

	err = env.client.Down(env.ctx, postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not down postgres")
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
