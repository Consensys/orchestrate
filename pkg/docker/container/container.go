package container

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

//go:generate mockgen -source=container.go -destination=mock/mock.go -package=mock

type ConfigGenerator interface {
	GenerateContainerConfig(ctx context.Context, configuration interface{}) (*container.Config, *container.HostConfig, *network.NetworkingConfig, error)
}
