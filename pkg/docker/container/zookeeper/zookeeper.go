package zookeeper

import (
	"context"
	"fmt"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const DefaultZookeeperImage = "confluentinc/cp-zookeeper:5.3.0"
const defaultHostPort = ""
const DefaultZookeeperClientPort = "32181"

type Zookeeper struct{}
type Config struct {
	Image string
	Port  string
}

func NewDefault() *Config {
	return &Config{
		Image: DefaultZookeeperImage,
		Port:  defaultHostPort,
	}
}

func (c *Config) SetHostPort(port string) *Config {
	c.Port = port
	return c
}

func (k *Zookeeper) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			"ZOOKEEPER_CLIENT_PORT=" + DefaultZookeeperClientPort,
			"ZOOKEEPER_TICK_TIME=2000",
			"ALLOW_ANONYMOUS_LOGIN=yes",
		},
		ExposedPorts: nat.PortSet{
			"2181/tcp": struct{}{},
		},
	}

	hostConfig := &dockercontainer.HostConfig{}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"2181/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (k *Zookeeper) WaitForService(configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}
	return nil
}
