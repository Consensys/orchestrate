package compose

import (
	"context"
	"fmt"
	goreflect "reflect"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/reflect"
)

type Compose struct {
	reflect *reflect.Reflect
}

func New() *Compose {
	generator := &Compose{
		reflect: reflect.New(),
	}

	generator.reflect.AddGenerator(goreflect.TypeOf(&config.Postgres{}), &postgres.Postgres{})

	return generator
}

func (gen *Compose) GenerateContainerConfig(ctx context.Context, configuration interface{}) (*dockercontainer.Config, *dockercontainer.HostConfig, *network.NetworkingConfig, error) {
	cfg, ok := configuration.(*config.Container)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return nil, nil, nil, err
	}

	return gen.reflect.GenerateContainerConfig(ctx, field)
}
