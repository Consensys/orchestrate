package hashicorp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	httputils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const defaultHashicorpVaultImage = "library/vault:1.1.1"
const defaultHostPort = "8200"
const defaultRootToken = "myRoot"
const defaultHost = "localhost"

type Vault struct{}

type Config struct {
	Image       string
	Port        string
	RootTokenID string
	Host        string
}

func NewDefault() *Config {
	return &Config{
		Image:       defaultHashicorpVaultImage,
		Port:        defaultHostPort,
		RootTokenID: defaultRootToken,
		Host:        defaultHost,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetRootTokenID(rootToken string) *Config {
	cfg.RootTokenID = rootToken
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	if host != "" {
		cfg.Host = host
	}

	return cfg
}

func (vault *Vault) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%v", cfg.RootTokenID),
			fmt.Sprintf("VAULT_API_ADDR=http://127.0.0.1:%v", cfg.Port),
		},
		ExposedPorts: nat.PortSet{
			"8200/tcp": struct{}{},
		},
		Tty: true,
	}

	hostConfig := &dockercontainer.HostConfig{
		CapAdd: []string{"IPC_LOCK"},
	}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"8200/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (vault *Vault) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	retryT := time.NewTicker(2 * time.Second)
	defer retryT.Stop()

	httpClient := httputils.NewClient(httputils.NewDefaultConfig())

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			resp, err := httpClient.Get(fmt.Sprintf("http://%v:%v/v1/sys/health", cfg.Host, cfg.Port))

			switch {
			case err != nil:
				log.WithContext(rctx).WithError(err).Warnf("waiting for Hashicorp Vault service to start")
			case resp.StatusCode != http.StatusOK:
				log.WithContext(rctx).WithField("status_code", resp.StatusCode).Warnf("waiting for Hashicorp Vault service to be ready")
			default:
				log.WithContext(rctx).Infof("Hashicorp Vault container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
