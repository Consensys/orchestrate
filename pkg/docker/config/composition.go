package config

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/kafka"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/docker/container/zookeeper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type Composition struct {
	Containers map[string]*Container
}

type Container struct {
	Postgres       *postgres.Config
	Zookeeper      *zookeeper.Config
	Kafka          *kafka.Config
	HashicorpVault *hashicorp.Config
}

func (c *Container) Field() (interface{}, error) {
	return utils.ExtractField(c)
}
