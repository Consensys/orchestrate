package reflect

import (
	"context"
	"fmt"
	"reflect"
	"time"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container"
)

type Reflect struct {
	generators map[reflect.Type]container.DockerContainerFactory
}

func New() *Reflect {
	return &Reflect{
		generators: make(map[reflect.Type]container.DockerContainerFactory),
	}
}

func (gen *Reflect) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	generator, ok := gen.generators[reflect.TypeOf(configuration)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("no container config generator for configuration of type %T (consider adding one)", configuration)
	}

	return generator.GenerateContainerConfig(ctx, configuration)
}

func (gen *Reflect) WaitForService(configuration interface{}, timeout time.Duration) error {
	generator, ok := gen.generators[reflect.TypeOf(configuration)]
	if !ok {
		return fmt.Errorf("no container config generator for configuration of type %T (consider adding one)", configuration)
	}

	return generator.WaitForService(configuration, timeout)
}

func (gen *Reflect) AddGenerator(typ reflect.Type, generator container.DockerContainerFactory) {
	gen.generators[typ] = generator
}
