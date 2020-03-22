package config

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type Composition struct {
	Containers map[string]*Container
}

type Container struct {
	Postgres *Postgres
}

func (c *Container) Field() (interface{}, error) {
	return utils.ExtractField(c)
}

const DefaultPostgresImage = "postgres:10.12-alpine"

type Postgres struct {
	Image    string
	Port     string
	Password string
}

func (p *Postgres) SetDefault() *Postgres {
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
