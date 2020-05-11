package postgres

import (
	"context"
	"fmt"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

const DefaultPostgresImage = "postgres:10.12-alpine"

type Postgres struct{}
type Config struct {
	Image    string
	Port     string
	Password string
}

func (p *Config) SetDefault() *Config {
	if p.Image == "" {
		p.Image = DefaultPostgresImage
	}

	if p.Port == "" {
		p.Port = "5432"
	}

	if p.Password == "" {
		p.Password = "postgres"
	}

	return p
}

func (g *Postgres) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*Config)
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
