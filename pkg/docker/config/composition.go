package config

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/ganache"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/kafka"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/docker/container/zookeeper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type Composition struct {
	Containers map[string]*Container
}

type Container struct {
	Postgres       *postgres.Config
	Zookeeper      *zookeeper.Config
	Kafka          *kafka.Config
	HashicorpVault *hashicorp.Config
	Ganache        *ganache.Config
}

func (c *Container) Field() (interface{}, error) {
	return utils.ExtractField(c)
}
