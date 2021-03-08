package ganache

import (
	"context"
	"fmt"
	"time"

	ethclient "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/rpc"
	log "github.com/sirupsen/logrus"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const defaultGanacheImage = "trufflesuite/ganache-cli:v6.12.1"
const defaultHostPort = "8545"
const defaultHost = "localhost"

type Ganache struct{}

type Config struct {
	Image string
	Port  string
	Host  string
}

func NewDefault() *Config {
	return &Config{
		Image: defaultGanacheImage,
		Port:  defaultHostPort,
		Host:  defaultHost,
	}
}

func (cfg *Config) SetHostPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetHost(host string) *Config {
	cfg.Host = host
	return cfg
}

func (*Ganache) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		ExposedPorts: nat.PortSet{
			"8545/tcp": struct{}{},
		},
		Cmd: []string{"ganache-cli", "--mnemonic", "surge arm pulse bus piano poet thrive erase angry dwarf cargo vanish", "--blockTime", "1"},
	}

	hostConfig := &dockercontainer.HostConfig{}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"8545/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (*Ganache) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	retryT := time.NewTicker(2 * time.Second)
	defer retryT.Stop()

	ethclient.Init(ctx)
	client := ethclient.GlobalClient()

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			_, err := client.Network(ctx, fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port))
			if err != nil {
				log.WithContext(rctx).WithError(err).Warnf("waiting for Ganache service to start")
			} else {
				log.WithContext(rctx).Info("ganache container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
