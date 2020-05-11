package zookeeper

import (
	"context"
	"fmt"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const DefaultZookeeperImage = "confluentinc/cp-zookeeper:5.3.0"

type Zookeeper struct{}
type Config struct {
	Image string
	Port  string
}

func (p *Config) SetDefault() *Config {
	if p.Image == "" {
		p.Image = DefaultZookeeperImage
	}

	if p.Port == "" {
		p.Port = "2181"
	}

	return p
}

func (k *Zookeeper) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			"ZOOKEEPER_CLIENT_PORT=32181",
			"ZOOKEEPER_TICK_TIME=2000",
			"ALLOW_ANONYMOUS_LOGIN=yes",
		},
		ExposedPorts: nat.PortSet{
			"2181/tcp": struct{}{},
		},
	}

	hostConfig := &dockercontainer.HostConfig{
		PortBindings: nat.PortMap{
			"2181/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		},
	}

	return containerCfg, hostConfig, nil, nil
}
