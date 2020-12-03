package compose

import (
	"context"
	"fmt"
	goreflect "reflect"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/ganache"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/hashicorp"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/kafka"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/zookeeper"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/config"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/reflect"
)

type Compose struct {
	reflect *reflect.Reflect
}

func New() *Compose {
	factory := &Compose{
		reflect: reflect.New(),
	}

	factory.reflect.AddGenerator(goreflect.TypeOf(&postgres.Config{}), &postgres.Postgres{})
	factory.reflect.AddGenerator(goreflect.TypeOf(&zookeeper.Config{}), &zookeeper.Zookeeper{})
	factory.reflect.AddGenerator(goreflect.TypeOf(&kafka.Config{}), &kafka.Kafka{})
	factory.reflect.AddGenerator(goreflect.TypeOf(&hashicorp.Config{}), &hashicorp.Vault{})
	factory.reflect.AddGenerator(goreflect.TypeOf(&ganache.Config{}), &ganache.Ganache{})

	return factory
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

func (gen *Compose) WaitForService(ctx context.Context, configuration interface{}, timeout time.Duration) error {
	cfg, ok := configuration.(*config.Container)
	if !ok {
		return fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	field, err := cfg.Field()
	if err != nil {
		return err
	}

	return gen.reflect.WaitForService(ctx, field, timeout)
}
