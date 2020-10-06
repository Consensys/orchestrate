package integrationtests

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	postgresDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/postgres/migrations"
	"k8s.io/apimachinery/pkg/util/rand"
)

const postgresContainerID = "postgres-chain-registry"

var envPGHostPort string
var envHTTPPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	client     *docker.Client
	pgmngr     postgres.Manager
	logger     log.Logger
	baseURL    string
	metricsURL string
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envPGHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envHTTPPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))

	// Initialize environment flags
	flgs := pflag.NewFlagSet("transaction-scheduler-integration-test", pflag.ContinueOnError)
	postgres.DBPort(flgs)
	httputils.MetricFlags(flgs)
	httputils.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--rest-port=" + envHTTPPort,
		"--db-port=" + envPGHostPort,
	}

	err := flgs.Parse(args)
	if err != nil {
		logger.WithError(err).Error("cannot parse environment flags")
		return nil, err
	}

	composition := &config.Composition{
		Containers: map[string]*config.Container{
			postgresContainerID: {Postgres: postgresDocker.NewDefault().SetHostPort(envPGHostPort)},
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
		client:     dockerClient,
		pgmngr:     postgres.NewManager(),
		logger:     logger,
		baseURL:    "http://localhost:" + envHTTPPort,
		metricsURL: "http://localhost:" + envMetricsPort,
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	// Start postgres database
	err := env.client.Up(ctx, postgresContainerID, "")
	if err != nil {
		env.logger.WithError(err).Error("could not up postgres")
		return err
	}
	err = env.client.WaitTillIsReady(ctx, postgresContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start postgres")
		return err
	}

	// Migrate database
	err = env.migrate(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not migrate postgres")
		return err
	}

	// Start chain registry API
	err = chainregistry.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start chain-registry")
		return err
	}

	integrationtest.WaitForServiceLive(ctx,
		fmt.Sprintf("http://localhost:%s/live", envMetricsPort),
		"chain-registry",
		10*time.Second)

	log.WithoutContext().Infof("chain-registry ready")
	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	log.WithoutContext().Infof("tearing test suite down")

	err := env.client.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down postgres")
		return
	}

	err = chainregistry.Stop(ctx)
	if err != nil {
		env.logger.WithError(err).Errorf("could not stop chain-registry")
		return
	}
}

func (env *IntegrationEnvironment) migrate(ctx context.Context) error {
	// Set database connection
	opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
	if err != nil {
		return err
	}

	db := env.pgmngr.Connect(ctx, opts)

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
