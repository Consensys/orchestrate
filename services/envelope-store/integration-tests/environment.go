package integrationtests

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/pflag"
	postgresDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
	"k8s.io/apimachinery/pkg/util/rand"

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
	client  *docker.Client
	pgmngr  postgres.Manager
	logger  log.Logger
	baseURL string
}

var envPGHostPort string
var envGRPCPort string
var envMetricsPort string

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	envPGHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envGRPCPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))
	logger := log.FromContext(ctx)

	// Initialize environment flags
	flgs := pflag.NewFlagSet("transaction-scheduler-integration-test", pflag.ContinueOnError)
	postgres.DBPort(flgs)
	httputils.MetricFlags(flgs)
	grpc.Flags(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
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

	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		return nil, err
	}

	return &IntegrationEnvironment{
		client:  dockerClient,
		pgmngr:  postgres.NewManager(),
		logger:  logger,
		baseURL: "localhost:" + envGRPCPort,
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

	// Start envelope store
	err = envelopestore.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start envelope-store")
		return err
	}

	integrationtest.WaitForServiceReady(ctx,
		fmt.Sprintf("http://localhost:%s/ready", envMetricsPort),
		"envelope-store",
		10*time.Second)

	env.logger.Infof("envelope-store ready")

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	log.WithoutContext().Infof("tearing test suite down")
	err := envelopestore.Stop(ctx)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not stop envelope-store")
		return
	}

	err = env.client.Down(ctx, postgresContainerID)
	if err != nil {
		log.WithoutContext().WithError(err).Errorf("could not down postgres")
		return
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
