package quorumkeymanager

import (
	"context"
	"fmt"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

const defaultDBHost = "localhost"
const defaultDBPort = "5432"

// nolint
type QuorumKeyManagerMigrate struct{}

type ConfigMigrate struct {
	Image  string
	DBPort string
	DBHost string
}

func NewDefaultMigrate() *ConfigMigrate {
	return &ConfigMigrate{
		Image:  defaultHashicorpVaultImage,
		DBPort: defaultDBPort,
		DBHost: defaultDBHost,
	}
}

func (cfg *ConfigMigrate) SetDBPort(port string) *ConfigMigrate {
	cfg.DBPort = port
	return cfg
}

func (cfg *ConfigMigrate) SetDBHost(host string) *ConfigMigrate {
	cfg.DBHost = host
	return cfg
}

func (q *QuorumKeyManagerMigrate) GenerateContainerConfig(_ context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*ConfigMigrate)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	containerCfg := &dockercontainer.Config{
		Image: cfg.Image,
		Env: []string{
			fmt.Sprintf("DB_PORT=%v", cfg.DBPort),
			fmt.Sprintf("DB_HOST=%v", cfg.DBHost),
		},
		Cmd: []string{"migrate", "up"},
	}

	return containerCfg, &dockercontainer.HostConfig{
		RestartPolicy: dockercontainer.RestartPolicy{
			Name: "no",
		},
	}, nil, nil
}

func (q *QuorumKeyManagerMigrate) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	return nil
}
