package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const DefaultPostgresImage = "postgres:13.1-alpine"

const defaultPassword = "postgres"
const defaultHostPort = "5432"

type Postgres struct{}

type Config struct {
	Image    string
	Port     string
	Password string
}

func NewDefault() *Config {
	cfg := &Config{
		Image:    DefaultPostgresImage,
		Port:     defaultHostPort,
		Password: defaultPassword,
	}

	return cfg
}

func (c *Config) SetHostPort(port string) *Config {
	c.Port = port
	return c
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

	hostConfig := &dockercontainer.HostConfig{}
	if cfg.Port != "" {
		hostConfig.PortBindings = nat.PortMap{
			"5432/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: cfg.Port}},
		}
	}

	return containerCfg, hostConfig, nil, nil
}

func (g *Postgres) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pgCfg, _ := postgres.NewConfig(viper.GetViper()).PGOptions()
	db := pg.Connect(pgCfg)
	defer db.Close()

	retryT := time.NewTicker(time.Second)
	defer retryT.Stop()

	var cerr error
waitForServiceLoop:
	for {
		select {
		case <-rctx.Done():
			cerr = rctx.Err()
			break waitForServiceLoop
		case <-retryT.C:
			_, err := db.Exec("SELECT 1")
			if err != nil {
				log.WithContext(rctx).WithError(err).Warnf("waiting for PostgreSQL service to start")
			} else {
				log.WithContext(rctx).Info("PostgreSQL container service is ready")
				break waitForServiceLoop
			}
		}
	}

	return cerr
}
