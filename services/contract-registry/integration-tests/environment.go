package integrationtests

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	postgresDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
	"k8s.io/apimachinery/pkg/util/rand"
)

const postgresContainerID = "postgres-contract-registry"

var envPGHostPort string
var envHTTPPort string
var envGRPCPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	client      *docker.Client
	pgmngr      postgres.Manager
	envContract *abi.Contract
	logger      log.Logger
	baseHTTP    string
	metricsURL  string
	baseGRPC    string
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envPGHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envHTTPPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envGRPCPort = strconv.Itoa(rand.IntnRange(30000, 38080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(40000, 48082))

	// Initialize environment flags
	flgs := pflag.NewFlagSet("transaction-scheduler-integration-test", pflag.ContinueOnError)
	postgres.DBPort(flgs)
	httputils.MetricFlags(flgs)
	httputils.Flags(flgs)
	grpc.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--rest-port=" + envHTTPPort,
		"--grpc-port=" + envGRPCPort,
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

	envContract := testutils.FakeContract()
	var re = regexp.MustCompile(`\s+`)
	contractAbi := fmt.Sprintf("%s:%s:%s:%s", envContract.Id.Name, re.ReplaceAllString(envContract.Abi, ""), envContract.Bytecode, envContract.DeployedBytecode)
	viper.SetDefault(contractregistry.ABIViperKey, contractAbi)

	client, err := docker.NewClient(composition)
	if err != nil {
		return nil, err
	}

	return &IntegrationEnvironment{
		client:      client,
		pgmngr:      postgres.NewManager(),
		envContract: envContract,
		logger:      logger,
		baseHTTP:    "http://localhost:" + envHTTPPort,
		metricsURL:  "http://localhost:" + envMetricsPort,
		baseGRPC:    "localhost:" + envGRPCPort,
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

	// Start contract registry
	err = contractregistry.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start contract-registry")
		return err
	}

	integrationtest.WaitForServiceLive(ctx,
		fmt.Sprintf("http://localhost:%s/live", envMetricsPort),
		"contract-registry",
		10*time.Second)

	env.logger.Infof("contract-registry ready")

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Infof("tearing test suite down")
	err := contractregistry.Stop(ctx)
	if err != nil {
		env.logger.Errorf("could not stop contract-registry")
	}

	err = env.client.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down postgres")
	}
}

func (env *IntegrationEnvironment) migrate(ctx context.Context) error {
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
