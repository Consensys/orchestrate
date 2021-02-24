package config

import (
	"github.com/ConsenSys/orchestrate/pkg/docker/container/ganache"
	"github.com/ConsenSys/orchestrate/pkg/docker/container/hashicorp"
	"github.com/ConsenSys/orchestrate/pkg/docker/container/kafka"
	"github.com/ConsenSys/orchestrate/pkg/docker/container/postgres"
	"github.com/ConsenSys/orchestrate/pkg/docker/container/zookeeper"
	"github.com/ConsenSys/orchestrate/pkg/utils"
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
