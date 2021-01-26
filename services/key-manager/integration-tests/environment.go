package integrationtests

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"

	"github.com/hashicorp/vault/api"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/hashicorp"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/spf13/pflag"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/config"
	hashicorpDocker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/hashicorp"
)

const vaultContainerID = "vault-key-manager"
const networkName = "key-manager"
const vaultTokenFilePrefix = "orchestrate_vault_token_"
const localhostPath = "http://localhost:"
const pluginFilename = "orchestrate-hashicorp-vault-plugin"
const pluginVersion = "v0.0.9"

var envVaultHostPort string
var envHTTPPort string
var envMetricsPort string

type IntegrationEnvironment struct {
	ctx        context.Context
	logger     *log.Logger
	keyManager *app.App
	client     *docker.Client
	baseURL    string
	metricsURL string
	rootToken  string
}

func NewIntegrationEnvironment(ctx context.Context) (*IntegrationEnvironment, error) {
	logger := log.NewLogger().WithContext(ctx)

	host := os.Getenv("VAULT_HOST")
	if host == "" {
		host = "localhost"
	}

	rootToken := fmt.Sprintf("root_token_%v", strconv.Itoa(rand.IntnRange(0, 10000)))
	tokenFileName, err := generateTokenFile(rootToken)
	if err != nil {
		logger.WithError(err).Error("cannot generate vault token file")
		return nil, err
	}

	envHTTPPort = strconv.Itoa(rand.IntnRange(20000, 28080))
	envMetricsPort = strconv.Itoa(rand.IntnRange(30000, 38082))
	envVaultHostPort = strconv.Itoa(rand.IntnRange(10000, 15235))

	// Initialize environment flags
	flgs := pflag.NewFlagSet("key-manager-integration-test", pflag.ContinueOnError)
	httputils.MetricFlags(flgs)
	httputils.Flags(flgs)
	hashicorp.InitFlags(flgs)
	log.Level(flgs)
	args := []string{
		"--metrics-port=" + envMetricsPort,
		"--rest-port=" + envHTTPPort,
		"--vault-addr=http://" + host + ":" + envVaultHostPort,
		"--vault-token-file=" + tokenFileName,
		"--log-level=panic",
	}

	err = flgs.Parse(args)
	if err != nil {
		logger.WithError(err).Error("cannot parse environment flags")
		return nil, err
	}

	pluginPath, err := getPluginPath(logger)
	if err != nil {
		return nil, err
	}

	vaultContainer := hashicorpDocker.
		NewDefault().
		SetHostPort(envVaultHostPort).
		SetRootToken(rootToken).
		SetHost(host).
		SetPluginSourceDirectory(pluginPath)

	err = vaultContainer.DownloadPlugin(pluginFilename, pluginVersion)
	if err != nil {
		return nil, err
	}

	// Initialize environment container setup
	composition := &config.Composition{
		Containers: map[string]*config.Container{
			vaultContainerID: {
				HashicorpVault: vaultContainer,
			},
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
		baseURL:    localhostPath + envHTTPPort,
		metricsURL: localhostPath + envMetricsPort,
		rootToken:  rootToken,
	}, nil
}

func (env *IntegrationEnvironment) Start(ctx context.Context) error {
	err := env.client.CreateNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not create network")
		return err
	}

	// Start Hashicorp Vault
	err = env.client.Up(ctx, vaultContainerID, networkName)
	if err != nil {
		env.logger.WithError(err).Error("could not up vault container")
		return err
	}

	err = env.client.WaitTillIsReady(ctx, vaultContainerID, 10*time.Second)
	if err != nil {
		env.logger.WithError(err).Error("could not start vault")
		return err
	}

	// Enable orchestrate secret engine
	vaultClient, err := api.NewClient(hashicorp.ToVaultConfig(hashicorp.ConfigFromViper()))
	if err != nil {
		env.logger.WithError(err).Error("failed to instantiate vault client")
		return err
	}
	vaultClient.SetToken(env.rootToken)
	err = vaultClient.Sys().Mount("orchestrate", &api.MountInput{
		Type:        "plugin",
		Description: "Orchestrate Wallets",
		Config: api.MountConfigInput{
			ForceNoCache:              true,
			PassthroughRequestHeaders: []string{"X-Vault-Namespace"},
		},
		PluginName: pluginFilename,
	})
	if err != nil {
		env.logger.WithError(err).Error("failed to mount (enable) orchestrate vault plugin")
		return err
	}

	env.keyManager, err = keymanager.NewKeyManager(ctx, keymanager.NewConfig(viper.GetViper()))
	if err != nil {
		env.logger.WithError(err).Error("could not initialize key manager")
		return err
	}

	// Start key-manager app
	err = env.keyManager.Start(ctx)
	if err != nil {
		env.logger.WithError(err).Error("could not start key-manager")
		return err
	}
	integrationtest.WaitForServiceLive(ctx, fmt.Sprintf("%s/live", env.metricsURL), "key-manager", 15*time.Second)

	return nil
}

func (env *IntegrationEnvironment) Teardown(ctx context.Context) {
	env.logger.Info("tearing test suite down")

	if env.keyManager.IsReady() {
		err := env.keyManager.Stop(ctx)
		if err != nil {
			env.logger.WithError(err).Error("could not stop key-manager")
		}
	}

	err := env.client.Down(ctx, vaultContainerID)
	if err != nil {
		env.logger.WithError(err).Errorf("could not down vault")
	}

	err = env.client.RemoveNetwork(ctx, networkName)
	if err != nil {
		env.logger.WithError(err).Errorf("could not remove network")
	}
}

func generateTokenFile(rootToken string) (string, error) {
	file, err := ioutil.TempFile("", vaultTokenFilePrefix)
	if err != nil {
		return "", err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.WriteString(rootToken)
	if err != nil {
		return "", err
	}

	err = w.Flush()
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func getPluginPath(logger *log.Logger) (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		logger.WithError(err).Error("failed to get the current directory path")
		return "", err
	}

	return fmt.Sprintf("%s/plugins", currDir), nil
}
