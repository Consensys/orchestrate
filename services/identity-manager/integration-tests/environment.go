package integrationtests

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager"
	keyManagerClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	"gopkg.in/h2non/gock.v1"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	postgresDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/postgres/migrations"
)

const postgresContainerID = "postgres-identity-manager"
const KeyManagerURL = "http://key-manager:8081"
const KeyManagerMetricsURL = "http://key-manager:8082"
const networkName = "identity-manager"

var envPGHostPort string
var envHTTPPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	ctx        context.Context
	logger     log.Logger
	app        *app.App
	client     *docker.Client
	pgmngr     postgres.Manager
	baseURL    string
	metricsURL string
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.FromContext(ctx)
	envPGHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))
	envHTTPPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))

	// Initialize environment flags
	flgs := pflag.NewFlagSet("identity-manager-integration-test", pflag.ContinueOnError)
	postgres.DBPort(flgs)
	httputils.MetricFlags(flgs)
	httputils.Flags(flgs)
	keyManagerClient.Flags(flgs)
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

	// Initialize environment container setup
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			postgresContainerID: {Postgres: postgresDocker.NewDefault().SetHostPort(envPGHostPort)},
		},
	}

	// Docker client
	dockerClient, err := docker.NewClient(composition)
	if err != nil {
		logger.WithError(err).Error("cannot initialize new environment")
		return nil, err
	}

	return &IntegrationEnvironment{
		ctx:        ctx,
		logger:     logger,
		client:     dockerClient,
		pgmngr:     postgres.NewManager(),
		baseURL:    "http://localhost:" + envHTTPPort,
		metricsURL: "http://localhost:" + envMetricsPort,
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	err := env.client.CreateNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not create network")
		return err
	}

	// Start postgres Database
	err = env.client.Up(ctx, postgresContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up postgres")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, postgresContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start postgres")
		return err
	}

	// Run postgres migrations
	err = env.migrate(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not migrate postgres")
		return err
	}

	env.app, err = newIdentityManager(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could initialize transaction scheduler")
		return err
	}

	err = env.app.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start identity-manager")
		return err
	}

	integrationtest.WaitForServiceLive(
		ctx,
		fmt.Sprintf("%s/live", env.metricsURL),
		"identity-manager",
		15*time.Second,
	)

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Infof("tearing test suite down")

	err := env.app.Stop(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not stop identity-manager")
	}

	err = env.client.Down(ctx, postgresContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down postgres")
	}

	err = env.client.RemoveNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Errorf("could not remove network")
	}
}

func (env *IntegrationEnvironment) migrate(ctx context.Context) error {
	// Set Database connection
	opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
	if err != nil {
		return err
	}

	db := env.pgmngr.Connect(ctx, opts)
	env.logger.Debugf("initializing Database migrations...")
	_, _, err = migrations.Run(db, "init")
	if err != nil {
		return err
	}

	env.logger.Debugf("running Database migrations...")
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

func newIdentityManager(
	ctx context.Context,
) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)

	// We mock the calls to the chain registry
	conf := keyManagerClient.NewConfig(KeyManagerURL, nil)
	conf.MetricsURL = KeyManagerMetricsURL
	httpClient := httputils.NewClient(httputils.NewDefaultConfig())
	gock.InterceptClient(httpClient)
	client := keyManagerClient.NewHTTPClient(httpClient, conf)

	pgmngr := postgres.GetManager()
	cfg := identitymanager.NewConfig(viper.GetViper())

	return identitymanager.NewIdentityManager(
		cfg,
		pgmngr,
		authjwt.GlobalChecker(), authkey.GlobalChecker(),
		client,
	)
}
