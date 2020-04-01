package integrationtests

import (
	"context"

	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"

	"net/http"
	"time"

	"github.com/docker/docker/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test-utils/containers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test-utils/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
)

type IntegrationEnvironment struct {
	client            *client.Client
	postgresContainer *containers.PostgresContainer
}

const (
	component  = "contract-registry-integration"
	metricsURL = "http://localhost:8082"
)

func NewIntegrationEnvironment() *IntegrationEnvironment {
	cli := utils.NewDockerClient()
	return &IntegrationEnvironment{
		client:            utils.NewDockerClient(),
		postgresContainer: containers.NewPostgresContainer(cli, component),
	}
}

func (env *IntegrationEnvironment) Start() {
	env.postgresContainer.Start(context.Background(), env.client)
	env.migrate()

	go contractregistry.StartService(context.Background(), "postgres")
	env.waitForService()
}

func (env *IntegrationEnvironment) Teardown() {
	contractregistry.StopService(context.Background())
	env.postgresContainer.Remove(context.Background(), env.client)
}

func (env *IntegrationEnvironment) migrate() {
	db := env.postgresContainer.GetDB()

	_, _, err := migrations.Run(db, "init")
	if err != nil {
		panic(err)
	}

	_, _, err = migrations.Run(db, "up")
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}
}

func (env *IntegrationEnvironment) waitForService() {
	for {
		resp, _ := http.Get(metricsURL + "/ready")
		if resp != nil && resp.StatusCode == 200 {
			return
		}

		time.Sleep(2 * time.Second)
	}
}
