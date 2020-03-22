package postgres

import (
	"context"
	"fmt"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
)

type Postgres struct{}

func (g *Postgres) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*config.Postgres)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%v", cfg.Password),
		},
		ExposedPorts: nat.PortSet{
			"5432/tcp": struct{}{},
		},
	}

	hostConfig := &dockercontainer.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		},
	}

	return containerCfg, hostConfig, nil, nil
}
